// Package models 提供应用程序的数据模型定义
package models

import (
	"lemon-tree-core/internal/base"

	"github.com/google/uuid"
)

// ChatAgentMcpServerTool ChatAgentMcp配置
type ChatAgentMcpServerTool struct {
	base.BaseModel                       // 继承基础模型，包含 ID、时间戳等通用字段
	ChatAgentID                uuid.UUID `json:"chat_agent_id" gorm:"type:char(36);not null;comment:所属的聊天智能体ID"`
	ApplicationMcpServerToolID uuid.UUID `json:"application_mcp_server_tool_id" gorm:"type:char(36);not null;comment:所属的mcp服务工具ID"`
	Enabled                    bool      `json:"enabled" gorm:"type:tinyint(1);not null;comment:是否启用"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ChatAgentMcpServerTool) TableName() string {
	return "ltc_chat_agent_mcp_server_tool"
}
