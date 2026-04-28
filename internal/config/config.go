package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	JWT   JWTConfig
}

type AppConfig struct {
	Env       string
	Port      string
	AdminPort string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type JWTConfig struct {
	AccessSecret             string
	RefreshSecret            string
	AccessExpireMinutes      int
	RefreshExpireDays        int
	AdminAccessExpireMinutes int
}

func Load() *Config {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用系统环境变量")
	}

	viper.AutomaticEnv()

	return &Config{
		App: AppConfig{
			Env:       viper.GetString("APP_ENV"),
			Port:      viper.GetString("APP_PORT"),
			AdminPort: viper.GetString("ADMIN_PORT"),
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
		JWT: JWTConfig{
			AccessSecret:             viper.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret:            viper.GetString("JWT_REFRESH_SECRET"),
			AccessExpireMinutes:      viper.GetInt("JWT_ACCESS_EXPIRE_MINUTES"),
			RefreshExpireDays:        viper.GetInt("JWT_REFRESH_EXPIRE_DAYS"),
			AdminAccessExpireMinutes: viper.GetInt("JWT_ADMIN_ACCESS_EXPIRE_MINUTES"),
		},
	}
}
