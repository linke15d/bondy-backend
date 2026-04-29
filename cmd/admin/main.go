// Package main 后台管理服务入口
//
//	@title			Bondy Admin API
//	@version		1.0
//	@description	Bondy 后台管理系统接口文档，仅供内部使用
//	@host			localhost:3001
//	@BasePath		/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				管理员登录后将 access_token 填入，格式：Bearer <token>
package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	_ "github.com/linke15d/bondy-backend/docs/admin"
	"github.com/linke15d/bondy-backend/internal/config"
	"github.com/linke15d/bondy-backend/internal/database"
	adminHandler "github.com/linke15d/bondy-backend/internal/handler/admin"
	"github.com/linke15d/bondy-backend/internal/repository"
	adminRouter "github.com/linke15d/bondy-backend/internal/router/admin"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
)

func main() {
	cfg := config.Load()
	db := database.Init(&cfg.DB)

	jwtManager := jwtpkg.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AdminAccessExpireMinutes,
		cfg.JWT.RefreshExpireDays,
	)

	// 初始化 repository
	userRepo := repository.NewUserRepository(db)
	coupleRepo := repository.NewCoupleRepository(db)
	recordRepo := repository.NewRecordRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	// 初始化 admin service
	adminAuthService := adminService.NewAdminAuthService(adminRepo, jwtManager, rdb)
	adminUserService := adminService.NewAdminUserService(userRepo, coupleRepo)
	adminRecordService := adminService.NewAdminRecordService(recordRepo)
	adminSubService := adminService.NewAdminSubService(subRepo)
	adminStatsService := adminService.NewAdminStatsService(db)
	adminContentService := adminService.NewAdminContentService(db)

	// 初始化 admin handler
	adminAuthHandler := adminHandler.NewAdminAuthHandler(adminAuthService)
	adminUserHandler := adminHandler.NewAdminUserHandler(adminUserService)
	adminRecordHandler := adminHandler.NewAdminRecordHandler(adminRecordService)
	adminSubHandler := adminHandler.NewAdminSubHandler(adminSubService)
	adminStatsHandler := adminHandler.NewAdminStatsHandler(adminStatsService)
	adminContentHandler := adminHandler.NewAdminContentHandler(adminContentService)

	i18nService := adminService.NewI18nService(db)
	i18nHandler := adminHandler.NewI18nHandler(i18nService)

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 开发环境先放开，生产环境改成具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "admin"})
	})

	adminRouter.Setup(r, jwtManager, adminAuthHandler, adminUserHandler,
		adminRecordHandler, adminSubHandler, adminStatsHandler, adminContentHandler, i18nHandler)

	addr := fmt.Sprintf(":%s", cfg.App.AdminPort)
	log.Printf("后台管理服务启动，监听端口 %s", addr)
	log.Printf("Swagger 文档: http://localhost%s/swagger/index.html", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
