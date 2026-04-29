// Package admin 后台管理业务逻辑层
package admin

import (
	"errors"
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

// CreatePositionInput 创建系统姿势请求参数
type CreatePositionInput struct {
	// Name 姿势名称
	Name string `json:"name" binding:"required,max=30" example:"传教士"`
	// Category 分类：CLASSIC / ADVENTURE / INTIMATE / FUN
	Category string `json:"category" binding:"required,oneof=CLASSIC ADVENTURE INTIMATE FUN" example:"CLASSIC"`
	// IconBase64 图标 base64 编码，格式：data:image/png;base64,xxx
	// 建议图标尺寸 64x64px 以内，大小不超过 50KB
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

// CreateCategoryInput 创建分类请求参数
type CreateCategoryInput struct {
	// Code 分类英文代码，唯一，建议全大写，如 CLASSIC
	Code string `json:"code" binding:"required,max=50" example:"CLASSIC"`

	// DefaultName 默认中文名称
	DefaultName string `json:"default_name" binding:"required,max=50" example:"经典"`

	// SortOrder 排序值，数字越小越靠前
	SortOrder int `json:"sort_order" example:"1"`
}

// UpdateCategoryInput 更新分类请求参数
type UpdateCategoryInput struct {
	// DefaultName 修改默认名称
	DefaultName *string `json:"default_name" binding:"omitempty,max=50"`

	// SortOrder 修改排序
	SortOrder *int `json:"sort_order"`

	// IsActive 修改启用状态
	IsActive *bool `json:"is_active"`
}

// CreatePositionCategory 创建姿势分类
func (s *AdminContentService) CreatePositionCategory(input CreateCategoryInput) (*model.PositionCategory, error) {
	// 检查 Code 是否已存在
	var count int64
	s.db.Model(&model.PositionCategory{}).Where("code = ?", input.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("该分类代码已存在")
	}

	category := &model.PositionCategory{
		Code:        input.Code,
		DefaultName: input.DefaultName,
		SortOrder:   input.SortOrder,
		IsActive:    true,
	}

	if err := s.db.Create(category).Error; err != nil {
		return nil, errors.New("创建分类失败")
	}

	return category, nil
}

// ListPositionCategories 获取分类列表（含多语言翻译）
func (s *AdminContentService) ListPositionCategories() ([]model.PositionCategory, error) {
	var categories []model.PositionCategory
	err := s.db.Order("sort_order ASC, created_at ASC").Find(&categories).Error
	if err != nil {
		return nil, errors.New("获取分类列表失败")
	}

	// 填充每个分类的多语言翻译
	for i := range categories {
		var translations []model.Translation
		s.db.Where("module = 'position_category' AND ref_id = ?", categories[i].ID).
			Find(&translations)
		categories[i].Translations = translations
	}

	return categories, nil
}

// UpdatePositionCategory 更新分类
func (s *AdminContentService) UpdatePositionCategory(id string, input UpdateCategoryInput) (*model.PositionCategory, error) {
	var category model.PositionCategory
	if err := s.db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, errors.New("分类不存在")
	}

	if input.DefaultName != nil {
		category.DefaultName = *input.DefaultName
	}
	if input.SortOrder != nil {
		category.SortOrder = *input.SortOrder
	}
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	if err := s.db.Save(&category).Error; err != nil {
		return nil, errors.New("更新失败")
	}

	return &category, nil
}

// DeletePositionCategory 删除分类
// 删除前检查是否有姿势关联此分类
func (s *AdminContentService) DeletePositionCategory(id string) error {
	var count int64
	s.db.Model(&model.Position{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该分类下还有姿势，请先删除或移动姿势后再删除分类")
	}

	return s.db.Where("id = ?", id).Delete(&model.PositionCategory{}).Error
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

// ListSystemPositions 获取系统预设姿势列表
func (s *AdminContentService) ListSystemPositions() ([]model.Position, error) {
	var positions []model.Position
	err := s.db.Where("is_system = true").Order("name ASC").Find(&positions).Error
	return positions, err
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

// CreateSystemPosition 创建系统预设姿势
func (s *AdminContentService) CreateSystemPosition(input CreatePositionInput) (*model.Position, error) {
	// 验证 base64 格式
	if input.IconBase64 != nil {
		if err := validateBase64Icon(*input.IconBase64); err != nil {
			return nil, err
		}
	}

	position := &model.Position{
		Name:       input.Name,
		IconBase64: input.IconBase64,
		IsSystem:   true,
	}
	if err := s.db.Create(position).Error; err != nil {
		return nil, errors.New("创建失败")
	}
	return position, nil
}

// DeleteSystemPosition 删除系统预设姿势
func (s *AdminContentService) DeleteSystemPosition(id string) error {
	return s.db.Where("id = ? AND is_system = true", id).Delete(&model.Position{}).Error
}
