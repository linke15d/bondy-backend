package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/linke15d/bondy-backend/internal/config"
	"github.com/linke15d/bondy-backend/internal/model"
)

var DB *gorm.DB

func Init(cfg *config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动建表
	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Couple{},
		&model.Tag{},
		&model.PositionCategory{},
		&model.Position{},
		&model.Record{},
		&model.Wishlist{},
		&model.HealthRecord{},
		&model.Subscription{},
		&model.Admin{},
		&model.SupportedLanguage{},
		&model.PositionCategory{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库连接成功")
	DB = db
	return db
}
