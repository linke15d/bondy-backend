// Package admin 后台管理中间件
package admin

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// AuthMiddleware 后台管理员鉴权中间件
// 和 App 端使用相同的 JWT 验证逻辑
// 但会额外把 adminID 写入 context 以区分普通用户
func AuthMiddleware(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, 401, "Authorization 格式错误")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseAccessToken(parts[1])
		if err != nil {
			response.Fail(c, 401, "token 无效或已过期")
			c.Abort()
			return
		}

		// 写入 adminID，后续 handler 通过 c.GetString("adminID") 获取
		c.Set("adminID", claims.UserID)
		c.Next()
	}
}
