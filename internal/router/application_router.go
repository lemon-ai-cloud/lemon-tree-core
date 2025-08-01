// Package router 提供路由管理功能
package router

import (
	"lemon-tree-core/internal/handler"
	"lemon-tree-core/internal/middleware"
	"lemon-tree-core/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupApplicationRoutes 设置 Application 相关路由
// 配置 Application 模块的所有 HTTP 路由
// 参数：api - API 路由组，appHandler - Application 处理器
func SetupApplicationRoutes(api *gin.RouterGroup, appHandler *handler.ApplicationHandler, userService service.UserService) {
	// Application 路由组
	// 所有 Application 相关的路由都以 /applications 为前缀
	applications := api.Group("/applications")
	{

		authenticated := applications.Group("")
		authenticated.Use(middleware.AuthMiddleware(userService))
		{
			// 获取应用列表
			// GET /api/v1/applications
			// 获取所有应用的列表
			applications.GET("", appHandler.GetAllApplications)

			// 根据ID获取应用
			// GET /api/v1/applications/:id
			// 根据 UUID 获取指定的应用信息
			applications.GET("/:id", appHandler.GetApplicationByID)

			// 保存应用（upsert）
			// POST /api/v1/applications/save
			// 如果应用存在则更新，不存在则创建
			applications.POST("/save", appHandler.SaveApplication)

			// 动态查询应用
			// POST /api/v1/applications/query
			// 根据查询条件动态查询应用
			applications.POST("/query", appHandler.QueryApplications)

			// 删除应用
			// DELETE /api/v1/applications/:id
			// 删除指定的应用（软删除）
			applications.DELETE("/:id", appHandler.DeleteApplication)
		}
	}
}
