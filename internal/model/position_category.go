// Package model 数据库模型定义
package model

import "github.com/linke15d/bondy-backend/pkg/timeformat"

// PositionCategory 姿势分类表
// 后台可配置分类，支持多语言
// 对应数据库表名: position_categories
type PositionCategory struct {
	// ID 分类唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// SortOrder 排序，数字越小越靠前
	SortOrder int `gorm:"default:0" json:"sort_order" example:"1"`

	// IsActive 是否启用
	IsActive bool `gorm:"default:true" json:"is_active"`

	// CreatedAt 创建时间
	CreatedAt timeformat.LocalTime `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt timeformat.LocalTime `json:"updated_at"`

	// Names 各语言名称，不存数据库，查询时填充
	Names []PositionCategoryName `gorm:"foreignKey:CategoryID" json:"names,omitempty"`

	// Name 当前语言名称，不存数据库，App 端查询时填充
	Name string `gorm:"-" json:"name,omitempty"`
}

// PositionCategoryName 分类多语言名称表
// 对应数据库表名: position_category_names
type PositionCategoryName struct {
	ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CategoryID string `gorm:"type:uuid;not null;index" json:"category_id"`
	// LanguageCode 语言代码，关联 supported_languages 表
	LanguageCode string `gorm:"size:10;not null" json:"language_code" example:"zh-CN"`
	// Name 该语言下的分类名称
	Name      string               `gorm:"size:50;not null" json:"name" example:"经典"`
	CreatedAt timeformat.LocalTime `json:"created_at"`
}
