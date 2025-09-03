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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LlmProviderHandler LlmProvider 控制器
// 处理 LlmProvider 相关的所有 HTTP 请求
// 相当于 Java Spring Boot 中的 Controller
type LlmProviderHandler struct {
	llmProviderService service.LlmProviderService // LlmProvider 业务逻辑层接口
}

// NewLlmProviderHandler 创建 LlmProvider Handler 实例
// 返回 LlmProviderHandler 的实例
// 参数：llmProviderService - LlmProvider 业务逻辑层接口
func NewLlmProviderHandler(llmProviderService service.LlmProviderService) *LlmProviderHandler {
	return &LlmProviderHandler{
		llmProviderService: llmProviderService,
	}
}

// GetLlmProviderByID 根据ID获取大语言模型提供商
// 处理 GET /api/v1/llm-providers/:id 请求
// 根据 UUID 获取指定的提供商信息
func (h *LlmProviderHandler) GetLlmProviderByID(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层获取提供商
	llmProvider, err := h.llmProviderService.GetLlmProviderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "LlmProvider not found"})
		return
	}

	// 转换为DTO返回
	llmProviderDto := converter.LlmProviderModelToLlmProviderDto(llmProvider)
	c.JSON(http.StatusOK, gin.H{
		"llm_provider": llmProviderDto,
	})
}

// GetAllLlmProviders 获取所有大语言模型提供商
// 处理 GET /api/v1/llm-providers 请求
// 获取所有提供商的列表
func (h *LlmProviderHandler) GetAllLlmProviders(c *gin.Context) {
	// 调用业务逻辑层获取所有提供商
	llmProviders, err := h.llmProviderService.GetAllLlmProviders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	llmProviderDtos := converter.LlmProviderModelListToLlmProviderDtoList(llmProviders)
	c.JSON(http.StatusOK, gin.H{
		"llm_providers": llmProviderDtos,
	})
}

// SaveLlmProvider 保存大语言模型提供商（upsert）
// 处理 POST /api/v1/llm-providers/save 请求
// 如果提供商存在则更新，不存在则创建
func (h *LlmProviderHandler) SaveLlmProvider(c *gin.Context) {
	// 绑定 JSON 请求体到 LlmProviderSaveDto 结构体
	var llmProviderSaveDto dto.LlmProviderSaveDto
	if err := c.ShouldBindJSON(&llmProviderSaveDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为模型
	llmProvider := converter.LlmProviderSaveDtoToLlmProviderModel(&llmProviderSaveDto)

	// 调用业务逻辑层保存提供商
	if err := h.llmProviderService.SaveLlmProvider(c.Request.Context(), llmProvider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	llmProviderDto := converter.LlmProviderModelToLlmProviderDto(llmProvider)
	c.JSON(http.StatusOK, gin.H{
		"llm_provider": llmProviderDto,
	})
}

// QueryLlmProviders 动态查询大语言模型提供商
// 处理 POST /api/v1/llm-providers/query 请求
// 根据查询条件动态查询提供商
func (h *LlmProviderHandler) QueryLlmProviders(c *gin.Context) {
	// 绑定 JSON 请求体到 LlmProviderQueryDto 结构体作为查询条件
	var queryDto dto.LlmProviderQueryDto
	if err := c.ShouldBindJSON(&queryDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为模型
	query := converter.LlmProviderQueryDtoToLlmProviderModel(&queryDto)

	// 调用业务逻辑层查询提供商
	llmProviders, err := h.llmProviderService.QueryLlmProviders(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	llmProviderDtos := converter.LlmProviderModelListToLlmProviderDtoList(llmProviders)
	c.JSON(http.StatusOK, gin.H{
		"llm_providers": llmProviderDtos,
	})
}

// DeleteLlmProvider 删除大语言模型提供商
// 处理 DELETE /api/v1/llm-providers/:id 请求
// 删除指定的提供商（软删除）
func (h *LlmProviderHandler) DeleteLlmProvider(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层删除提供商
	if err := h.llmProviderService.DeleteLlmProvider(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回删除成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "LlmProvider deleted successfully"})
}

// GetLlmProvidersByApplicationID 根据应用ID获取大语言模型提供商列表
// 处理 GET /api/v1/llm-providers/application/:applicationId 请求
// 返回指定应用下的所有提供商
func (h *LlmProviderHandler) GetLlmProvidersByApplicationID(c *gin.Context) {
	// 从 URL 参数中获取应用 ID
	applicationIDStr := c.Param("applicationId")

	// 解析 UUID 格式的应用 ID
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application UUID format"})
		return
	}

	// 调用业务逻辑层获取指定应用下的提供商
	llmProviders, err := h.llmProviderService.GetLlmProvidersByApplicationID(c.Request.Context(), applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	llmProviderDtos := converter.LlmProviderModelListToLlmProviderDtoList(llmProviders)
	c.JSON(http.StatusOK, gin.H{
		"llm_providers": llmProviderDtos,
	})
}

// UploadLlmProviderIcon 上传大语言模型提供商图标
// 处理 POST /api/v1/llm-providers/upload-icon 请求
// 上传图标文件并返回可用的 URL
func (h *LlmProviderHandler) UploadLlmProviderIcon(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("icon")
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
	saveDir := filepath.Join(workspacePath, define.WorkspaceDirNameLlmProviderIcon)

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
		"message": "图片上传成功",
		"data": gin.H{
			"file_name": fileName,
			"file_path": define.WorkspaceDirNameLlmProviderIcon + fileName,
			"file_size": file.Size,
			"mime_type": contentType,
		},
	})
}
