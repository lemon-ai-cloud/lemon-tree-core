// Package di 提供依赖注入容器功能
// 使用 Uber FX 框架管理应用程序的依赖关系
// 负责组件的生命周期管理和依赖注入
package core

import (
	"context"
	"lemon-tree-core/internal/config"
	"lemon-tree-core/internal/handler"
	"lemon-tree-core/internal/repository"
	"lemon-tree-core/internal/router"
	"lemon-tree-core/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewContainer 创建依赖注入容器
// 配置所有组件的依赖关系和生命周期
// 按类型分组注册，便于管理和维护
// 返回配置完成的 FX 应用程序实例
func NewContainer() *fx.App {
	return fx.New(
		// 基础设施提供者（Infrastructure Providers）
		// 包含配置、数据库、日志等基础组件
		fx.Provide(
			config.LoadConfig, // 加载配置文件
			NewDatabase,       // 创建数据库连接
			NewLogger,         // 创建日志记录器
		),

		// Repository 层提供者（Repository Providers）
		// 包含所有数据访问层的组件
		fx.Provide(
			repository.NewApplicationRepository,                            // 创建 Application Repository
			repository.NewSystemUserRepository,                             // 创建 SystemUser Repository
			repository.NewSystemUserSessionRepository,                      // 创建 SystemUserSession Repository
			repository.NewApplicationModelRepository,                       // 创建 ApplicationModel Repository
			repository.NewChatAgentRepository,                              // 创建 ChatAgent Repository
			repository.NewChatAgentConversationRepository,                  // 创建 ChatAgentConversation Repository
			repository.NewChatAgentMessageRepository,                       // 创建 ChatAgentMessage Repository
			repository.NewChatAgentAttachmentRepository,                    // 创建 ChatAgentAttachment Repository
			repository.NewChatAgentApiKeyRepository,                        // 创建 ChatAgentApiKey Repository
			repository.NewChatConversationRepository,                       // 创建 ChatConversation Repository
			repository.NewLlmProviderRepository,                            // 创建 LlmProvider Repository
			repository.NewLlmProviderDefineRepository,                      // 创建 LlmProviderDefine Repository
			repository.NewApplicationStorageConfigRepository,               // 创建 ApplicationStorageConfig Repository
			repository.NewApplicationInternalToolNetSearchConfigRepository, // 创建 ApplicationInternalToolNetSearchConfig Repository
			repository.NewApplicationMcpConfigConfigRepository,             // 创建 ApplicationMcpConfigConfig Repository
			// 未来可以在这里添加更多 Repository
			// repository.NewUserRepository,
			// repository.NewOrderRepository,
		),

		// Service 层提供者（Service Providers）
		// 包含所有业务逻辑层的组件
		fx.Provide(
			service.NewApplicationService, // 创建 Application Service
			service.NewUserService,        // 创建 User Service
			// 未来可以在这里添加更多 Service
			// service.NewOrderService,
		),

		// Handler 层提供者（Handler Providers）
		// 包含所有 HTTP 请求处理层的组件
		fx.Provide(
			handler.NewApplicationHandler, // 创建 Application Handler
			handler.NewUserHandler,        // 创建 User Handler
			// 未来可以在这里添加更多 Handler
			// handler.NewOrderHandler,
		),

		// 路由和 Web 层提供者（Router & Web Providers）
		// 包含路由管理和 Web 服务器相关组件
		fx.Provide(
			router.NewRouterManager, // 创建路由管理器
			// 创建 Gin 引擎实例
			func(rm *router.RouterManager) *gin.Engine {
				return rm.SetupAllRoutes()
			},
		),

		// 启动钩子（Invokes）
		// 在应用程序启动时执行的函数
		fx.Invoke(StartServer),
	)
}

// NewLogger 创建日志记录器
// 使用 Zap 库创建生产环境的日志记录器
// 返回 Zap 日志记录器实例和错误信息
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

// StartServer 启动服务器
// 配置 HTTP 服务器的启动和关闭逻辑
// 使用 FX 的生命周期管理功能
// 参数：lifecycle - FX 生命周期管理器，router - Gin 路由引擎，config - 应用程序配置，logger - 日志记录器
func StartServer(
	lifecycle fx.Lifecycle,
	router *gin.Engine,
	config *config.Config,
	logger *zap.Logger,
) {
	// 创建 HTTP 服务器实例
	server := &http.Server{
		Addr:    config.Server.Port, // 服务器监听地址
		Handler: router,             // HTTP 请求处理器
	}

	// 配置应用程序生命周期钩子
	lifecycle.Append(fx.Hook{
		// OnStart 在应用程序启动时执行
		OnStart: func(context.Context) error {
			logger.Info("Starting server", zap.String("port", config.Server.Port))
			// 在后台 goroutine 中启动服务器
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		// OnStop 在应用程序停止时执行
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server")
			// 设置关闭超时时间
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			// 优雅关闭服务器
			return server.Shutdown(ctx)
		},
	})
}
