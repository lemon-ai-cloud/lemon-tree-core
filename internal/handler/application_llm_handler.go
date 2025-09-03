// Package handler 提供 HTTP 请求处理层功能
// 负责处理 HTTP 请求、参数验证、调用业务逻辑和返回响应
package handler

import (
	"lemon-tree-core/internal/converter"
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ApplicationLlmHandler ApplicationLLM 控制器
// 处理 ApplicationLLM 相关的所有 HTTP 请求
// 相当于 Java Spring Boot 中的 Controller
type ApplicationLlmHandler struct {
	applicationLlmService service.ApplicationLlmService // ApplicationLLM 业务逻辑层接口
}

// NewApplicationLlmHandler 创建 ApplicationLLM Handler 实例
// 返回 ApplicationLlmHandler 的实例
// 参数：applicationLlmService - ApplicationLLM 业务逻辑层接口
func NewApplicationLlmHandler(applicationLlmService service.ApplicationLlmService) *ApplicationLlmHandler {
	return &ApplicationLlmHandler{
		applicationLlmService: applicationLlmService,
	}
}

// SaveApplicationLlm 保存应用模型信息
// 处理 POST /api/v1/application-llms/save 请求
// 如果模型存在则更新，不存在则创建
func (h *ApplicationLlmHandler) SaveApplicationLlm(c *gin.Context) {
	// 绑定 JSON 请求体到 SaveApplicationLlmRequest 结构体
	var saveRequest dto.SaveApplicationLlmRequest
	if err := c.ShouldBindJSON(&saveRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为模型
	applicationLlm := converter.SaveApplicationLlmRequestToApplicationLlmModel(&saveRequest)

	// 调用业务逻辑层保存模型
	if err := h.applicationLlmService.SaveApplicationLlm(c.Request.Context(), applicationLlm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	applicationLlmDto := converter.ApplicationLlmModelToApplicationLlmDto(applicationLlm)
	c.JSON(http.StatusOK, gin.H{
		"application_llm": applicationLlmDto,
	})
}

// UpdateEnabledStatus 更新模型启用状态
// 处理 PUT /api/v1/application-llms/:id/enabled 请求
// 只更新 Enabled 字段
func (h *ApplicationLlmHandler) UpdateEnabledStatus(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 绑定 JSON 请求体到 UpdateEnabledStatusRequest 结构体
	var updateRequest dto.UpdateEnabledStatusRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证ID一致性
	if updateRequest.ID != idStr {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID in URL and request body do not match"})
		return
	}

	// 调用业务逻辑层更新启用状态
	if err := h.applicationLlmService.UpdateEnabledStatus(c.Request.Context(), id, updateRequest.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "Model enabled status updated successfully"})
}

// GetModelsByProviderID 根据提供商ID获取模型列表
// 处理 GET /api/v1/application-llms/provider/:providerId 请求
// 返回指定提供商下的所有模型
func (h *ApplicationLlmHandler) GetModelsByProviderID(c *gin.Context) {
	// 从 URL 参数中获取提供商 ID
	providerIDStr := c.Param("providerId")

	// 解析 UUID 格式的提供商 ID
	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider UUID format"})
		return
	}

	// 调用业务逻辑层获取指定提供商下的模型
	models, err := h.applicationLlmService.GetModelsByProviderID(c.Request.Context(), providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	modelDtos := converter.ApplicationLlmModelListToApplicationLlmDtoList(models)
	c.JSON(http.StatusOK, gin.H{
		"application_llm": modelDtos,
	})
}

// GetModelsByApplicationID 根据应用ID获取模型列表
// 处理 GET /api/v1/application-llms/application/:applicationId 请求
// 返回指定应用下的所有模型
func (h *ApplicationLlmHandler) GetModelsByApplicationID(c *gin.Context) {
	// 从 URL 参数中获取应用 ID
	applicationIDStr := c.Param("applicationId")

	// 解析 UUID 格式的应用 ID
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application UUID format"})
		return
	}

	// 调用业务逻辑层获取指定应用下的模型
	models, err := h.applicationLlmService.GetModelsByApplicationID(c.Request.Context(), applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	modelDtos := converter.ApplicationLlmModelListToApplicationLlmDtoList(models)
	c.JSON(http.StatusOK, gin.H{
		"application_llm": modelDtos,
	})
}
