// Package repository 数据访问层
// 负责所有数据库操作，不包含业务逻辑
// 上层 service 通过 repository 接口访问数据
package repository

import (
	"github.com/linke15d/bondy-backend/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问对象
// 封装所有与 users 表和 refresh_tokens 表相关的数据库操作
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建 UserRepository 实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建新用户
// 插入一条用户记录，ID 由数据库 gen_random_uuid() 自动生成
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByEmail 通过邮箱查找用户
// 只返回未被软删除的用户
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone 通过手机号查找用户
// 只返回未被软删除的用户
func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ? AND deleted_at IS NULL", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 通过用户 ID 查找用户
// 只返回未被软删除的用户
func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SaveRefreshToken 保存 Refresh Token 到数据库
// 每次登录都会生成新的 Refresh Token 并持久化
func (r *UserRepository) SaveRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

// FindRefreshToken 通过 token 字符串查找 RefreshToken 记录
func (r *UserRepository) FindRefreshToken(token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.Where("token = ?", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// DeleteRefreshToken 删除指定的 Refresh Token
// 用于登出时使 token 失效
func (r *UserRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}

// DeleteAllUserRefreshTokens 删除某用户的所有 Refresh Token
// 用于修改密码、封号等需要强制下线的场景
func (r *UserRepository) DeleteAllUserRefreshTokens(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}

// Update 更新用户信息
// 只更新模型中非零值字段（使用 Save 全量更新）
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// AdminList 管理端获取用户列表，支持关键词搜索和封禁状态过滤
func (r *UserRepository) AdminList(keyword string, isBlocked *bool, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{}).Where("deleted_at IS NULL")

	if keyword != "" {
		query = query.Where("email LIKE ? OR nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	if isBlocked != nil {
		query = query.Where("is_blocked = ?", *isBlocked)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	return users, total, err
}
