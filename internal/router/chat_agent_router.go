// Package router 提供路由管理功能
package router

import (
	"lemon-tree-core/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupChatAgentRoutes 设置智能体模块的路由
// 配置 ChatAgent 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - ChatAgent 处理器
func SetupChatAgentRoutes(api *gin.RouterGroup, handler *handler.ChatAgentHandler) {
	// 智能体路由组
	chatAgents := api.Group("/chat-agents")
	{
		// 保存智能体信息（创建或更新）
		// POST /api/v1/chat-agents/save
		// 保存智能体信息，如果存在则更新，不存在则创建
		chatAgents.POST("/save", handler.SaveChatAgent)

		// 删除智能体
		// DELETE /api/v1/chat-agents/:id
		// 删除指定的智能体
		chatAgents.DELETE("/:id", handler.DeleteChatAgent)

		// 根据应用ID获取智能体列表（分页）
		// GET /api/v1/chat-agents/application/:applicationId
		// 根据应用ID获取该应用下的所有智能体列表
		chatAgents.GET("/application/:applicationId", handler.GetChatAgentsByApplicationID)

		// 上传智能体头像
		// POST /api/v1/chat-agents/upload-avatar
		// 上传智能体头像文件
		chatAgents.POST("/upload-avatar", handler.UploadChatAgentAvatar)
	}
}
