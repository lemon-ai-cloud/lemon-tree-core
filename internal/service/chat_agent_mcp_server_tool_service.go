// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"fmt"
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"
	"log"

	"github.com/google/uuid"
)

// ChatAgentMcpServerToolService ChatAgentMcpServerTool 业务逻辑层接口
// 定义 ChatAgentMcpServerTool 相关的业务逻辑方法
type ChatAgentMcpServerToolService interface {
	// SaveChatAgentMcpServerToolSettings 保存聊天智能体的MCP工具设置
	SaveChatAgentMcpServerToolSettings(ctx context.Context, chatAgentID uuid.UUID, toolSettings []dto.ChatAgentMcpServerToolSettingDto) error

	// GetChatAgentMcpServerToolSettings 获取聊天智能体的MCP工具设置
	GetChatAgentMcpServerToolSettings(ctx context.Context, chatAgentID uuid.UUID) ([]dto.ChatAgentMcpServerToolSettingDto, error)

	// GetChatAgentAvailableMcpServerTools 获取聊天智能体可用的MCP工具列表
	// 返回指定ChatAgentID下所有MCP工具及其启用状态，按MCP Server分组
	GetChatAgentAvailableMcpServerTools(ctx context.Context, chatAgentID uuid.UUID) ([]dto.McpServerToolGroupDto, error)
}

// chatAgentMcpServerToolService ChatAgentMcpServerTool 业务逻辑层实现
// 实现 ChatAgentMcpServerToolService 接口
type chatAgentMcpServerToolService struct {
	chatAgentMcpServerToolRepo     repository.ChatAgentMcpServerToolRepository
	applicationMcpServerToolRepo   repository.ApplicationMcpServerToolRepository
	applicationMcpServerConfigRepo repository.ApplicationMcpServerConfigRepository
	chatAgentRepo                  repository.ChatAgentRepository
}

// NewChatAgentMcpServerToolService 创建 ChatAgentMcpServerTool 服务实例
// 返回 ChatAgentMcpServerToolService 接口的实现
func NewChatAgentMcpServerToolService(
	chatAgentMcpServerToolRepo repository.ChatAgentMcpServerToolRepository,
	applicationMcpServerToolRepo repository.ApplicationMcpServerToolRepository,
	applicationMcpServerConfigRepo repository.ApplicationMcpServerConfigRepository,
	chatAgentRepo repository.ChatAgentRepository,
) ChatAgentMcpServerToolService {
	return &chatAgentMcpServerToolService{
		chatAgentMcpServerToolRepo:     chatAgentMcpServerToolRepo,
		applicationMcpServerToolRepo:   applicationMcpServerToolRepo,
		applicationMcpServerConfigRepo: applicationMcpServerConfigRepo,
		chatAgentRepo:                  chatAgentRepo,
	}
}

// SaveChatAgentMcpServerToolSettings 保存聊天智能体的MCP工具设置
func (s *chatAgentMcpServerToolService) SaveChatAgentMcpServerToolSettings(ctx context.Context, chatAgentID uuid.UUID, toolSettings []dto.ChatAgentMcpServerToolSettingDto) error {
	// 验证聊天智能体是否存在
	_, err := s.chatAgentRepo.GetByID(ctx, chatAgentID)
	if err != nil {
		return fmt.Errorf("聊天智能体不存在: %w", err)
	}

	// 获取当前智能体的所有工具配置
	existingSettings, err := s.chatAgentMcpServerToolRepo.GetByChatAgentID(ctx, chatAgentID)
	if err != nil {
		return fmt.Errorf("获取现有工具配置失败: %w", err)
	}

	// 创建现有配置的映射，便于查找
	existingMap := make(map[uuid.UUID]*models.ChatAgentMcpServerTool)
	for _, setting := range existingSettings {
		existingMap[setting.ApplicationMcpServerToolID] = setting
	}

	// 处理新的工具设置
	var toCreate []*models.ChatAgentMcpServerTool
	var toUpdate []*models.ChatAgentMcpServerTool
	var toDelete []uuid.UUID

	// 记录新设置中的工具ID
	newToolIDs := make(map[uuid.UUID]bool)

	for _, toolSetting := range toolSettings {
		toolID, err := uuid.Parse(toolSetting.ApplicationMcpServerToolID)
		if err != nil {
			return fmt.Errorf("无效的工具ID: %s", toolSetting.ApplicationMcpServerToolID)
		}

		newToolIDs[toolID] = true

		if existingSetting, exists := existingMap[toolID]; exists {
			// 更新现有配置
			existingSetting.Enabled = toolSetting.Enabled
			toUpdate = append(toUpdate, existingSetting)
		} else {
			// 创建新配置
			newSetting := &models.ChatAgentMcpServerTool{
				ChatAgentID:                chatAgentID,
				ApplicationMcpServerToolID: toolID,
				Enabled:                    toolSetting.Enabled,
			}
			toCreate = append(toCreate, newSetting)
		}
	}

	// 找出需要删除的配置（在新设置中不存在的）
	for toolID, existingSetting := range existingMap {
		if !newToolIDs[toolID] {
			toDelete = append(toDelete, existingSetting.ID)
		}
	}

	// 执行数据库操作
	// 1. 删除不需要的配置
	for _, id := range toDelete {
		if err := s.chatAgentMcpServerToolRepo.DeleteByID(ctx, id); err != nil {
			return fmt.Errorf("删除工具配置失败: %w", err)
		}
	}

	// 2. 批量创建新配置
	if len(toCreate) > 0 {
		if err := s.chatAgentMcpServerToolRepo.BatchCreate(ctx, toCreate); err != nil {
			return fmt.Errorf("创建工具配置失败: %w", err)
		}
	}

	// 3. 批量更新现有配置
	if len(toUpdate) > 0 {
		if err := s.chatAgentMcpServerToolRepo.BatchUpdate(ctx, toUpdate); err != nil {
			return fmt.Errorf("更新工具配置失败: %w", err)
		}
	}

	return nil
}

