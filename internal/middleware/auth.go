// Package middleware Gin 中间件集合
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// AuthMiddleware JWT 鉴权中间件
// 从请求 Header 的 Authorization 字段提取并验证 Bearer Token
// 验证通过后将 userID 注入到 gin.Context 中
// 后续 handler 通过 c.GetString("userID") 获取当前登录用户 ID
func AuthMiddleware(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 格式校验：必须是 "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, 401, "Authorization 格式错误，应为：Bearer <token>")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseAccessToken(parts[1])
		if err != nil {
			response.Fail(c, 401, "token 无效或已过期，请重新登录")
			c.Abort()
			return
		}

		// 将 userID 写入 context，供后续 handler 使用
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
