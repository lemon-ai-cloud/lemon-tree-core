// Package handler 提供 HTTP 请求处理层功能
// 负责处理 HTTP 请求、参数验证、调用业务逻辑和返回响应
package handler

import (
	"lemon-tree-core/internal/converter"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ChatAgentHandler 智能体 控制器
// 处理 智能体 相关的所有 HTTP 请求
// 相当于 Java Spring Boot 中的 Controller
type ChatAgentHandler struct {
	chatAgentService service.ChatAgentService // 智能体 业务逻辑层接口
}

// NewChatAgentHandler 创建 智能体 Handler 实例
// 返回 ChatAgentHandler 的实例
// 参数：chatAgentService - 智能体 业务逻辑层接口
func NewChatAgentHandler(chatAgentService service.ChatAgentService) *ChatAgentHandler {
	return &ChatAgentHandler{
		chatAgentService: chatAgentService,
	}
}

// SaveChatAgent 保存智能体信息
// 处理 POST /api/v1/chat-agents/save 请求
// 如果智能体存在则更新，不存在则创建
func (h *ChatAgentHandler) SaveChatAgent(c *gin.Context) {
	// 绑定 JSON 请求体到 SaveChatAgentRequest 结构体
	var saveRequest dto.SaveChatAgentRequest
	if err := c.ShouldBindJSON(&saveRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为模型
	agent := converter.SaveChatAgentRequestToChatAgentModel(&saveRequest)

	// 调用业务逻辑层保存智能体
	if err := h.chatAgentService.SaveChatAgent(c.Request.Context(), agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	agentDto := converter.ChatAgentModelToChatAgentDto(agent)
	c.JSON(http.StatusOK, gin.H{
		"chat_agent": agentDto,
	})
}

// DeleteChatAgent 删除智能体
// 处理 DELETE /api/v1/chat-agents/:id 请求
// 根据ID删除指定的智能体记录
func (h *ChatAgentHandler) DeleteChatAgent(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层删除智能体
	if err := h.chatAgentService.DeleteChatAgent(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回删除成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "智能体删除成功"})
}

// GetChatAgentsByApplicationID 根据应用ID获取智能体列表
// 处理 GET /api/v1/chat-agents/application/:applicationId 请求
// 返回指定应用下的所有智能体，支持分页
func (h *ChatAgentHandler) GetChatAgentsByApplicationID(c *gin.Context) {
	// 从 URL 参数中获取应用 ID
	applicationIDStr := c.Param("applicationId")

	// 解析 UUID 格式的应用 ID
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application UUID format"})
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用业务逻辑层获取指定应用下的智能体
	agents, total, err := h.chatAgentService.GetChatAgentsByApplicationID(c.Request.Context(), applicationID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	agentDtos := converter.ChatAgentModelListToChatAgentDtoList(agents)
	c.JSON(http.StatusOK, gin.H{
		"chat_agents": agentDtos,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

// UploadChatAgentAvatar 上传智能体头像
// 处理 POST /api/v1/chat-agents/upload-avatar 请求
// 上传头像文件并返回可用的 URL
func (h *ChatAgentHandler) UploadChatAgentAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的图片文件"})
		return
	}

	// 验证文件类型
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持图片文件上传"})
		return
	}

	// 验证文件大小（限制为 5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片文件大小不能超过 5MB"})
		return
	}

	// 生成临时文件名
	tempID := uuid.New()

	// 确定文件扩展名
	var ext string
	switch contentType {
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的图片格式"})
		return
	}

	// 获取工作区路径
	workspacePath := os.Getenv("WORKSPACE_PUBLIC_PATH")
	if workspacePath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "环境变量 WORKSPACE_PUBLIC_PATH 未设置"})
		return
	}

	// 构建保存路径
	saveDir := filepath.Join(workspacePath, define.WorkspaceDirNameChatAgentAvatar)

	// 确保目录存在
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	// 构建文件路径
	fileName := tempID.String() + ext
	filePath := filepath.Join(saveDir, fileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 返回文件信息
	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"data": gin.H{
			"file_name": fileName,
			"file_path": define.WorkspaceDirNameChatAgentAvatar + fileName,
			"file_size": file.Size,
			"mime_type": contentType,
		},
	})
}
