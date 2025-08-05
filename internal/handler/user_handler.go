// Package handler 提供 HTTP 请求处理层功能
// 负责处理 HTTP 请求、参数验证、调用业务逻辑和返回响应
package handler

import (
	"lemon-tree-core/internal/converter"
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/service"
	"lemon-tree-core/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler User 控制器
// 处理 User 相关的所有 HTTP 请求
type UserHandler struct {
	userService service.UserService // User 业务逻辑层接口
}

// NewUserHandler 创建 User Handler 实例
// 返回 UserHandler 的实例
// 参数：userService - User 业务逻辑层接口
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login 用户登录
// 处理 POST /api/v1/users/login 请求
// 验证用户账号密码，创建会话并返回Token
func (h *UserHandler) Login(c *gin.Context) {
	// 绑定登录请求参数
	var loginRequest dto.SystemUserLoginDto

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 调用业务逻辑层进行登录
	user, token, err := h.userService.Login(c.Request.Context(), loginRequest.Number, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	userDto := converter.SystemUserModelToSystemUserDto(user)
	utils.JsonResponse(c, http.StatusOK, gin.H{"user": userDto, "token": token})
}

// SaveUser 保存用户（创建或更新）
// 处理 POST /api/v1/users/save 请求
// 如果用户存在则更新，不存在则创建
func (h *UserHandler) SaveUser(c *gin.Context) {
	// 绑定用户信息
	var userSaveDto dto.SystemUserSaveDto
	if err := c.ShouldBindJSON(&userSaveDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 转换为模型
	user := converter.SystemUserSaveDtoToSystemUserModel(&userSaveDto)

	// 调用业务逻辑层保存用户
	if err := h.userService.SaveUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	userDto := converter.SystemUserModelToSystemUserDto(user)
	c.JSON(http.StatusOK, gin.H{
		"user": userDto,
	})
}

// GetAllUsers 获取所有用户
// 处理 GET /api/v1/users 请求
// 获取所有用户的列表
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// 调用业务逻辑层获取所有用户
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	userDtos := converter.SystemUserModelListToSystemUserDtoList(users)
	c.JSON(http.StatusOK, gin.H{
		"users": userDtos,
	})
}

// GetUserByID 根据ID获取用户详情
// 处理 GET /api/v1/users/:id 请求
// 根据 UUID 获取指定的用户信息
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	// 调用业务逻辑层获取用户
	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 转换为DTO返回
	userDto := converter.SystemUserModelToSystemUserDto(user)
	c.JSON(http.StatusOK, gin.H{
		"user": userDto,
	})
}

// GetCurrentUser 获取当前登录用户信息
// 处理 GET /api/v1/users/current 请求
// 根据Token获取当前登录用户信息
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// 调用业务逻辑层获取当前用户
	user, err := h.userService.GetCurrentUser(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	userDto := converter.SystemUserModelToSystemUserDto(user)
	c.JSON(http.StatusOK, gin.H{
		"user": userDto,
	})
}

// Logout 用户登出
// 处理 POST /api/v1/users/logout 请求
// 删除用户的会话记录
func (h *UserHandler) Logout(c *gin.Context) {
	// 从请求头中获取Token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证Token"})
		return
	}

	// 移除Bearer前缀（如果存在）
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 调用业务逻辑层登出
	if err := h.userService.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回登出成功的响应
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// DeleteUser 删除用户
// 处理 DELETE /api/v1/users/:id 请求
// 删除指定用户及其所有会话记录
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	// 调用业务逻辑层删除用户
	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回删除成功的响应
	c.JSON(http.StatusOK, gin.H{
		"message": "用户删除成功",
	})
}
