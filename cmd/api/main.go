package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/linke15d/bondy-backend/internal/config"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置 gin 模式
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    cfg.App.Env,
		})
	})

	// 启动服务
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("服务启动，监听端口 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
