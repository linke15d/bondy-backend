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
	coupleHandler *handler.CoupleHandler, // 新增
) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 无需登录
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// 需要登录
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.POST("/auth/logout", authHandler.Logout)

		// 用户信息
		user := protected.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
		}

		// 伴侣配对
		couple := protected.Group("/couple")
		{
			couple.POST("/invite", coupleHandler.GenerateInviteCode) // 生成邀请码
			couple.POST("/bind", coupleHandler.BindPartner)          // 绑定伴侣
			couple.GET("/info", coupleHandler.GetCoupleInfo)         // 获取伴侣信息
			couple.DELETE("/unbind", coupleHandler.Unlink)           // 解除绑定
		}
	}
}
