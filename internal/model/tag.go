package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID          string         `gorm:"type:uuid;primaryKey" json:"id"`
	IconBase64  *string        `gorm:"type:text" json:"icon_base64"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	IsSystem    bool           `gorm:"default:false" json:"is_system"`
	DefaultName string         `gorm:"type:varchar(100)" json:"default_name"`
	Name        string         `gorm:"-" json:"name"` // 不存库，运行时填充
	Names       []TagName      `gorm:"foreignKey:TagID" json:"names"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type TagName struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	TagID        string    `gorm:"type:uuid;index" json:"tag_id"`
	LanguageCode string    `gorm:"type:varchar(10)" json:"language_code"`
	Name         string    `gorm:"type:varchar(100)" json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (tn *TagName) BeforeCreate(tx *gorm.DB) error {
	if tn.ID == "" {
		tn.ID = uuid.New().String()
	}
	return nil
}
