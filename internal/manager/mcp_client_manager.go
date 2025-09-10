package manager

import (
	"context"
	"fmt"
	"lemon-tree-core/internal/models"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetMcpClient(ctx context.Context, config *models.ApplicationMcpServerConfig) (*client.Client, error) {
	// 根据连接方式创建MCP客户端
	var c *client.Client
	var err error

	switch config.McpServerConnectType {
	case "streamable-http":
		httpTransport, err := transport.NewStreamableHTTP(config.McpServerUrl)
		if err != nil {
			return nil, fmt.Errorf("创建Streamable HTTP传输失败: %w", err)
		}
		c = client.NewClient(httpTransport)
	case "sse":
		sse, err := transport.NewSSE(config.McpServerUrl)
		if err != nil {
			return nil, fmt.Errorf("创建SSE传输失败: %w", err)
		}
		c = client.NewClient(sse)
	case "stdio":
		studio := transport.NewStdio(config.McpServerCommand, []string{config.McpServerEnv}, "")
		c = client.NewClient(studio)
	default:
		return nil, fmt.Errorf("不支持的连接方式: %s", config.McpServerConnectType)
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
		return nil, fmt.Errorf("初始化MCP客户端失败: %w", err)
	}

	log.Printf("连接到MCP服务器: %s (版本 %s)", serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)

	// 健康检查
	if err := c.Ping(ctx); err != nil {
		return nil, fmt.Errorf("MCP服务器健康检查失败: %w", err)
	}

	return c, nil
}
