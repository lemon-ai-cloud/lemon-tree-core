// Package models 提供应用程序的数据模型定义
package models

import (
	"lemon-tree-core/internal/base"

	"github.com/google/uuid"
)

// ApplicationMcpServerConfig 应用mcp配置
type ApplicationMcpServerConfig struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	ConfigID       string    `json:"config_id" gorm:"type:varchar(64);not null;comment:配置ID"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	Name           string    `json:"name" gorm:"type:varchar(64);not null;comment:名称"`
	Description    string    `json:"description" gorm:"type:varchar(512);not null;comment:描述"`
	Version        string    `json:"version" gorm:"type:varchar(64);not null;comment:版本"`
	// MCP连接方式 sse stdio streamable-http
	McpServerConnectType string `json:"mcp_server_protocol" gorm:"type:varchar(64);not null;comment:MCP服务连接方式"`
	McpServerTimeout     int    `json:"mcp_server_timeout" gorm:"type:int;not null;comment:MCP服务超时时间"`
	// sse / streamable-http使用
	McpServerUrl    string `json:"mcp_server_url" gorm:"type:varchar(512);not null;comment:MCP服务URL"`
	McpServerHeader string `json:"mcp_server_header" gorm:"type:text;not null;comment:MCP服务请求头"`
	// stdio 使用
	McpServerCommand string `json:"mcp_server_command" gorm:"type:varchar(512);not null;comment:MCP服务命令"`
	McpServerArgs    string `json:"mcp_server_args" gorm:"type:varchar(512);not null;comment:MCP服务参数"`
	McpServerEnv     string `json:"mcp_server_env" gorm:"type:varchar(512);not null;comment:MCP服务环境变量"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ApplicationMcpServerConfig) TableName() string {
	return "ltc_application_mcp_server_config"
}
