// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
)

// ChatAgent AI聊天智能体
// 智能体是一个可以有自主思考能力且可以自主调用各种工具的AI助手
type ChatAgent struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	Name           string    `json:"name" gorm:"type:varchar(64);not null;comment:Agent名称"`
	Description    string    `json:"description" gorm:"type:varchar(512);not null;comment:Agent描述"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	AvatarUrl      string    `json:"avatar_url" gorm:"type:varchar(512);not null;comment:Agent的头像URL"`
	// LLM设置
	ChatSystemPrompt               string    `json:"system_prompt" gorm:"type:text;not null;comment:系统提示"`
	ChatModelID                    uuid.UUID `json:"chat_model_id" gorm:"type:char(36);not null;comment:聊天模型ID"`
	ConversationNamingPrompt       string    `json:"conversation_naming_prompt" gorm:"type:text;not null;comment:会话命名提示词"`
	ConversationNamingModelID      uuid.UUID `json:"conversation_naming_model_id" gorm:"type:char(36);not null;comment:会话命名模型ID"`
	ModelParamTemperature          float64   `json:"model_temperature" gorm:"type:decimal(10,2);not null;comment:模型温度"`
	ModelParamTopP                 float64   `json:"model_top_p" gorm:"type:decimal(10,2);not null;comment:模型TopP"`
	EnableContextLengthLimit       bool      `json:"enable_context_length_limit" gorm:"type:tinyint(1);not null;comment:是否启用上下文长度限制，单位是消息数量"`
	ContextLengthLimit             int       `json:"context_length_limit" gorm:"type:int;not null;comment:上下文长度限制，单位是消息数量"`
	EnableMaxOutputTokenCountLimit bool      `json:"enable_max_output_token_count_limit" gorm:"type:tinyint(1);not null;comment:是否启用最大输出Token数量限制"`
	MaxOutputTokenCountLimit       int       `json:"max_output_token_count_limit" gorm:"type:int;not null;comment:最大输出Token数量"`
	// 这个流式返回只是针对默认的Lemon Tree UI界面，通过API访问时可以通过传参来控制是否流式返回
	DefaultStreamable bool `json:"default_streamable" gorm:"type:tinyint(1);not null;comment:是否默认流式返回"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ChatAgent) TableName() string {
	return "ltc_chat_agent"
}
