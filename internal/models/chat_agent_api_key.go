// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
)

// ChatAgentApiKey 聊天智能体API Key模型结构体
type ChatAgentApiKey struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	Name           string    `json:"name" gorm:"type:varchar(64);not null;comment:Key名称"`
	Description    string    `json:"description" gorm:"type:varchar(512);not null;comment:Key描述"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:应用ID"`
	ChatAgentID    uuid.UUID `json:"chat_agent_id" gorm:"type:char(36);not null;comment:智能体ID"`
	ApiKey         string    `json:"api_key" gorm:"type:varchar(512);not null;comment:API Key"`
	ApiSecret      string    `json:"api_secret" gorm:"type:varchar(512);not null;comment:API Secret"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ChatAgentApiKey) TableName() string {
	return "ltc_chat_agent_api_key"
}
