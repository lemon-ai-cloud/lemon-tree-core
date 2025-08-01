// Package middleware 提供 HTTP 中间件功能
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware 恢复中间件
// 捕获和处理 panic，防止程序崩溃
// 记录错误信息并返回友好的错误响应
// 参数：logger - Zap 日志记录器
// 返回 Gin 中间件函数
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	// 使用 Gin 的自定义恢复函数
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 尝试将 recovered 转换为字符串
		if err, ok := recovered.(string); ok {
			// 记录 panic 错误信息
			logger.Error("Panic recovered",
				zap.String("error", err),               // 错误信息
				zap.String("path", c.Request.URL.Path), // 请求路径
				zap.String("method", c.Request.Method), // HTTP 方法
			)
		}

		// 返回友好的错误响应
		// 避免向客户端暴露敏感的错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	})
}
