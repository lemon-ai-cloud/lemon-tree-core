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
	appHandler                        *handler.ApplicationHandler                // Application 处理器
	userHandler                       *handler.UserHandler                       // User 处理器
	llmProviderHandler                *handler.LlmProviderHandler                // LlmProvider 处理器
	applicationLlmHandler             *handler.ApplicationLlmHandler             // ApplicationLLM 处理器
	applicationMcpServerConfigHandler *handler.ApplicationMcpServerConfigHandler // ApplicationMCP配置 处理器
	chatAgentHandler                  *handler.ChatAgentHandler                  // ChatAgent 处理器
	chatAgentConversationHandler      *handler.ChatAgentConversationHandler      // ChatAgentConversation 处理器
	applicationStorageConfigHandler   *handler.ApplicationStorageConfigHandler   // ApplicationStorageConfig 处理器
	resourceHandler                   *handler.ResourceHandler                   // Resource 处理器
	chatAgentMcpServerToolHandler     *handler.ChatAgentMcpServerToolHandler     // ChatAgentMcpServerTool 处理器
	userService                       service.UserService                        // User 服务
	chatAgentService                  service.ChatAgentService                   // ChatAgent 服务
	applicationService                service.ApplicationService                 // Application 服务
	logger                            *zap.Logger                                // 日志记录器
}

// NewRouterManager 创建路由管理器实例
// 返回 RouterManager 的实例
// 参数：appHandler - Application 处理器，userHandler - User 处理器，llmProviderHandler - LlmProvider 处理器，applicationLlmHandler - ApplicationLLM 处理器，applicationMcpServerConfigHandler - ApplicationMCP配置 处理器，chatAgentHandler - ChatAgent 处理器，chatAgentConversationHandler - ChatAgentConversation 处理器，applicationStorageConfigHandler - ApplicationStorageConfig 处理器，resourceHandler - Resource 处理器，chatAgentMcpServerToolHandler - ChatAgentMcpServerTool 处理器，userService - User 服务，chatAgentService - ChatAgent 服务，applicationService - Application 服务，logger - 日志记录器
func NewRouterManager(appHandler *handler.ApplicationHandler, userHandler *handler.UserHandler, llmProviderHandler *handler.LlmProviderHandler, applicationLlmHandler *handler.ApplicationLlmHandler, applicationMcpServerConfigHandler *handler.ApplicationMcpServerConfigHandler, chatAgentHandler *handler.ChatAgentHandler, chatAgentConversationHandler *handler.ChatAgentConversationHandler, applicationStorageConfigHandler *handler.ApplicationStorageConfigHandler, resourceHandler *handler.ResourceHandler, chatAgentMcpServerToolHandler *handler.ChatAgentMcpServerToolHandler, userService service.UserService, chatAgentService service.ChatAgentService, applicationService service.ApplicationService, logger *zap.Logger) *RouterManager {
	return &RouterManager{
		appHandler:                        appHandler,
		userHandler:                       userHandler,
		llmProviderHandler:                llmProviderHandler,
		applicationLlmHandler:             applicationLlmHandler,
		applicationMcpServerConfigHandler: applicationMcpServerConfigHandler,
		chatAgentHandler:                  chatAgentHandler,
		chatAgentConversationHandler:      chatAgentConversationHandler,
		applicationStorageConfigHandler:   applicationStorageConfigHandler,
		resourceHandler:                   resourceHandler,
		chatAgentMcpServerToolHandler:     chatAgentMcpServerToolHandler,
		userService:                       userService,
		chatAgentService:                  chatAgentService,
		applicationService:                applicationService,
		logger:                            logger,
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

		// 设置 ApplicationMCP配置 模块的路由
		SetupApplicationMcpServerConfigRoutes(api, rm.applicationMcpServerConfigHandler)

		// 设置 ChatAgent 模块的路由
		SetupChatAgentRoutes(api, rm.chatAgentHandler)

		// 设置 ChatAgentConversation 模块的路由
		SetupChatAgentConversationRoutes(api, rm.chatAgentConversationHandler, rm.chatAgentService, rm.applicationService)

		// 设置 ApplicationStorageConfig 模块的路由
		SetupApplicationStorageConfigRoutes(api, rm.applicationStorageConfigHandler)

		// 设置 Resource 模块的路由
		SetupResourceRoutes(api, rm.resourceHandler)

		// 设置 ChatAgentMcpServerTool 模块的路由
		SetupChatAgentMcpServerToolRoutes(api, rm.chatAgentMcpServerToolHandler, rm.userService)

		// 未来可以在这里添加更多模块的路由
		// SetupAuthRoutes(api, authHandler)
	}

	return r
}
