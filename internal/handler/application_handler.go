// Package handler 提供 HTTP 请求处理层功能
// 负责处理 HTTP 请求、参数验证、调用业务逻辑和返回响应
package handler

import (
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ApplicationHandler Application 控制器
// 处理 Application 相关的所有 HTTP 请求
// 相当于 Java Spring Boot 中的 Controller
type ApplicationHandler struct {
	appService service.ApplicationService // Application 业务逻辑层接口
}

// NewApplicationHandler 创建 Application Handler 实例
// 返回 ApplicationHandler 的实例
// 参数：appService - Application 业务逻辑层接口
func NewApplicationHandler(appService service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		appService: appService,
	}
}

// GetApplicationByID 根据ID获取应用
// 处理 GET /api/v1/applications/:id 请求
// 根据 UUID 获取指定的应用信息
func (h *ApplicationHandler) GetApplicationByID(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层获取应用
	application, err := h.appService.GetApplicationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// 返回应用信息
	c.JSON(http.StatusOK, application)
}

// GetAllApplications 获取所有应用
// 处理 GET /api/v1/applications 请求
// 获取所有应用的列表
func (h *ApplicationHandler) GetAllApplications(c *gin.Context) {
	// 调用业务逻辑层获取所有应用
	applications, err := h.appService.GetAllApplications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回应用列表
	c.JSON(http.StatusOK, applications)
}

// SaveApplication 保存应用（upsert）
// 处理 POST /api/v1/applications/save 请求
// 如果应用存在则更新，不存在则创建
func (h *ApplicationHandler) SaveApplication(c *gin.Context) {
	// 绑定 JSON 请求体到 Application 结构体
	var application models.Application
	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用业务逻辑层保存应用
	if err := h.appService.SaveApplication(c.Request.Context(), &application); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回保存成功的响应
	c.JSON(http.StatusOK, application)
}

// QueryApplications 动态查询应用
// 处理 POST /api/v1/applications/query 请求
// 根据查询条件动态查询应用
func (h *ApplicationHandler) QueryApplications(c *gin.Context) {
	// 绑定 JSON 请求体到 Application 结构体作为查询条件
	var query models.Application
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用业务逻辑层查询应用
	applications, err := h.appService.QueryApplications(c.Request.Context(), &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, applications)
}

// DeleteApplication 删除应用
// 处理 DELETE /api/v1/applications/:id 请求
// 删除指定的应用（软删除）
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	// 从 URL 参数中获取 ID
	idStr := c.Param("id")

	// 解析 UUID 格式的 ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// 调用业务逻辑层删除应用
	if err := h.appService.DeleteApplication(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回删除成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}
