// Package model 数据库模型定义
package model

import "time"

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

	// Name 标签名称，如"家里"、"酒店"、"浪漫"
	Name string `gorm:"size:30;not null" json:"name" example:"家里"`

	// Type 标签类型：LOCATION（地点）或 ACTIVITY（活动）
	Type string `gorm:"size:20;not null" json:"type" example:"LOCATION"`

	// IsSystem 是否为系统预设标签，系统标签不可删除
	IsSystem bool `gorm:"default:false" json:"is_system"`

	// CoupleID 所属伴侣 ID，系统标签此字段为空
	CoupleID *string `gorm:"type:uuid;index" json:"couple_id,omitempty"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// Position 姿势表
// 系统预设姿势 + 用户自定义姿势
// 对应数据库表名: positions
type Position struct {
	// ID 姿势唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// Name 姿势名称
	Name string `gorm:"size:30;not null" json:"name" example:"传教士"`

	// IsSystem 是否为系统预设姿势
	IsSystem bool `gorm:"default:false" json:"is_system"`

	// CoupleID 所属伴侣 ID，系统姿势此字段为空
	CoupleID *string `gorm:"type:uuid;index" json:"couple_id,omitempty"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}
