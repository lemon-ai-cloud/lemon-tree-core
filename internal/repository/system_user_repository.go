// Package repository 提供数据访问层功能
// 负责与数据库的交互，实现数据持久化操作
package repository

import (
	"context"
	"lemon-tree-core/internal/base"
	"lemon-tree-core/internal/models"

	"gorm.io/gorm"
)

// SystemUserRepository SystemUser 数据访问层接口
// 定义了 SystemUser 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type SystemUserRepository interface {
	base.BaseRepository[models.SystemUser]                                      // 继承基础仓库接口
	GetByNumber(ctx context.Context, number string) (*models.SystemUser, error) // 根据用户账号获取用户
	GetByEmail(ctx context.Context, email string) (*models.SystemUser, error)   // 根据邮箱获取用户
}

// systemUserRepository SystemUser 数据访问层实现
// 实现了 SystemUserRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type systemUserRepository struct {
	base.BaseRepository[models.SystemUser]          // 组合基础仓库实现
	db                                     *gorm.DB // GORM 数据库连接实例
}

// NewSystemUserRepository 创建 SystemUser Repository 实例
// 返回 SystemUserRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewSystemUserRepository(db *gorm.DB) SystemUserRepository {
	return &systemUserRepository{
		BaseRepository: base.NewBaseRepository[models.SystemUser](db),
		db:             db,
	}
}

// GetByNumber 根据用户账号获取用户
// 根据用户账号查找并返回指定的用户信息
// 参数：ctx - 上下文，number - 用户账号
// 返回：用户对象和错误信息
func (r *systemUserRepository) GetByNumber(ctx context.Context, number string) (*models.SystemUser, error) {
	var user models.SystemUser
	err := r.db.WithContext(ctx).Where("number = ?", number).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
// 根据用户邮箱查找并返回指定的用户信息
// 参数：ctx - 上下文，email - 用户邮箱
// 返回：用户对象和错误信息
func (r *systemUserRepository) GetByEmail(ctx context.Context, email string) (*models.SystemUser, error) {
	var user models.SystemUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
