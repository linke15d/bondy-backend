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
