// Package handler HTTP 请求处理层
// 负责解析请求参数、调用 service 层、返回统一格式的 JSON 响应
// 不包含任何业务逻辑
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/linke15d/bondy-backend/internal/service"
	"github.com/linke15d/bondy-backend/pkg/i18n"
	"github.com/linke15d/bondy-backend/pkg/response"
)

// AuthHandler 认证相关请求处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建 AuthHandler 实例
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RefreshTokenInput 刷新/登出请求参数
type RefreshTokenInput struct {
	// RefreshToken 登录时获得的长期刷新令牌
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Register 用户注册
//
//	@Summary		用户注册
//	@Description	使用邮箱和密码创建新账号。注册成功后直接返回 token 对，无需再次调用登录接口。
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.RegisterInput						true	"注册信息"
//	@Success		201		{object}	response.Response{data=service.AuthResult}	"注册成功，返回 token 对和用户信息"
//	@Failure		400		{object}	response.Response							"参数格式错误 或 邮箱已被注册"
//	@Failure		500		{object}	response.Response							"服务器内部错误"
//	@Router			/api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	lang := c.GetString("lang")

	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok && len(validationErrors) > 0 {
			switch validationErrors[0].Field() {
			case "Email":
				response.BadRequest(c, i18n.Get("email_invalid", lang))
			case "Password":
				response.BadRequest(c, i18n.Get("password_too_short", lang))
			case "Nickname":
				response.BadRequest(c, i18n.Get("nickname_required", lang))
			case "Gender":
				response.BadRequest(c, i18n.Get("gender_invalid", lang))
			default:
				response.BadRequest(c, i18n.Get("invalid_params", lang))
			}
			return
		}
		response.BadRequest(c, i18n.Get("invalid_params", lang))
		return
	}

	result, err := h.authService.Register(input)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, i18n.Get(err.Error(), lang))
		return
	}

	response.Created(c, result)
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	使用邮箱和密码登录。返回 access_token（有效期15分钟）和 refresh_token（有效期30天）。access_token 用于调用其他接口，refresh_token 用于在 access_token 过期后换取新令牌。
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.LoginInput							true	"登录信息"
//	@Success		200		{object}	response.Response{data=service.AuthResult}	"登录成功"
//	@Failure		400		{object}	response.Response							"参数格式错误"
//	@Failure		401		{object}	response.Response							"邮箱或密码错误 或 账号已被封禁"
//	@Router			/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	lang := c.GetString("lang")

	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok && len(validationErrors) > 0 {
			switch validationErrors[0].Field() {
			case "Email":
				response.BadRequest(c, i18n.Get("email_required", lang))
			case "Password":
				response.BadRequest(c, i18n.Get("password_required", lang))
			default:
				response.BadRequest(c, i18n.Get("invalid_params", lang))
			}
			return
		}
		response.BadRequest(c, i18n.Get("invalid_params", lang))
		return
	}

	result, err := h.authService.Login(input)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, i18n.Get(err.Error(), lang))
		return
	}

	response.Success(c, result)
}

// RefreshToken 刷新令牌
//
//	@Summary		刷新 Access Token
//	@Description	当 access_token 过期时，使用 refresh_token 换取新的 token 对。采用 Rotation 策略，每次调用后旧的 refresh_token 立即失效，请使用新返回的 refresh_token 替换本地存储。
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RefreshTokenInput							true	"Refresh Token"
//	@Success		200		{object}	response.Response{data=service.AuthResult}	"刷新成功，返回新的 token 对"
//	@Failure		400		{object}	response.Response							"参数格式错误"
//	@Failure		401		{object}	response.Response							"refresh_token 无效、已过期或已被使用"
//	@Router			/api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	response.Success(c, result)
}

// Logout 用户登出
//
//	@Summary		用户登出
//	@Description	使当前 refresh_token 永久失效。登出后 access_token 仍在有效期内可使用（最多15分钟），客户端应同时清除本地存储的所有 token。
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {access_token}"
//	@Param			body			body		RefreshTokenInput	true	"要失效的 Refresh Token"
//	@Success		200				{object}	response.Response	"登出成功"
//	@Failure		400				{object}	response.Response	"参数格式错误"
//	@Failure		401				{object}	response.Response	"未登录或 token 已过期"
//	@Security		BearerAuth
//	@Router			/api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var input RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	_ = h.authService.Logout(input.RefreshToken)
	response.Success(c, nil)
}
