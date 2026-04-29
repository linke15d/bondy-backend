// Package admin 后台管理业务逻辑层
package admin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// AdminContentService 内容管理业务逻辑
// 负责系统预设标签和姿势的维护
type AdminContentService struct {
	db *gorm.DB
}

// NewAdminContentService 创建 AdminContentService 实例
func NewAdminContentService(db *gorm.DB) *AdminContentService {
	return &AdminContentService{db: db}
}

// CreateTagInput 创建系统标签请求参数
type CreateTagInput struct {
	// Name 标签名称
	Name string `json:"name" binding:"required,max=30" example:"酒店"`

	// Type 标签类型：LOCATION 或 ACTIVITY
	Type string `json:"type" binding:"required,oneof=LOCATION ACTIVITY" example:"LOCATION"`

	IconBase64 *string `json:"icon_base64"`
}

// ListSystemTags 获取系统预设标签列表
func (s *AdminContentService) ListSystemTags(tagType string) ([]model.Tag, error) {
	var tags []model.Tag
	query := s.db.Where("is_system = true")
	if tagType != "" {
		query = query.Where("type = ?", tagType)
	}
	err := query.Order("type ASC, name ASC").Find(&tags).Error
	return tags, err
}

// CategoryNameInput 单个语言的分类名称
type CategoryNameInput struct {
	// LanguageCode 语言代码，如 zh-CN / en
	LanguageCode string `json:"language_code" binding:"required" example:"zh-CN"`
	// Name 该语言下的名称
	Name string `json:"name" binding:"required,max=50" example:"经典"`
}

// CreateCategoryInput 创建分类请求参数
type CreateCategoryInput struct {
	// Names 各语言名称列表，至少传一种语言
	Names []CategoryNameInput `json:"names" binding:"required,min=1"`

	// SortOrder 排序值，数字越小越靠前
	SortOrder int `json:"sort_order" example:"1"`
}

// UpdateCategoryInput 更新分类请求参数
type UpdateCategoryInput struct {
	// Names 更新各语言名称，传入的语言会覆盖，未传的语言保持不变
	Names []CategoryNameInput `json:"names"`

	// SortOrder 修改排序
	SortOrder *int `json:"sort_order"`

	// IsActive 修改启用状态
	IsActive *bool `json:"is_active"`
}

