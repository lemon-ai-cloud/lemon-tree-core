package al_client

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// OpenAIChatCompletionsClient OpenAI聊天完成客户端实现
type OpenAIChatCompletionsClient struct {
	client *openai.Client
}

// NewOpenAIChatCompletionsClient 创建OpenAI聊天完成客户端
func NewOpenAIChatCompletionsClient(apiKey string) *OpenAIChatCompletionsClient {
	return &OpenAIChatCompletionsClient{
		client: openai.NewClient(apiKey),
	}
}

// SendMessage 发送消息
func (c *OpenAIChatCompletionsClient) SendMessage(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 转换请求格式
	openaiReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    convertToOpenAIMessages(req.Messages),
		Stream:      req.Stream,
		Tools:       convertToOpenAITools(req.Tools),
		Temperature: req.Temperature,
		TopP:        req.TopP,
		ToolChoice:  req.ToolChoice,
		MaxTokens:   req.MaxTokens,
	}

	// 调用OpenAI API
	response, err := c.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	return &ChatCompletionResponse{
		Choices: convertToLemonChoices(response.Choices),
	}, nil
}

// SendMessageStream 发送流式消息
func (c *OpenAIChatCompletionsClient) SendMessageStream(ctx context.Context, req ChatCompletionRequest) (ChatCompletionStream, error) {
	// 转换请求格式
	openaiReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    convertToOpenAIMessages(req.Messages),
		Stream:      req.Stream,
		Tools:       convertToOpenAITools(req.Tools),
		Temperature: req.Temperature,
		TopP:        req.TopP,
		ToolChoice:  req.ToolChoice,
		MaxTokens:   req.MaxTokens,
	}

	// 调用OpenAI流式API
	stream, err := c.client.CreateChatCompletionStream(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	// 包装流式响应
	return &OpenAIStreamWrapper{stream: stream}, nil
}

// OpenAIStreamWrapper OpenAI流式响应包装器
type OpenAIStreamWrapper struct {
	stream *openai.ChatCompletionStream
}

// Recv 接收流式数据
func (w *OpenAIStreamWrapper) Recv() (*ChatCompletionStreamResponse, error) {
	chunk, err := w.stream.Recv()
	if err != nil {
		return nil, err
	}

	return &ChatCompletionStreamResponse{
		Choices: convertToLemonStreamChoices(chunk.Choices),
	}, nil
}

// Close 关闭流
func (w *OpenAIStreamWrapper) Close() {
	w.stream.Close()
}

// 转换函数

// convertToOpenAIMessages 转换消息格式
func convertToOpenAIMessages(messages []ChatMessage) []openai.ChatCompletionMessage {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMsg := openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}

		// 处理工具调用
		if len(msg.ToolCalls) > 0 {
			openaiMsg.ToolCalls = convertToOpenAIToolCalls(msg.ToolCalls)
		}

		// 处理工具调用ID
		if msg.ToolCallID != "" {
			openaiMsg.ToolCallID = msg.ToolCallID
		}

		openaiMessages[i] = openaiMsg
	}
	return openaiMessages
}

// convertToOpenAITools 转换工具格式
func convertToOpenAITools(tools []Tool) []openai.Tool {
	openaiTools := make([]openai.Tool, len(tools))
	for i, tool := range tools {
		openaiTools[i] = openai.Tool{
			Type: openai.ToolType(tool.Type),
			Function: &openai.FunctionDefinition{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters:  tool.Function.Parameters,
			},
		}
	}
	return openaiTools
}

// convertToLemonChoices 转换选择格式
func convertToLemonChoices(choices []openai.ChatCompletionChoice) []ChatCompletionChoice {
	lemonChoices := make([]ChatCompletionChoice, len(choices))
	for i, choice := range choices {
		lemonMsg := ChatMessage{
			Role:    choice.Message.Role,
			Content: choice.Message.Content,
		}

		// 处理工具调用
		if choice.Message.ToolCalls != nil {
			lemonMsg.ToolCalls = convertToLemonToolCalls(choice.Message.ToolCalls)
		}

		// 处理工具调用ID
		if choice.Message.ToolCallID != "" {
			lemonMsg.ToolCallID = choice.Message.ToolCallID
		}

		lemonChoices[i] = ChatCompletionChoice{
			Message:      lemonMsg,
			FinishReason: string(choice.FinishReason),
		}
	}
	return lemonChoices
}

// convertToLemonStreamChoices 转换流式选择格式
func convertToLemonStreamChoices(choices []openai.ChatCompletionStreamChoice) []ChatCompletionStreamChoice {
	lemonChoices := make([]ChatCompletionStreamChoice, len(choices))
	for i, choice := range choices {
		lemonChoices[i] = ChatCompletionStreamChoice{
			Delta: ChatCompletionStreamDelta{
				Content:   choice.Delta.Content,
				ToolCalls: convertToLemonToolCalls(choice.Delta.ToolCalls),
			},
			FinishReason: string(choice.FinishReason),
		}
	}
	return lemonChoices
}

// convertToOpenAIToolCalls 转换工具调用格式到OpenAI
func convertToOpenAIToolCalls(toolCalls []ToolCall) []openai.ToolCall {
	openaiToolCalls := make([]openai.ToolCall, len(toolCalls))
	for i, toolCall := range toolCalls {
		openaiToolCalls[i] = openai.ToolCall{
			ID:   toolCall.ID,
			Type: openai.ToolType(toolCall.Type),
			Function: openai.FunctionCall{
				Name:      toolCall.Function.Name,
				Arguments: toolCall.Function.Arguments,
			},
		}
	}
	return openaiToolCalls
}

// convertToLemonToolCalls 转换工具调用格式
func convertToLemonToolCalls(toolCalls []openai.ToolCall) []ToolCall {
	lemonToolCalls := make([]ToolCall, len(toolCalls))
	for i, toolCall := range toolCalls {
		lemonToolCalls[i] = ToolCall{
			ID:   toolCall.ID,
			Type: string(toolCall.Type),
			Function: FunctionCall{
				Name:      toolCall.Function.Name,
				Arguments: toolCall.Function.Arguments,
			},
		}
	}
	return lemonToolCalls
}
