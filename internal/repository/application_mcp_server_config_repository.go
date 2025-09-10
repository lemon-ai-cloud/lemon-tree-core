// Package repository 提供数据访问层功能
// 负责与数据库交互，执行 CRUD 操作
package repository

import (
	"context"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApplicationMcpServerConfigRepository ApplicationMCP配置 数据访问层接口
// 定义 ApplicationMCP配置 相关的数据库操作方法
type ApplicationMcpServerConfigRepository interface {
	// Create 创建新的 ApplicationMCP配置 记录
	Create(ctx context.Context, config *models.ApplicationMcpServerConfig) error

	// GetByID 根据ID获取 ApplicationMCP配置 记录
	GetByID(ctx context.Context, id uuid.UUID) (*models.ApplicationMcpServerConfig, error)

	// Update 更新 ApplicationMCP配置 记录
	Update(ctx context.Context, config *models.ApplicationMcpServerConfig) error

	// Delete 删除 ApplicationMCP配置 记录
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByApplicationID 根据应用ID获取 ApplicationMCP配置 列表
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationMcpServerConfig, error)

	GetByConfigID(ctx context.Context, configID string) (*models.ApplicationMcpServerConfig, error)
}

// applicationMcpServerConfigRepository ApplicationMCP配置 数据访问层实现
// 实现 ApplicationMcpServerConfigRepository 接口
type applicationMcpServerConfigRepository struct {
	db *gorm.DB // 数据库连接
}

// NewApplicationMcpServerConfigRepository 创建 ApplicationMCP配置 Repository 实例
// 返回 ApplicationMcpServerConfigRepository 接口的实现
// 参数：db - 数据库连接
func NewApplicationMcpServerConfigRepository(db *gorm.DB) ApplicationMcpServerConfigRepository {
	return &applicationMcpServerConfigRepository{
		db: db,
	}
}

// Create 创建新的 ApplicationMCP配置 记录
// 在数据库中插入新的 ApplicationMCP配置 记录
// 参数：ctx - 上下文，config - ApplicationMCP配置 模型
// 返回：错误信息
func (r *applicationMcpServerConfigRepository) Create(ctx context.Context, config *models.ApplicationMcpServerConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// GetByID 根据ID获取 ApplicationMCP配置 记录
// 从数据库中查询指定ID的 ApplicationMCP配置 记录
// 参数：ctx - 上下文，id - ApplicationMCP配置 ID
// 返回：ApplicationMCP配置 模型和错误信息
func (r *applicationMcpServerConfigRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ApplicationMcpServerConfig, error) {
	var config models.ApplicationMcpServerConfig
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Update 更新 ApplicationMCP配置 记录
// 在数据库中更新指定的 ApplicationMCP配置 记录
// 参数：ctx - 上下文，config - ApplicationMCP配置 模型
// 返回：错误信息
func (r *applicationMcpServerConfigRepository) Update(ctx context.Context, config *models.ApplicationMcpServerConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// Delete 删除 ApplicationMCP配置 记录
// 从数据库中删除指定ID的 ApplicationMCP配置 记录
// 参数：ctx - 上下文，id - ApplicationMCP配置 ID
// 返回：错误信息
func (r *applicationMcpServerConfigRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ApplicationMcpServerConfig{}, id).Error
}

// GetByApplicationID 根据应用ID获取 ApplicationMCP配置 列表
// 从数据库中查询指定应用下的所有 ApplicationMCP配置 记录
// 参数：ctx - 上下文，applicationID - 应用ID
// 返回：ApplicationMCP配置 列表和错误信息
func (r *applicationMcpServerConfigRepository) GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationMcpServerConfig, error) {
	var configs []*models.ApplicationMcpServerConfig
	err := r.db.WithContext(ctx).Where("application_id = ?", applicationID).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// GetByConfigID 根据配置ID获取 ApplicationMCP配置
// 从数据库中查询指定配置ID的 ApplicationMCP配置
// 参数：ctx - 上下文，configID - 配置ID
// 返回：ApplicationMCP配置和错误信息
func (r *applicationMcpServerConfigRepository) GetByConfigID(ctx context.Context, configID string) (*models.ApplicationMcpServerConfig, error) {
	var config models.ApplicationMcpServerConfig
	err := r.db.WithContext(ctx).Where("config_id = ?", configID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}
