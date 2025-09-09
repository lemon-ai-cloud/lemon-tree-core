// Package router 提供路由配置功能
// 负责定义HTTP路由和中间件配置
package router

import (
	"lemon-tree-core/internal/handler"
	"lemon-tree-core/internal/middleware"
	"lemon-tree-core/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupChatAgentMcpServerToolRoutes 设置聊天智能体MCP工具相关路由
// 参数：api - API 路由组，chatAgentMcpServerToolHandler - 聊天智能体MCP工具处理器，userService - 用户服务
func SetupChatAgentMcpServerToolRoutes(api *gin.RouterGroup, chatAgentMcpServerToolHandler *handler.ChatAgentMcpServerToolHandler, userService service.UserService) {
	// 创建聊天智能体MCP工具路由组
	chatAgentMcpServerToolGroup := api.Group("/chat-agents")

	// 应用认证中间件
	chatAgentMcpServerToolGroup.Use(middleware.UserAuthMiddleware(userService))

	// 保存聊天智能体的MCP工具设置
	// PUT /api/v1/chat-agents/:chatAgentID/mcp-tools
	chatAgentMcpServerToolGroup.PUT("/:chatAgentID/mcp-tools", chatAgentMcpServerToolHandler.SaveChatAgentMcpServerToolSettings)

	// 获取聊天智能体的MCP工具设置
	// GET /api/v1/chat-agents/:chatAgentID/mcp-tools
	chatAgentMcpServerToolGroup.GET("/:chatAgentID/mcp-tools", chatAgentMcpServerToolHandler.GetChatAgentAvailableMcpServerTools)
}
