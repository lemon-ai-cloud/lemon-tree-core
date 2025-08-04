// Package middleware 提供 HTTP 中间件功能
// 包含各种 HTTP 请求处理中间件
package middleware

import (
	"context"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type myKey string

// AuthMiddleware 认证中间件
// 验证请求中的Token，确保用户已登录
// 相当于 Java Spring Boot 中的拦截器
// 返回 Gin 中间件函数
func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
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

// OptionalAuthMiddleware 可选认证中间件
// 验证Token但不强制要求，如果Token有效则设置用户信息
// 适用于一些可选认证的接口
// 返回 Gin 中间件函数
func OptionalAuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Token
		token := c.GetHeader("Authorization")
		if token == "" {
			// Token不存在，继续处理（不中断请求）
			c.Next()
			return
		}

		// 移除Bearer前缀（如果存在）
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// 验证Token并获取当前用户
		user, err := userService.GetUserByToken(c.Request.Context(), token)
		if err != nil {
			// Token无效，但不中断请求，继续处理
			c.Next()
			return
		}

		// 将用户信息存储到上下文中，供后续处理器使用
		c.Set(define.AppContextKeyCurrentUser, user)

		// 继续处理下一个中间件或路由处理器
		c.Next()
	}
}

// GetCurrentUser 从上下文中获取当前用户
// 辅助函数，用于在处理器中获取当前登录用户信息
// 参数：c - Gin上下文
// 返回：当前用户对象和是否存在
func GetCurrentUser(c *gin.Context) (*models.SystemUser, bool) {
	user, exists := c.Get(define.AppContextKeyCurrentUser)
	if !exists {
		return nil, false
	}

	if userObj, ok := user.(*models.SystemUser); ok {
		return userObj, true
	}

	return nil, false
}
