// Package dto 提供数据传输对象定义
package dto

// LlmProviderDefineDto 大语言模型提供商定义DTO
// 用于前后端数据传输，避免直接暴露内部模型结构
type LlmProviderDefineDto struct {
	BaseModelDto        // 继承基础DTO，包含 ID、时间戳等通用字段
	Name         string `json:"name"`        // 提供商名称
	Description  string `json:"description"` // 提供商描述
	IconUrl      string `json:"icon_url"`    // 提供商图标URL
	Type         string `json:"type"`        // 提供商类型
}

// LlmProviderDefineSaveDto 大语言模型提供商定义保存DTO（创建或更新）
// 用于创建或更新提供商定义时的数据传输
type LlmProviderDefineSaveDto struct {
	ID          string `json:"id,omitempty"`                   // 提供商定义ID（更新时提供）
	Name        string `json:"name" binding:"required"`        // 提供商名称
	Description string `json:"description" binding:"required"` // 提供商描述
	IconUrl     string `json:"icon_url" binding:"required"`    // 提供商图标URL
	Type        string `json:"type" binding:"required"`        // 提供商类型
}

// LlmProviderDefineQueryDto 大语言模型提供商定义查询DTO
// 用于查询提供商定义时的数据传输
type LlmProviderDefineQueryDto struct {
	Name        string `json:"name,omitempty"`        // 提供商名称（模糊查询）
	Description string `json:"description,omitempty"` // 提供商描述（模糊查询）
	Type        string `json:"type,omitempty"`        // 提供商类型（精确查询）
}
