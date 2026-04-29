// Package model 数据库模型定义
package model

import "github.com/linke15d/bondy-backend/pkg/timeformat"

// PositionCategory 姿势分类表
// 后台可配置分类，支持多语言
// 对应数据库表名: position_categories
type PositionCategory struct {
	// ID 分类唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// Code 分类英文代码，唯一，App 端用此值做多语言 Key
	// 例如：CLASSIC / ADVENTURE / INTIMATE / FUN
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code" example:"CLASSIC"`

	// DefaultName 默认名称（中文），后台展示用
	DefaultName string `gorm:"size:50;not null" json:"default_name" example:"经典"`

	// SortOrder 排序，数字越小越靠前
	SortOrder int `gorm:"default:0" json:"sort_order" example:"1"`

	// IsActive 是否启用
	IsActive bool `gorm:"default:true" json:"is_active"`

	// CreatedAt 创建时间
	CreatedAt timeformat.LocalTime `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt timeformat.LocalTime `json:"updated_at"`

	// Translations 多语言翻译，不存数据库，查询时填充
	Translations []Translation `gorm:"-" json:"translations,omitempty"`
}
