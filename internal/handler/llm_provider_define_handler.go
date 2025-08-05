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

// LlmProviderDefineHandler LlmProviderDefine 控制器
// 处理 LlmProviderDefine 相关的所有 HTTP 请求
// 相当于 Java Spring Boot 中的 Controller
type LlmProviderDefineHandler struct {
	llmProviderDefineService service.LlmProviderDefineService // LlmProviderDefine 业务逻辑层接口
}

// NewLlmProviderDefineHandler 创建 LlmProviderDefine Handler 实例
// 返回 LlmProviderDefineHandler 的实例
// 参数：llmProviderDefineService - LlmProviderDefine 业务逻辑层接口
func NewLlmProviderDefineHandler(llmProviderDefineService service.LlmProviderDefineService) *LlmProviderDefineHandler {
	return &LlmProviderDefineHandler{
		llmProviderDefineService: llmProviderDefineService,
	}
}

// GetLlmProviderDefineByID 根据ID获取大语言模型提供商定义
// 处理 GET /api/v1/llm-provider-defines/:id 请求
// 根据 UUID 获取指定的提供商定义信息
func (h *LlmProviderDefineHandler) GetLlmProviderDefineByID(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层获取提供商定义
	llmProviderDefine, err := h.llmProviderDefineService.GetLlmProviderDefineByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "LlmProviderDefine not found"})
		return
	}

	// 转换为DTO返回
	llmProviderDefineDto := converter.LlmProviderDefineModelToLlmProviderDefineDto(llmProviderDefine)
	c.JSON(http.StatusOK, gin.H{
		"llm_provider_define": llmProviderDefineDto,
	})
}

// GetAllLlmProviderDefines 获取所有大语言模型提供商定义
// 处理 GET /api/v1/llm-provider-defines 请求
// 获取所有提供商定义的列表
func (h *LlmProviderDefineHandler) GetAllLlmProviderDefines(c *gin.Context) {
	// 调用业务逻辑层获取所有提供商定义
	llmProviderDefines, err := h.llmProviderDefineService.GetAllLlmProviderDefines(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	llmProviderDefineDtos := converter.LlmProviderDefineModelListToLlmProviderDefineDtoList(llmProviderDefines)
	c.JSON(http.StatusOK, gin.H{
		"llm_provider_defines": llmProviderDefineDtos,
	})
}

// SaveLlmProviderDefine 保存大语言模型提供商定义（upsert）
// 处理 POST /api/v1/llm-provider-defines/save 请求
// 如果提供商定义存在则更新，不存在则创建
func (h *LlmProviderDefineHandler) SaveLlmProviderDefine(c *gin.Context) {
	// 绑定 JSON 请求体到 LlmProviderDefineSaveDto 结构体
	var llmProviderDefineSaveDto dto.LlmProviderDefineSaveDto
	if err := c.ShouldBindJSON(&llmProviderDefineSaveDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为模型
	llmProviderDefine := converter.LlmProviderDefineSaveDtoToLlmProviderDefineModel(&llmProviderDefineSaveDto)

	// 调用业务逻辑层保存提供商定义
	if err := h.llmProviderDefineService.SaveLlmProviderDefine(c.Request.Context(), llmProviderDefine); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	llmProviderDefineDto := converter.LlmProviderDefineModelToLlmProviderDefineDto(llmProviderDefine)
	c.JSON(http.StatusOK, gin.H{
		"llm_provider_define": llmProviderDefineDto,
	})
}

// QueryLlmProviderDefines 动态查询大语言模型提供商定义
// 处理 POST /api/v1/llm-provider-defines/query 请求
// 根据查询条件动态查询提供商定义
func (h *LlmProviderDefineHandler) QueryLlmProviderDefines(c *gin.Context) {
	// 绑定 JSON 请求体到 LlmProviderDefineQueryDto 结构体作为查询条件
	var queryDto dto.LlmProviderDefineQueryDto
	if err := c.ShouldBindJSON(&queryDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为模型
	query := converter.LlmProviderDefineQueryDtoToLlmProviderDefineModel(&queryDto)

	// 调用业务逻辑层查询提供商定义
	llmProviderDefines, err := h.llmProviderDefineService.QueryLlmProviderDefines(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	llmProviderDefineDtos := converter.LlmProviderDefineModelListToLlmProviderDefineDtoList(llmProviderDefines)
	c.JSON(http.StatusOK, gin.H{
		"llm_provider_defines": llmProviderDefineDtos,
	})
}

// DeleteLlmProviderDefine 删除大语言模型提供商定义
// 处理 DELETE /api/v1/llm-provider-defines/:id 请求
// 删除指定的提供商定义（软删除）
func (h *LlmProviderDefineHandler) DeleteLlmProviderDefine(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层删除提供商定义
	if err := h.llmProviderDefineService.DeleteLlmProviderDefine(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回删除成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "LlmProviderDefine deleted successfully"})
}

// UploadLlmProviderDefineIcon 上传大语言模型提供商定义图标
// 处理 POST /api/v1/llm-provider-defines/upload-icon 请求
// 上传图标文件并返回可用的 URL
func (h *LlmProviderDefineHandler) UploadLlmProviderDefineIcon(c *gin.Context) {
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
	saveDir := filepath.Join(workspacePath, define.WorkspaceDirNameLlmProviderDefineIcon)

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
			"file_path": define.WorkspaceDirNameLlmProviderDefineIcon + fileName,
			"file_size": file.Size,
			"mime_type": contentType,
		},
	})
}
