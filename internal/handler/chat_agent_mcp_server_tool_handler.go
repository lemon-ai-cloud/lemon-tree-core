// Package handler 提供HTTP请求处理功能
// 负责接收HTTP请求、参数验证、调用业务逻辑层和返回HTTP响应
package handler

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ChatAgentMcpServerToolHandler ChatAgentMcpServerTool HTTP处理器
// 负责处理与ChatAgentMcpServerTool相关的HTTP请求
type ChatAgentMcpServerToolHandler struct {
	chatAgentMcpServerToolService service.ChatAgentMcpServerToolService
}

// NewChatAgentMcpServerToolHandler 创建 ChatAgentMcpServerTool HTTP处理器实例
// 参数：chatAgentMcpServerToolService - ChatAgentMcpServerTool业务逻辑层服务
func NewChatAgentMcpServerToolHandler(chatAgentMcpServerToolService service.ChatAgentMcpServerToolService) *ChatAgentMcpServerToolHandler {
	return &ChatAgentMcpServerToolHandler{
		chatAgentMcpServerToolService: chatAgentMcpServerToolService,
	}
}

// SaveChatAgentMcpServerToolSettings 保存聊天智能体的MCP工具设置
// @Summary 保存聊天智能体的MCP工具设置
// @Description 保存指定聊天智能体的MCP工具启用/禁用设置
// @Tags ChatAgentMcpServerTool
// @Accept json
// @Produce json
// @Param request body dto.SaveChatAgentMcpServerToolSettingsRequest true "保存工具设置请求"
// @Success 200 {object} dto.SaveChatAgentMcpServerToolSettingsResponse "保存成功"
// @Failure 400 {object} dto.SaveChatAgentMcpServerToolSettingsResponse "请求参数错误"
// @Failure 500 {object} dto.SaveChatAgentMcpServerToolSettingsResponse "服务器内部错误"
// @Router /api/v1/chat-agent-mcp-server-tool/settings [post]
func (h *ChatAgentMcpServerToolHandler) SaveChatAgentMcpServerToolSettings(c *gin.Context) {
	// 从URL路径参数获取chatAgentID
	chatAgentIDStr := c.Param("chatAgentID")
	chatAgentID, err := uuid.Parse(chatAgentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.SaveChatAgentMcpServerToolSettingsResponse{
			Success: false,
			Message: "无效的聊天智能体ID",
		})
		return
	}

	// 从请求体获取工具设置
	var req dto.SaveChatAgentMcpServerToolSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.SaveChatAgentMcpServerToolSettingsResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用业务逻辑层
	err = h.chatAgentMcpServerToolService.SaveChatAgentMcpServerToolSettings(c.Request.Context(), chatAgentID, req.ToolSettings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.SaveChatAgentMcpServerToolSettingsResponse{
			Success: false,
			Message: "保存工具设置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SaveChatAgentMcpServerToolSettingsResponse{
		Success: true,
		Message: "保存成功",
	})
}

// GetChatAgentMcpServerToolSettings 获取聊天智能体的MCP工具设置
// @Summary 获取聊天智能体的MCP工具设置
// @Description 获取指定聊天智能体的MCP工具启用/禁用设置
// @Tags ChatAgentMcpServerTool
// @Accept json
// @Produce json
// @Param chat_agent_id path string true "聊天智能体ID"
// @Success 200 {object} dto.GetChatAgentMcpServerToolSettingsResponse "获取成功"
// @Failure 400 {object} dto.GetChatAgentMcpServerToolSettingsResponse "请求参数错误"
// @Failure 500 {object} dto.GetChatAgentMcpServerToolSettingsResponse "服务器内部错误"
// @Router /api/v1/chat-agent-mcp-server-tool/settings/{chat_agent_id} [get]
func (h *ChatAgentMcpServerToolHandler) GetChatAgentMcpServerToolSettings(c *gin.Context) {
	chatAgentIDStr := c.Param("chatAgentID")
	if chatAgentIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.GetChatAgentMcpServerToolSettingsResponse{
			Success: false,
			Message: "聊天智能体ID不能为空",
		})
		return
	}

	// 验证聊天智能体ID
	chatAgentID, err := uuid.Parse(chatAgentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetChatAgentMcpServerToolSettingsResponse{
			Success: false,
			Message: "无效的聊天智能体ID",
		})
		return
	}

	// 调用业务逻辑层
	settings, err := h.chatAgentMcpServerToolService.GetChatAgentMcpServerToolSettings(c.Request.Context(), chatAgentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetChatAgentMcpServerToolSettingsResponse{
			Success: false,
			Message: "获取工具设置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetChatAgentMcpServerToolSettingsResponse{
		Success: true,
		Data:    settings,
		Message: "获取成功",
	})
}

// GetChatAgentAvailableMcpServerTools 获取聊天智能体可用的MCP工具列表
// @Summary 获取聊天智能体可用的MCP工具列表
// @Description 获取指定聊天智能体可用的所有MCP工具及其启用状态
// @Tags ChatAgentMcpServerTool
// @Accept json
// @Produce json
// @Param chat_agent_id path string true "聊天智能体ID"
// @Success 200 {object} dto.GetChatAgentAvailableMcpServerToolsResponse "获取成功"
// @Failure 400 {object} dto.GetChatAgentAvailableMcpServerToolsResponse "请求参数错误"
// @Failure 500 {object} dto.GetChatAgentAvailableMcpServerToolsResponse "服务器内部错误"
// @Router /api/v1/chat-agent-mcp-server-tool/available/{chat_agent_id} [get]
func (h *ChatAgentMcpServerToolHandler) GetChatAgentAvailableMcpServerTools(c *gin.Context) {
	chatAgentIDStr := c.Param("chatAgentID")
	if chatAgentIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.GetChatAgentAvailableMcpServerToolsResponse{
			Success: false,
			Message: "聊天智能体ID不能为空",
		})
		return
	}

	// 验证聊天智能体ID
	chatAgentID, err := uuid.Parse(chatAgentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetChatAgentAvailableMcpServerToolsResponse{
			Success: false,
			Message: "无效的聊天智能体ID",
		})
		return
	}

	// 调用业务逻辑层
	tools, err := h.chatAgentMcpServerToolService.GetChatAgentAvailableMcpServerTools(c.Request.Context(), chatAgentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetChatAgentAvailableMcpServerToolsResponse{
			Success: false,
			Message: "获取可用工具失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetChatAgentAvailableMcpServerToolsResponse{
		Success: true,
		Data:    tools,
		Message: "获取成功",
	})
}
