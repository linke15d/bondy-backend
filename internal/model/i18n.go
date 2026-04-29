// Package model 数据库模型定义
package model

import "github.com/linke15d/bondy-backend/pkg/timeformat"

// SupportedLanguage 支持的语言列表
// 对应数据库表名: supported_languages
type SupportedLanguage struct {
	// ID 语言唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// Code 语言代码，如 zh-CN / en / ja / ko
	Code string `gorm:"size:10;not null;uniqueIndex" json:"code" example:"zh-CN"`

	// Name 语言名称，如 简体中文 / English / 日本語
	Name string `gorm:"size:50;not null" json:"name" example:"简体中文"`

	// IsDefault 是否为默认语言
	IsDefault bool `gorm:"default:false" json:"is_default"`

	// IsActive 是否启用
	IsActive bool `gorm:"default:true" json:"is_active"`

	// SortOrder 排序，数字越小越靠前
	SortOrder int `gorm:"default:0" json:"sort_order"`

	// CreatedAt 创建时间
	CreatedAt timeformat.LocalTime `json:"created_at"`
}

// Translation 多语言翻译表
// 统一存储所有模块的多语言内容
// 对应数据库表名: translations
type Translation struct {
	// ID 翻译记录唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// Module 所属模块：position / tag / common
	Module string `gorm:"size:50;not null;index" json:"module" example:"position"`

	// RefID 关联的记录 ID，如姿势 ID、标签 ID
	RefID string `gorm:"type:uuid;not null;index" json:"ref_id"`

	// Field 翻译的字段名：name / category_name / description
	Field string `gorm:"size:50;not null" json:"field" example:"name"`

	// LanguageCode 语言代码，关联 supported_languages 表
	LanguageCode string `gorm:"size:10;not null;index" json:"language_code" example:"zh-CN"`

	// Content 翻译内容
	Content string `gorm:"size:500;not null" json:"content" example:"传教士"`

	// CreatedAt 创建时间
	CreatedAt timeformat.LocalTime `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt timeformat.LocalTime `json:"updated_at"`
}
