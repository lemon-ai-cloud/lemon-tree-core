// Package handler 提供 HTTP 请求处理层功能
// 负责处理 HTTP 请求、参数验证、调用业务逻辑和返回响应
package handler

import (
	"encoding/json"
	"io"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ChatAgentConversationHandler 聊天会话 控制器
// 处理 聊天会话 相关的所有 HTTP 请求
// 相当于 Java Spring Boot 中的 Controller
type ChatAgentConversationHandler struct {
	chatAgentConversationService service.ChatAgentConversationService // 聊天会话 业务逻辑层接口
}

// NewChatAgentConversationHandler 创建 聊天会话 Handler 实例
// 返回 ChatAgentConversationHandler 的实例
// 参数：chatAgentConversationService - 聊天会话 业务逻辑层接口
func NewChatAgentConversationHandler(chatAgentConversationService service.ChatAgentConversationService) *ChatAgentConversationHandler {
	return &ChatAgentConversationHandler{
		chatAgentConversationService: chatAgentConversationService,
	}
}

// GetConversationList 获取会话列表
// 处理 GET /api/v1/chat-agent-conversations/conversation-list 请求
func (h *ChatAgentConversationHandler) GetConversationList(c *gin.Context) {
	// 获取查询参数
	serviceUserID := c.Query("service_user_id")
	if serviceUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_user_id 参数不能为空"})
		return
	}

	lastID := c.Query("last_id")
	sizeStr := c.DefaultQuery("size", "10")

	// 解析size参数
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// 从上下文获取智能体信息（通过中间件设置）
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	chatAgent, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 调用业务逻辑层获取会话列表
	conversations, err := h.chatAgentConversationService.GetConversationList(
		c.Request.Context(),
		chatAgent.ID.String(),
		serviceUserID,
		lastID,
		size,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	conversationList := make([]dto.ConversationInfoDto, 0, len(conversations))
	for _, conv := range conversations {
		createdAt := conv.CreatedAt.UnixMilli()
		updatedAt := conv.UpdatedAt.UnixMilli()
		conversationList = append(conversationList, dto.ConversationInfoDto{
			ID:            conv.ID.String(),
			Title:         conv.Title,
			ApplicationID: conv.ApplicationID.String(),
			ServiceUserID: conv.ServiceUserID,
			CreatedAt:     &createdAt,
			UpdatedAt:     &updatedAt,
		})
	}

	response := dto.GetConversationListResponse{
		Conversations: conversationList,
		TotalCount:    len(conversationList),
	}

	c.JSON(http.StatusOK, response)
}

// GetChatMessageList 获取聊天消息列表
// 处理 GET /api/v1/chat-agent-conversations/message-list 请求
func (h *ChatAgentConversationHandler) GetChatMessageList(c *gin.Context) {
	// 获取查询参数
	conversationID := c.Query("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversation_id 参数不能为空"})
		return
	}

	lastID := c.Query("last_id")
	sizeStr := c.DefaultQuery("size", "10")

	// 解析size参数
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// 从上下文获取智能体信息
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	chatAgent, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 调用业务逻辑层获取消息列表
	messages, err := h.chatAgentConversationService.GetChatMessageList(
		c.Request.Context(),
		chatAgent.ID.String(),
		conversationID,
		lastID,
		size,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	messageList := make([]dto.ChatMessageInfoDto, 0, len(messages))
	for _, msg := range messages {
		attachmentInfoList := make([]dto.ChatMessageAttachmentInfoDto, 0)
		if msg.AttachmentsInfo != "" {
			// 解析附件信息JSON
			var attachmentDictList []map[string]interface{}
			if err := json.Unmarshal([]byte(msg.AttachmentsInfo), &attachmentDictList); err == nil {
				for _, item := range attachmentDictList {
					if id, ok := item["id"].(string); ok {
						if name, ok := item["name"].(string); ok {
							attachmentInfoList = append(attachmentInfoList, dto.ChatMessageAttachmentInfoDto{
								ID:   id,
								Name: name,
							})
						}
					}
				}
			}
		}

		createdAt := msg.CreatedAt.UnixMilli()
		updatedAt := msg.UpdatedAt.UnixMilli()
		messageList = append(messageList, dto.ChatMessageInfoDto{
			ID:                    msg.ID.String(),
			ApplicationID:         msg.ApplicationID.String(),
			ConversationID:        msg.ConversationID.String(),
			RequestID:             msg.RequestID,
			Type:                  msg.Type,
			Role:                  &msg.Role,
			Content:               &msg.Content,
			FunctionCallID:        &msg.FunctionCallID,
			FunctionCallName:      &msg.FunctionCallName,
			FunctionCallArguments: &msg.FunctionCallArguments,
			FunctionCallOutput:    &msg.FunctionCallOutput,
			PromptTokenCount:      msg.PromptTokenCount,
			CompletionTokenCount:  msg.CompletionTokenCount,
			TotalTokenCount:       msg.TotalTokenCount,
			CreatedAt:             &createdAt,
			UpdatedAt:             &updatedAt,
			AttachmentInfoList:    attachmentInfoList,
		})
	}

	response := dto.GetChatMessageListResponse{
		Messages:   messageList,
		TotalCount: len(messageList),
	}

	c.JSON(http.StatusOK, response)
}

