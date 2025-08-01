// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
)

// ApplicationInternalToolNetSearchConfig 应用内部工具 - 网络搜索配置
type ApplicationInternalToolNetSearchConfig struct {
	base.BaseModel              // 继承基础模型，包含 ID、时间戳等通用字段
	ApplicationID     uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	Type              string    `json:"type" gorm:"type:varchar(64);not null;comment:网络搜索工具类型"`
	ApiUrl            string    `json:"api_url" gorm:"type:varchar(512);not null;comment:网络搜索API URL"`
	ApiKey            string    `json:"api_key" gorm:"type:varchar(512);not null;comment:网络搜索API Key"`
	SearchResultCount int       `json:"search_result_count" gorm:"type:int;not null;comment:网络搜索结果数量"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ApplicationInternalToolNetSearchConfig) TableName() string {
	return "ltc_application_internal_tool_netsearch_config"
}
