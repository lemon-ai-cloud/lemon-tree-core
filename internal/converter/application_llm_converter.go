// Package converter 提供数据转换功能
// 用于在不同层之间转换数据格式，如模型与DTO之间的转换
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
)

// ApplicationLlmModelToApplicationLlmDto 将 ApplicationLlm 模型转换为 ApplicationLlmDto
// 参数：applicationLlm - 应用模型
// 返回：应用模型DTO
func ApplicationLlmModelToApplicationLlmDto(applicationLlm *models.ApplicationLlm) *dto.ApplicationLlmDto {
	if applicationLlm == nil {
		return nil
	}

	return &dto.ApplicationLlmDto{
		ID:                    applicationLlm.ID.String(),
		Name:                  applicationLlm.Name,
		Alias:                 applicationLlm.Alias,
		ApplicationID:         applicationLlm.ApplicationID.String(),
		LlmProviderID:         applicationLlm.LlmProviderID.String(),
		Enabled:               applicationLlm.Enabled,
		AbilityVision:         applicationLlm.AbilityVision,
		AbilityNetwork:        applicationLlm.AbilityNetwork,
		AbilityTextEmbeddings: applicationLlm.AbilityTextEmbeddings,
		AbilityThinking:       applicationLlm.AbilityThinking,
		AbilityCallTools:      applicationLlm.AbilityCallTools,
		AbilityReranking:      applicationLlm.AbilityReranking,
		BillingCurrency:       applicationLlm.BillingCurrency,
		BillingPriceInput:     applicationLlm.BillingPriceInput,
		BillingPriceOutput:    applicationLlm.BillingPriceOutput,
		CreatedAt:             applicationLlm.CreatedAt.String(),
		UpdatedAt:             applicationLlm.UpdatedAt.String(),
	}
}

// ApplicationLlmDtoToApplicationLlmModel 将 ApplicationLlmDto 转换为 ApplicationLlm 模型
// 参数：applicationLlmDto - 应用模型DTO
// 返回：应用模型
func ApplicationLlmDtoToApplicationLlmModel(applicationLlmDto *dto.ApplicationLlmDto) *models.ApplicationLlm {
	if applicationLlmDto == nil {
		return nil
	}

	// 解析UUID字符串
	var applicationID, llmProviderID uuid.UUID
	if applicationLlmDto.ApplicationID != "" {
		if parsedID, err := uuid.Parse(applicationLlmDto.ApplicationID); err == nil {
			applicationID = parsedID
		}
	}
	if applicationLlmDto.LlmProviderID != "" {
		if parsedID, err := uuid.Parse(applicationLlmDto.LlmProviderID); err == nil {
			llmProviderID = parsedID
		}
	}

	return &models.ApplicationLlm{
		Name:                  applicationLlmDto.Name,
		Alias:                 applicationLlmDto.Alias,
		ApplicationID:         applicationID,
		LlmProviderID:         llmProviderID,
		Enabled:               applicationLlmDto.Enabled,
		AbilityVision:         applicationLlmDto.AbilityVision,
		AbilityNetwork:        applicationLlmDto.AbilityNetwork,
		AbilityTextEmbeddings: applicationLlmDto.AbilityTextEmbeddings,
		AbilityThinking:       applicationLlmDto.AbilityThinking,
		AbilityCallTools:      applicationLlmDto.AbilityCallTools,
		AbilityReranking:      applicationLlmDto.AbilityReranking,
		BillingCurrency:       applicationLlmDto.BillingCurrency,
		BillingPriceInput:     applicationLlmDto.BillingPriceInput,
		BillingPriceOutput:    applicationLlmDto.BillingPriceOutput,
	}
}

// SaveApplicationLlmRequestToApplicationLlmModel 将 SaveApplicationLlmRequest 转换为 ApplicationLlm 模型
// 参数：saveRequest - 保存应用模型请求
// 返回：应用模型
func SaveApplicationLlmRequestToApplicationLlmModel(saveRequest *dto.SaveApplicationLlmRequest) *models.ApplicationLlm {
	if saveRequest == nil {
		return nil
	}

	// 解析UUID字符串
	var id, applicationID, llmProviderID uuid.UUID
	if saveRequest.ID != nil && *saveRequest.ID != "" {
		if parsedID, err := uuid.Parse(*saveRequest.ID); err == nil {
			id = parsedID
		}
	}
	if saveRequest.ApplicationID != "" {
		if parsedID, err := uuid.Parse(saveRequest.ApplicationID); err == nil {
			applicationID = parsedID
		}
	}
	if saveRequest.LlmProviderID != "" {
		if parsedID, err := uuid.Parse(saveRequest.LlmProviderID); err == nil {
			llmProviderID = parsedID
		}
	}

	applicationLlm := &models.ApplicationLlm{
		Name:                  saveRequest.Name,
		Alias:                 saveRequest.Alias,
		ApplicationID:         applicationID,
		LlmProviderID:         llmProviderID,
		Enabled:               saveRequest.Enabled,
		AbilityVision:         saveRequest.AbilityVision,
		AbilityNetwork:        saveRequest.AbilityNetwork,
		AbilityTextEmbeddings: saveRequest.AbilityTextEmbeddings,
		AbilityThinking:       saveRequest.AbilityThinking,
		AbilityCallTools:      saveRequest.AbilityCallTools,
		AbilityReranking:      saveRequest.AbilityReranking,
		BillingCurrency:       saveRequest.BillingCurrency,
		BillingPriceInput:     saveRequest.BillingPriceInput,
		BillingPriceOutput:    saveRequest.BillingPriceOutput,
	}

	// 设置ID字段（如果存在）
	if id != uuid.Nil {
		applicationLlm.ID = id
	}

	return applicationLlm
}

// ApplicationLlmModelListToApplicationLlmDtoList 将 ApplicationLlm 模型列表转换为 ApplicationLlmDto 列表
// 参数：applicationLlms - 应用模型列表
// 返回：应用模型DTO列表
func ApplicationLlmModelListToApplicationLlmDtoList(applicationLlms []*models.ApplicationLlm) []*dto.ApplicationLlmDto {
	if applicationLlms == nil {
		return nil
	}

	result := make([]*dto.ApplicationLlmDto, len(applicationLlms))
	for i, applicationLlm := range applicationLlms {
		result[i] = ApplicationLlmModelToApplicationLlmDto(applicationLlm)
	}
	return result
}
