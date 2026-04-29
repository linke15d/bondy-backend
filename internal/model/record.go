// Package model 数据库模型定义
package model

import (
	"time"

	"gorm.io/gorm"
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

// Tag 标签表
// 分为地点标签和活动标签两种类型
// 系统预设标签（is_system=true）所有用户共享
// 用户自定义标签（is_system=false）仅属于当前伴侣
// 对应数据库表名: tags
type Tag struct {
	// ID 标签唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// DefaultName 默认名称（中文）
	DefaultName string `gorm:"size:30;not null" json:"default_name" example:"家里"`

	// Name 标签名称，如"家里"、"酒店"、"浪漫"
	Name string `gorm:"size:30;not null" json:"name" example:"家里"`

	// Type 标签类型：LOCATION（地点）或 ACTIVITY（活动）
	Type string `gorm:"size:20;not null" json:"type" example:"LOCATION"`

	// IconBase64 标签图标 base64
	IconBase64 *string `gorm:"type:text" json:"icon_base64,omitempty"`

	// IsSystem 是否为系统预设标签，系统标签不可删除
	IsSystem bool `gorm:"default:false" json:"is_system"`

	// CoupleID 所属伴侣 ID，系统标签此字段为空
	CoupleID *string `gorm:"type:uuid;index" json:"couple_id,omitempty"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// Position 姿势表
// 对应数据库表名: positions
type Position struct {
	// ID 姿势唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	// DefaultName 默认名称（中文），用于后台展示和搜索
	DefaultName string `gorm:"size:30;not null" json:"default_name" example:"传教士"`
	// Name 姿势名称
	Name string `gorm:"size:30;not null" json:"name" example:"传教士"`
	// CategoryID 所属分类 ID，关联 position_categories 表
	CategoryID string `gorm:"type:uuid;not null;index" json:"category_id"`
	// Category 关联的分类对象，查询时填充
	Category string `gorm:"size:20;not null;default:'CLASSIC'" json:"category" example:"CLASSIC"`
	// CategoryName 分类中文名，由 category 字段转换而来，不存数据库
	CategoryName string `gorm:"-" json:"category_name" example:"经典"`
	// IconBase64 图标的 base64 编码，格式：data:image/png;base64,xxx
	IconBase64 *string `gorm:"type:text" json:"icon_base64,omitempty"`
	// IsSystem 是否为系统预设姿势
	IsSystem bool `gorm:"default:false" json:"is_system"`
	// CoupleID 所属伴侣 ID，系统姿势此字段为空
	CoupleID *string `gorm:"type:uuid;index" json:"couple_id,omitempty"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// Names 多语言名称列表，后台管理时返回
	Names []PositionName `gorm:"foreignKey:PositionID" json:"names,omitempty"`
}

// CategoryNameMap 分类枚举值到中文的映射
var CategoryNameMap = map[string]string{
	"CLASSIC":   "经典",
	"ADVENTURE": "探险",
	"INTIMATE":  "亲密",
	"FUN":       "趣味",
}

// AfterFind GORM 查询后自动填充 CategoryName
func (p *Position) AfterFind(tx *gorm.DB) error {
	if name, ok := CategoryNameMap[p.Category]; ok {
		p.CategoryName = name
	}
	return nil
}
