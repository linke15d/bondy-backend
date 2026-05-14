// Package model 数据库模型定义
// 对应数据库表结构，使用 GORM tag 描述字段约束
package model

import (
	"time"

	"github.com/linke15d/bondy-backend/pkg/timeformat"
)

// User 用户表
// 存储所有注册用户的基本信息
// 对应数据库表名: users
type User struct {
	// ID 用户唯一标识，UUID 格式，由数据库自动生成
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Email 用户邮箱，唯一，可用于登录，注册时必填
	Email *string `gorm:"uniqueIndex;size:255" json:"email,omitempty" example:"user@example.com"`

	// Phone 用户手机号，唯一，可用于登录，可选
	Phone *string `gorm:"uniqueIndex;size:20" json:"phone,omitempty" example:"+8613800138000"`

	// PasswordHash 密码的 bcrypt 哈希值，不在 API 响应中返回
	PasswordHash *string `gorm:"size:255" json:"-"`

	// Nickname 用户昵称，显示在界面上的名称
	Nickname *string `gorm:"size:50" json:"nickname,omitempty" example:"小明"`

	// AvatarURL 用户头像图片地址
	AvatarURL *string `gorm:"size:500" json:"avatar_url,omitempty" example:"https://cdn.example.com/avatar.jpg"`

	// Birthday 用户生日，用于计算星座、年龄等
	Birthday *time.Time `json:"birthday,omitempty" example:"1995-06-15T00:00:00Z"`

	// IsVerified 邮箱或手机号是否已通过验证
	IsVerified bool `gorm:"default:false" json:"is_verified" example:"false"`

	// IsBlocked 账号是否被管理员封禁，封禁后无法登录
	IsBlocked bool `gorm:"default:false" json:"is_blocked" example:"false"`

	// FCMToken Firebase 推送通知 token，用于向设备发送消息，不在 API 响应中返回
	FCMToken *string `gorm:"size:500" json:"-"`

	//性别
	Gender string `gorm:"type:varchar(10)" json:"gender"`

	// CreatedAt 账号创建时间
	CreatedAt timeformat.LocalTime `json:"created_at" example:"2024-01-01 12:00:00"`

	// UpdatedAt 账号最后更新时间
	UpdatedAt timeformat.LocalTime `json:"updated_at" example:"2024-01-01 12:00:00"`

	// DeletedAt 软删除时间
	DeletedAt *timeformat.LocalTime `gorm:"index" json:"-"`
}
