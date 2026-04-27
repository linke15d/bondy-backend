// Package admin 后台管理业务逻辑层
package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/linke15d/bondy-backend/internal/model"
	"github.com/linke15d/bondy-backend/internal/repository"
	jwtpkg "github.com/linke15d/bondy-backend/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

// AdminAuthService 管理员认证业务逻辑
type AdminAuthService struct {
	adminRepo  *repository.AdminRepository
	jwtManager *jwtpkg.Manager
	redis      *redis.Client
}

// NewAdminAuthService 创建 AdminAuthService 实例
func NewAdminAuthService(adminRepo *repository.AdminRepository, jwtManager *jwtpkg.Manager, redis *redis.Client) *AdminAuthService {
	return &AdminAuthService{
		adminRepo:  adminRepo,
		jwtManager: jwtManager,
		redis:      redis,
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
	// 检查是否被锁定
	lockKey := fmt.Sprintf("admin:login:lock:%s", input.Username)
	locked, _ := s.redis.Get(context.Background(), lockKey).Result()
	if locked == "1" {
		return nil, errors.New("账号已被锁定，请15分钟后再试")
	}

	admin, err := s.adminRepo.FindByUsername(input.Username)
	if err != nil {
		s.incrLoginFail(input.Username)
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		s.incrLoginFail(input.Username)
		return nil, errors.New("用户名或密码错误")
	}

	// 登录成功，清除失败记录
	failKey := fmt.Sprintf("admin:login:fail:%s", input.Username)
	s.redis.Del(context.Background(), failKey)

	token, err := s.jwtManager.GenerateAccessToken(admin.ID)
	if err != nil {
		return nil, errors.New("token 生成失败")
	}

	_ = s.adminRepo.UpdateLastLogin(admin.ID)

	return &AdminLoginResult{
		AccessToken: token,
		Admin:       *admin,
	}, nil
}

// incrLoginFail 记录登录失败次数，超过5次锁定15分钟
func (s *AdminAuthService) incrLoginFail(username string) {
	ctx := context.Background()
	failKey := fmt.Sprintf("admin:login:fail:%s", username)
	lockKey := fmt.Sprintf("admin:login:lock:%s", username)

	count, _ := s.redis.Incr(ctx, failKey).Result()
	s.redis.Expire(ctx, failKey, 15*time.Minute)

	if count >= 5 {
		s.redis.Set(ctx, lockKey, "1", 15*time.Minute)
	}
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
