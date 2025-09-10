// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"fmt"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"
	"lemon-tree-core/internal/utils"
	"log"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// ApplicationMcpServerConfigService ApplicationMCP配置 业务逻辑层接口
// 定义 ApplicationMCP配置 相关的业务逻辑方法
type ApplicationMcpServerConfigService interface {
	// SaveApplicationMcpServerConfig 保存应用MCP配置信息
	// 如果ID为空则新增，否则更新现有记录
	SaveApplicationMcpServerConfig(ctx context.Context, config *models.ApplicationMcpServerConfig) error

	// DeleteApplicationMcpServerConfig 删除MCP配置
	// 根据ID删除指定的MCP配置记录
	DeleteApplicationMcpServerConfig(ctx context.Context, id uuid.UUID) error

	// GetMcpServerConfigsByApplicationID 根据应用ID获取MCP配置列表
	// 返回指定应用下的所有MCP配置
	GetMcpServerConfigsByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationMcpServerConfig, error)

	// GetMcpServerTools 获取MCP服务器的所有工具
	// 根据MCP配置ID从数据库获取工具列表，如果为空则自动同步
	GetMcpServerTools(ctx context.Context, configID uuid.UUID) ([]*models.ApplicationMcpServerTool, error)

	// SyncMcpServerTools 同步MCP服务器的工具列表
	// 从MCP服务器获取工具列表并同步到数据库
	SyncMcpServerTools(ctx context.Context, configID uuid.UUID) ([]*models.ApplicationMcpServerTool, error)
}

// applicationMcpServerConfigService ApplicationMCP配置 业务逻辑层实现
// 实现 ApplicationMcpServerConfigService 接口
type applicationMcpServerConfigService struct {
	applicationMcpServerConfigRepo repository.ApplicationMcpServerConfigRepository // 数据访问层接口
	applicationMcpServerToolRepo   repository.ApplicationMcpServerToolRepository   // 工具数据访问层接口
}

// NewApplicationMcpServerConfigService 创建 ApplicationMCP配置 服务实例
// 返回 ApplicationMcpServerConfigService 接口的实现
// 参数：applicationMcpServerConfigRepo - ApplicationMCP配置 数据访问层接口
// 参数：applicationMcpServerToolRepo - ApplicationMCP工具 数据访问层接口
func NewApplicationMcpServerConfigService(applicationMcpServerConfigRepo repository.ApplicationMcpServerConfigRepository, applicationMcpServerToolRepo repository.ApplicationMcpServerToolRepository) ApplicationMcpServerConfigService {
	return &applicationMcpServerConfigService{
		applicationMcpServerConfigRepo: applicationMcpServerConfigRepo,
		applicationMcpServerToolRepo:   applicationMcpServerToolRepo,
	}
}

// SaveApplicationMcpServerConfig 保存应用MCP配置信息
// 如果ID为空则新增，否则更新现有记录
func (s *applicationMcpServerConfigService) SaveApplicationMcpServerConfig(ctx context.Context, config *models.ApplicationMcpServerConfig) error {
	// 数据验证
	if err := s.validateApplicationMcpServerConfig(config); err != nil {
		return err
	}

	if config.ID == uuid.Nil {
		// 新增：生成新的UUID
		config.ID = uuid.New()
		shortID, _ := utils.ShortUUID(config.ID.String())
		config.ConfigID = shortID
		return s.applicationMcpServerConfigRepo.Create(ctx, config)
	} else {
		// 更新：检查记录是否存在
		existing, err := s.applicationMcpServerConfigRepo.GetByID(ctx, config.ID)
		if err != nil {
			return fmt.Errorf("MCP配置不存在: %w", err)
		}
		if existing == nil {
			return fmt.Errorf("MCP配置不存在")
		}
		if config.ConfigID == "" {
			shortID, _ := utils.ShortUUID(config.ID.String())
			config.ConfigID = shortID
		}
		return s.applicationMcpServerConfigRepo.Update(ctx, config)
	}
}

// DeleteApplicationMcpServerConfig 删除MCP配置
// 根据ID删除指定的MCP配置记录
func (s *applicationMcpServerConfigService) DeleteApplicationMcpServerConfig(ctx context.Context, id uuid.UUID) error {
	// 检查记录是否存在
	existing, err := s.applicationMcpServerConfigRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("MCP配置不存在: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("MCP配置不存在")
	}

	return s.applicationMcpServerConfigRepo.Delete(ctx, id)
}

