// Package dto 提供数据传输对象定义
// 用于在不同层之间传递数据，确保数据格式的一致性
package dto

// ChatAgentDto 智能体数据传输对象
// 用于在业务逻辑层和HTTP处理层之间传递数据
type ChatAgentDto struct {
	ID                             string  `json:"id"`                                  // 主键ID
	Name                           string  `json:"name"`                                // Agent名称
	Description                    string  `json:"description"`                         // Agent描述
	ApplicationID                  string  `json:"application_id"`                      // 所属应用ID
	AvatarUrl                      string  `json:"avatar_url"`                          // Agent的头像URL
	ChatSystemPrompt               string  `json:"system_prompt"`                       // 系统提示
	ChatModelID                    string  `json:"chat_model_id"`                       // 聊天模型ID
	ConversationNamingPrompt       string  `json:"conversation_naming_prompt"`          // 会话命名提示词
	ConversationNamingModelID      string  `json:"conversation_naming_model_id"`        // 会话命名模型ID
	ModelParamTemperature          float64 `json:"model_temperature"`                   // 模型温度
	ModelParamTopP                 float64 `json:"model_top_p"`                         // 模型TopP
	EnableContextLengthLimit       bool    `json:"enable_context_length_limit"`         // 是否启用上下文长度限制
	ContextLengthLimit             int     `json:"context_length_limit"`                // 上下文长度限制
	EnableMaxOutputTokenCountLimit bool    `json:"enable_max_output_token_count_limit"` // 是否启用最大输出Token数量限制
	MaxOutputTokenCountLimit       int     `json:"max_output_token_count_limit"`        // 最大输出Token数量
	DefaultStreamable              bool    `json:"default_streamable"`                  // 是否默认流式返回
	CreatedAt                      string  `json:"created_at"`                          // 创建时间
	UpdatedAt                      string  `json:"updated_at"`                          // 更新时间
}

// SaveChatAgentRequest 保存智能体请求
// 用于前端保存智能体的请求数据
type SaveChatAgentRequest struct {
	ID                             *string `json:"id,omitempty"`                        // 主键ID（更新时提供）
	Name                           string  `json:"name"`                                // Agent名称
	Description                    string  `json:"description"`                         // Agent描述
	ApplicationID                  string  `json:"application_id"`                      // 所属应用ID
	AvatarUrl                      string  `json:"avatar_url"`                          // Agent的头像URL
	ChatSystemPrompt               string  `json:"system_prompt"`                       // 系统提示
	ChatModelID                    string  `json:"chat_model_id"`                       // 聊天模型ID
	ConversationNamingPrompt       string  `json:"conversation_naming_prompt"`          // 会话命名提示词
	ConversationNamingModelID      string  `json:"conversation_naming_model_id"`        // 会话命名模型ID
	ModelParamTemperature          float64 `json:"model_temperature"`                   // 模型温度
	ModelParamTopP                 float64 `json:"model_top_p"`                         // 模型TopP
	EnableContextLengthLimit       bool    `json:"enable_context_length_limit"`         // 是否启用上下文长度限制
	ContextLengthLimit             int     `json:"context_length_limit"`                // 上下文长度限制
	EnableMaxOutputTokenCountLimit bool    `json:"enable_max_output_token_count_limit"` // 是否启用最大输出Token数量限制
	MaxOutputTokenCountLimit       int     `json:"max_output_token_count_limit"`        // 最大输出Token数量
	DefaultStreamable              bool    `json:"default_streamable"`                  // 是否默认流式返回
}

// ChatAgentListResponse 智能体列表响应
// 用于返回分页的智能体列表
type ChatAgentListResponse struct {
	ChatAgents []ChatAgentDto `json:"chat_agents"` // 智能体列表
	Total      int64          `json:"total"`       // 总数量
	Page       int            `json:"page"`        // 当前页码
	PageSize   int            `json:"page_size"`   // 每页大小
}

// SingleChatAgentResponse 单个智能体响应
// 用于返回单个智能体数据
type SingleChatAgentResponse struct {
	ChatAgent ChatAgentDto `json:"chat_agent"` // 智能体数据
}
