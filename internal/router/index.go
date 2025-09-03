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
	appHandler            *handler.ApplicationHandler    // Application 处理器
	userHandler           *handler.UserHandler           // User 处理器
	llmProviderHandler    *handler.LlmProviderHandler    // LlmProvider 处理器
	applicationLlmHandler *handler.ApplicationLlmHandler // ApplicationLLM 处理器
	resourceHandler       *handler.ResourceHandler       // Resource 处理器
	userService           service.UserService            // User 服务
	logger                *zap.Logger                    // 日志记录器
}

// NewRouterManager 创建路由管理器实例
// 返回 RouterManager 的实例
// 参数：appHandler - Application 处理器，userHandler - User 处理器，llmProviderHandler - LlmProvider 处理器，applicationLlmHandler - ApplicationLLM 处理器，resourceHandler - Resource 处理器，userService - User 服务，logger - 日志记录器
func NewRouterManager(appHandler *handler.ApplicationHandler, userHandler *handler.UserHandler, llmProviderHandler *handler.LlmProviderHandler, applicationLlmHandler *handler.ApplicationLlmHandler, resourceHandler *handler.ResourceHandler, userService service.UserService, logger *zap.Logger) *RouterManager {
	return &RouterManager{
		appHandler:            appHandler,
		userHandler:           userHandler,
		llmProviderHandler:    llmProviderHandler,
		applicationLlmHandler: applicationLlmHandler,
		resourceHandler:       resourceHandler,
		userService:           userService,
		logger:                logger,
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

		// 设置 LlmProvider 模块的路由
		SetupLlmProviderRoutes(api, rm.llmProviderHandler)

		// 设置 ApplicationLLM 模块的路由
		SetupApplicationLlmRoutes(api, rm.applicationLlmHandler)

		// 设置 Resource 模块的路由
		SetupResourceRoutes(api, rm.resourceHandler)

		// 未来可以在这里添加更多模块的路由
		// SetupAuthRoutes(api, authHandler)
	}

	return r
}

// SetupLlmProviderRoutes 设置大语言模型提供商模块的路由
// 配置 LlmProvider 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - LlmProvider 处理器
func SetupLlmProviderRoutes(api *gin.RouterGroup, handler *handler.LlmProviderHandler) {
	// 大语言模型提供商路由组
	llmProviders := api.Group("/llm-providers")
	{
		// 获取所有提供商
		llmProviders.GET("", handler.GetAllLlmProviders)

		// 根据ID获取提供商
		llmProviders.GET("/:id", handler.GetLlmProviderByID)

		// 保存提供商（创建或更新）
		llmProviders.POST("/save", handler.SaveLlmProvider)

		// 上传提供商图标
		llmProviders.POST("/upload-icon", handler.UploadLlmProviderIcon)

		// 动态查询提供商
		llmProviders.POST("/query", handler.QueryLlmProviders)

		// 删除提供商
		llmProviders.DELETE("/:id", handler.DeleteLlmProvider)

		// 根据应用ID获取提供商列表
		llmProviders.GET("/application/:applicationId", handler.GetLlmProvidersByApplicationID)
	}
}

// SetupApplicationLlmRoutes 设置应用模型模块的路由
// 配置 ApplicationLLM 相关的所有 HTTP 路由
// 参数：api - API 路由组，handler - ApplicationLLM 处理器
func SetupApplicationLlmRoutes(api *gin.RouterGroup, handler *handler.ApplicationLlmHandler) {
	// 应用模型路由组
	applicationLlms := api.Group("/application-llms")
	{
		// 保存应用模型信息（创建或更新）
		applicationLlms.POST("/save", handler.SaveApplicationLlm)

		// 更新模型启用状态
		applicationLlms.PUT("/:id/enabled", handler.UpdateEnabledStatus)

		// 根据提供商ID获取模型列表
		applicationLlms.GET("/provider/:providerId", handler.GetModelsByProviderID)

		// 根据应用ID获取模型列表
		applicationLlms.GET("/application/:applicationId", handler.GetModelsByApplicationID)
	}
}
