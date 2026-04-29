// Package admin 后台管理业务逻辑层
package admin

import (
	"errors"

	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// I18nService 多语言管理业务逻辑
type I18nService struct {
	db *gorm.DB
}

// NewI18nService 创建 I18nService 实例
func NewI18nService(db *gorm.DB) *I18nService {
	return &I18nService{db: db}
}

// ─────────────────────────────────────────
// 语言管理
// ─────────────────────────────────────────

// CreateLanguageInput 创建语言请求参数
type CreateLanguageInput struct {
	// Code 语言代码，如 zh-CN / en / ja / ko
	Code string `json:"code" binding:"required,max=10" example:"en"`

	// Name 语言名称
	Name string `json:"name" binding:"required,max=50" example:"English"`

	// IsDefault 是否设为默认语言
	IsDefault bool `json:"is_default" example:"false"`

	// SortOrder 排序
	SortOrder int `json:"sort_order" example:"1"`
}

// CreateLanguage 创建语言
func (s *I18nService) CreateLanguage(input CreateLanguageInput) (*model.SupportedLanguage, error) {
	// 如果设为默认，先取消其他默认
	if input.IsDefault {
		s.db.Model(&model.SupportedLanguage{}).
			Where("is_default = true").
			Update("is_default", false)
	}

	lang := &model.SupportedLanguage{
		Code:      input.Code,
		Name:      input.Name,
		IsDefault: input.IsDefault,
		SortOrder: input.SortOrder,
	}

	if err := s.db.Create(lang).Error; err != nil {
		return nil, errors.New("创建语言失败，语言代码可能已存在")
	}
	return lang, nil
}

// ListLanguages 获取所有语言列表
func (s *I18nService) ListLanguages() ([]model.SupportedLanguage, error) {
	var languages []model.SupportedLanguage
	err := s.db.Order("sort_order ASC, created_at ASC").Find(&languages).Error
	return languages, err
}

// UpdateLanguageStatus 启用/禁用语言
func (s *I18nService) UpdateLanguageStatus(id string, isActive bool) error {
	return s.db.Model(&model.SupportedLanguage{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

// UpdateLanguageInput 更新语言请求参数
type UpdateLanguageInput struct {
	// Name 修改语言名称
	Name *string `json:"name" binding:"omitempty,max=50" example:"English"`

	// IsDefault 是否设为默认语言
	IsDefault *bool `json:"is_default"`

	// IsActive 是否启用
	IsActive *bool `json:"is_active"`

	// SortOrder 修改排序
	SortOrder *int `json:"sort_order"`
}

// UpdateLanguage 更新语言
func (s *I18nService) UpdateLanguage(id string, input UpdateLanguageInput) (*model.SupportedLanguage, error) {
	var lang model.SupportedLanguage
	if err := s.db.Where("id = ?", id).First(&lang).Error; err != nil {
		return nil, errors.New("语言不存在")
	}

	// 如果设为默认，先取消其他默认
	if input.IsDefault != nil && *input.IsDefault {
		s.db.Model(&model.SupportedLanguage{}).
			Where("is_default = true AND id != ?", id).
			Update("is_default", false)
		lang.IsDefault = true
	}

	if input.Name != nil {
		lang.Name = *input.Name
	}
	if input.IsActive != nil {
		lang.IsActive = *input.IsActive
	}
	if input.SortOrder != nil {
		lang.SortOrder = *input.SortOrder
	}

	if err := s.db.Save(&lang).Error; err != nil {
		return nil, errors.New("更新失败")
	}

	return &lang, nil
}

// DeleteLanguage 删除语言
// 删除前检查是否有翻译内容使用此语言
func (s *I18nService) DeleteLanguage(id string) error {
	var lang model.SupportedLanguage
	if err := s.db.Where("id = ?", id).First(&lang).Error; err != nil {
		return errors.New("语言不存在")
	}

	// 默认语言不允许删除
	if lang.IsDefault {
		return errors.New("默认语言不允许删除")
	}

	// 检查是否有翻译内容使用此语言
	var count int64
	s.db.Model(&model.Translation{}).
		Where("language_code = ?", lang.Code).
		Count(&count)
	if count > 0 {
		return errors.New("该语言下还有翻译内容，请先删除相关翻译")
	}

	return s.db.Where("id = ?", id).Delete(&model.SupportedLanguage{}).Error
}

// TranslationResult 翻译结果
type TranslationResult struct {
	// RefID 关联记录 ID
	RefID string `json:"ref_id"`
	// Field 翻译字段
	Field string `json:"field"`
	// Translations 各语言翻译
	Translations []model.Translation `json:"translations"`
}

// DeleteTranslationsByRef 删除某条记录的所有翻译
// 用于删除姿势/标签时级联删除翻译
func (s *I18nService) DeleteTranslationsByRef(module, refID string) error {
	return s.db.
		Where("module = ? AND ref_id = ?", module, refID).
		Delete(&model.Translation{}).Error
}
