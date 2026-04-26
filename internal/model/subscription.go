// Package model 数据库模型定义
package model

import "time"

// Subscription 订阅会员表
// 记录用户的会员购买记录，一个用户同时只有一条有效记录
// 对应数据库表名: subscriptions
type Subscription struct {
	// ID 订阅记录唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// UserID 所属用户 ID，唯一索引保证一个用户只有一条记录
	UserID string `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Plan 订阅套餐：MONTHLY（月付）或 LIFETIME（买断）
	Plan string `gorm:"size:20;not null" json:"plan" example:"MONTHLY"`

	// Status 订阅状态：ACTIVE（有效）、EXPIRED（已过期）、CANCELLED（已取消）
	Status string `gorm:"size:20;not null;default:'ACTIVE'" json:"status" example:"ACTIVE"`

	// StartAt 订阅开始时间
	StartAt time.Time `gorm:"not null" json:"start_at"`

	// ExpiresAt 订阅过期时间，买断套餐此字段为空（永久有效）
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// Provider 支付渠道：apple（App Store）、google（Google Play）、stripe
	Provider string `gorm:"size:20" json:"provider" example:"apple"`

	// ProviderSubID 第三方平台的订阅 ID，用于核销和退款
	ProviderSubID *string `gorm:"size:200" json:"provider_sub_id,omitempty"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`

	// User 关联用户，仅用于联表查询
	User User `gorm:"foreignKey:UserID" json:"-"`
}
