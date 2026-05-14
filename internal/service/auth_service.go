// Package service 业务逻辑层
// 处理所有业务规则，调用 repository 层访问数据
// handler 层只做参数解析和响应格式化，不写业务逻辑
package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
)

// AuthService 认证业务逻辑服务
// 负责注册、登录、登出、刷新 token 等功能
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *jwtpkg.Manager
}

// NewAuthService 创建 AuthService 实例
func NewAuthService(userRepo *repository.UserRepository, jwtManager *jwtpkg.Manager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// RegisterInput 注册接口请求参数
type RegisterInput struct {
	// Email 注册邮箱，必须是合法的邮箱格式，注册后作为登录账号
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	// Password 登录密码，最少8位，建议包含字母和数字
	Password string `json:"password" binding:"required,min=8" example:"password123"`
	// Nickname 用户昵称，显示在 App 界面，2到20个字符
	Nickname string `json:"nickname" binding:"required,min=2,max=20" example:"小明"`
	// Gender 用户性别，app端选择
	Gender string `json:"gender" binding:"required,oneof=male female other"`
}

// LoginInput 登录接口请求参数
type LoginInput struct {
	// Email 登录邮箱
	Email string `json:"email" binding:"required,email" example:"user@example.com"`

	// Password 登录密码
	Password string `json:"password" binding:"required" example:"password123"`
}

// AuthResult 认证成功后的统一返回结构
// 注册和登录都返回此结构
type AuthResult struct {
	// AccessToken 短期访问令牌，有效期15分钟
	// 调用需要登录的接口时放在 Header: Authorization: Bearer <access_token>
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// RefreshToken 长期刷新令牌，有效期30天
	// Access Token 过期后，用此 token 换取新的 token 对
	// 请安全存储，不要暴露
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// User 当前登录的用户信息
	User model.User `json:"user"`
}

// Register 用户注册
func (s *AuthService) Register(input RegisterInput) (*AuthResult, error) {
	existing, err := s.userRepo.FindByEmail(input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("server_error")
	}
	if existing != nil {
		return nil, errors.New("email_registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, errors.New("server_error")
	}

	hashStr := string(hash)
	user := &model.User{
		Email:        &input.Email,
		PasswordHash: &hashStr,
		Nickname:     &input.Nickname,
		Gender:       input.Gender,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("register_failed")
	}

	return s.generateTokenPair(user)
}

// Login 用户登录
// 流程：查找用户 → 检查封禁状态 → 验证密码 → 生成 token 对
func (s *AuthService) Login(input LoginInput) (*AuthResult, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		// 返回模糊错误，防止攻击者通过错误信息枚举已注册邮箱
		return nil, errors.New("邮箱或密码错误")
	}

	if user.IsBlocked {
		return nil, errors.New("账号已被禁用，请联系客服")
	}

	if user.PasswordHash == nil {
		return nil, errors.New("邮箱或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	return s.generateTokenPair(user)
}

// RefreshToken 刷新 Access Token
// 采用 Rotation 策略：每次刷新都删除旧 token 并生成新 token 对
// 这样即使 Refresh Token 泄露，使用一次后旧 token 立即失效
func (s *AuthService) RefreshToken(refreshToken string) (*AuthResult, error) {
	claims, err := s.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("token 无效或已过期")
	}

	rt, err := s.userRepo.FindRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("token 不存在或已被使用")
	}

	if rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token 已过期，请重新登录")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 删除旧 token
	_ = s.userRepo.DeleteRefreshToken(refreshToken)

	return s.generateTokenPair(user)
}

// Logout 用户登出
// 删除数据库中的 Refresh Token，使其永久失效
// Access Token 本身无法撤销，依赖其短暂的有效期（15分钟）自然过期
func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(refreshToken)
}

// generateTokenPair 生成 Access Token + Refresh Token 对
// 并将 Refresh Token 持久化到数据库
func (s *AuthService) generateTokenPair(user *model.User) (*AuthResult, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, errors.New("token 生成失败")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("token 生成失败")
	}

	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.userRepo.SaveRefreshToken(rt); err != nil {
		return nil, errors.New("token 保存失败")
	}

	return &AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}
