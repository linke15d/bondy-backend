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
	coupleHandler *handler.CoupleHandler,
	recordHandler *handler.RecordHandler,
	statsHandler *handler.StatsHandler,
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
			user.POST("/profile/get", userHandler.GetProfile)         // 获取个人信息
			user.POST("/profile/update", userHandler.UpdateProfile)   // 更新个人信息
			user.POST("/password/change", userHandler.ChangePassword) // 修改密码
		}

		// 伴侣配对
		couple := protected.Group("/couple")
		{
			couple.POST("/invite", coupleHandler.GenerateInviteCode) // 生成邀请码
			couple.POST("/bind", coupleHandler.BindPartner)          // 绑定伴侣
			couple.POST("/info", coupleHandler.GetCoupleInfo)        // 获取伴侣信息
			couple.POST("/unbind", coupleHandler.Unlink)             // 解除绑定
		}

		// 亲密记录
		records := protected.Group("/records")
		{
			records.POST("/create", recordHandler.CreateRecord)    // 创建记录
			records.POST("/list", recordHandler.ListRecords)       // 获取列表
			records.POST("/detail", recordHandler.GetRecord)       // 获取详情
			records.POST("/update", recordHandler.UpdateRecord)    // 更新记录
			records.POST("/delete", recordHandler.DeleteRecord)    // 删除记录
			records.POST("/tags", recordHandler.GetTags)           // 获取标签列表
			records.POST("/positions", recordHandler.GetPositions) // 获取姿势列表
		}

		// 数据统计
		stats := protected.Group("/stats")
		{
			stats.POST("/overview", statsHandler.GetOverview)    // 总览统计
			stats.POST("/yearly", statsHandler.GetYearlyStats)   // 年度统计
			stats.POST("/monthly", statsHandler.GetMonthlyStats) // 月度统计
		}
	}
}
