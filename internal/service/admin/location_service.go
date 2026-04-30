// Package admin 后台管理业务逻辑层
package admin

import (
	"errors"
	"fmt"

	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// LocationService 地点管理业务逻辑
type LocationService struct {
	db *gorm.DB
}

// NewLocationService 创建 LocationService 实例
func NewLocationService(db *gorm.DB) *LocationService {
	return &LocationService{db: db}
}

// LocationNameInput 单个语言的地点名称
type LocationNameInput struct {
	// LanguageCode 语言代码，如 zh-CN / en
	LanguageCode string `json:"language_code" binding:"required" example:"zh-CN"`
	// Name 该语言下的名称
	Name string `json:"name" binding:"required,max=50" example:"家里"`
}

// CreateLocationInput 创建地点请求参数
type CreateLocationInput struct {
	// Names 各语言名称列表，至少传一种语言
	Names []LocationNameInput `json:"names" binding:"required,min=1"`
	// IconBase64 图标 base64，可选
	IconBase64 *string `json:"icon_base64"`
	// SortOrder 排序值
	SortOrder int `json:"sort_order" example:"1"`
}

// UpdateLocationInput 更新地点请求参数
type UpdateLocationInput struct {
	// Names 更新各语言名称，传入的语言会覆盖，未传的保持不变
	Names []LocationNameInput `json:"names"`
	// IconBase64 修改图标
	IconBase64 *string `json:"icon_base64"`
	// SortOrder 修改排序
	SortOrder *int `json:"sort_order"`
	// IsActive 修改启用状态
	IsActive *bool `json:"is_active"`
}

// LocationListInput 地点列表查询参数
type LocationListInput struct {
	// Keyword 搜索关键词，匹配默认名称
	Keyword string `json:"keyword"`
	// IsActive 按启用状态过滤，不传返回全部
	IsActive *bool `json:"is_active"`
	// Page 页码，默认 1
	Page int `json:"page" binding:"omitempty,min=1" example:"1"`
	// PageSize 每页数量，默认 20
	PageSize int `json:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

// LocationListResult 地点列表返回结构
type LocationListResult struct {
	// List 地点列表
	List []model.Location `json:"list"`
	// Total 总数量
	Total int64 `json:"total"`
	// Page 当前页码
	Page int `json:"page"`
	// PageSize 每页数量
	PageSize int `json:"page_size"`
}

// CreateLocation 创建地点
func (s *LocationService) CreateLocation(input CreateLocationInput) (*model.Location, error) {
	// 检查重复名称
	for _, n := range input.Names {
		var count int64
		s.db.Model(&model.LocationName{}).
			Where("language_code = ? AND name = ?", n.LanguageCode, n.Name).
			Count(&count)
		if count > 0 {
			return nil, fmt.Errorf("「%s」名称已存在，请勿重复添加", n.Name)
		}
	}

	// 取中文名作为 DefaultName
	defaultName := input.Names[0].Name
	for _, n := range input.Names {
		if n.LanguageCode == "zh-CN" {
			defaultName = n.Name
			break
		}
	}

	location := &model.Location{
		DefaultName: defaultName,
		IconBase64:  input.IconBase64,
		SortOrder:   input.SortOrder,
		IsSystem:    true,
		IsActive:    true,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(location).Error; err != nil {
			return errors.New("创建地点失败")
		}

		for _, n := range input.Names {
			name := model.LocationName{
				LocationID:   location.ID,
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

	s.db.Where("location_id = ?", location.ID).Find(&location.Names)
	return location, nil
}

// ListLocations 获取地点列表（分页）
func (s *LocationService) ListLocations(input LocationListInput) (*LocationListResult, error) {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	query := s.db.Model(&model.Location{}).
		Preload("Names").
		Where("is_system = true")

	if input.Keyword != "" {
		query = query.Where("default_name LIKE ?", "%"+input.Keyword+"%")
	}

	if input.IsActive != nil {
		query = query.Where("is_active = ?", *input.IsActive)
	}

	var total int64
	query.Count(&total)

	var locations []model.Location
	offset := (input.Page - 1) * input.PageSize
	err := query.
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(input.PageSize).
		Find(&locations).Error
	if err != nil {
		return nil, errors.New("获取地点列表失败")
	}

	return &LocationListResult{
		List:     locations,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

// UpdateLocation 更新地点
func (s *LocationService) UpdateLocation(id string, input UpdateLocationInput) (*model.Location, error) {
	var location model.Location
	if err := s.db.Where("id = ? AND is_system = true", id).First(&location).Error; err != nil {
		return nil, errors.New("地点不存在")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{}

		if input.IconBase64 != nil {
			updates["icon_base64"] = input.IconBase64
		}
		if input.SortOrder != nil {
			updates["sort_order"] = *input.SortOrder
		}
		if input.IsActive != nil {
			updates["is_active"] = *input.IsActive
		}

		// 如果更新了 zh-CN 名称，同步更新 DefaultName
		for _, n := range input.Names {
			if n.LanguageCode == "zh-CN" {
				updates["default_name"] = n.Name
				break
			}
		}

		if len(updates) > 0 {
			if err := tx.Model(&location).Updates(updates).Error; err != nil {
				return errors.New("更新失败")
			}
		}

		// 更新语言名称（upsert）
		for _, n := range input.Names {
			var existing model.LocationName
			err := tx.Where("location_id = ? AND language_code = ?", id, n.LanguageCode).
				First(&existing).Error
			if err == gorm.ErrRecordNotFound {
				newName := model.LocationName{
					LocationID:   id,
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

	s.db.Preload("Names").First(&location, "id = ?", id)
	return &location, nil
}

// DeleteLocation 删除地点
func (s *LocationService) DeleteLocation(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先删除多语言名称
		if err := tx.Where("location_id = ?", id).
			Delete(&model.LocationName{}).Error; err != nil {
			return errors.New("删除语言名称失败")
		}
		// 再删除地点
		if err := tx.Where("id = ? AND is_system = true", id).
			Delete(&model.Location{}).Error; err != nil {
			return errors.New("删除失败")
		}
		return nil
	})
}
