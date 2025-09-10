package al_client

import (
	"context"
)

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// Tool 工具定义结构
type Tool struct {
	Type     string              `json:"type"`
	Function *FunctionDefinition `json:"function"`
}

// FunctionDefinition 函数定义结构
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolCall 工具调用结构
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall 函数调用结构
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// SendMessageRequest 聊天完成请求结构
type SendMessageRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	Tools       []Tool        `json:"tools,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
	ToolChoice  string        `json:"tool_choice,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// SendMessageResponse 聊天完成响应结构
type SendMessageResponse struct {
	Choices []SendMessageChoice `json:"choices"`
}

// SendMessageChoice 聊天完成选择结构
type SendMessageChoice struct {
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// SendMessageStreamResponse 流式聊天完成响应结构
type SendMessageStreamResponse struct {
	Choices []SendMessageStreamChoice `json:"choices"`
}

// SendMessageStreamChoice 流式聊天完成选择结构
type SendMessageStreamChoice struct {
	Delta        SendMessageStreamDelta `json:"delta"`
	FinishReason string                 `json:"finish_reason"`
}

// SendMessageStreamDelta 流式聊天完成增量结构
type SendMessageStreamDelta struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// SendMessageStream 流式聊天完成接口
type SendMessageStream interface {
	Recv() (*SendMessageStreamResponse, error)
	Close()
}

// LemonAiClient AI客户端接口
type LemonAiClient interface {
	// SendMessage 发送消息
	SendMessage(ctx context.Context, req SendMessageRequest) (*SendMessageResponse, error)

	// SendMessageStream 发送流式消息
	SendMessageStream(ctx context.Context, req SendMessageRequest) (SendMessageStream, error)
}
