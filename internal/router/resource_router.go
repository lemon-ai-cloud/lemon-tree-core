// Package router 提供路由管理功能
package router

import (
	"lemon-tree-core/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupResourceRoutes 设置资源文件模块的路由
// 配置 Resource 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - Resource 处理器
func SetupResourceRoutes(api *gin.RouterGroup, handler *handler.ResourceHandler) {
	// 资源文件路由组
	resources := api.Group("/resources")
	{
		// 下载文件
		// GET /api/v1/resources/download
		// 下载指定的资源文件
		resources.GET("/download", handler.DownloadFile)

		// 列出目录文件
		// GET /api/v1/resources/list
		// 列出指定目录下的所有文件
		resources.GET("/list", handler.ListFiles)

		// 获取文件信息
		// GET /api/v1/resources/info
		// 获取指定文件的详细信息
		resources.GET("/info", handler.GetFileInfo)
	}
}
