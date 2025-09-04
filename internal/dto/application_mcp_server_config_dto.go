// Package dto 提供数据传输对象定义
// 用于在不同层之间传递数据，确保数据格式的一致性
package dto

// ApplicationMcpServerConfigDto ApplicationMCP配置 数据传输对象
// 用于在业务逻辑层和HTTP处理层之间传递数据
type ApplicationMcpServerConfigDto struct {
	ID                   string `json:"id"`                      // 主键ID
	ApplicationID        string `json:"application_id"`          // 所属应用ID
	Name                 string `json:"name"`                    // 名称
	Description          string `json:"description"`             // 描述
	Version              string `json:"version"`                 // 版本
	McpServerConnectType string `json:"mcp_server_connect_type"` // MCP服务连接方式
	McpServerTimeout     int    `json:"mcp_server_timeout"`      // MCP服务超时时间
	McpServerUrl         string `json:"mcp_server_url"`          // MCP服务URL
	McpServerHeader      string `json:"mcp_server_header"`       // MCP服务请求头
	McpServerCommand     string `json:"mcp_server_command"`      // MCP服务命令
	McpServerArgs        string `json:"mcp_server_args"`         // MCP服务参数
	McpServerEnv         string `json:"mcp_server_env"`          // MCP服务环境变量
	CreatedAt            string `json:"created_at"`              // 创建时间
	UpdatedAt            string `json:"updated_at"`              // 更新时间
}

// SaveApplicationMcpServerConfigRequest 保存应用MCP配置请求
// 用于接收前端保存MCP配置的请求数据
type SaveApplicationMcpServerConfigRequest struct {
	ID                   *string `json:"id,omitempty"`            // 主键ID，为空时新增，有值时更新
	ApplicationID        string  `json:"application_id"`          // 所属应用ID
	Name                 string  `json:"name"`                    // 名称
	Description          string  `json:"description"`             // 描述
	Version              string  `json:"version"`                 // 版本
	McpServerConnectType string  `json:"mcp_server_connect_type"` // MCP服务连接方式
	McpServerTimeout     int     `json:"mcp_server_timeout"`      // MCP服务超时时间
	McpServerUrl         string  `json:"mcp_server_url"`          // MCP服务URL
	McpServerHeader      string  `json:"mcp_server_header"`       // MCP服务请求头
	McpServerCommand     string  `json:"mcp_server_command"`      // MCP服务命令
	McpServerArgs        string  `json:"mcp_server_args"`         // MCP服务参数
	McpServerEnv         string  `json:"mcp_server_env"`          // MCP服务环境变量
}

// ApplicationMcpServerConfigListResponse ApplicationMCP配置列表响应
// 用于返回MCP配置列表的响应数据
type ApplicationMcpServerConfigListResponse struct {
	ApplicationMcpServerConfigs []ApplicationMcpServerConfigDto `json:"application_mcp_server_configs"` // MCP配置列表
}

// SingleApplicationMcpServerConfigResponse 单个ApplicationMCP配置响应
// 用于返回单个MCP配置的响应数据
type SingleApplicationMcpServerConfigResponse struct {
	ApplicationMcpServerConfig ApplicationMcpServerConfigDto `json:"application_mcp_server_config"` // MCP配置
}