// GetMcpServerConfigsByApplicationID 根据应用ID获取MCP配置列表
// 返回指定应用下的所有MCP配置
func (s *applicationMcpServerConfigService) GetMcpServerConfigsByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationMcpServerConfig, error) {
	return s.applicationMcpServerConfigRepo.GetByApplicationID(ctx, applicationID)
}

// validateApplicationMcpServerConfig 验证应用MCP配置数据
// 检查必填字段是否为空
func (s *applicationMcpServerConfigService) validateApplicationMcpServerConfig(config *models.ApplicationMcpServerConfig) error {
	if config == nil {
		return fmt.Errorf("MCP配置不能为空")
	}

	if config.Name == "" {
		return fmt.Errorf("MCP配置名称不能为空")
	}

	if config.Description == "" {
		return fmt.Errorf("MCP配置描述不能为空")
	}

	if config.Version == "" {
		return fmt.Errorf("MCP配置版本不能为空")
	}

	if config.McpServerConnectType == "" {
		return fmt.Errorf("MCP服务连接方式不能为空")
	}

	if config.ApplicationID == uuid.Nil {
		return fmt.Errorf("所属应用ID不能为空")
	}

	// 根据连接方式验证必填字段
	switch config.McpServerConnectType {
	case "sse", "streamable-http":
		if config.McpServerUrl == "" {
			return fmt.Errorf("MCP服务URL不能为空")
		}
	case "stdio":
		if config.McpServerCommand == "" {
			return fmt.Errorf("MCP服务命令不能为空")
		}
	default:
		return fmt.Errorf("不支持的MCP服务连接方式: %s", config.McpServerConnectType)
	}

	return nil
}

// GetMcpServerTools 获取MCP服务器的所有工具
// 根据MCP配置ID从数据库获取工具列表，如果为空则自动同步
func (s *applicationMcpServerConfigService) GetMcpServerTools(ctx context.Context, configID uuid.UUID) ([]*models.ApplicationMcpServerTool, error) {
	// 从数据库获取工具列表
	tools, err := s.applicationMcpServerToolRepo.GetByApplicationMcpServerConfigID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("从数据库获取工具失败: %w", err)
	}

	// 如果工具列表为空，自动执行一次同步
	if len(tools) == 0 {
		log.Printf("工具列表为空，自动执行同步: configID=%s", configID)
		return s.SyncMcpServerTools(ctx, configID)
	}

	return tools, nil
}

// SyncMcpServerTools 同步MCP服务器的工具列表
// 从MCP服务器获取工具列表并同步到数据库
func (s *applicationMcpServerConfigService) SyncMcpServerTools(ctx context.Context, configID uuid.UUID) ([]*models.ApplicationMcpServerTool, error) {
	// 获取MCP配置信息
	config, err := s.applicationMcpServerConfigRepo.GetByID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("获取MCP配置失败: %w", err)
	}
	if config == nil {
		return nil, fmt.Errorf("MCP配置不存在")
	}

	// 根据连接方式创建MCP客户端并获取工具
	var tools []mcp.Tool
	switch config.McpServerConnectType {
	case "sse", "streamable-http", "stdio":
		tools, err = s.getToolsFromMcpClient(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("从HTTP客户端获取工具失败: %w", err)
		}

	default:
		return nil, fmt.Errorf("不支持的连接方式: %s", config.McpServerConnectType)
	}

	// 同步工具到数据库
	if err := s.syncMcpToolsToDatabase(ctx, config, tools); err != nil {
		log.Printf("同步工具到数据库失败: %v", err)
		// 不返回错误，继续执行
	}

	// 从数据库获取同步后的工具列表
	syncedTools, err := s.applicationMcpServerToolRepo.GetByApplicationMcpServerConfigID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("从数据库获取工具失败: %w", err)
	}

	return syncedTools, nil
}

