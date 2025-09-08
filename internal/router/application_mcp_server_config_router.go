// Package router 提供路由管理功能
package router

import (
	"lemon-tree-core/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupApplicationMcpServerConfigRoutes 设置应用MCP配置模块的路由
// 配置 ApplicationMCP配置 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - ApplicationMCP配置 处理器
func SetupApplicationMcpServerConfigRoutes(api *gin.RouterGroup, handler *handler.ApplicationMcpServerConfigHandler) {
	// 应用MCP配置路由组
	applicationMcpServerConfigs := api.Group("/application-mcp-server-configs")
	{
		// 保存应用MCP配置信息（创建或更新）
		// POST /api/v1/application-mcp-server-configs/save
		// 保存应用MCP配置信息，如果存在则更新，不存在则创建
		applicationMcpServerConfigs.POST("/save", handler.SaveApplicationMcpServerConfig)

		// 删除MCP配置
		// DELETE /api/v1/application-mcp-server-configs/:id
		// 删除指定的MCP配置
		applicationMcpServerConfigs.DELETE("/:id", handler.DeleteApplicationMcpServerConfig)

		// 根据应用ID获取MCP配置列表
		// GET /api/v1/application-mcp-server-configs/application/:applicationId
		// 根据应用ID获取该应用下的所有MCP配置列表
		applicationMcpServerConfigs.GET("/application/:applicationId", handler.GetMcpServerConfigsByApplicationID)

		// 获取MCP服务器的所有工具
		// GET /api/v1/application-mcp-server-configs/:id/tools
		// 获取指定MCP服务器的所有工具列表
		applicationMcpServerConfigs.GET("/:id/tools", handler.GetMcpServerTools)

		// 同步MCP服务器的工具列表
		// POST /api/v1/application-mcp-server-configs/:id/sync-tools
		// 同步指定MCP服务器的工具列表到数据库
		applicationMcpServerConfigs.POST("/:id/sync-tools", handler.SyncMcpServerTools)
	}
}
