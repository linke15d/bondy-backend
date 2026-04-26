// Package model 数据库模型定义
package model

import "time"

// Wishlist 心愿清单表
// 伴侣双方都可以添加心愿，支持匿名提案（隐藏提案人）
// 对应数据库表名: wishlists
type Wishlist struct {
	// ID 心愿唯一标识，UUID 格式
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// CoupleID 所属伴侣关系 ID
	CoupleID string `gorm:"type:uuid;not null;index" json:"couple_id"`

	// CreatedByID 提案人用户 ID
	CreatedByID string `gorm:"type:uuid;not null" json:"created_by_id"`

	// Title 心愿标题
	Title string `gorm:"size:100;not null" json:"title" example:"去海边看日出"`

	// Description 心愿详细描述，可选
	Description *string `gorm:"type:text" json:"description,omitempty" example:"找一个好天气，一起去看日出"`

	// IsAnonymous 是否匿名提案
	// 匿名时对方看不到提案人是谁，增加神秘感
	IsAnonymous bool `gorm:"default:false" json:"is_anonymous"`

	// IsCompleted 是否已完成
	IsCompleted bool `gorm:"default:false;index" json:"is_completed"`

	// CompletedAt 完成时间
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Heat 热度值，双方都可以点赞提高热度，热度越高越优先展示
	Heat int `gorm:"default:0" json:"heat"`

	// Scope 心愿范围：COUPLE（共同心愿）或 PERSONAL（个人心愿）
	Scope string `gorm:"size:20;default:'COUPLE'" json:"scope" example:"COUPLE"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}
