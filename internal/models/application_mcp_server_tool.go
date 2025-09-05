// Package models 提供应用程序的数据模型定义
package models

import (
	"lemon-tree-core/internal/base"

	"github.com/google/uuid"
)

// ApplicationMcpServerConfig 应用mcp配置
type ApplicationMcpServerTool struct {
	base.BaseModel                         // 继承基础模型，包含 ID、时间戳等通用字段
	ApplicationID                uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	ApplicationMcpServerConfigID uuid.UUID `json:"application_mcp_server_config_id" gorm:"type:char(36);not null;comment:所属mcp服务配置ID"`
	Name                         string    `json:"name" gorm:"type:varchar(64);not null;comment:工具名称"`
	Title                        string    `json:"title" gorm:"type:varchar(64);not null;comment:标题"`
	Description                  string    `json:"description" gorm:"type:text;not null;comment:描述"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ApplicationMcpServerTool) TableName() string {
	return "ltc_application_mcp_server_tool"
}
