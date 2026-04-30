// Package model 数据库模型定义
package model

import (
	"time"

	"github.com/linke15d/bondy-backend/pkg/timeformat"
)

// Record 亲密记录表
// 每条记录代表一次亲密行为，由伴侣中的任意一方创建
// 对应数据库表名: records
type Record struct {
	// ID 记录唯一标识，UUID 格式
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// CoupleID 所属伴侣关系 ID，关联 couples 表
	CoupleID string `gorm:"type:uuid;not null;index" json:"couple_id"`

	// CreatedByID 创建此记录的用户 ID
	CreatedByID string `gorm:"type:uuid;not null" json:"created_by_id"`

	// HappenedAt 实际发生时间，支持回填历史记录
	HappenedAt time.Time `gorm:"not null;index" json:"happened_at"`

	// DurationMins 持续时长（分钟），可选
	DurationMins *int `json:"duration_mins,omitempty" example:"30"`

	// Mood 心情评分 1-5，1最差5最好，可选
	Mood *int `gorm:"check:mood >= 1 AND mood <= 5" json:"mood,omitempty" example:"4"`

	// Satisfaction 满意度评分 1-5，1最差5最好，可选
	Satisfaction *int `gorm:"check:satisfaction >= 1 AND satisfaction <= 5" json:"satisfaction,omitempty" example:"5"`

	// NoteEncrypted 备注内容，客户端加密后存储的密文，后端不解密
	NoteEncrypted *string `gorm:"type:text" json:"note_encrypted,omitempty"`

	// IsDeleted 软删除标记
	IsDeleted bool `gorm:"default:false;index" json:"-"`

	// CreatedAt 记录创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 记录最后更新时间
	UpdatedAt time.Time `json:"updated_at"`

	// Couple 关联的伴侣关系，仅用于联表查询
	Couple Couple `gorm:"foreignKey:CoupleID" json:"-"`

	// CreatedBy 创建人用户信息，仅用于联表查询
	CreatedBy User `gorm:"foreignKey:CreatedByID" json:"-"`

	// Tags 关联的标签列表（地点、活动）
	Tags []Tag `gorm:"many2many:record_tags;" json:"tags,omitempty"`

	// Positions 关联的姿势列表
	Positions []Position `gorm:"many2many:record_positions;" json:"positions,omitempty"`
}

// Position 姿势表
// Position 姿势表
type Position struct {
	// ID 姿势唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// DefaultName 默认名称（中文），用于后台展示和搜索
	DefaultName string `gorm:"size:30;not null" json:"default_name" example:"传教士"`

	// Name 当前语言名称，不存数据库，查询时填充
	Name string `gorm:"-" json:"name" example:"传教士"`

	// CategoryID 所属分类 ID，关联 position_categories 表
	CategoryID string `gorm:"type:uuid;not null;index" json:"category_id"`

	// Category 关联的分类对象，查询时填充
	Category *PositionCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`

	// IconBase64 图标的 base64 编码
	IconBase64 *string `gorm:"type:text" json:"icon_base64,omitempty"`

	// IsSystem 是否为系统预设姿势
	IsSystem bool `gorm:"default:false" json:"is_system"`

	// IsActive 是否启用，禁用后 App 端不显示此姿势
	IsActive bool `gorm:"default:true" json:"is_active"`

	// CoupleID 所属伴侣 ID，系统姿势此字段为空
	CoupleID *string `gorm:"type:uuid;index" json:"couple_id,omitempty"`

	// CreatedAt 创建时间
	CreatedAt timeformat.LocalTime `json:"created_at"`

	// Names 多语言名称列表，后台管理时返回
	Names []PositionName `gorm:"foreignKey:PositionID" json:"names,omitempty"`
}
