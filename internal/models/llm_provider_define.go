// Package models 提供应用程序的数据模型定义
package models

import (
	"lemon-tree-core/internal/base"
)

// LlmProviderDefine 大模型供应商预定义
// 只是定义，没有API KEY 和 URL，因为只是定义没有真实落地
type LlmProviderDefine struct {
	base.BaseModel        // 继承基础模型，包含 ID、时间戳等通用字段
	Name           string `json:"name" gorm:"type:varchar(64);not null;comment:大语言模型供应商名称"` // 应用名称，最大长度64字符
	Description    string `json:"description" gorm:"type:varchar(512);not null;comment:大语言模型供应商描述"`
	IconUrl        string `json:"icon_url" gorm:"type:varchar(512);not null;comment:大语言模型供应商图标URL"`
	Type           string `json:"type" gorm:"type:varchar(64);not null;comment:大语言模型供应商类型"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (LlmProviderDefine) TableName() string {
	return "ltc_llm_provider_define"
}
