// Package dto 提供数据传输对象定义
// 用于在不同层之间传递数据，确保数据格式的一致性
package dto

// ChatAgentMcpServerToolSettingDto 聊天智能体MCP工具设置
type ChatAgentMcpServerToolSettingDto struct {
	ID                         string `json:"id"`                             // 配置ID
	ApplicationMcpServerToolID string `json:"application_mcp_server_tool_id"` // 应用MCP工具ID
	Enabled                    bool   `json:"enabled"`                        // 是否启用
}

// ChatAgentAvailableMcpServerToolDto 聊天智能体可用的MCP工具
type ChatAgentAvailableMcpServerToolDto struct {
	ID                           string `json:"id"`                               // 工具ID
	ApplicationMcpServerConfigID string `json:"application_mcp_server_config_id"` // MCP服务配置ID
	Name                         string `json:"name"`                             // 工具名称
	Title                        string `json:"title"`                            // 工具标题
	Description                  string `json:"description"`                      // 工具描述
	Enabled                      bool   `json:"enabled"`                          // 是否启用
}

// McpServerConfigDto MCP服务配置信息
type McpServerConfigDto struct {
	ID          string `json:"id"`          // 配置ID
	Name        string `json:"name"`        // 服务名称
	Description string `json:"description"` // 服务描述
}

// McpServerToolGroupDto 按MCP服务分组的工具
type McpServerToolGroupDto struct {
	Server McpServerConfigDto                   `json:"server"` // MCP服务信息
	Tools  []ChatAgentAvailableMcpServerToolDto `json:"tools"`  // 该服务下的工具列表
}

// SaveChatAgentMcpServerToolSettingsRequest 保存聊天智能体MCP工具设置请求
type SaveChatAgentMcpServerToolSettingsRequest struct {
	ToolSettings []ChatAgentMcpServerToolSettingDto `json:"tool_settings"` // 工具设置列表
}

// SaveChatAgentMcpServerToolSettingsResponse 保存聊天智能体MCP工具设置响应
type SaveChatAgentMcpServerToolSettingsResponse struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
}

// GetChatAgentMcpServerToolSettingsRequest 获取聊天智能体MCP工具设置请求
type GetChatAgentMcpServerToolSettingsRequest struct {
	ChatAgentID string `json:"chat_agent_id"` // 聊天智能体ID
}

// GetChatAgentMcpServerToolSettingsResponse 获取聊天智能体MCP工具设置响应
type GetChatAgentMcpServerToolSettingsResponse struct {
	Success bool                               `json:"success"` // 是否成功
	Data    []ChatAgentMcpServerToolSettingDto `json:"data"`    // 工具设置列表
	Message string                             `json:"message"` // 消息
}

// GetChatAgentAvailableMcpServerToolsRequest 获取聊天智能体可用MCP工具请求
type GetChatAgentAvailableMcpServerToolsRequest struct {
	ChatAgentID string `json:"chat_agent_id"` // 聊天智能体ID
}

// GetChatAgentAvailableMcpServerToolsResponse 获取聊天智能体可用MCP工具响应
type GetChatAgentAvailableMcpServerToolsResponse struct {
	Success bool                    `json:"success"` // 是否成功
	Data    []McpServerToolGroupDto `json:"data"`    // 按MCP服务分组的工具列表
	Message string                  `json:"message"` // 消息
}
