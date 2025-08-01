// Package router 提供路由管理功能
package router

import (
	"lemon-tree-core/internal/handler"
	"lemon-tree-core/internal/middleware"
	"lemon-tree-core/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户相关路由
// 配置用户模块的所有 HTTP 路由
// 参数：api - API 路由组，userHandler - 用户处理器，userService - 用户服务
func SetupUserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, userService service.UserService) {
	// 用户路由组
	// 所有用户相关的路由都以 /users 为前缀
	users := api.Group("/users")
	{
		// 用户登录（无需认证）
		// POST /api/v1/users/login
		// 用户登录，验证账号密码，返回Token
		users.POST("/login", userHandler.Login)

		// 需要认证的路由组
		authenticated := users.Group("")
		authenticated.Use(middleware.AuthMiddleware(userService))
		{
			// 保存用户（创建或更新）
			// POST /api/v1/users/save
			// 保存用户信息，如果存在则更新，不存在则创建
			authenticated.POST("/save", userHandler.SaveUser)

			// 获取所有用户
			// GET /api/v1/users
			// 获取所有用户的列表
			authenticated.GET("", userHandler.GetAllUsers)

			// 根据ID获取用户详情
			// GET /api/v1/users/:id
			// 根据 ID 获取指定的用户信息
			authenticated.GET("/:id", userHandler.GetUserByID)

			// 获取当前登录用户信息
			// GET /api/v1/users/current
			// 根据Token获取当前登录用户信息
			authenticated.GET("/current", userHandler.GetCurrentUser)

			// 用户登出
			// POST /api/v1/users/logout
			// 用户登出，删除会话记录
			authenticated.POST("/logout", userHandler.Logout)
		}
	}
}
