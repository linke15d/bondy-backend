package model

import "time"

// RefreshToken 刷新令牌表
// 每次登录生成一条记录，登出或刷新后删除
// 对应数据库表名: refresh_tokens
type RefreshToken struct {
	// ID 记录唯一标识，UUID 格式
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// UserID 所属用户 ID，关联 users 表
	UserID string `gorm:"type:uuid;not null;index" json:"user_id"`

	// Token Refresh Token 字符串，JWT 格式，唯一，不在 API 响应中返回
	Token string `gorm:"uniqueIndex;size:500" json:"-"`

	// ExpiresAt token 过期时间，过期后即使存在数据库中也无法使用
	ExpiresAt time.Time `json:"expires_at"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// User 关联的用户对象，仅用于 GORM 联表查询，不在 API 响应中返回
	User User `gorm:"foreignKey:UserID" json:"-"`
}
