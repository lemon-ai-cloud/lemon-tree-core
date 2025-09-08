// Package dto 提供数据传输对象定义
// 用于在不同层之间传递数据，确保数据格式的一致性
package dto

// ChatMessageResponseEventDto 聊天消息响应事件
// 用于流式返回聊天消息更新
type ChatMessageResponseEventDto struct {
	ConversationID string `json:"conversation_id"` // 会话ID
	RequestID      string `json:"request_id"`      // 请求ID
	MessageType    string `json:"message_type"`    // 消息类型：answer answer_delta 消息回复，tool_call tool_call_processing tool_call_end 工具调用
	Content        string `json:"content"`         // 内容：answer，内容为消息内容；当消息类型为tool_call时，内容为调用的工具名字
}

// ChatMessageUseToolDto 聊天消息使用工具
type ChatMessageUseToolDto struct {
	ApplicationMcpConfigID string `json:"application_mcp_config_id"` // 应用mcp配置ID
	ToolName               string `json:"tool_name"`                 // 工具名称
}

// ChatUserSendMessageRequest 用户发送消息请求
type ChatUserSendMessageRequest struct {
	ServiceUserID        string                  `json:"service_user_id"`         // 业务侧用户ID
	SystemPrompt         string                  `json:"system_prompt"`           // 系统提示词
	UserMessage          string                  `json:"user_message"`            // 用户消息
	PredefinedAnswer     *string                 `json:"predefined_answer"`       // 预制答案（可选）
	UsedMcpToolList      []ChatMessageUseToolDto `json:"used_mcp_tool_list"`      // 使用的MCP工具列表
	UsedInternalToolList []string                `json:"used_internal_tool_list"` // 使用的内部工具列表
	ConversationID       *string                 `json:"conversation_id"`         // 会话ID（可选）
	Attachments          []string                `json:"attachments"`             // 附件ID列表（可选）
}

// GetConversationListRequest 获取会话列表请求
type GetConversationListRequest struct {
	ServiceUserID string  `json:"service_user_id"` // 业务侧用户ID
	LastID        *string `json:"last_id"`         // 最后一个会话的ID，用于游标分页
	Size          *int    `json:"size"`            // 返回数量
	Sort          *string `json:"sort"`            // 排序方式，默认按创建时间倒序
}

// ConversationInfoDto 会话信息
type ConversationInfoDto struct {
	ID            string `json:"id"`              // 会话ID
	Title         string `json:"title"`           // 会话标题
	ApplicationID string `json:"application_id"`  // 应用ID
	ServiceUserID string `json:"service_user_id"` // 业务侧用户ID
	CreatedAt     *int64 `json:"created_at"`      // 创建时间（时间戳）
	UpdatedAt     *int64 `json:"updated_at"`      // 更新时间（时间戳）
}

// GetConversationListResponse 获取会话列表响应
type GetConversationListResponse struct {
	Conversations []ConversationInfoDto `json:"conversations"` // 会话列表
	TotalCount    int                   `json:"total_count"`   // 总数量
}

// GetChatMessageListRequest 获取聊天消息列表请求
type GetChatMessageListRequest struct {
	ConversationID string  `json:"conversation_id"` // 会话ID
	LastID         *string `json:"last_id"`         // 最后一个消息的ID，用于游标分页
	Size           *int    `json:"size"`            // 返回数量
	Sort           *string `json:"sort"`            // 排序方式，默认按创建时间倒序
}

// ChatMessageAttachmentInfoDto 聊天附件信息
type ChatMessageAttachmentInfoDto struct {
	ID   string `json:"id"`   // 附件ID
	Name string `json:"name"` // 附件名称
}

// ChatMessageInfoDto 聊天消息信息
type ChatMessageInfoDto struct {
	ID                    string                         `json:"id"`                      // 消息ID
	ApplicationID         string                         `json:"application_id"`          // 应用ID
	ConversationID        string                         `json:"conversation_id"`         // 会话ID
	RequestID             string                         `json:"request_id"`              // 请求ID
	Type                  string                         `json:"type"`                    // 消息类型
	Role                  *string                        `json:"role"`                    // 消息角色
	Content               *string                        `json:"content"`                 // 消息内容
	FunctionCallID        *string                        `json:"function_call_id"`        // 函数调用ID
	FunctionCallName      *string                        `json:"function_call_name"`      // 函数调用名称
	FunctionCallArguments *string                        `json:"function_call_arguments"` // 函数调用参数
	FunctionCallOutput    *string                        `json:"function_call_output"`    // 函数调用返回值
	PromptTokenCount      int                            `json:"prompt_token_count"`      // 提示词token数
	CompletionTokenCount  int                            `json:"completion_token_count"`  // 回复token数
	TotalTokenCount       int                            `json:"total_token_count"`       // 总token数
	CreatedAt             *int64                         `json:"created_at"`              // 创建时间（时间戳）
	UpdatedAt             *int64                         `json:"updated_at"`              // 更新时间（时间戳）
	AttachmentInfoList    []ChatMessageAttachmentInfoDto `json:"attachment_info_list"`    // 附件信息列表
}

// GetChatMessageListResponse 获取聊天消息列表响应
type GetChatMessageListResponse struct {
	Messages   []ChatMessageInfoDto `json:"messages"`    // 消息列表
	TotalCount int                  `json:"total_count"` // 总数量
}

// DeleteConversationRequest 删除会话请求
type DeleteConversationRequest struct {
	ConversationID string `json:"conversation_id"` // 会话ID
}

// DeleteConversationResponse 删除会话响应
type DeleteConversationResponse struct {
	Success                 bool    `json:"success"`                   // 是否成功
	Message                 *string `json:"message"`                   // 成功消息
	Error                   *string `json:"error"`                     // 错误消息
	DeletedMessagesCount    *int    `json:"deleted_messages_count"`    // 删除的消息数量
	DeletedAttachmentsCount *int    `json:"deleted_attachments_count"` // 删除的附件数量
}

// RenameConversationRequest 重命名会话请求
type RenameConversationRequest struct {
	ConversationID string `json:"conversation_id"` // 会话ID
	NewTitle       string `json:"new_title"`       // 新标题
}

// RenameConversationResponse 重命名会话响应
type RenameConversationResponse struct {
	Success        bool    `json:"success"`         // 是否成功
	Message        *string `json:"message"`         // 成功消息
	Error          *string `json:"error"`           // 错误消息
	ConversationID *string `json:"conversation_id"` // 会话ID
	NewTitle       *string `json:"new_title"`       // 新标题
}

// UploadAttachmentResponse 上传附件响应
type UploadAttachmentResponse struct {
	Success          bool    `json:"success"`           // 是否成功
	AttachmentID     *string `json:"attachment_id"`     // 附件ID
	OriginalFileName *string `json:"original_filename"` // 原始文件名
	FileSize         *int64  `json:"file_size"`         // 文件大小
	AttachmentType   *string `json:"attachment_type"`   // 附件类型
	IsProcessed      *bool   `json:"is_processed"`      // 是否已处理
	ProcessingError  *string `json:"processing_error"`  // 处理错误信息
	MarkdownContent  *string `json:"markdown_content"`  // Markdown内容
	Error            *string `json:"error"`             // 错误消息
}
