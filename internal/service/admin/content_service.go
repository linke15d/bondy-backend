// Package admin 后台管理业务逻辑层
package admin

import (
	"errors"

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
}

// CreatePositionInput 创建系统姿势请求参数
type CreatePositionInput struct {
	// Name 姿势名称
	Name string `json:"name" binding:"required,max=30" example:"传教士"`
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

// CreateSystemPosition 创建系统预设姿势
func (s *AdminContentService) CreateSystemPosition(input CreatePositionInput) (*model.Position, error) {
	position := &model.Position{
		Name:     input.Name,
		IsSystem: true,
	}
	if err := s.db.Create(position).Error; err != nil {
		return nil, errors.New("创建姿势失败")
	}
	return position, nil
}

// DeleteSystemPosition 删除系统预设姿势
func (s *AdminContentService) DeleteSystemPosition(id string) error {
	return s.db.Where("id = ? AND is_system = true", id).Delete(&model.Position{}).Error
}
