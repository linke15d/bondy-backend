//go:build ignore

package main

import (
	"fmt"
	"log"

	"github.com/linke15d/bondy-backend/internal/config"
	"github.com/linke15d/bondy-backend/internal/database"
	"github.com/linke15d/bondy-backend/internal/repository"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
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

	adminRepo := repository.NewAdminRepository(db)
	authService := adminService.NewAdminAuthService(adminRepo, jwtManager)

	if err := authService.CreateFirstAdmin("admin", "Admin@123456"); err != nil {
		log.Fatalf("创建管理员失败: %v", err)
	}

	fmt.Println("管理员账号创建成功")
	fmt.Println("用户名: admin")
	fmt.Println("密码: admin@1234")
	fmt.Println("请登录后立即修改密码！")
}
