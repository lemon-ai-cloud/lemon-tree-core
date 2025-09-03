package dto

// Package dto 提供数据传输对象定义
// 用于在不同层之间传递数据，如前端与后端之间的数据交换

// ApplicationLlmDto 应用模型 DTO
// 用于返回给前端的应用模型信息
type ApplicationLlmDto struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	Alias                 string  `json:"alias"`
	ApplicationID         string  `json:"application_id"`
	LlmProviderID         string  `json:"llm_provider_id"`
	Enabled               bool    `json:"enabled"`
	AbilityVision         bool    `json:"ability_vision"`
	AbilityNetwork        bool    `json:"ability_network"`
	AbilityTextEmbeddings bool    `json:"ability_text_embeddings"`
	AbilityThinking       bool    `json:"ability_thinking"`
	AbilityCallTools      bool    `json:"ability_call_tools"`
	AbilityReranking      bool    `json:"ability_reranking"`
	BillingCurrency       string  `json:"billing_currency"`
	BillingPriceInput     float64 `json:"billing_price_input"`
	BillingPriceOutput    float64 `json:"billing_price_output"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
}

// SaveApplicationLlmRequest 保存应用模型请求
// 用于前端提交的应用模型信息
type SaveApplicationLlmRequest struct {
	ID                    *string `json:"id,omitempty"` // 更新时必填
	Name                  string  `json:"name"`
	Alias                 string  `json:"alias"`
	ApplicationID         string  `json:"application_id"`
	LlmProviderID         string  `json:"llm_provider_id"`
	Enabled               bool    `json:"enabled"`
	AbilityVision         bool    `json:"ability_vision"`
	AbilityNetwork        bool    `json:"ability_network"`
	AbilityTextEmbeddings bool    `json:"ability_text_embeddings"`
	AbilityThinking       bool    `json:"ability_thinking"`
	AbilityCallTools      bool    `json:"ability_call_tools"`
	AbilityReranking      bool    `json:"ability_reranking"`
	BillingCurrency       string  `json:"billing_currency"`
	BillingPriceInput     float64 `json:"billing_price_input"`
	BillingPriceOutput    float64 `json:"billing_price_output"`
}

// UpdateEnabledStatusRequest 更新启用状态请求
// 用于前端提交的启用状态更新
type UpdateEnabledStatusRequest struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

// ApplicationLlmListResponse 应用模型列表响应
// 用于返回给前端的应用模型列表
type ApplicationLlmListResponse struct {
	ApplicationLlm []ApplicationLlmDto `json:"application_llm"`
}

// SingleApplicationLlmResponse 单个应用模型响应
// 用于返回给前端的单个应用模型
type SingleApplicationLlmResponse struct {
	ApplicationLlm ApplicationLlmDto `json:"application_llm"`
}
