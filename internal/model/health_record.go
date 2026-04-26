// Package model 数据库模型定义
package model

import "time"

// HealthRecord 健康记录表
// 记录用户的 STI 检测和疫苗接种情况
// 健康记录属于个人私密数据，只有本人可见
// 对应数据库表名: health_records
type HealthRecord struct {
	// ID 记录唯一标识，UUID 格式
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// UserID 所属用户 ID，健康记录只属于个人，不共享给伴侣
	UserID string `gorm:"type:uuid;not null;index" json:"user_id"`

	// Type 记录类型：STI_TEST（性病检测）或 VACCINE（疫苗接种）
	Type string `gorm:"size:20;not null" json:"type" example:"STI_TEST"`

	// ItemName 检测/接种项目名称，如"HIV检测"、"HPV疫苗"
	ItemName string `gorm:"size:100;not null" json:"item_name" example:"HIV检测"`

	// Result 检测结果：NEGATIVE（阴性/正常）、POSITIVE（阳性）、UNKNOWN（未知）
	// 疫苗记录此字段填 COMPLETED（已完成）
	Result string `gorm:"size:20;not null" json:"result" example:"NEGATIVE"`

	// TestedAt 检测/接种日期
	TestedAt time.Time `gorm:"not null" json:"tested_at" example:"2024-01-15T00:00:00Z"`

	// NextRemindAt 下次提醒时间，到期后系统发送推送提醒
	NextRemindAt *time.Time `json:"next_remind_at,omitempty" example:"2024-07-15T00:00:00Z"`

	// NoteEncrypted 备注，客户端加密后存储，后端不解密
	NoteEncrypted *string `gorm:"type:text" json:"note_encrypted,omitempty"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`

	// User 关联用户，仅用于联表查询
	User User `gorm:"foreignKey:UserID" json:"-"`
}
