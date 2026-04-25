// Package jwt 提供 JWT token 的生成与解析功能
// 支持 Access Token（短期）和 Refresh Token（长期）两种令牌
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT 载荷
// 在标准 RegisteredClaims 基础上增加 UserID 字段
type Claims struct {
	// UserID 用户唯一标识
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Manager JWT 管理器，持有签名密钥和过期配置
type Manager struct {
	accessSecret  string        // Access Token 签名密钥
	refreshSecret string        // Refresh Token 签名密钥
	accessExpire  time.Duration // Access Token 有效期（默认15分钟）
	refreshExpire time.Duration // Refresh Token 有效期（默认30天）
}

// NewManager 创建 JWT 管理器
// accessSecret: Access Token 签名密钥
// refreshSecret: Refresh Token 签名密钥
// accessMinutes: Access Token 有效分钟数
// refreshDays: Refresh Token 有效天数
func NewManager(accessSecret, refreshSecret string, accessMinutes, refreshDays int) *Manager {
	return &Manager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpire:  time.Duration(accessMinutes) * time.Minute,
		refreshExpire: time.Duration(refreshDays) * 24 * time.Hour,
	}
}

// GenerateAccessToken 生成短期 Access Token
// 用于 API 请求鉴权，有效期较短（默认15分钟）
func (m *Manager) GenerateAccessToken(userID string) (string, error) {
	return m.generate(userID, m.accessSecret, m.accessExpire)
}

// GenerateRefreshToken 生成长期 Refresh Token
// 用于刷新 Access Token，有效期较长（默认30天）
func (m *Manager) GenerateRefreshToken(userID string) (string, error) {
	return m.generate(userID, m.refreshSecret, m.refreshExpire)
}

// generate 内部通用 token 生成方法
func (m *Manager) generate(userID, secret string, expire time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseAccessToken 解析并验证 Access Token
// 返回 Claims 或错误（token 过期、签名无效等）
func (m *Manager) ParseAccessToken(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, m.accessSecret)
}

// ParseRefreshToken 解析并验证 Refresh Token
func (m *Manager) ParseRefreshToken(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, m.refreshSecret)
}

// parse 内部通用 token 解析方法
func (m *Manager) parse(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的 token")
	}
	return claims, nil
}
