// Package model 数据库模型定义
package model

import "time"

// Admin 管理员表
// 后台管理系统的登录账号，与普通用户完全隔离
// 对应数据库表名: admins
type Admin struct {
	// ID 管理员唯一标识
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// Username 管理员登录用户名，唯一
	Username string `gorm:"size:50;not null;uniqueIndex" json:"username"`

	// PasswordHash 密码 bcrypt 哈希，不在响应中返回
	PasswordHash string `gorm:"size:255;not null" json:"-"`

	// Role 角色：SUPER_ADMIN（超级管理员）、ADMIN（普通管理员）
	// 超级管理员可以管理其他管理员账号
	Role string `gorm:"size:20;not null;default:'ADMIN'" json:"role"`

	// IsActive 是否启用，禁用后无法登录
	IsActive bool `gorm:"default:true" json:"is_active"`

	// LastLoginAt 最后登录时间
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}
