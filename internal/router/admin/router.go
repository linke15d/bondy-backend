// Package admin 后台管理路由
package admin

import (
	"github.com/gin-gonic/gin"
	_ "github.com/linke15d/bondy-backend/docs/admin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	adminHandler "github.com/linke15d/bondy-backend/internal/handler/admin"
	adminMiddleware "github.com/linke15d/bondy-backend/internal/middleware/admin"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
)

// Setup 注册后台管理所有路由
func Setup(
	r *gin.Engine,
	jwtManager *jwtpkg.Manager,
	authHandler *adminHandler.AdminAuthHandler,
	userHandler *adminHandler.AdminUserHandler,
	recordHandler *adminHandler.AdminRecordHandler,
	subHandler *adminHandler.AdminSubHandler,
	statsHandler *adminHandler.AdminStatsHandler,
	contentHandler *adminHandler.AdminContentHandler,
) {
	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.InstanceName("admin"),
	))

	v1 := r.Group("/admin/v1")

	// 无需登录
	v1.POST("/auth/login", authHandler.Login)

	// 需要管理员登录
	protected := v1.Group("")
	protected.Use(adminMiddleware.AuthMiddleware(jwtManager))
	{
		// 用户管理
		users := protected.Group("/users")
		{
			users.POST("/list", userHandler.ListUsers)
			users.POST("/detail", userHandler.GetUserDetail)
			users.POST("/block", userHandler.BlockUser)
			users.POST("/unblock", userHandler.UnblockUser)
			users.POST("/delete", userHandler.DeleteUser)
		}

		// 记录管理
		records := protected.Group("/records")
		{
			records.POST("/list", recordHandler.ListRecords)
		}

		// 内容管理
		content := protected.Group("/content")
		{
			// 姿势分类
			content.POST("/categories/create", contentHandler.CreatePositionCategory)
			content.POST("/categories/list", contentHandler.ListPositionCategories)
			content.POST("/categories/update", contentHandler.UpdatePositionCategory)
			content.POST("/categories/delete", contentHandler.DeletePositionCategory)

			// 标签
			content.POST("/tags/list", contentHandler.ListSystemTags)
			content.POST("/tags/create", contentHandler.CreateSystemTag)
			content.POST("/tags/delete", contentHandler.DeleteSystemTag)

			// 姿势
			content.POST("/positions/list", contentHandler.ListSystemPositions)
			content.POST("/positions/create", contentHandler.CreateSystemPosition)
			content.POST("/positions/delete", contentHandler.DeleteSystemPosition)

			// 图标上传
			// content.POST("/upload/icon", contentHandler.UploadIcon)
		}

		// 数据统计
		stats := protected.Group("/stats")
		{
			stats.POST("/dashboard", statsHandler.GetDashboard)
		}
	}
}
