// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"fmt"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"
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
	// 根据MCP配置ID连接服务器并返回可用工具列表
	GetMcpServerTools(ctx context.Context, configID uuid.UUID) ([]*models.ApplicationMcpServerTool, error)
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
// 根据MCP配置ID连接服务器并返回可用工具列表
func (s *applicationMcpServerConfigService) GetMcpServerTools(ctx context.Context, configID uuid.UUID) ([]*models.ApplicationMcpServerTool, error) {
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
	case "sse", "streamable-http":
		tools, err = s.getToolsFromHTTPClient(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("从HTTP客户端获取工具失败: %w", err)
		}

	case "stdio":
		tools, err = s.getToolsFromStdioClient(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("从STDIO客户端获取工具失败: %w", err)
		}

	default:
		return nil, fmt.Errorf("不支持的连接方式: %s", config.McpServerConnectType)
	}

	// 保存工具到数据库
	if err := s.saveMcpToolsToDatabase(ctx, config, tools); err != nil {
		log.Printf("保存工具到数据库失败: %v", err)
		// 不返回错误，继续执行
	}

	// 从数据库获取保存的工具列表
	savedTools, err := s.applicationMcpServerToolRepo.GetByApplicationMcpServerConfigID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("从数据库获取工具失败: %w", err)
	}

	return savedTools, nil
}

// getToolsFromHTTPClient 从HTTP/SSE客户端获取工具
func (s *applicationMcpServerConfigService) getToolsFromHTTPClient(ctx context.Context, config *models.ApplicationMcpServerConfig) ([]mcp.Tool, error) {
	// 创建HTTP传输
	httpTransport, err := transport.NewStreamableHTTP(config.McpServerUrl)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP传输失败: %w", err)
	}

	// 创建客户端
	c := client.NewClient(httpTransport)

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

// getToolsFromStdioClient 从STDIO客户端获取工具
func (s *applicationMcpServerConfigService) getToolsFromStdioClient(ctx context.Context, config *models.ApplicationMcpServerConfig) ([]mcp.Tool, error) {
	// TODO: 实现STDIO客户端连接
	// 这里需要根据config.McpServerCommand, config.McpServerArgs, config.McpServerEnv来创建STDIO传输
	// 暂时返回空列表，等STDIO传输实现后再完善
	log.Println("STDIO客户端连接功能待实现")
	return []mcp.Tool{}, nil
}

// saveMcpToolsToDatabase 保存MCP工具到数据库
// 先删除该配置下的所有工具，然后重新插入新的工具列表
func (s *applicationMcpServerConfigService) saveMcpToolsToDatabase(ctx context.Context, config *models.ApplicationMcpServerConfig, tools []mcp.Tool) error {
	// 先删除该配置下的所有工具
	if err := s.applicationMcpServerToolRepo.DeleteByApplicationMcpServerConfigID(ctx, config.ID); err != nil {
		return fmt.Errorf("删除现有工具失败: %w", err)
	}

	// 准备批量插入的工具数据
	var toolsToCreate []*models.ApplicationMcpServerTool
	for _, tool := range tools {
		title := ""
		if tool.Annotations.Title != "" {
			title = tool.Annotations.Title
		}

		// 如果title为空，使用name作为title
		if title == "" {
			title = tool.Name
		}

		toolModel := &models.ApplicationMcpServerTool{
			ApplicationID:                config.ApplicationID,
			ApplicationMcpServerConfigID: config.ID,
			Name:                         tool.Name,
			Title:                        title,
			Description:                  tool.Description,
		}

		toolsToCreate = append(toolsToCreate, toolModel)
	}

	// 批量创建工具记录
	if len(toolsToCreate) > 0 {
		if err := s.applicationMcpServerToolRepo.BatchCreate(ctx, toolsToCreate); err != nil {
			return fmt.Errorf("批量创建工具失败: %w", err)
		}
	}

	return nil
}
