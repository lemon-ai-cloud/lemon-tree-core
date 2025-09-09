// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"fmt"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"

	"github.com/google/uuid"
)

// ChatAgentService 智能体 业务逻辑层接口
// 定义 智能体 相关的业务逻辑方法
type ChatAgentService interface {
	// SaveChatAgent 保存智能体信息
	// 如果ID为空则新增，否则更新现有记录
	SaveChatAgent(ctx context.Context, agent *models.ChatAgent) error

	// DeleteChatAgent 删除智能体
	// 根据ID删除指定的智能体记录
	DeleteChatAgent(ctx context.Context, id uuid.UUID) error

	// GetChatAgentsByApplicationID 根据应用ID获取智能体列表
	// 返回指定应用下的所有智能体，支持分页
	GetChatAgentsByApplicationID(ctx context.Context, applicationID uuid.UUID, page, pageSize int) ([]*models.ChatAgent, int64, error)

	// GetChatAgentByApiKey 根据API Key获取聊天智能体
	// 返回指定API Key的聊天智能体
	GetChatAgentByApiKey(ctx context.Context, apiKey string) (*models.ChatAgent, error) // 根据API Key获取应用
}

// chatAgentService 智能体 业务逻辑层实现
// 实现 ChatAgentService 接口
type chatAgentService struct {
	chatAgentRepo       repository.ChatAgentRepository // 数据访问层接口
	chatAgentApiKeyRepo repository.ChatAgentApiKeyRepository
}

// NewChatAgentService 创建 智能体 服务实例
// 返回 ChatAgentService 接口的实现
// 参数：chatAgentRepo - 智能体 数据访问层接口
func NewChatAgentService(chatAgentRepo repository.ChatAgentRepository, chatAgentApiKeyRepo repository.ChatAgentApiKeyRepository) ChatAgentService {
	return &chatAgentService{
		chatAgentRepo:       chatAgentRepo,
		chatAgentApiKeyRepo: chatAgentApiKeyRepo,
	}
}

// SaveChatAgent 保存智能体信息
// 如果ID为空则新增，否则更新现有记录
func (s *chatAgentService) SaveChatAgent(ctx context.Context, agent *models.ChatAgent) error {
	// 数据验证
	if err := s.validateChatAgent(agent); err != nil {
		return err
	}

	if agent.ID == uuid.Nil {
		// 新增：生成新的UUID
		agent.ID = uuid.New()
		return s.chatAgentRepo.Create(ctx, agent)
	} else {
		// 更新：检查记录是否存在
		existing, err := s.chatAgentRepo.GetByID(ctx, agent.ID)
		if err != nil {
			return fmt.Errorf("智能体不存在: %w", err)
		}
		if existing == nil {
			return fmt.Errorf("智能体不存在")
		}
		return s.chatAgentRepo.Update(ctx, agent)
	}
}

// DeleteChatAgent 删除智能体
// 根据ID删除指定的智能体记录
func (s *chatAgentService) DeleteChatAgent(ctx context.Context, id uuid.UUID) error {
	// 检查记录是否存在
	existing, err := s.chatAgentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("智能体不存在: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("智能体不存在")
	}

	return s.chatAgentRepo.DeleteByID(ctx, id)
}

// GetChatAgentsByApplicationID 根据应用ID获取智能体列表
// 返回指定应用下的所有智能体，支持分页
func (s *chatAgentService) GetChatAgentsByApplicationID(ctx context.Context, applicationID uuid.UUID, page, pageSize int) ([]*models.ChatAgent, int64, error) {
	return s.chatAgentRepo.GetByApplicationIDWithPagination(ctx, applicationID, page, pageSize)
}

// GetChatAgentByApiKey 根据API Key获取聊天智能体
// 返回指定API Key的聊天智能体
func (s *chatAgentService) GetChatAgentByApiKey(ctx context.Context, apiKey string) (*models.ChatAgent, error) {
	apiKeyObj, getApiKeyErr := s.chatAgentApiKeyRepo.GetByApiKey(ctx, apiKey)
	if getApiKeyErr != nil {
		return nil, getApiKeyErr
	} else {
		return s.chatAgentRepo.GetByID(ctx, apiKeyObj.ChatAgentID)
	}
}

// validateChatAgent 验证智能体数据
// 检查必填字段是否为空
func (s *chatAgentService) validateChatAgent(agent *models.ChatAgent) error {
	if agent == nil {
		return fmt.Errorf("智能体不能为空")
	}

	if agent.Name == "" {
		return fmt.Errorf("智能体名称不能为空")
	}

	if agent.Description == "" {
		return fmt.Errorf("智能体描述不能为空")
	}

	if agent.ApplicationID == uuid.Nil {
		return fmt.Errorf("所属应用ID不能为空")
	}

	if agent.ChatSystemPrompt == "" {
		return fmt.Errorf("系统提示词不能为空")
	}

	if agent.ChatModelID == uuid.Nil {
		return fmt.Errorf("聊天模型ID不能为空")
	}

	if agent.ConversationNamingModelID == uuid.Nil {
		return fmt.Errorf("会话命名模型ID不能为空")
	}

	// 验证数值范围
	if agent.ModelParamTemperature < 0 || agent.ModelParamTemperature > 2 {
		return fmt.Errorf("模型温度必须在0-2之间")
	}

	if agent.ModelParamTopP < 0 || agent.ModelParamTopP > 1 {
		return fmt.Errorf("模型TopP必须在0-1之间")
	}

	if agent.EnableContextLengthLimit && agent.ContextLengthLimit <= 0 {
		return fmt.Errorf("启用上下文长度限制时，限制值必须大于0")
	}

	if agent.EnableMaxOutputTokenCountLimit && agent.MaxOutputTokenCountLimit <= 0 {
		return fmt.Errorf("启用最大输出Token限制时，限制值必须大于0")
	}

	return nil
}
