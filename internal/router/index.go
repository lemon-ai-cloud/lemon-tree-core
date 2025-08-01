// Package router 提供路由管理功能
// 负责设置和管理 HTTP 路由，包括中间件配置和模块路由注册
package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lemon-tree-core/internal/handler"
	middleware2 "lemon-tree-core/internal/middleware"
)

// RouterManager 路由管理器
// 负责管理所有路由的配置和中间件
// 提供统一的路由设置接口
type RouterManager struct {
	appHandler *handler.ApplicationHandler // Application 处理器
	logger     *zap.Logger                 // 日志记录器
}

// NewRouterManager 创建路由管理器实例
// 返回 RouterManager 的实例
// 参数：appHandler - Application 处理器，logger - 日志记录器
func NewRouterManager(appHandler *handler.ApplicationHandler, logger *zap.Logger) *RouterManager {
	return &RouterManager{
		appHandler: appHandler,
		logger:     logger,
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
		SetupApplicationRoutes(api, rm.appHandler)

		// 未来可以在这里添加更多模块的路由
		// SetupUserRoutes(api, userHandler)
		// SetupAuthRoutes(api, authHandler)
	}

	return r
}