// GetChatAgentMcpServerToolSettings 获取聊天智能体的MCP工具设置
func (s *chatAgentMcpServerToolService) GetChatAgentMcpServerToolSettings(ctx context.Context, chatAgentID uuid.UUID) ([]dto.ChatAgentMcpServerToolSettingDto, error) {
	// 验证聊天智能体是否存在
	_, err := s.chatAgentRepo.GetByID(ctx, chatAgentID)
	if err != nil {
		return nil, fmt.Errorf("聊天智能体不存在: %w", err)
	}

	// 获取智能体的工具配置
	settings, err := s.chatAgentMcpServerToolRepo.GetByChatAgentID(ctx, chatAgentID)
	if err != nil {
		return nil, fmt.Errorf("获取工具配置失败: %w", err)
	}

	// 转换为DTO
	var result []dto.ChatAgentMcpServerToolSettingDto
	for _, setting := range settings {
		result = append(result, dto.ChatAgentMcpServerToolSettingDto{
			ID:                         setting.ID.String(),
			ApplicationMcpServerToolID: setting.ApplicationMcpServerToolID.String(),
			Enabled:                    setting.Enabled,
		})
	}

	return result, nil
}

// GetChatAgentAvailableMcpServerTools 获取聊天智能体可用的MCP工具列表
func (s *chatAgentMcpServerToolService) GetChatAgentAvailableMcpServerTools(ctx context.Context, chatAgentID uuid.UUID) ([]dto.McpServerToolGroupDto, error) {
	// 验证聊天智能体是否存在
	chatAgent, err := s.chatAgentRepo.GetByID(ctx, chatAgentID)
	if err != nil {
		return nil, fmt.Errorf("聊天智能体不存在: %w", err)
	}

	// 获取应用下的所有MCP工具
	allTools, err := s.applicationMcpServerToolRepo.GetByApplicationID(ctx, chatAgent.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("获取应用MCP工具失败: %w", err)
	}

	// 获取智能体当前的工具配置
	agentSettings, err := s.chatAgentMcpServerToolRepo.GetByChatAgentID(ctx, chatAgentID)
	if err != nil {
		return nil, fmt.Errorf("获取智能体工具配置失败: %w", err)
	}

	// 创建配置映射
	settingMap := make(map[uuid.UUID]bool)
	for _, setting := range agentSettings {
		settingMap[setting.ApplicationMcpServerToolID] = setting.Enabled
	}

	// 按MCP Server分组
	serverGroups := make(map[uuid.UUID]*dto.McpServerToolGroupDto)

	for _, tool := range allTools {
		enabled := settingMap[tool.ID] // 如果不存在，默认为false

		toolDto := dto.ChatAgentAvailableMcpServerToolDto{
			ID:                           tool.ID.String(),
			ApplicationMcpServerConfigID: tool.ApplicationMcpServerConfigID.String(),
			Name:                         tool.Name,
			Title:                        tool.Title,
			Description:                  tool.Description,
			Enabled:                      enabled,
		}

		// 如果该MCP Server还没有分组，创建新分组
		if _, exists := serverGroups[tool.ApplicationMcpServerConfigID]; !exists {
			// 获取MCP Server配置信息
			serverConfig, err := s.applicationMcpServerConfigRepo.GetByID(ctx, tool.ApplicationMcpServerConfigID)
			if err != nil {
				log.Printf("获取MCP Server配置失败: %v", err)
				continue
			}

			serverGroups[tool.ApplicationMcpServerConfigID] = &dto.McpServerToolGroupDto{
				Server: dto.McpServerConfigDto{
					ID:          serverConfig.ID.String(),
					Name:        serverConfig.Name,
					Description: serverConfig.Description,
				},
				Tools: []dto.ChatAgentAvailableMcpServerToolDto{},
			}
		}

		// 将工具添加到对应的分组
		serverGroups[tool.ApplicationMcpServerConfigID].Tools = append(serverGroups[tool.ApplicationMcpServerConfigID].Tools, toolDto)
	}

	// 转换为切片
	var result []dto.McpServerToolGroupDto
	for _, group := range serverGroups {
		result = append(result, *group)
	}

	return result, nil
}

func (s *chatAgentMcpServerToolService) CallTool(ctx context.Context, mcpServerConfigID string, toolName, toolArgs string) (string, error) {
	// TODO: 实现MCP工具调用逻辑
	return "", fmt.Errorf("CallTool方法尚未实现")
}
