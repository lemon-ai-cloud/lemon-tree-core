// Package repository 提供数据访问层功能
// 负责与数据库的交互，实现数据持久化操作
package repository

import (
	"context"
	"lemon-tree-core/internal/base"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SystemUserSessionRepository SystemUserSession 数据访问层接口
// 定义了 SystemUserSession 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type SystemUserSessionRepository interface {
	base.BaseRepository[models.SystemUserSession]                                           // 继承基础仓库接口
	GetByToken(ctx context.Context, token string) (*models.SystemUserSession, error)        // 根据Token获取会话
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.SystemUserSession, error) // 根据用户ID获取会话列表
	DeleteExpiredSessions(ctx context.Context) error                                        // 删除过期会话
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error                             // 根据用户ID删除所有会话
}

// systemUserSessionRepository SystemUserSession 数据访问层实现
// 实现了 SystemUserSessionRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type systemUserSessionRepository struct {
	base.BaseRepository[models.SystemUserSession]          // 组合基础仓库实现
	db                                            *gorm.DB // GORM 数据库连接实例
}

// NewSystemUserSessionRepository 创建 SystemUserSession Repository 实例
// 返回 SystemUserSessionRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewSystemUserSessionRepository(db *gorm.DB) SystemUserSessionRepository {
	return &systemUserSessionRepository{
		BaseRepository: base.NewBaseRepository[models.SystemUserSession](db),
		db:             db,
	}
}

// GetByToken 根据Token获取会话
// 根据Token查找并返回指定的会话信息
// 参数：ctx - 上下文，token - 会话Token
// 返回：会话对象和错误信息
func (r *systemUserSessionRepository) GetByToken(ctx context.Context, token string) (*models.SystemUserSession, error) {
	var session models.SystemUserSession
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByUserID 根据用户ID获取会话列表
// 根据用户ID查找并返回该用户的所有会话信息
// 参数：ctx - 上下文，userID - 用户ID
// 返回：会话列表和错误信息
func (r *systemUserSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.SystemUserSession, error) {
	var sessions []*models.SystemUserSession
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

// DeleteExpiredSessions 删除过期会话
// 删除所有已过期的会话记录
// 参数：ctx - 上下文
// 返回：错误信息
func (r *systemUserSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("login_expired_at < NOW()").Delete(&models.SystemUserSession{}).Error
}

// DeleteByUserID 根据用户ID删除所有会话
// 删除指定用户的所有会话记录
// 参数：ctx - 上下文，userID - 用户ID
// 返回：错误信息
func (r *systemUserSessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.SystemUserSession{}).Error
}
