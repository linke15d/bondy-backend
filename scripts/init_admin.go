//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/linke15d/bondy-backend/internal/config"
	"github.com/linke15d/bondy-backend/internal/database"
	"github.com/linke15d/bondy-backend/internal/repository"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()
	db := database.Init(&cfg.DB)

	jwtManager := jwtpkg.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpireMinutes,
		cfg.JWT.RefreshExpireDays,
	)

	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})
	defer rdb.Close()

	// 测试 Redis 连接
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("Redis 连接失败（非致命）: %v", err)
	}

	adminRepo := repository.NewAdminRepository(db)
	authService := adminService.NewAdminAuthService(adminRepo, jwtManager, rdb)

	if err := authService.CreateFirstAdmin("admin", "Admin@123456"); err != nil {
		log.Fatalf("创建管理员失败: %v", err)
	}

	fmt.Println("管理员账号创建成功")
	fmt.Println("用户名: admin")
	fmt.Println("密码: Admin@123456")
	fmt.Println("请登录后立即修改密码！")
}
