// Package dto 提供数据传输对象定义
// 用于在不同层之间传输数据，避免直接暴露内部模型
package dto

import (
	"time"
)

// LlmProviderDto 大语言模型提供商数据传输对象
// 用于向前端返回提供商信息
type LlmProviderDto struct {
	ID            string    `json:"id"`             // 提供商ID
	Name          string    `json:"name"`           // 提供商名称
	Description   string    `json:"description"`    // 提供商描述
	Type          string    `json:"type"`           // 提供商类型
	IconUrl       string    `json:"icon_url"`       // 提供商图标URL
	ApplicationID string    `json:"application_id"` // 所属应用ID
	ApiUrl        string    `json:"api_url"`        // API URL
	ApiKey        string    `json:"api_key"`        // API Key
	CreatedAt     time.Time `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`     // 更新时间
}

// LlmProviderSaveDto 大语言模型提供商保存数据传输对象
// 用于接收前端提交的提供商信息
type LlmProviderSaveDto struct {
	ID            string `json:"id"`             // 提供商ID（更新时必填）
	Name          string `json:"name"`           // 提供商名称
	Description   string `json:"description"`    // 提供商描述
	Type          string `json:"type"`           // 提供商类型
	IconUrl       string `json:"icon_url"`       // 提供商图标URL
	ApplicationID string `json:"application_id"` // 所属应用ID
	ApiUrl        string `json:"api_url"`        // API URL
	ApiKey        string `json:"api_key"`        // API Key
}

// LlmProviderQueryDto 大语言模型提供商查询数据传输对象
// 用于接收前端的查询条件
type LlmProviderQueryDto struct {
	ID            string `json:"id"`             // 提供商ID
	Name          string `json:"name"`           // 提供商名称
	Description   string `json:"description"`    // 提供商描述
	Type          string `json:"type"`           // 提供商类型
	IconUrl       string `json:"icon_url"`       // 提供商图标URL
	ApplicationID string `json:"application_id"` // 所属应用ID
	ApiUrl        string `json:"api_url"`        // API URL
	ApiKey        string `json:"api_key"`        // API Key
}
