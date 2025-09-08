// Package router 提供路由管理功能
// 负责设置和管理 HTTP 路由，包括中间件配置和模块路由注册
package router

import (
	"lemon-tree-core/internal/handler"
	"lemon-tree-core/internal/middleware"
	"lemon-tree-core/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupChatAgentConversationRoutes 设置聊天会话模块的路由
// 配置 ChatAgentConversation 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - ChatAgentConversation 处理器，chatAgentService - ChatAgent 服务
func SetupChatAgentConversationRoutes(api *gin.RouterGroup, handler *handler.ChatAgentConversationHandler,
	chatAgentService service.ChatAgentService, applicationService service.ApplicationService) {
	// 聊天会话路由组
	chatAgentConversations := api.Group("/chat-agent-conversations")
	chatAgentConversations.Use(middleware.ChatAgentAuthMiddleware(chatAgentService, applicationService))
	{
		// 获取会话列表
		// GET /api/v1/chat-agent-conversations/conversation-list
		// 获取指定智能体的会话列表
		chatAgentConversations.GET("/conversation-list", handler.GetConversationList)

		// 获取聊天消息列表
		// GET /api/v1/chat-agent-conversations/message-list
		// 获取指定会话的消息列表
		chatAgentConversations.GET("/message-list", handler.GetChatMessageList)

		// 发送消息（非流式）
		// POST /api/v1/chat-agent-conversations/send-message
		// 发送消息并等待完整回复
		chatAgentConversations.POST("/send-message", handler.SendMessage)

		// 发送消息（流式）
		// POST /api/v1/chat-agent-conversations/send-message-streamable
		// 发送消息并流式返回回复
		chatAgentConversations.POST("/send-message-streamable", handler.SendMessageStreamable)

		// 发送预制答案（非流式）
		// POST /api/v1/chat-agent-conversations/send-message-predefined
		// 发送预制答案并等待完整回复
		chatAgentConversations.POST("/send-message-predefined", handler.SendMessagePredefined)

		// 发送预制答案（流式）
		// POST /api/v1/chat-agent-conversations/send-message-predefined-streamable
		// 发送预制答案并流式返回回复
		chatAgentConversations.POST("/send-message-predefined-streamable", handler.SendMessagePredefinedStreamable)

		// 上传附件
		// POST /api/v1/chat-agent-conversations/upload-attachment
		// 上传聊天附件文件
		chatAgentConversations.POST("/upload-attachment", handler.UploadAttachment)

		// 删除会话
		// DELETE /api/v1/chat-agent-conversations/conversation
		// 删除指定的会话及其所有消息
		chatAgentConversations.DELETE("/conversation", handler.DeleteConversation)

		// 重命名会话
		// PUT /api/v1/chat-agent-conversations/conversation-title
		// 重命名指定的会话标题
		chatAgentConversations.PUT("/conversation-title", handler.RenameConversationTitle)
	}
}
