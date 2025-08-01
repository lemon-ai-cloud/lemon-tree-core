// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
)

// ApplicationModel 应用模型
// 记录当前应用从模型提供商中选择确认使用的模型
type ApplicationModel struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	Name           string    `json:"name" gorm:"type:varchar(64);not null;comment:模型名称"`
	Alias          string    `json:"alias" gorm:"type:varchar(64);not null;comment:模型别名"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	LlmProviderID  uuid.UUID `json:"llm_provider_id" gorm:"type:char(36);not null;comment:所属模型供应商ID"`
	// 大模型能力相关
	AbilityVision         bool `json:"ability_vision" gorm:"type:tinyint(1);not null;comment:是否为视觉能力"`
	AbilityNetwork        bool `json:"ability_network" gorm:"type:tinyint(1);not null;comment:是否有联网能力能力"`
	AbilityTextEmbeddings bool `json:"ability_text_embeddings" gorm:"type:tinyint(1);not null;comment:是否文本有嵌入能力"`
	AbilityThinking       bool `json:"ability_thinking" gorm:"type:tinyint(1);not null;comment:是否为思考能力"`
	AbilityCallTools      bool `json:"ability_call_tools" gorm:"type:tinyint(1);not null;comment:是否为调用工具能力"`
	AbilityReranking      bool `json:"ability_reranking" gorm:"type:tinyint(1);not null;comment:是否拥有重排能力"`
	// 计费相关
	BillingCurrency    string  `json:"billing_currency" gorm:"type:varchar(64);not null;comment:计费币种"`
	BillingPriceInput  float64 `json:"billing_price_input" gorm:"type:decimal(10,2);not null;comment:输入计费价格，单位每百万Token"`
	BillingPriceOutput float64 `json:"billing_price_output" gorm:"type:decimal(10,2);not null;comment:输出计费价格，单位每百万Token"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ApplicationModel) TableName() string {
	return "ltc_application_model"
}
