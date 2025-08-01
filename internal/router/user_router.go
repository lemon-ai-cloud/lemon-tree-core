// Package router 提供路由管理功能
package router

import (
	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器接口（示例）
// 定义了用户相关的所有 HTTP 请求处理方法
// 这是一个示例接口，用于演示如何添加新的模块路由
type UserHandler interface {
	CreateUser(c *gin.Context)  // 创建用户
	GetUserByID(c *gin.Context) // 根据ID获取用户
	GetAllUsers(c *gin.Context) // 获取所有用户
	UpdateUser(c *gin.Context)  // 更新用户
	DeleteUser(c *gin.Context)  // 删除用户
}

// SetupUserRoutes 设置用户相关路由（示例）
// 配置用户模块的所有 HTTP 路由
// 这是一个示例函数，演示如何添加新的模块路由
// 参数：api - API 路由组，userHandler - 用户处理器
func SetupUserRoutes(api *gin.RouterGroup, userHandler UserHandler) {
	// 用户路由组
	// 所有用户相关的路由都以 /users 为前缀
	users := api.Group("/users")
	{
		// 创建用户
		// POST /api/v1/users
		// 创建新的用户
		users.POST("", userHandler.CreateUser)

		// 获取用户列表
		// GET /api/v1/users
		// 获取所有用户的列表
		users.GET("", userHandler.GetAllUsers)

		// 根据ID获取用户
		// GET /api/v1/users/:id
		// 根据 ID 获取指定的用户信息
		users.GET("/:id", userHandler.GetUserByID)

		// 更新用户
		// PUT /api/v1/users/:id
		// 更新指定用户的信息
		users.PUT("/:id", userHandler.UpdateUser)

		// 删除用户
		// DELETE /api/v1/users/:id
		// 删除指定的用户
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
