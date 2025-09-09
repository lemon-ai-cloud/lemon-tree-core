// Package repository 提供数据访问层功能
// 负责与数据库交互，执行CRUD操作
package repository

import (
	"context"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatAgentMcpServerToolRepository ChatAgentMcpServerTool 数据访问层接口
// 定义 ChatAgentMcpServerTool 相关的数据库操作方法
type ChatAgentMcpServerToolRepository interface {
	// Create 创建 ChatAgentMcpServerTool 记录
	Create(ctx context.Context, chatAgentMcpServerTool *models.ChatAgentMcpServerTool) error

	// GetByID 根据ID获取 ChatAgentMcpServerTool 记录
	GetByID(ctx context.Context, id uuid.UUID) (*models.ChatAgentMcpServerTool, error)

	// Update 更新 ChatAgentMcpServerTool 记录
	Update(ctx context.Context, chatAgentMcpServerTool *models.ChatAgentMcpServerTool) error

	// DeleteByID 根据ID删除 ChatAgentMcpServerTool 记录
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// GetByChatAgentID 根据ChatAgentID获取所有工具配置
	GetByChatAgentID(ctx context.Context, chatAgentID uuid.UUID) ([]*models.ChatAgentMcpServerTool, error)

	// GetByChatAgentIDAndApplicationMcpServerToolID 根据ChatAgentID和ApplicationMcpServerToolID获取配置
	GetByChatAgentIDAndApplicationMcpServerToolID(ctx context.Context, chatAgentID, applicationMcpServerToolID uuid.UUID) (*models.ChatAgentMcpServerTool, error)

	// BatchCreate 批量创建 ChatAgentMcpServerTool 记录
	BatchCreate(ctx context.Context, chatAgentMcpServerTools []*models.ChatAgentMcpServerTool) error

	// BatchUpdate 批量更新 ChatAgentMcpServerTool 记录
	BatchUpdate(ctx context.Context, chatAgentMcpServerTools []*models.ChatAgentMcpServerTool) error

	// DeleteByChatAgentID 根据ChatAgentID删除所有相关记录
	DeleteByChatAgentID(ctx context.Context, chatAgentID uuid.UUID) error
}

// chatAgentMcpServerToolRepository ChatAgentMcpServerTool 数据访问层实现
// 实现 ChatAgentMcpServerToolRepository 接口
type chatAgentMcpServerToolRepository struct {
	db *gorm.DB // 数据库连接
}

// NewChatAgentMcpServerToolRepository 创建 ChatAgentMcpServerTool 数据访问层实例
// 返回 ChatAgentMcpServerToolRepository 接口的实现
// 参数：db - 数据库连接
func NewChatAgentMcpServerToolRepository(db *gorm.DB) ChatAgentMcpServerToolRepository {
	return &chatAgentMcpServerToolRepository{
		db: db,
	}
}

// Create 创建 ChatAgentMcpServerTool 记录
func (r *chatAgentMcpServerToolRepository) Create(ctx context.Context, chatAgentMcpServerTool *models.ChatAgentMcpServerTool) error {
	return r.db.WithContext(ctx).Create(chatAgentMcpServerTool).Error
}

// GetByID 根据ID获取 ChatAgentMcpServerTool 记录
func (r *chatAgentMcpServerToolRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ChatAgentMcpServerTool, error) {
	var chatAgentMcpServerTool models.ChatAgentMcpServerTool
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&chatAgentMcpServerTool).Error
	if err != nil {
		return nil, err
	}
	return &chatAgentMcpServerTool, nil
}

// Update 更新 ChatAgentMcpServerTool 记录
func (r *chatAgentMcpServerToolRepository) Update(ctx context.Context, chatAgentMcpServerTool *models.ChatAgentMcpServerTool) error {
	return r.db.WithContext(ctx).Save(chatAgentMcpServerTool).Error
}

// DeleteByID 根据ID删除 ChatAgentMcpServerTool 记录
func (r *chatAgentMcpServerToolRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ChatAgentMcpServerTool{}, id).Error
}

// GetByChatAgentID 根据ChatAgentID获取所有工具配置
func (r *chatAgentMcpServerToolRepository) GetByChatAgentID(ctx context.Context, chatAgentID uuid.UUID) ([]*models.ChatAgentMcpServerTool, error) {
	var chatAgentMcpServerTools []*models.ChatAgentMcpServerTool
	err := r.db.WithContext(ctx).Where("chat_agent_id = ?", chatAgentID).Find(&chatAgentMcpServerTools).Error
	if err != nil {
		return nil, err
	}
	return chatAgentMcpServerTools, nil
}

// GetByChatAgentIDAndApplicationMcpServerToolID 根据ChatAgentID和ApplicationMcpServerToolID获取配置
func (r *chatAgentMcpServerToolRepository) GetByChatAgentIDAndApplicationMcpServerToolID(ctx context.Context, chatAgentID, applicationMcpServerToolID uuid.UUID) (*models.ChatAgentMcpServerTool, error) {
	var chatAgentMcpServerTool models.ChatAgentMcpServerTool
	err := r.db.WithContext(ctx).Where("chat_agent_id = ? AND application_mcp_server_tool_id = ?", chatAgentID, applicationMcpServerToolID).First(&chatAgentMcpServerTool).Error
	if err != nil {
		return nil, err
	}
	return &chatAgentMcpServerTool, nil
}

// BatchCreate 批量创建 ChatAgentMcpServerTool 记录
func (r *chatAgentMcpServerToolRepository) BatchCreate(ctx context.Context, chatAgentMcpServerTools []*models.ChatAgentMcpServerTool) error {
	if len(chatAgentMcpServerTools) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(chatAgentMcpServerTools, 100).Error
}

// BatchUpdate 批量更新 ChatAgentMcpServerTool 记录
func (r *chatAgentMcpServerToolRepository) BatchUpdate(ctx context.Context, chatAgentMcpServerTools []*models.ChatAgentMcpServerTool) error {
	if len(chatAgentMcpServerTools) == 0 {
		return nil
	}

	// 使用事务进行批量更新
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, tool := range chatAgentMcpServerTools {
			if err := tx.Save(tool).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteByChatAgentID 根据ChatAgentID删除所有相关记录
func (r *chatAgentMcpServerToolRepository) DeleteByChatAgentID(ctx context.Context, chatAgentID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("chat_agent_id = ?", chatAgentID).Delete(&models.ChatAgentMcpServerTool{}).Error
}
