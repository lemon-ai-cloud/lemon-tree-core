// Package models 提供应用程序的数据模型定义
package models

import "lemon-tree-core/internal/base"

// Application 应用模型结构体
type Application struct {
	base.BaseModel        // 继承基础模型，包含 ID、时间戳等通用字段
	Name           string `json:"name" gorm:"type:varchar(64);not null;comment:应用名称"` // 应用名称，最大长度64字符
	Description    string `json:"description" gorm:"type:varchar(512);not null;comment:应用描述"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (Application) TableName() string {
	return "ltc_application"
}
