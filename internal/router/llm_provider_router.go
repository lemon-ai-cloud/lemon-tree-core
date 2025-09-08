// Package router 提供路由管理功能
package router

import (
	"lemon-tree-core/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupLlmProviderRoutes 设置大语言模型提供商模块的路由
// 配置 LlmProvider 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - LlmProvider 处理器
func SetupLlmProviderRoutes(api *gin.RouterGroup, handler *handler.LlmProviderHandler) {
	// 大语言模型提供商路由组
	llmProviders := api.Group("/llm-providers")
	{
		// 获取所有提供商
		// GET /api/v1/llm-providers
		// 获取所有大语言模型提供商的列表
		llmProviders.GET("", handler.GetAllLlmProviders)

		// 根据ID获取提供商
		// GET /api/v1/llm-providers/:id
		// 根据 UUID 获取指定的提供商信息
		llmProviders.GET("/:id", handler.GetLlmProviderByID)

		// 保存提供商（创建或更新）
		// POST /api/v1/llm-providers/save
		// 如果提供商存在则更新，不存在则创建
		llmProviders.POST("/save", handler.SaveLlmProvider)

		// 上传提供商图标
		// POST /api/v1/llm-providers/upload-icon
		// 上传提供商图标文件
		llmProviders.POST("/upload-icon", handler.UploadLlmProviderIcon)

		// 动态查询提供商
		// POST /api/v1/llm-providers/query
		// 根据查询条件动态查询提供商
		llmProviders.POST("/query", handler.QueryLlmProviders)

		// 删除提供商
		// DELETE /api/v1/llm-providers/:id
		// 删除指定的提供商（软删除）
		llmProviders.DELETE("/:id", handler.DeleteLlmProvider)

		// 根据应用ID获取提供商列表
		// GET /api/v1/llm-providers/application/:applicationId
		// 根据应用ID获取该应用下的所有提供商列表
		llmProviders.GET("/application/:applicationId", handler.GetLlmProvidersByApplicationID)
	}
}
