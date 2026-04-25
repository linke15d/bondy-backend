// Package router 路由注册与分组管理
package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/linke15d/bondy-backend/internal/handler"
	"github.com/linke15d/bondy-backend/internal/middleware"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
)

// Setup 注册所有路由和中间件
func Setup(r *gin.Engine, jwtManager *jwtpkg.Manager, authHandler *handler.AuthHandler) {
	// Swagger 文档（开发环境访问 /swagger/index.html）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 无需登录的路由
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// 需要登录的路由
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.POST("/auth/logout", authHandler.Logout)
	}
}
