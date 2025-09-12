// Package models 提供应用程序的数据模型定义
package models

import (
	"lemon-tree-core/internal/base"

	"github.com/google/uuid"
)

// ChatAgentConversation 聊天智能体的会话
type ChatAgentConversation struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	Title          string    `json:"title" gorm:"type:varchar(64);not null;comment:会话标题"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	ChatAgentID    uuid.UUID `json:"chat_agent_id" gorm:"type:char(36);not null;comment:所属Chat Agent ID"`
	ServiceUserID  string    `json:"service_user_id" gorm:"type:varchar(256);not null;comment:业务侧的用户ID"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ChatAgentConversation) TableName() string {
	return "ltc_chat_agent_conversation"
}
