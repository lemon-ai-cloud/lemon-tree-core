// Package repository 提供数据访问层功能
// 负责与数据库交互，执行 CRUD 操作
package repository

import (
	"context"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApplicationMcpServerToolRepository ApplicationMCP工具 数据访问层接口
// 定义 ApplicationMCP工具 相关的数据库操作方法
type ApplicationMcpServerToolRepository interface {
	// Create 创建新的 ApplicationMCP工具 记录
	Create(ctx context.Context, tool *models.ApplicationMcpServerTool) error

	// GetByID 根据ID获取 ApplicationMCP工具 记录
	GetByID(ctx context.Context, id uuid.UUID) (*models.ApplicationMcpServerTool, error)

	// Update 更新 ApplicationMCP工具 记录
	Update(ctx context.Context, tool *models.ApplicationMcpServerTool) error

	// Delete 删除 ApplicationMCP工具 记录
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByApplicationMcpServerConfigID 根据MCP配置ID获取工具列表
	GetByApplicationMcpServerConfigID(ctx context.Context, configID uuid.UUID) ([]*models.ApplicationMcpServerTool, error)

	// GetByConfigIDAndName 根据MCP配置ID和工具名称获取工具
	GetByConfigIDAndName(ctx context.Context, configID uuid.UUID, name string) (*models.ApplicationMcpServerTool, error)

	// DeleteByApplicationMcpServerConfigID 根据MCP配置ID删除所有工具
	DeleteByApplicationMcpServerConfigID(ctx context.Context, configID uuid.UUID) error

	// BatchCreate 批量创建工具记录
	BatchCreate(ctx context.Context, tools []*models.ApplicationMcpServerTool) error

	// GetByApplicationID 根据应用ID获取所有MCP工具
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationMcpServerTool, error)
}

// applicationMcpServerToolRepository ApplicationMCP工具 数据访问层实现
// 实现 ApplicationMcpServerToolRepository 接口
type applicationMcpServerToolRepository struct {
	db *gorm.DB // 数据库连接
}

// NewApplicationMcpServerToolRepository 创建 ApplicationMCP工具 Repository 实例
// 返回 ApplicationMcpServerToolRepository 接口的实现
// 参数：db - 数据库连接
func NewApplicationMcpServerToolRepository(db *gorm.DB) ApplicationMcpServerToolRepository {
	return &applicationMcpServerToolRepository{
		db: db,
	}
}

// Create 创建新的 ApplicationMCP工具 记录
// 在数据库中插入新的 ApplicationMCP工具 记录
// 参数：ctx - 上下文，tool - ApplicationMCP工具 模型
// 返回：错误信息
func (r *applicationMcpServerToolRepository) Create(ctx context.Context, tool *models.ApplicationMcpServerTool) error {
	return r.db.WithContext(ctx).Create(tool).Error
}

// GetByID 根据ID获取 ApplicationMCP工具 记录
// 从数据库中查询指定ID的 ApplicationMCP工具 记录
// 参数：ctx - 上下文，id - ApplicationMCP工具 ID
// 返回：ApplicationMCP工具 模型和错误信息
func (r *applicationMcpServerToolRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ApplicationMcpServerTool, error) {
	var tool models.ApplicationMcpServerTool
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tool).Error
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

// Update 更新 ApplicationMCP工具 记录
// 在数据库中更新指定的 ApplicationMCP工具 记录
// 参数：ctx - 上下文，tool - ApplicationMCP工具 模型
// 返回：错误信息
func (r *applicationMcpServerToolRepository) Update(ctx context.Context, tool *models.ApplicationMcpServerTool) error {
	return r.db.WithContext(ctx).Save(tool).Error
}

// Delete 删除 ApplicationMCP工具 记录
// 从数据库中删除指定ID的 ApplicationMCP工具 记录
// 参数：ctx - 上下文，id - ApplicationMCP工具 ID
// 返回：错误信息
func (r *applicationMcpServerToolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ApplicationMcpServerTool{}, id).Error
}

// GetByApplicationMcpServerConfigID 根据MCP配置ID获取工具列表
// 从数据库中查询指定MCP配置下的所有工具记录
// 参数：ctx - 上下文，configID - MCP配置ID
// 返回：工具列表和错误信息
func (r *applicationMcpServerToolRepository) GetByApplicationMcpServerConfigID(ctx context.Context, configID uuid.UUID) ([]*models.ApplicationMcpServerTool, error) {
	var tools []*models.ApplicationMcpServerTool
	err := r.db.WithContext(ctx).Where("application_mcp_server_config_id = ?", configID).Find(&tools).Error
	if err != nil {
		return nil, err
	}
	return tools, nil
}

// GetByConfigIDAndName 根据MCP配置ID和工具名称获取工具
// 从数据库中查询指定MCP配置和工具名称的工具记录
// 参数：ctx - 上下文，configID - MCP配置ID，name - 工具名称
// 返回：工具模型和错误信息
func (r *applicationMcpServerToolRepository) GetByConfigIDAndName(ctx context.Context, configID uuid.UUID, name string) (*models.ApplicationMcpServerTool, error) {
	var tool models.ApplicationMcpServerTool
	err := r.db.WithContext(ctx).Where("application_mcp_server_config_id = ? AND name = ?", configID, name).First(&tool).Error
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

// DeleteByApplicationMcpServerConfigID 根据MCP配置ID删除所有工具
// 从数据库中删除指定MCP配置下的所有工具记录
// 参数：ctx - 上下文，configID - MCP配置ID
// 返回：错误信息
func (r *applicationMcpServerToolRepository) DeleteByApplicationMcpServerConfigID(ctx context.Context, configID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("application_mcp_server_config_id = ?", configID).Delete(&models.ApplicationMcpServerTool{}).Error
}

// BatchCreate 批量创建工具记录
// 在数据库中批量插入工具记录
// 参数：ctx - 上下文，tools - 工具模型列表
// 返回：错误信息
func (r *applicationMcpServerToolRepository) BatchCreate(ctx context.Context, tools []*models.ApplicationMcpServerTool) error {
	if len(tools) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(tools, 100).Error
}

// GetByApplicationID 根据应用ID获取所有MCP工具
// 从数据库中查询指定应用下的所有MCP工具
// 参数：ctx - 上下文，applicationID - 应用ID
// 返回：工具列表和错误信息
func (r *applicationMcpServerToolRepository) GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationMcpServerTool, error) {
	var tools []*models.ApplicationMcpServerTool
	err := r.db.WithContext(ctx).Where("application_id = ?", applicationID).Find(&tools).Error
	if err != nil {
		return nil, err
	}
	return tools, nil
}
