// Package main Bondy API 服务入口
//
//	@title			Bondy API
//	@version		1.0
//	@description	Bondy 伴侣亲密记录 App 后端接口文档。调用需要登录的接口时，在 Header 中加入 Authorization: Bearer <access_token>
//	@host			localhost:8080
//	@BasePath		/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				登录后将 access_token 填入此处，格式：Bearer <token>
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/linke15d/bondy-backend/docs"
	"github.com/linke15d/bondy-backend/internal/config"
	"github.com/linke15d/bondy-backend/internal/database"
	"github.com/linke15d/bondy-backend/internal/handler"
	"github.com/linke15d/bondy-backend/internal/repository"
	"github.com/linke15d/bondy-backend/internal/router"
	"github.com/linke15d/bondy-backend/internal/service"
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

	userRepo := repository.NewUserRepository(db)
	coupleRepo := repository.NewCoupleRepository(db)
	recordRepo := repository.NewRecordRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo)
	coupleService := service.NewCoupleService(coupleRepo, userRepo)
	recordService := service.NewRecordService(recordRepo, coupleRepo)
	statsService := service.NewStatsService(statsRepo, coupleRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	coupleHandler := handler.NewCoupleHandler(coupleService)
	recordHandler := handler.NewRecordHandler(recordService)
	statsHandler := handler.NewStatsHandler(statsService)

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "env": cfg.App.Env})
	})
	router.Setup(r, jwtManager, authHandler, userHandler, coupleHandler, recordHandler, statsHandler)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("服务启动，监听端口 %s", addr)
	log.Printf("Swagger 文档: http://localhost%s/swagger/index.html", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
