// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
)

// ChatAgentMessage 聊天智能体的聊天具体消息
type ChatAgentMessage struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	Title          string    `json:"title" gorm:"type:varchar(64);not null;comment:会话标题"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	ChatAgentID    uuid.UUID `json:"chat_agent_id" gorm:"type:char(36);not null;comment:所属Chat Agent ID"`
	ConversationID uuid.UUID `json:"conversation_id" gorm:"type:char(36);not null;comment:所属会话ID"`
	RequestID      string    `json:"request_id" gorm:"type:varchar(64);not null;comment:请求ID，同一条消息相关的子消息请求ID一致"`
	// 消息类型： message 普通消息 function_call 函数调用 function_call_output 函数调用返回值
	Type string `json:"type" gorm:"type:varchar(32);not null;comment:消息类型"`

	// 下面字段仅在type为message时有用
	Role    string `json:"role" gorm:"type:varchar(32);not null;comment:消息角色"`
	Content string `json:"content" gorm:"type:text;not null;comment:消息内容"`

	// 下面字段仅在消息类型是function_call 和 function_call_output时有用
	FunctionCallID        string `json:"function_call_id" gorm:"type:varchar(64);not null;comment:函数调用ID"`
	FunctionCallName      string `json:"function_call_name" gorm:"type:varchar(128);not null;comment:函数调用名称"`
	FunctionCallArguments string `json:"function_call_arguments" gorm:"type:text;not null;comment:函数调用参数"`
	FunctionCallOutput    string `json:"function_call_output" gorm:"type:text;not null;comment:函数调用返回值"`

	// token数统计，在type是message，且role是system 和 user时都为0，或者function_call_output时为0，其他情况下有值
	// 总之就是在服务器端回复的消息才有值
	PromptTokenCount     int `json:"prompt_token_count" gorm:"type:int;not null;comment:提示词token数"`
	CompletionTokenCount int `json:"completion_token_count" gorm:"type:int;not null;comment:回复token数"`
	TotalTokenCount      int `json:"total_token_count" gorm:"type:int;not null;comment:总token数"`

	// 附件消息 {id: 'xxx', name: 'xxx.docx'}[]这种格式的json
	AttachmentsInfo string `json:"attachments_info" gorm:"type:text;not null;comment:附件信息"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ChatAgentMessage) TableName() string {
	return "ltc_chat_agent_message"
}