// getToolsFromMcpClient 从HTTP/SSE客户端获取工具
func (s *applicationMcpServerConfigService) getToolsFromMcpClient(ctx context.Context, config *models.ApplicationMcpServerConfig) ([]mcp.Tool, error) {
	var c *client.Client
	switch config.McpServerConnectType {
	case "streamable-http":
		// 创建HTTP传输
		httpTransport, err := transport.NewStreamableHTTP(config.McpServerUrl)
		if err != nil {
			return nil, fmt.Errorf("创建Streamable HTTP传输失败: %w", err)
		}
		// 创建客户端
		c = client.NewClient(httpTransport)
	case "sse":
		// 创建HTTP传输
		sse, err := transport.NewSSE(config.McpServerUrl)
		if err != nil {
			return nil, fmt.Errorf("创建SSE传输失败: %w", err)
		}
		c = client.NewClient(sse)
	case "studio":
		studio := transport.NewStdio(config.McpServerCommand, []string{config.McpServerEnv}, "")
		c = client.NewClient(studio)
	}

	// 初始化客户端
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Lemon-Tree MCP Client",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, initRequest)
	if err != nil {
		return nil, fmt.Errorf("初始化客户端失败: %w", err)
	}

	log.Printf("连接到服务器: %s (版本 %s)", serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)

	// 健康检查
	if err := c.Ping(ctx); err != nil {
		return nil, fmt.Errorf("健康检查失败: %w", err)
	}

	// 获取工具列表
	if serverInfo.Capabilities.Tools == nil {
		log.Println("服务器不支持工具功能")
		return []mcp.Tool{}, nil
	}

	toolsRequest := mcp.ListToolsRequest{}
	toolsResult, err := c.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, fmt.Errorf("获取工具列表失败: %w", err)
	}

	log.Printf("服务器有 %d 个可用工具", len(toolsResult.Tools))
	return toolsResult.Tools, nil
}

// syncMcpToolsToDatabase 同步MCP工具到数据库
// 实现增量更新：保留现有记录的ID和创建时间，只更新变化的字段
func (s *applicationMcpServerConfigService) syncMcpToolsToDatabase(ctx context.Context, config *models.ApplicationMcpServerConfig, tools []mcp.Tool) error {
	// 获取现有的工具列表
	existingTools, err := s.applicationMcpServerToolRepo.GetByApplicationMcpServerConfigID(ctx, config.ID)
	if err != nil {
		return fmt.Errorf("获取现有工具失败: %w", err)
	}

	// 创建现有工具的映射，以name为key
	existingToolsMap := make(map[string]*models.ApplicationMcpServerTool)
	for _, tool := range existingTools {
		existingToolsMap[tool.Name] = tool
	}

	// 创建新工具列表的映射，以name为key
	newToolsMap := make(map[string]mcp.Tool)
	for _, tool := range tools {
		newToolsMap[tool.Name] = tool
	}

	// 处理每个新工具
	for _, newTool := range tools {
		title := ""
		if newTool.Annotations.Title != "" {
			title = newTool.Annotations.Title
		}
		// 如果title为空，使用name作为title
		if title == "" {
			title = newTool.Name
		}

		if existingTool, exists := existingToolsMap[newTool.Name]; exists {
			// 工具已存在，检查是否需要更新
			needsUpdate := false
			if existingTool.Description != newTool.Description {
				existingTool.Description = newTool.Description
				needsUpdate = true
			}
			if existingTool.Title != title {
				existingTool.Title = title
				needsUpdate = true
			}

			if needsUpdate {
				if err := s.applicationMcpServerToolRepo.Update(ctx, existingTool); err != nil {
					log.Printf("更新工具失败: %s, error: %v", newTool.Name, err)
				}
			}
		} else {
			// 工具不存在，创建新记录
			newToolModel := &models.ApplicationMcpServerTool{
				ApplicationID:                config.ApplicationID,
				ApplicationMcpServerConfigID: config.ID,
				Name:                         newTool.Name,
				Title:                        title,
				Description:                  newTool.Description,
			}

			if err := s.applicationMcpServerToolRepo.Create(ctx, newToolModel); err != nil {
				log.Printf("创建工具失败: %s, error: %v", newTool.Name, err)
			}
		}
	}

	// 删除不再存在的工具
	for toolName, existingTool := range existingToolsMap {
		if _, exists := newToolsMap[toolName]; !exists {
			if err := s.applicationMcpServerToolRepo.Delete(ctx, existingTool.ID); err != nil {
				log.Printf("删除工具失败: %s, error: %v", toolName, err)
			}
		}
	}

	return nil
}
