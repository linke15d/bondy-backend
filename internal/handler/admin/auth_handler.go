// Package admin 后台管理 Handler 层
package admin

import (
	"github.com/gin-gonic/gin"
	_ "github.com/linke15d/bondy-backend/internal/model"
	adminService "github.com/linke15d/bondy-backend/internal/service/admin"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// AdminAuthHandler 管理员认证请求处理器
type AdminAuthHandler struct {
	authService *adminService.AdminAuthService
}

// NewAdminAuthHandler 创建 AdminAuthHandler 实例
func NewAdminAuthHandler(authService *adminService.AdminAuthService) *AdminAuthHandler {
	return &AdminAuthHandler{authService: authService}
}

// Login 管理员登录
//
//	@Summary		管理员登录
//	@Description	使用用户名和密码登录后台管理系统，返回 access_token
//	@Tags			后台管理-认证
//	@Accept			json
//	@Produce		json
//	@Param			body	body		adminService.AdminLoginInput						true	"登录信息"
//	@Success		200		{object}	response.Response{data=adminService.AdminLoginResult}	"登录成功"
//	@Failure		401		{object}	response.Response										"用户名或密码错误"
//	@Router			/admin/v1/auth/login [post]
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var input adminService.AdminLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Login(input)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}

	response.Success(c, result)
}
