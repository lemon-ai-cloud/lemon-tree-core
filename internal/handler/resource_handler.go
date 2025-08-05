// Package handler 提供 HTTP 请求处理层功能
// 负责处理 HTTP 请求、参数验证、调用业务逻辑和返回响应
package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ResourceHandler 资源文件处理器
// 处理静态资源文件的下载请求
type ResourceHandler struct{}

// NewResourceHandler 创建资源处理器实例
// 返回 ResourceHandler 的实例
func NewResourceHandler() *ResourceHandler {
	return &ResourceHandler{}
}

// DownloadFile 下载文件
// 处理 GET /api/v1/resources/download 请求
// 根据子路径下载 WORKSPACE_PUBLIC_PATH 下的文件
func (h *ResourceHandler) DownloadFile(c *gin.Context) {
	// 从查询参数获取子路径
	subPath := c.Query("path")
	if subPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 path 参数"})
		return
	}

	// 安全检查：防止路径遍历攻击
	if strings.Contains(subPath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件路径"})
		return
	}

	// 处理以 / 开头的路径，去掉开头的 /
	if strings.HasPrefix(subPath, "/") {
		subPath = subPath[1:]
	}

	// 获取工作区公共路径
	workspacePublicPath := os.Getenv("WORKSPACE_PUBLIC_PATH")
	if workspacePublicPath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "环境变量 WORKSPACE_PUBLIC_PATH 未设置"})
		return
	}

	// 构建完整文件路径
	fullPath := filepath.Join(workspacePublicPath, subPath)

	// 安全检查：确保文件路径在工作区公共目录内
	absWorkspacePath, err := filepath.Abs(workspacePublicPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析工作区路径"})
		return
	}

	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析文件路径"})
		return
	}

	if !strings.HasPrefix(absFilePath, absWorkspacePath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "访问路径超出允许范围"})
		return
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法访问文件"})
		}
		return
	}

	// 检查是否为目录
	if fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能下载目录"})
		return
	}

	// 设置响应头
	filename := filepath.Base(fullPath)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(fullPath)
}

// ListFiles 列出目录下的文件
// 处理 GET /api/v1/resources/list 请求
// 列出 WORKSPACE_PUBLIC_PATH 下指定子目录的文件列表
func (h *ResourceHandler) ListFiles(c *gin.Context) {
	// 从查询参数获取子路径
	subPath := c.Query("path")
	if subPath == "" {
		subPath = "." // 默认为根目录
	}

	// 安全检查：防止路径遍历攻击
	if strings.Contains(subPath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目录路径"})
		return
	}

	// 处理以 / 开头的路径，去掉开头的 /
	if strings.HasPrefix(subPath, "/") {
		subPath = subPath[1:]
	}

	// 获取工作区公共路径
	workspacePublicPath := os.Getenv("WORKSPACE_PUBLIC_PATH")
	if workspacePublicPath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "环境变量 WORKSPACE_PUBLIC_PATH 未设置"})
		return
	}

	// 构建完整目录路径
	fullPath := filepath.Join(workspacePublicPath, subPath)

	// 安全检查：确保目录路径在工作区公共目录内
	absWorkspacePath, err := filepath.Abs(workspacePublicPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析工作区路径"})
		return
	}

	absDirPath, err := filepath.Abs(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析目录路径"})
		return
	}

	if !strings.HasPrefix(absDirPath, absWorkspacePath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "访问路径超出允许范围"})
		return
	}

	// 检查目录是否存在
	dirInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "目录不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法访问目录"})
		}
		return
	}

	// 检查是否为目录
	if !dirInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "指定路径不是目录"})
		return
	}

	// 读取目录内容
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取目录内容"})
		return
	}

	// 构建文件列表
	var files []gin.H
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileType := "file"
		if entry.IsDir() {
			fileType = "directory"
		}

		files = append(files, gin.H{
			"name":          entry.Name(),
			"type":          fileType,
			"size":          info.Size(),
			"modified_time": info.ModTime().Unix(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"path":  subPath,
		"files": files,
	})
}

// GetFileInfo 获取文件信息
// 处理 GET /api/v1/resources/info 请求
// 获取指定文件的详细信息
func (h *ResourceHandler) GetFileInfo(c *gin.Context) {
	// 从查询参数获取子路径
	subPath := c.Query("path")
	if subPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 path 参数"})
		return
	}

	// 安全检查：防止路径遍历攻击
	if strings.Contains(subPath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件路径"})
		return
	}

	// 处理以 / 开头的路径，去掉开头的 /
	if strings.HasPrefix(subPath, "/") {
		subPath = subPath[1:]
	}

	// 获取工作区公共路径
	workspacePublicPath := os.Getenv("WORKSPACE_PUBLIC_PATH")
	if workspacePublicPath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "环境变量 WORKSPACE_PUBLIC_PATH 未设置"})
		return
	}

	// 构建完整文件路径
	fullPath := filepath.Join(workspacePublicPath, subPath)

	// 安全检查：确保文件路径在工作区公共目录内
	absWorkspacePath, err := filepath.Abs(workspacePublicPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析工作区路径"})
		return
	}

	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析文件路径"})
		return
	}

	if !strings.HasPrefix(absFilePath, absWorkspacePath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "访问路径超出允许范围"})
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法访问文件"})
		}
		return
	}

	// 确定文件类型
	fileType := "file"
	if fileInfo.IsDir() {
		fileType = "directory"
	}

	// 获取文件扩展名
	ext := filepath.Ext(fileInfo.Name())

	c.JSON(http.StatusOK, gin.H{
		"name":          fileInfo.Name(),
		"path":          subPath,
		"type":          fileType,
		"size":          fileInfo.Size(),
		"extension":     ext,
		"modified_time": fileInfo.ModTime().Unix(),
		"created_time":  fileInfo.ModTime().Unix(), // 注意：Go 的 os.Stat 不提供创建时间
	})
}