// CreatePositionCategory 创建姿势分类
func (s *AdminContentService) CreatePositionCategory(input CreateCategoryInput) (*model.PositionCategory, error) {
	// 检查是否有重复的名称（任意语言下名称相同都不允许）
	for _, n := range input.Names {
		var count int64
		s.db.Model(&model.PositionCategoryName{}).
			Where("language_code = ? AND name = ?", n.LanguageCode, n.Name).
			Count(&count)
		if count > 0 {
			return nil, fmt.Errorf("「%s」名称已存在，请勿重复添加", n.Name)
		}
	}

	// 取中文名作为 DefaultName，没有中文则取第一个
	defaultName := input.Names[0].Name
	for _, n := range input.Names {
		if n.LanguageCode == "zh-CN" {
			defaultName = n.Name
			break
		}
	}

	category := &model.PositionCategory{
		DefaultName: defaultName,
		SortOrder:   input.SortOrder,
		IsActive:    true,
	}

	// 开启事务
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(category).Error; err != nil {
			return errors.New("创建分类失败")
		}

		// 批量插入各语言名称
		for _, n := range input.Names {
			name := model.PositionCategoryName{
				CategoryID:   category.ID,
				LanguageCode: n.LanguageCode,
				Name:         n.Name,
			}
			if err := tx.Create(&name).Error; err != nil {
				return errors.New("保存语言名称失败")
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 填充 Names
	s.db.Where("category_id = ?", category.ID).Find(&category.Names)
	return category, nil
}

// ListPositionCategories 获取分类列表（含所有语言名称）
func (s *AdminContentService) ListPositionCategories() ([]model.PositionCategory, error) {
	var categories []model.PositionCategory
	err := s.db.Order("sort_order ASC, created_at ASC").Find(&categories).Error
	if err != nil {
		return nil, errors.New("获取分类列表失败")
	}

	// 填充每个分类的多语言名称
	for i := range categories {
		s.db.Where("category_id = ?", categories[i].ID).Find(&categories[i].Names)
	}

	return categories, nil
}

// UpdatePositionCategory 更新分类
func (s *AdminContentService) UpdatePositionCategory(id string, input UpdateCategoryInput) (*model.PositionCategory, error) {
	var category model.PositionCategory
	if err := s.db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, errors.New("分类不存在")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if input.SortOrder != nil {
			category.SortOrder = *input.SortOrder
		}
		if input.IsActive != nil {
			category.IsActive = *input.IsActive
		}
		// 如果更新了 zh-CN 的名称，同步更新 DefaultName
		for _, n := range input.Names {
			if n.LanguageCode == "zh-CN" {
				category.DefaultName = n.Name
				break
			}
		}
		if err := tx.Save(&category).Error; err != nil {
			return errors.New("更新失败")
		}
		// 更新语言名称（upsert）
		for _, n := range input.Names {
			var existing model.PositionCategoryName
			err := tx.Where("category_id = ? AND language_code = ?", id, n.LanguageCode).
				First(&existing).Error
			if err == gorm.ErrRecordNotFound {
				newName := model.PositionCategoryName{
					CategoryID:   id,
					LanguageCode: n.LanguageCode,
					Name:         n.Name,
				}
				if err := tx.Create(&newName).Error; err != nil {
					return errors.New("保存语言名称失败")
				}
			} else {
				if err := tx.Model(&existing).Update("name", n.Name).Error; err != nil {
					return errors.New("更新语言名称失败")
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.db.Where("category_id = ?", id).Find(&category.Names)
	return &category, nil
}

// DeletePositionCategory 删除分类
// 删除前检查是否有姿势关联此分类，同时级联删除多语言名称
func (s *AdminContentService) DeletePositionCategory(id string) error {
	// 检查是否有姿势关联此分类
	var count int64
	s.db.Model(&model.Position{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该分类下还有姿势，请先删除或移动姿势后再删除分类")
	}

	// 开启事务，先删除多语言名称，再删除分类
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的多语言名称
		if err := tx.Where("category_id = ?", id).
			Delete(&model.PositionCategoryName{}).Error; err != nil {
			return errors.New("删除语言名称失败")
		}

		// 删除分类
		if err := tx.Where("id = ?", id).
			Delete(&model.PositionCategory{}).Error; err != nil {
			return errors.New("删除分类失败")
		}

		return nil
	})
}

// CreateSystemTag 创建系统预设标签
func (s *AdminContentService) CreateSystemTag(input CreateTagInput) (*model.Tag, error) {
	tag := &model.Tag{
		Name:     input.Name,
		Type:     input.Type,
		IsSystem: true,
	}
	if err := s.db.Create(tag).Error; err != nil {
		return nil, errors.New("创建标签失败")
	}
	return tag, nil
}

// DeleteSystemTag 删除系统预设标签
func (s *AdminContentService) DeleteSystemTag(id string) error {
	return s.db.Where("id = ? AND is_system = true", id).Delete(&model.Tag{}).Error
}

// PositionListInput 姿势列表查询参数
type PositionListInput struct {
	// CategoryID 按分类过滤，不传返回全部
	CategoryID string `json:"category_id"`

	// Keyword 搜索关键词，匹配姿势默认名称
	Keyword string `json:"keyword"`

	// IsActive 按启用状态过滤，不传返回全部
	IsActive *bool `json:"is_active"`

	// Page 页码，默认 1
	Page int `json:"page" binding:"omitempty,min=1" example:"1"`

	// PageSize 每页数量，默认 20
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

// PositionListResult 姿势列表返回结构
type PositionListResult struct {
	// List 姿势列表
	List []model.Position `json:"list"`

	// Total 总数量
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页数量
	PageSize int `json:"page_size"`
}

// ListSystemPositions 获取系统姿势列表（分页）
func (s *AdminContentService) ListSystemPositions(input PositionListInput) (*PositionListResult, error) {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	query := s.db.Model(&model.Position{}).
		Preload("Category").
		Where("is_system = true")

	// 按分类过滤
	if input.CategoryID != "" {
		query = query.Where("category_id = ?", input.CategoryID)
	}

	// 关键词搜索
	if input.Keyword != "" {
		query = query.Where("default_name LIKE ?", "%"+input.Keyword+"%")
	}

	var total int64
	query.Count(&total)

	var positions []model.Position
	offset := (input.Page - 1) * input.PageSize
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(input.PageSize).
		Find(&positions).Error
	if err != nil {
		return nil, errors.New("获取姿势列表失败")
	}

	return &PositionListResult{
		List:     positions,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

// validateBase64Icon 验证 base64 图标格式和大小
func validateBase64Icon(base64Str string) error {
	// 必须是 data:image/xxx;base64, 格式
	if !strings.HasPrefix(base64Str, "data:image/") {
		return errors.New("图标格式错误，必须是 base64 编码的图片")
	}

	// 限制大小，base64 字符串不超过 100KB
	if len(base64Str) > 100*1024 {
		return errors.New("图标太大，base64 编码后不能超过 100KB，建议使用 64x64px 的小图标")
	}

	return nil
}

// PositionNameInput 单个语言的姿势名称
type PositionNameInput struct {
	// LanguageCode 语言代码，如 zh-CN / en
	LanguageCode string `json:"language_code" binding:"required" example:"zh-CN"`
	// Name 该语言下的名称
	Name string `json:"name" binding:"required,max=50" example:"传教士"`
}

// CreatePositionInput 创建系统姿势请求参数
type CreatePositionInput struct {
	// Names 各语言名称列表，至少传一种语言
	Names []PositionNameInput `json:"names" binding:"required,min=1"`

	// CategoryID 所属分类 ID
	CategoryID string `json:"category_id" binding:"required"`

	// IconBase64 图标 base64，可选
	IconBase64 *string `json:"icon_base64"`
}

// UpdatePositionInput 更新姿势请求参数
type UpdatePositionInput struct {
	// Names 更新各语言名称，传入的语言会覆盖，未传的保持不变
	Names []PositionNameInput `json:"names"`

	// CategoryID 修改所属分类
	CategoryID *string `json:"category_id"`

	// IconBase64 修改图标
	IconBase64 *string `json:"icon_base64"`
}

// CreateSystemPosition 创建系统预设姿势
func (s *AdminContentService) CreateSystemPosition(input CreatePositionInput) (*model.Position, error) {
	// 验证分类是否存在
	var category model.PositionCategory
	if err := s.db.Where("id = ? AND is_active = true", input.CategoryID).First(&category).Error; err != nil {
		return nil, errors.New("分类不存在或已禁用")
	}

	// 检查重复名称
	for _, n := range input.Names {
		var count int64
		s.db.Model(&model.PositionName{}).
			Where("language_code = ? AND name = ?", n.LanguageCode, n.Name).
			Count(&count)
		if count > 0 {
			return nil, fmt.Errorf("「%s」名称已存在，请勿重复添加", n.Name)
		}
	}

	// 验证 base64
	if input.IconBase64 != nil {
		if err := validateBase64Icon(*input.IconBase64); err != nil {
			return nil, err
		}
	}

	// 取中文名作为 DefaultName，没有中文则取第一个
	defaultName := input.Names[0].Name
	for _, n := range input.Names {
		if n.LanguageCode == "zh-CN" {
			defaultName = n.Name
			break
		}
	}

	position := &model.Position{
		DefaultName: defaultName,
		CategoryID:  input.CategoryID,
		IconBase64:  input.IconBase64,
		IsSystem:    true,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(position).Error; err != nil {
			return errors.New("创建姿势失败")
		}

		for _, n := range input.Names {
			name := model.PositionName{
				PositionID:   position.ID,
				LanguageCode: n.LanguageCode,
				Name:         n.Name,
			}
			if err := tx.Create(&name).Error; err != nil {
				return errors.New("保存语言名称失败")
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.db.Where("position_id = ?", position.ID).Find(&position.Names)
	return position, nil
}

// UpdateSystemPosition 更新系统姿势
func (s *AdminContentService) UpdateSystemPosition(id string, input UpdatePositionInput) (*model.Position, error) {
	var position model.Position
	if err := s.db.Where("id = ? AND is_system = true", id).First(&position).Error; err != nil {
		return nil, errors.New("姿势不存在")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if input.CategoryID != nil {
			position.CategoryID = *input.CategoryID
		}
		if input.IconBase64 != nil {
			position.IconBase64 = input.IconBase64
		}

		if err := tx.Save(&position).Error; err != nil {
			return errors.New("更新失败")
		}

		// 更新语言名称（upsert）
		for _, n := range input.Names {
			var existing model.PositionName
			err := tx.Where("position_id = ? AND language_code = ?", id, n.LanguageCode).
				First(&existing).Error
			if err == gorm.ErrRecordNotFound {
				newName := model.PositionName{
					PositionID:   id,
					LanguageCode: n.LanguageCode,
					Name:         n.Name,
				}
				if err := tx.Create(&newName).Error; err != nil {
					return errors.New("保存语言名称失败")
				}
			} else {
				if err := tx.Model(&existing).Update("name", n.Name).Error; err != nil {
					return errors.New("更新语言名称失败")
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.db.Where("position_id = ?", id).Find(&position.Names)
	s.db.Preload("Category").First(&position, "id = ?", id)
	return &position, nil
}

// DeleteSystemPosition 删除系统姿势
func (s *AdminContentService) DeleteSystemPosition(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先删除多语言名称
		if err := tx.Where("position_id = ?", id).
			Delete(&model.PositionName{}).Error; err != nil {
			return errors.New("删除语言名称失败")
		}
		// 再删除姿势
		if err := tx.Where("id = ? AND is_system = true", id).
			Delete(&model.Position{}).Error; err != nil {
			return errors.New("删除失败")
		}
		return nil
	})
}
