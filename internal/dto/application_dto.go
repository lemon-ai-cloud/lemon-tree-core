// Package dto 提供数据传输对象定义
package dto

// ApplicationDto 应用DTO
// 用于前后端数据传输，避免直接暴露内部模型结构
type ApplicationDto struct {
	BaseModelDto        // 继承基础DTO，包含 ID、时间戳等通用字段
	Name         string `json:"name"`        // 应用名称
	Description  string `json:"description"` // 应用描述
}

// ApplicationSaveDto 应用保存DTO（创建或更新）
// 用于创建或更新应用时的数据传输
type ApplicationSaveDto struct {
	ID          string `json:"id,omitempty"`                   // 应用ID（更新时提供）
	Name        string `json:"name" binding:"required"`        // 应用名称
	Description string `json:"description" binding:"required"` // 应用描述
}

// ApplicationQueryDto 应用查询DTO
// 用于查询应用时的数据传输
type ApplicationQueryDto struct {
	Name        string `json:"name,omitempty"`        // 应用名称（模糊查询）
	Description string `json:"description,omitempty"` // 应用描述（模糊查询）
}
