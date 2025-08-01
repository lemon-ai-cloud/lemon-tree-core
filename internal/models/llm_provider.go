// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
)

// Application 应用模型结构体
type LlmProvider struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	Name           string    `json:"name" gorm:"type:varchar(64);not null;comment:大语言模型供应商名称"` // 应用名称，最大长度64字符
	Description    string    `json:"description" gorm:"type:varchar(512);not null;comment:大语言模型供应商描述"`
	Type           string    `json:"type" gorm:"type:varchar(64);not null;comment:大语言模型供应商类型"`
	IconUrl        string    `json:"icon_url" gorm:"type:varchar(512);not null;comment:大语言模型供应商图标URL"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	ApiUrl         string    `json:"api_url" gorm:"type:varchar(512);not null;comment:大语言模型供应商API URL"`
	ApiKey         string    `json:"api_key" gorm:"type:varchar(512);not null;comment:大语言模型供应商API Key"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (LlmProvider) TableName() string {
	return "ltc_llm_provider"
}
