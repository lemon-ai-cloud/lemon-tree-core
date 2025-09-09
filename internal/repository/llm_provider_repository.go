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

// LlmProviderRepository ApplicationLlmProvider 数据访问层接口
// 定义了 ApplicationLlmProvider 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type LlmProviderRepository interface {
	base.BaseRepository[models.ApplicationLlmProvider] // 继承基础仓库接口

	// GetByApplicationID 根据应用ID获取大语言模型提供商列表
	// 返回指定应用下的所有提供商
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationLlmProvider, error)
}

// llmProviderRepository ApplicationLlmProvider 数据访问层实现
// 实现了 LlmProviderRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type llmProviderRepository struct {
	base.BaseRepository[models.ApplicationLlmProvider]          // 组合基础仓库实现
	db                                                 *gorm.DB // 直接保存数据库连接引用
}

// NewLlmProviderRepository 创建 ApplicationLlmProvider Repository 实例
// 返回 LlmProviderRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewLlmProviderRepository(db *gorm.DB) LlmProviderRepository {
	return &llmProviderRepository{
		BaseRepository: base.NewBaseRepository[models.ApplicationLlmProvider](db),
		db:             db,
	}
}

// GetByApplicationID 根据应用ID获取大语言模型提供商列表
// 返回指定应用下的所有提供商
func (r *llmProviderRepository) GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationLlmProvider, error) {
	var llmProviders []*models.ApplicationLlmProvider
	err := r.db.WithContext(ctx).Where("application_id = ?", applicationID).Find(&llmProviders).Error
	return llmProviders, err
}
