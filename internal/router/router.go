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
func Setup(
	r *gin.Engine,
	jwtManager *jwtpkg.Manager,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// ── 无需登录 ──
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// ── 需要登录 ──
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// 认证
		protected.POST("/auth/logout", authHandler.Logout)

		// 用户信息
		user := protected.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)      // 获取个人信息
			user.PUT("/profile", userHandler.UpdateProfile)   // 更新个人信息
			user.PUT("/password", userHandler.ChangePassword) // 修改密码
		}
	}
}
