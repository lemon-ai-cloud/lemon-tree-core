// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
)

// ChatAgentAttachment 聊天智能体的聊天附件
type ChatAgentAttachment struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	Title          string    `json:"title" gorm:"type:varchar(64);not null;comment:会话标题"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	ChatAgentID    uuid.UUID `json:"chat_agent_id" gorm:"type:char(36);not null;comment:所属Chat Agent ID"`
	ConversationID uuid.UUID `json:"conversation_id" gorm:"type:char(36);not null;comment:所属会话ID"`
	MessageID      uuid.UUID `json:"message_id" gorm:"type:char(36);not null;comment:所属消息ID"`

	// 文件信息
	OriginalFileName string `json:"original_file_name" gorm:"type:varchar(256);not null;comment:原始文件名"`
	FileExtension    string `json:"file_extension" gorm:"type:varchar(32);not null;comment:文件扩展名"`
	FileSize         int64  `json:"file_size" gorm:"type:bigint;not null;comment:文件大小"`
	MimeType         string `json:"mime_type" gorm:"type:varchar(128);not null;comment:MIME类型"`

	// 存储路径
	FilePath string `json:"file_path" gorm:"type:varchar(512);not null;comment:文件存储路径"`
	// 仅文档类型需要生成markdown时候有值
	MarkdownPath    string `json:"markdown_path" gorm:"type:varchar(512);not null;comment:Markdown文件存储路径"`
	AttachmentType  string `json:"attachment_type" gorm:"type:varchar(64);not null;comment:附件类型"`
	MarkdownContent string `json:"markdown_content" gorm:"type:text;not null;comment:Markdown内容"`
	IsProcessed     bool   `json:"is_processed" gorm:"type:tinyint(1);not null;comment:是否处理完毕"`
	ProcessingError string `json:"processing_error" gorm:"type:text;not null;comment:处理错误信息"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ChatAgentAttachment) TableName() string {
	return "ltc_chat_agent_attachment"
}
