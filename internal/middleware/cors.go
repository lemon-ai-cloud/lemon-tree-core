// Package middleware 提供 HTTP 中间件功能
// 包含各种 HTTP 请求处理中间件
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS 中间件
// 处理跨域资源共享（Cross-Origin Resource Sharing）
// 允许来自不同域的前端应用访问 API
// 返回 Gin 中间件函数
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源（Origin）
		// "*" 表示允许所有域，生产环境建议设置具体的域名
		c.Header("Access-Control-Allow-Origin", "*")

		// 设置是否允许发送 Cookie
		c.Header("Access-Control-Allow-Credentials", "true")

		// 设置允许的请求头
		// 包含常用的 HTTP 请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		// 设置允许的 HTTP 方法
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理预检请求（Preflight Request）
		// OPTIONS 请求是浏览器在发送实际请求前的预检请求
		if c.Request.Method == "OPTIONS" {
			// 返回 204 状态码，表示预检请求成功
			c.AbortWithStatus(204)
			return
		}

		// 继续处理下一个中间件或路由处理器
		c.Next()
	}
}
