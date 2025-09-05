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

// ApplicationStorageConfigHandler 应用存储配置 控制器
// 处理 应用存储配置 相关的所有 HTTP 请求
// 相当于 Java Spring Boot 中的 Controller
type ApplicationStorageConfigHandler struct {
	applicationStorageConfigService service.ApplicationStorageConfigService // 应用存储配置 业务逻辑层接口
}

// NewApplicationStorageConfigHandler 创建 应用存储配置 Handler 实例
// 返回 ApplicationStorageConfigHandler 的实例
// 参数：applicationStorageConfigService - 应用存储配置 业务逻辑层接口
func NewApplicationStorageConfigHandler(applicationStorageConfigService service.ApplicationStorageConfigService) *ApplicationStorageConfigHandler {
	return &ApplicationStorageConfigHandler{
		applicationStorageConfigService: applicationStorageConfigService,
	}
}

// SaveApplicationStorageConfig 保存应用存储配置
// 处理 POST /api/v1/application-storage-configs/save 请求
// 根据ApplicationID保存配置，如果存在则覆盖，不存在则创建
func (h *ApplicationStorageConfigHandler) SaveApplicationStorageConfig(c *gin.Context) {
	// 绑定 JSON 请求体到 SaveApplicationStorageConfigRequest 结构体
	var saveRequest dto.SaveApplicationStorageConfigRequest
	if err := c.ShouldBindJSON(&saveRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为模型
	config := converter.SaveApplicationStorageConfigRequestToApplicationStorageConfigModel(&saveRequest)

	// 调用业务逻辑层保存存储配置
	if err := h.applicationStorageConfigService.SaveApplicationStorageConfig(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	configDto := converter.ApplicationStorageConfigModelToApplicationStorageConfigDto(config)
	c.JSON(http.StatusOK, gin.H{
		"application_storage_config": configDto,
	})
}

// GetApplicationStorageConfigByApplicationID 根据应用ID获取存储配置
// 处理 GET /api/v1/application-storage-configs/application/:applicationId 请求
// 返回指定应用的存储配置
func (h *ApplicationStorageConfigHandler) GetApplicationStorageConfigByApplicationID(c *gin.Context) {
	// 从 URL 参数中获取应用 ID
	applicationIDStr := c.Param("applicationId")

	// 解析 UUID 格式的应用 ID
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application UUID format"})
		return
	}

	// 调用业务逻辑层获取指定应用的存储配置
	config, err := h.applicationStorageConfigService.GetApplicationStorageConfigByApplicationID(c.Request.Context(), applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果配置不存在，返回空对象
	if config == nil {
		c.JSON(http.StatusOK, gin.H{
			"application_storage_config": nil,
		})
		return
	}

	// 转换为DTO返回
	configDto := converter.ApplicationStorageConfigModelToApplicationStorageConfigDto(config)
	c.JSON(http.StatusOK, gin.H{
		"application_storage_config": configDto,
	})
}
