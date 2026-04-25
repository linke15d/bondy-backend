// Package response 统一 HTTP 响应格式
// 所有接口返回值都通过此包的方法输出，保证格式一致
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	// Code 业务状态码，0 表示成功，其他值表示错误
	Code int `json:"code" example:"0"`
	// Message 提示信息
	Message string `json:"message" example:"success"`
	// Data 响应数据，失败时为 null
	Data interface{} `json:"data,omitempty"`
}

// Success 返回 200 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created 返回 201 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Fail 返回指定 HTTP 状态码的错误响应
func Fail(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Message: message,
	})
}

// BadRequest 返回 400 参数错误
func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, message)
}

// Unauthorized 返回 401 未授权
func Unauthorized(c *gin.Context) {
	Fail(c, http.StatusUnauthorized, "未授权，请先登录")
}

// Forbidden 返回 403 无权限
func Forbidden(c *gin.Context) {
	Fail(c, http.StatusForbidden, "无权限访问")
}

// NotFound 返回 404 资源不存在
func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, message)
}

// ServerError 返回 500 服务器错误
func ServerError(c *gin.Context) {
	Fail(c, http.StatusInternalServerError, "服务器内部错误")
}
