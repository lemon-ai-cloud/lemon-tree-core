// Package middleware 提供 HTTP 中间件功能
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
// 记录 HTTP 请求的详细信息，包括请求方法、路径、状态码、延迟等
// 参数：logger - Zap 日志记录器
// 返回 Gin 中间件函数
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	// 使用 Gin 的自定义日志格式化器
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 使用结构化日志记录请求信息
		logger.Info("HTTP Request",
			zap.String("method", param.Method),                  // HTTP 方法（GET、POST 等）
			zap.String("path", param.Path),                      // 请求路径
			zap.Int("status", param.StatusCode),                 // HTTP 状态码
			zap.Duration("latency", param.Latency),              // 请求处理延迟
			zap.String("client_ip", param.ClientIP),             // 客户端 IP 地址
			zap.String("user_agent", param.Request.UserAgent()), // 用户代理（浏览器信息）
		)
		// 返回空字符串，因为我们使用 Zap 进行日志记录
		return ""
	})
}
