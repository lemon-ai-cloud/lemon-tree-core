// Package router 提供路由管理功能
package router

import (
	"lemon-tree-core/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupApplicationStorageConfigRoutes 设置应用存储配置模块的路由
// 配置 ApplicationStorageConfig 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - ApplicationStorageConfig 处理器
func SetupApplicationStorageConfigRoutes(api *gin.RouterGroup, handler *handler.ApplicationStorageConfigHandler) {
	// 应用存储配置路由组
	applicationStorageConfigs := api.Group("/application-storage-configs")
	{
		// 保存应用存储配置
		// POST /api/v1/application-storage-configs/save
		// 保存应用存储配置信息，如果存在则更新，不存在则创建
		applicationStorageConfigs.POST("/save", handler.SaveApplicationStorageConfig)

		// 根据应用ID获取存储配置
		// GET /api/v1/application-storage-configs/application/:applicationId
		// 根据应用ID获取该应用的存储配置
		applicationStorageConfigs.GET("/application/:applicationId", handler.GetApplicationStorageConfigByApplicationID)
	}
}
