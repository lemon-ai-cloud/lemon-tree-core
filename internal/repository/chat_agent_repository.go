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

// ChatAgentRepository ChatAgent 数据访问层接口
// 定义了 ChatAgent 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type ChatAgentRepository interface {
	base.BaseRepository[models.ChatAgent] // 继承基础仓库接口

	// GetByApplicationIDWithPagination 根据应用ID获取智能体列表（分页）
	GetByApplicationIDWithPagination(ctx context.Context, applicationID uuid.UUID, page, pageSize int) ([]*models.ChatAgent, int64, error)
}

// chatAgentRepository ChatAgent 数据访问层实现
// 实现了 ChatAgentRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type chatAgentRepository struct {
	base.BaseRepository[models.ChatAgent]          // 组合基础仓库实现
	db                                    *gorm.DB // 数据库连接
}

// NewChatAgentRepository 创建 ChatAgent Repository 实例
// 返回 ChatAgentRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewChatAgentRepository(db *gorm.DB) ChatAgentRepository {
	return &chatAgentRepository{
		BaseRepository: base.NewBaseRepository[models.ChatAgent](db),
		db:             db,
	}
}

// GetByApplicationIDWithPagination 根据应用ID获取智能体列表（分页）
// 从数据库中查询指定应用下的智能体记录，支持分页
// 参数：ctx - 上下文，applicationID - 应用ID，page - 页码（从1开始），pageSize - 每页大小
// 返回：智能体列表、总数量和错误信息
func (r *chatAgentRepository) GetByApplicationIDWithPagination(ctx context.Context, applicationID uuid.UUID, page, pageSize int) ([]*models.ChatAgent, int64, error) {
	var agents []*models.ChatAgent
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.ChatAgent{}).Where("application_id = ?", applicationID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := r.db.WithContext(ctx).Where("application_id = ?", applicationID).Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}
