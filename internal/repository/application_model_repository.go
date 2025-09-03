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

// ApplicationLlmRepository ApplicationLlm 数据访问层接口
// 定义了 ApplicationLlm 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type ApplicationLlmRepository interface {
	base.BaseRepository[models.ApplicationLlm] // 继承基础仓库实现

	// DeleteByProviderID 根据提供商ID删除模型记录
	// 删除指定提供商下的所有模型
	DeleteByProviderID(ctx context.Context, providerID uuid.UUID) error

	// GetByProviderID 根据提供商ID获取模型列表
	// 返回指定提供商下的所有模型
	GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.ApplicationLlm, error)

	// GetByApplicationID 根据应用ID获取模型列表
	// 返回指定应用下的所有模型
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationLlm, error)
}

// applicationLlmRepository ApplicationLlm 数据访问层实现
// 实现了 ApplicationLlmRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type applicationLlmRepository struct {
	base.BaseRepository[models.ApplicationLlm]          // 组合基础仓库实现
	db                                         *gorm.DB // 直接保存数据库连接引用
}

// NewApplicationLlmRepository 创建 ApplicationLlm Repository 实例
// 返回 ApplicationLlmRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewApplicationLlmRepository(db *gorm.DB) ApplicationLlmRepository {
	return &applicationLlmRepository{
		BaseRepository: base.NewBaseRepository[models.ApplicationLlm](db),
		db:             db,
	}
}

// DeleteByProviderID 根据提供商ID删除模型记录
// 删除指定提供商下的所有模型
func (r *applicationLlmRepository) DeleteByProviderID(ctx context.Context, providerID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("llm_provider_id = ?", providerID).Delete(&models.ApplicationLlm{}).Error
}

// GetByProviderID 根据提供商ID获取模型列表
// 返回指定提供商下的所有模型
func (r *applicationLlmRepository) GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.ApplicationLlm, error) {
	var models []*models.ApplicationLlm
	err := r.db.WithContext(ctx).Where("llm_provider_id = ?", providerID).Find(&models).Error
	return models, err
}

// GetByApplicationID 根据应用ID获取模型列表
// 返回指定应用下的所有模型
func (r *applicationLlmRepository) GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationLlm, error) {
	var models []*models.ApplicationLlm
	err := r.db.WithContext(ctx).Where("application_id = ?", applicationID).Find(&models).Error
	return models, err
}
