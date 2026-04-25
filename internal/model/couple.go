// Package model 数据库模型定义
package model

import "time"

// Couple 伴侣关系表
// 记录两个用户之间的绑定关系，一个用户同时只能有一段有效的伴侣关系
// 对应数据库表名: couples
type Couple struct {
	// ID 关系唯一标识，UUID 格式
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// User1ID 发起绑定的用户 ID（生成邀请码的一方）
	User1ID string `gorm:"type:uuid;not null;uniqueIndex" json:"user1_id"`

	// User2ID 接受绑定的用户 ID（输入邀请码的一方）
	User2ID string `gorm:"type:uuid;not null;uniqueIndex" json:"user2_id"`

	// InviteCode 邀请码，6位大写字母+数字组合，有效期15分钟，使用后失效
	InviteCode *string `gorm:"size:10;uniqueIndex" json:"invite_code,omitempty"`

	// InviteExpiresAt 邀请码过期时间，超过此时间邀请码自动失效
	InviteExpiresAt *time.Time `json:"invite_expires_at,omitempty"`

	// CreatedAt 绑定成功时间
	CreatedAt time.Time `json:"created_at"`

	// UnlinkedAt 解除绑定时间，不为 null 表示已解绑（保留历史记录）
	UnlinkedAt *time.Time `json:"unlinked_at,omitempty"`

	// User1 关联的用户1对象，仅用于联表查询
	User1 User `gorm:"foreignKey:User1ID" json:"user1,omitempty"`

	// User2 关联的用户2对象，仅用于联表查询
	User2 User `gorm:"foreignKey:User2ID" json:"user2,omitempty"`
}
