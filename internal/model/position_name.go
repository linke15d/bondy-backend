package model

import "github.com/linke15d/bondy-backend/pkg/timeformat"

// PositionName 姿势多语言名称表
// 对应数据库表名: position_names
type PositionName struct {
	ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PositionID string `gorm:"type:uuid;not null;index" json:"position_id"`
	// LanguageCode 语言代码，如 zh-CN / en
	LanguageCode string `gorm:"size:10;not null" json:"language_code" example:"zh-CN"`
	// Name 该语言下的姿势名称
	Name      string               `gorm:"size:50;not null" json:"name" example:"传教士"`
	CreatedAt timeformat.LocalTime `json:"created_at"`
}
