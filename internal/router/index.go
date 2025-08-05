// Package router 提供路由管理功能
// 负责设置和管理 HTTP 路由，包括中间件配置和模块路由注册
package router

import (
	"lemon-tree-core/internal/handler"
	middleware2 "lemon-tree-core/internal/middleware"
	"lemon-tree-core/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RouterManager 路由管理器
// 负责管理所有路由的配置和中间件
// 提供统一的路由设置接口
type RouterManager struct {
	appHandler               *handler.ApplicationHandler       // Application 处理器
	userHandler              *handler.UserHandler              // User 处理器
	llmProviderDefineHandler *handler.LlmProviderDefineHandler // LlmProviderDefine 处理器
	resourceHandler          *handler.ResourceHandler          // Resource 处理器
	userService              service.UserService               // User 服务
	logger                   *zap.Logger                       // 日志记录器
}

// NewRouterManager 创建路由管理器实例
// 返回 RouterManager 的实例
// 参数：appHandler - Application 处理器，userHandler - User 处理器，llmProviderDefineHandler - LlmProviderDefine 处理器，userService - User 服务，logger - 日志记录器
func NewRouterManager(appHandler *handler.ApplicationHandler, userHandler *handler.UserHandler, llmProviderDefineHandler *handler.LlmProviderDefineHandler, resourceHandler *handler.ResourceHandler, userService service.UserService, logger *zap.Logger) *RouterManager {
	return &RouterManager{
		appHandler:               appHandler,
		userHandler:              userHandler,
		llmProviderDefineHandler: llmProviderDefineHandler,
		resourceHandler:          resourceHandler,
		userService:              userService,
		logger:                   logger,
	}
}

// SetupResourceRoutes 设置资源文件模块的路由
// 配置 Resource 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - Resource 处理器
func SetupResourceRoutes(api *gin.RouterGroup, handler *handler.ResourceHandler) {
	// 资源文件路由组
	resources := api.Group("/resources")
	{
		// 下载文件
		resources.GET("/download", handler.DownloadFile)

		// 列出目录文件
		resources.GET("/list", handler.ListFiles)

		// 获取文件信息
		resources.GET("/info", handler.GetFileInfo)
	}
}

// SetupAllRoutes 设置所有路由
// 配置中间件、API 路由组和各模块的路由
// 返回配置完成的 Gin 引擎实例
func (rm *RouterManager) SetupAllRoutes() *gin.Engine {
	// 创建新的 Gin 引擎实例
	r := gin.New()

	// 添加中间件
	// 恢复中间件：处理 panic 并记录错误
	r.Use(middleware2.RecoveryMiddleware(rm.logger))
	// 日志中间件：记录 HTTP 请求日志
	r.Use(middleware2.LoggerMiddleware(rm.logger))
	// CORS 中间件：处理跨域请求
	r.Use(middleware2.CORSMiddleware())

	// API 路由组
	// 所有 API 路由都以 /api/v1 为前缀
	api := r.Group("/api/v1")
	{
		// 注册各个模块的路由
		// 设置 Application 模块的路由
		SetupApplicationRoutes(api, rm.appHandler, rm.userService)

		// 设置 User 模块的路由
		SetupUserRoutes(api, rm.userHandler, rm.userService)

		// 设置 LlmProviderDefine 模块的路由
		SetupLlmProviderDefineRoutes(api, rm.llmProviderDefineHandler)

		// 设置 Resource 模块的路由
		SetupResourceRoutes(api, rm.resourceHandler)

		// 未来可以在这里添加更多模块的路由
		// SetupAuthRoutes(api, authHandler)
	}

	return r
}

// SetupLlmProviderDefineRoutes 设置大语言模型提供商定义模块的路由
// 配置 LlmProviderDefine 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - LlmProviderDefine 处理器
func SetupLlmProviderDefineRoutes(api *gin.RouterGroup, handler *handler.LlmProviderDefineHandler) {
	// 大语言模型提供商定义路由组
	llmProviderDefines := api.Group("/llm-provider-defines")
	{
		// 获取所有提供商定义
		llmProviderDefines.GET("", handler.GetAllLlmProviderDefines)

		// 根据ID获取提供商定义
		llmProviderDefines.GET("/:id", handler.GetLlmProviderDefineByID)

		// 保存提供商定义（创建或更新）
		llmProviderDefines.POST("/save", handler.SaveLlmProviderDefine)

		// 上传提供商定义图标
		llmProviderDefines.POST("/upload-icon", handler.UploadLlmProviderDefineIcon)

		// 动态查询提供商定义
		llmProviderDefines.POST("/query", handler.QueryLlmProviderDefines)

		// 删除提供商定义
		llmProviderDefines.DELETE("/:id", handler.DeleteLlmProviderDefine)
	}
}
