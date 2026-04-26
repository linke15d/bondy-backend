// Package admin 后台管理业务逻辑层
package admin

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
)

// AdminAuthService 管理员认证业务逻辑
type AdminAuthService struct {
	adminRepo  *repository.AdminRepository
	jwtManager *jwtpkg.Manager
}

// NewAdminAuthService 创建 AdminAuthService 实例
func NewAdminAuthService(adminRepo *repository.AdminRepository, jwtManager *jwtpkg.Manager) *AdminAuthService {
	return &AdminAuthService{
		adminRepo:  adminRepo,
		jwtManager: jwtManager,
	}
}

// AdminLoginInput 管理员登录请求参数
type AdminLoginInput struct {
	// Username 管理员用户名
	Username string `json:"username" binding:"required" example:"admin"`

	// Password 管理员密码
	Password string `json:"password" binding:"required" example:"admin123456"`
}

// AdminLoginResult 管理员登录返回结构
type AdminLoginResult struct {
	// AccessToken 访问令牌，有效期15分钟
	AccessToken string `json:"access_token"`

	// Admin 当前登录的管理员信息
	Admin model.Admin `json:"admin"`
}

// Login 管理员登录
func (s *AdminAuthService) Login(input AdminLoginInput) (*AdminLoginResult, error) {
	admin, err := s.adminRepo.FindByUsername(input.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 token
	token, err := s.jwtManager.GenerateAccessToken(admin.ID)
	if err != nil {
		return nil, errors.New("token 生成失败")
	}

	// 更新最后登录时间
	_ = s.adminRepo.UpdateLastLogin(admin.ID)

	return &AdminLoginResult{
		AccessToken: token,
		Admin:       *admin,
	}, nil
}

// CreateFirstAdmin 创建初始超级管理员
// 仅在系统初始化时调用一次
func (s *AdminAuthService) CreateFirstAdmin(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	admin := &model.Admin{
		Username:     username,
		PasswordHash: string(hash),
		Role:         "SUPER_ADMIN",
		IsActive:     true,
	}

	return s.adminRepo.Create(admin)
}
