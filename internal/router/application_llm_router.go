// Package router 提供路由管理功能
package router

import (
	"lemon-tree-core/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupApplicationLlmRoutes 设置应用模型模块的路由
// 配置 ApplicationLLM 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - ApplicationLLM 处理器
func SetupApplicationLlmRoutes(api *gin.RouterGroup, handler *handler.ApplicationLlmHandler) {
	// 应用模型路由组
	applicationLlms := api.Group("/application-llms")
	{
		// 保存应用模型信息（创建或更新）
		// POST /api/v1/application-llms/save
		// 保存应用模型信息，如果存在则更新，不存在则创建
		applicationLlms.POST("/save", handler.SaveApplicationLlm)

		// 更新模型启用状态
		// PUT /api/v1/application-llms/:id/enabled
		// 更新指定模型的启用状态
		applicationLlms.PUT("/:id/enabled", handler.UpdateEnabledStatus)

		// 根据提供商ID获取模型列表
		// GET /api/v1/application-llms/provider/:providerId
		// 根据提供商ID获取该提供商下的所有模型列表
		applicationLlms.GET("/provider/:providerId", handler.GetModelsByProviderID)

		// 根据应用ID获取模型列表
		// GET /api/v1/application-llms/application/:applicationId
		// 根据应用ID获取该应用下的所有模型列表
		applicationLlms.GET("/application/:applicationId", handler.GetModelsByApplicationID)

		// 获取并保存模型列表
		// POST /api/v1/application-llms/provider/:providerId/fetch
		// 从提供商获取模型列表并保存到数据库
		applicationLlms.POST("/provider/:providerId/fetch", handler.FetchAndSaveModels)
	}
}
