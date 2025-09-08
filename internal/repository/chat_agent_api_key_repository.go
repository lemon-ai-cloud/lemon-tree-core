// Package repository 提供数据访问层功能
// 负责与数据库的交互，实现数据持久化操作
package repository

import (
	"context"
	"lemon-tree-core/internal/base"
	"lemon-tree-core/internal/models"

	"gorm.io/gorm"
)

// ChatAgentApiKeyRepository ChatAgentApiKey 数据访问层接口
// 定义了 ChatAgentApiKey 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type ChatAgentApiKeyRepository interface {
	base.BaseRepository[models.ChatAgentApiKey]                                      // 继承基础仓库接口
	GetByApiKey(ctx context.Context, apiKey string) (*models.ChatAgentApiKey, error) // 通过 apiKey 获取 ChatAgentApiKey
}

// chatAgentApiKeyRepository ChatAgentApiKey 数据访问层实现
// 实现了 ChatAgentApiKeyRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type chatAgentApiKeyRepository struct {
	base.BaseRepository[models.ChatAgentApiKey]          // 组合基础仓库实现
	db                                          *gorm.DB // GORM 数据库连接实例
}

// NewChatAgentApiKeyRepository 创建 ChatAgentApiKey Repository 实例
// 返回 ChatAgentApiKeyRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewChatAgentApiKeyRepository(db *gorm.DB) ChatAgentApiKeyRepository {
	return &chatAgentApiKeyRepository{
		BaseRepository: base.NewBaseRepository[models.ChatAgentApiKey](db),
	}
}

// GetByApiKey 通过 apiKey 获取 ChatAgentApiKey
// 参数：ctx - 上下文对象, apiKey - 待查询的 apiKey
func (r *chatAgentApiKeyRepository) GetByApiKey(ctx context.Context, apiKey string) (*models.ChatAgentApiKey, error) {
	var chatAgentApiKey models.ChatAgentApiKey
	err := r.db.WithContext(ctx).Where("api_key = ? AND deleted_at IS NULL", apiKey).First(&chatAgentApiKey).Error
	if err != nil {
		return nil, err
	}
	return &chatAgentApiKey, nil
}
