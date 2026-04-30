// Package model 数据库模型定义
package model

import "github.com/linke15d/bondy-backend/pkg/timeformat"

// Location 地点表
// 系统预设地点，支持多语言名称
// 对应数据库表名: locations
type Location struct {
	// ID 地点唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// DefaultName 默认名称（中文），用于后台展示和搜索
	DefaultName string `gorm:"size:50;not null" json:"default_name" example:"家里"`

	// Name 当前语言名称，不存数据库，查询时填充
	Name string `gorm:"-" json:"name" example:"家里"`

	// IconBase64 图标 base64 编码
	IconBase64 *string `gorm:"type:text" json:"icon_base64,omitempty"`

	// IsSystem 是否为系统预设地点
	IsSystem bool `gorm:"default:false" json:"is_system"`

	// IsActive 是否启用，禁用后 App 端不显示
	IsActive bool `gorm:"default:true" json:"is_active"`

	// SortOrder 排序，数字越小越靠前
	SortOrder int `gorm:"default:0" json:"sort_order"`

	// CoupleID 所属伴侣 ID，系统地点此字段为空
	CoupleID *string `gorm:"type:uuid;index" json:"couple_id,omitempty"`

	// CreatedAt 创建时间
	CreatedAt timeformat.LocalTime `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt timeformat.LocalTime `json:"updated_at"`

	// Names 多语言名称列表，后台管理时返回
	Names []LocationName `gorm:"foreignKey:LocationID" json:"names,omitempty"`
}

// LocationName 地点多语言名称表
// 对应数据库表名: location_names
type LocationName struct {
	// ID 唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// LocationID 所属地点 ID
	LocationID string `gorm:"type:uuid;not null;index" json:"location_id"`

	// LanguageCode 语言代码，如 zh-CN / en
	LanguageCode string `gorm:"size:10;not null" json:"language_code" example:"zh-CN"`

	// Name 该语言下的地点名称
	Name string `gorm:"size:50;not null" json:"name" example:"家里"`

	// CreatedAt 创建时间
	CreatedAt timeformat.LocalTime `json:"created_at"`
}
