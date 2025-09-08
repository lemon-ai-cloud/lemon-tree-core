// Package middleware 提供 HTTP 中间件功能
// 包含各种 HTTP 请求处理中间件
package middleware

import (
	"context"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type myKey string

// UserAuthMiddleware 认证中间件
// 验证请求中的Token，确保用户已登录
// 返回 Gin 中间件函数
func UserAuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证Token"})
			c.Abort()
			return
		}

		// 移除Bearer前缀（如果存在）
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// 验证Token并获取当前用户
		user, err := userService.GetUserByToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，供后续处理器使用
		c.Set(define.AppContextKeyCurrentUser, user)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, define.AppContextKeyCurrentUser, user)
		c.Request = c.Request.WithContext(ctx)
		// 继续处理下一个中间件或路由处理器
		c.Next()
	}
}

// ChatAgentAuthMiddleware 认证中间件
// 验证请求中的ApiKey，确认是哪个应用
// 返回 Gin 中间件函数
func ChatAgentAuthMiddleware(chatAgentService service.ChatAgentService, applicationService service.ApplicationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Token
		apiKey := c.GetHeader("lemon-ai-api-key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Lemon AI ApiKey Not Found"})
			c.Abort()
			return
		}

		// 验证Api Key获取ChatAgent
		chatAgent, err := chatAgentService.GetChatAgentByApiKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		application, getAppErr := applicationService.GetApplicationByID(c.Request.Context(), chatAgent.ApplicationID)

		if getAppErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": getAppErr.Error()})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，供后续处理器使用
		c.Set(define.AppContextKeyCurrentChatAgent, chatAgent)
		c.Set(define.AppContextKeyCurrentApplication, application)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, define.AppContextKeyCurrentChatAgent, chatAgent)
		ctx = context.WithValue(ctx, define.AppContextKeyCurrentApplication, application)
		c.Request = c.Request.WithContext(ctx)
		// 继续处理下一个中间件或路由处理器
		c.Next()
	}
}
