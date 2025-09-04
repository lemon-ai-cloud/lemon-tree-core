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

// ApplicationMcpServerConfigHandler ApplicationMCP配置 控制器
// 处理 ApplicationMCP配置 相关的所有 HTTP 请求
// 相当于 Java Spring Boot 中的 Controller
type ApplicationMcpServerConfigHandler struct {
	applicationMcpServerConfigService service.ApplicationMcpServerConfigService // ApplicationMCP配置 业务逻辑层接口
}

// NewApplicationMcpServerConfigHandler 创建 ApplicationMCP配置 Handler 实例
// 返回 ApplicationMcpServerConfigHandler 的实例
// 参数：applicationMcpServerConfigService - ApplicationMCP配置 业务逻辑层接口
func NewApplicationMcpServerConfigHandler(applicationMcpServerConfigService service.ApplicationMcpServerConfigService) *ApplicationMcpServerConfigHandler {
	return &ApplicationMcpServerConfigHandler{
		applicationMcpServerConfigService: applicationMcpServerConfigService,
	}
}

// SaveApplicationMcpServerConfig 保存应用MCP配置信息
// 处理 POST /api/v1/application-mcp-server-configs/save 请求
// 如果配置存在则更新，不存在则创建
func (h *ApplicationMcpServerConfigHandler) SaveApplicationMcpServerConfig(c *gin.Context) {
	// 绑定 JSON 请求体到 SaveApplicationMcpServerConfigRequest 结构体
	var saveRequest dto.SaveApplicationMcpServerConfigRequest
	if err := c.ShouldBindJSON(&saveRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为模型
	config := converter.SaveApplicationMcpServerConfigRequestToApplicationMcpServerConfigModel(&saveRequest)

	// 调用业务逻辑层保存配置
	if err := h.applicationMcpServerConfigService.SaveApplicationMcpServerConfig(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO返回
	configDto := converter.ApplicationMcpServerConfigModelToApplicationMcpServerConfigDto(config)
	c.JSON(http.StatusOK, gin.H{
		"application_mcp_server_config": configDto,
	})
}

// DeleteApplicationMcpServerConfig 删除MCP配置
// 处理 DELETE /api/v1/application-mcp-server-configs/:id 请求
// 根据ID删除指定的MCP配置记录
func (h *ApplicationMcpServerConfigHandler) DeleteApplicationMcpServerConfig(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层删除配置
	if err := h.applicationMcpServerConfigService.DeleteApplicationMcpServerConfig(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回删除成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "MCP配置删除成功"})
}

// GetMcpServerConfigsByApplicationID 根据应用ID获取MCP配置列表
// 处理 GET /api/v1/application-mcp-server-configs/application/:applicationId 请求
// 返回指定应用下的所有MCP配置
func (h *ApplicationMcpServerConfigHandler) GetMcpServerConfigsByApplicationID(c *gin.Context) {
	// 从 URL 参数中获取应用 ID
	applicationIDStr := c.Param("applicationId")

	// 解析 UUID 格式的应用 ID
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application UUID format"})
		return
	}

	// 调用业务逻辑层获取指定应用下的MCP配置
	configs, err := h.applicationMcpServerConfigService.GetMcpServerConfigsByApplicationID(c.Request.Context(), applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	configDtos := converter.ApplicationMcpServerConfigModelListToApplicationMcpServerConfigDtoList(configs)
	c.JSON(http.StatusOK, gin.H{
		"application_mcp_server_configs": configDtos,
	})
}

// GetMcpServerTools 获取MCP服务器的所有工具
// 处理 GET /api/v1/application-mcp-server-configs/:id/tools 请求
// 根据MCP配置ID连接服务器并返回可用工具列表
func (h *ApplicationMcpServerConfigHandler) GetMcpServerTools(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层获取工具列表
	tools, err := h.applicationMcpServerConfigService.GetMcpServerTools(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表返回
	toolDtos := converter.ApplicationMcpServerToolModelListToApplicationMcpServerToolDtoList(tools)
	c.JSON(http.StatusOK, gin.H{
		"tools": toolDtos,
	})
}
