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

// ─────────────────────────────────────────
// 翻译管理
// ─────────────────────────────────────────

// TranslationItem 单条翻译内容
type TranslationItem struct {
	// LanguageCode 语言代码
	LanguageCode string `json:"language_code" binding:"required" example:"en"`
	// Content 翻译内容
	Content string `json:"content" binding:"required" example:"Missionary"`
}

// SaveTranslationsInput 批量保存翻译请求参数
type SaveTranslationsInput struct {
	// Module 模块名：position / tag
	Module string `json:"module" binding:"required,oneof=position tag" example:"position"`

	// RefID 关联记录 ID（姿势 ID 或标签 ID）
	RefID string `json:"ref_id" binding:"required"`

	// Field 翻译字段：name / category_name
	Field string `json:"field" binding:"required,oneof=name category_name" example:"name"`

	// Translations 各语言翻译内容列表
	Translations []TranslationItem `json:"translations" binding:"required,min=1"`
}

// GetTranslationsInput 获取翻译请求参数
type GetTranslationsInput struct {
	// Module 模块名
	Module string `json:"module" binding:"required"`
	// RefID 关联记录 ID
	RefID string `json:"ref_id" binding:"required"`
	// Field 翻译字段
	Field string `json:"field" binding:"required"`
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

// SaveTranslations 批量保存翻译（存在则更新，不存在则创建）
func (s *I18nService) SaveTranslations(input SaveTranslationsInput) error {
	for _, t := range input.Translations {
		translation := model.Translation{
			Module:       input.Module,
			RefID:        input.RefID,
			Field:        input.Field,
			LanguageCode: t.LanguageCode,
			Content:      t.Content,
		}

		// upsert：有则更新，无则创建
		var existing model.Translation
		err := s.db.Where(
			"module = ? AND ref_id = ? AND field = ? AND language_code = ?",
			input.Module, input.RefID, input.Field, t.LanguageCode,
		).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&translation).Error; err != nil {
				return errors.New("保存翻译失败")
			}
		} else {
			if err := s.db.Model(&existing).Update("content", t.Content).Error; err != nil {
				return errors.New("更新翻译失败")
			}
		}
	}
	return nil
}

// GetTranslations 获取某条记录的所有翻译
func (s *I18nService) GetTranslations(input GetTranslationsInput) (*TranslationResult, error) {
	var translations []model.Translation
	err := s.db.
		Where("module = ? AND ref_id = ? AND field = ?",
			input.Module, input.RefID, input.Field).
		Find(&translations).Error
	if err != nil {
		return nil, errors.New("获取翻译失败")
	}

	return &TranslationResult{
		RefID:        input.RefID,
		Field:        input.Field,
		Translations: translations,
	}, nil
}

// GetAllTranslationsByRef 获取某条记录的所有字段的所有翻译
// 用于编辑页面一次性加载所有翻译
func (s *I18nService) GetAllTranslationsByRef(module, refID string) ([]model.Translation, error) {
	var translations []model.Translation
	err := s.db.
		Where("module = ? AND ref_id = ?", module, refID).
		Order("field ASC, language_code ASC").
		Find(&translations).Error
	return translations, err
}

// DeleteTranslationsByRef 删除某条记录的所有翻译
// 用于删除姿势/标签时级联删除翻译
func (s *I18nService) DeleteTranslationsByRef(module, refID string) error {
	return s.db.
		Where("module = ? AND ref_id = ?", module, refID).
		Delete(&model.Translation{}).Error
}