// SendMessage 自然语言对话
// 处理 POST /api/v1/chat-agent-conversations/send-message 请求
func (h *ChatAgentConversationHandler) SendMessage(c *gin.Context) {
	// 绑定JSON请求体
	var req dto.ChatUserSendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取智能体信息
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	_, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 调用业务逻辑层处理消息
	stream, err := h.chatAgentConversationService.UserSendMessage(
		c.Request.Context(),
		&req,
		false, // 非流式
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 流式返回响应
	c.Stream(func(w io.Writer) bool {
		buffer := make([]byte, 1024)
		n, err := stream.Read(buffer)
		if err != nil {
			return false
		}
		w.Write(buffer[:n])
		return true
	})
}

// SendMessageStreamable 自然语言对话-流式回复
// 处理 POST /api/v1/chat-agent-conversations/send-message-streamable 请求
func (h *ChatAgentConversationHandler) SendMessageStreamable(c *gin.Context) {
	// 绑定JSON请求体
	var req dto.ChatUserSendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取智能体信息
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	_, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 调用业务逻辑层处理消息
	stream, err := h.chatAgentConversationService.UserSendMessage(
		c.Request.Context(),
		&req,
		true, // 流式
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 流式返回响应
	c.Stream(func(w io.Writer) bool {
		buffer := make([]byte, 1024)
		n, err := stream.Read(buffer)
		if err != nil {
			return false
		}
		w.Write(buffer[:n])
		return true
	})
}

// SendMessagePredefined 自然语言对话-预制答案
// 处理 POST /api/v1/chat-agent-conversations/send-message-predefined 请求
func (h *ChatAgentConversationHandler) SendMessagePredefined(c *gin.Context) {
	// 绑定JSON请求体
	var req dto.ChatUserSendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证预制答案
	if req.PredefinedAnswer == nil || *req.PredefinedAnswer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预制答案不能为空"})
		return
	}

	// 从上下文获取智能体信息
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	_, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 调用业务逻辑层处理消息
	stream, err := h.chatAgentConversationService.UserSendMessagePredefinedAnswer(
		c.Request.Context(),
		&req,
		false, // 非流式
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 流式返回响应
	c.Stream(func(w io.Writer) bool {
		buffer := make([]byte, 1024)
		n, err := stream.Read(buffer)
		if err != nil {
			return false
		}
		w.Write(buffer[:n])
		return true
	})
}

// SendMessagePredefinedStreamable 自然语言对话-预制答案-流式回复
// 处理 POST /api/v1/chat-agent-conversations/send-message-predefined-streamable 请求
func (h *ChatAgentConversationHandler) SendMessagePredefinedStreamable(c *gin.Context) {
	// 绑定JSON请求体
	var req dto.ChatUserSendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证预制答案
	if req.PredefinedAnswer == nil || *req.PredefinedAnswer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预制答案不能为空"})
		return
	}

	// 从上下文获取智能体信息
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	_, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 调用业务逻辑层处理消息
	stream, err := h.chatAgentConversationService.UserSendMessagePredefinedAnswer(
		c.Request.Context(),
		&req,
		true, // 流式
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 流式返回响应
	c.Stream(func(w io.Writer) bool {
		buffer := make([]byte, 1024)
		n, err := stream.Read(buffer)
		if err != nil {
			return false
		}
		w.Write(buffer[:n])
		return true
	})
}

// UploadAttachment 上传聊天附件
// 处理 POST /api/v1/chat-agent-conversations/upload-attachment 请求
func (h *ChatAgentConversationHandler) UploadAttachment(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	// 从上下文获取智能体信息
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	chatAgent, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败"})
		return
	}
	defer src.Close()

	// 调用业务逻辑层上传附件
	result, err := h.chatAgentConversationService.UploadAttachment(
		c.Request.Context(),
		chatAgent.ID.String(),
		src,
		file.Filename,
		file.Size,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteConversation 删除会话
// 处理 DELETE /api/v1/chat-agent-conversations/conversation 请求
func (h *ChatAgentConversationHandler) DeleteConversation(c *gin.Context) {
	// 获取查询参数
	conversationID := c.Query("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversation_id 参数不能为空"})
		return
	}

	serviceUserID := c.Query("service_user_id")
	if serviceUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_user_id 参数不能为空"})
		return
	}

	// 从上下文获取智能体信息
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	chatAgent, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 调用业务逻辑层删除会话
	result, err := h.chatAgentConversationService.DeleteConversation(
		c.Request.Context(),
		chatAgent.ID.String(),
		serviceUserID,
		conversationID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RenameConversationTitle 重命名会话
// 处理 PUT /api/v1/chat-agent-conversations/conversation-title 请求
func (h *ChatAgentConversationHandler) RenameConversationTitle(c *gin.Context) {
	// 获取查询参数
	conversationID := c.Query("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversation_id 参数不能为空"})
		return
	}

	newTitle := c.Query("new_title")
	if newTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new_title 参数不能为空"})
		return
	}

	serviceUserID := c.Query("service_user_id")
	if serviceUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_user_id 参数不能为空"})
		return
	}

	// 从上下文获取智能体信息
	chatAgentValue, exists := c.Get(define.AppContextKeyCurrentChatAgent)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息未找到"})
		return
	}

	chatAgent, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体信息类型错误"})
		return
	}

	// 调用业务逻辑层重命名会话
	result, err := h.chatAgentConversationService.RenameConversationTitle(
		c.Request.Context(),
		chatAgent.ID.String(),
		serviceUserID,
		conversationID,
		newTitle,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
