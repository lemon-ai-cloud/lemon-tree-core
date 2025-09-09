// Package converter 提供数据转换功能
// 用于在不同层之间转换数据格式，如模型与DTO之间的转换
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
)

// LlmProviderModelToLlmProviderDto 将 ApplicationLlmProvider 模型转换为 LlmProviderDto
// 参数：llmProvider - 大语言模型提供商模型
// 返回：大语言模型提供商DTO
func LlmProviderModelToLlmProviderDto(llmProvider *models.ApplicationLlmProvider) *dto.LlmProviderDto {
	if llmProvider == nil {
		return nil
	}

	return &dto.LlmProviderDto{
		ID:            llmProvider.ID.String(),
		Name:          llmProvider.Name,
		Description:   llmProvider.Description,
		Type:          llmProvider.Type,
		IconUrl:       llmProvider.IconUrl,
		ApplicationID: llmProvider.ApplicationID.String(),
		ApiUrl:        llmProvider.ApiUrl,
		ApiKey:        llmProvider.ApiKey,
		CreatedAt:     llmProvider.CreatedAt,
		UpdatedAt:     llmProvider.UpdatedAt,
	}
}

// LlmProviderDtoToLlmProviderModel 将 LlmProviderDto 转换为 ApplicationLlmProvider 模型
// 参数：llmProviderDto - 大语言模型提供商DTO
// 返回：大语言模型提供商模型
func LlmProviderDtoToLlmProviderModel(llmProviderDto *dto.LlmProviderDto) *models.ApplicationLlmProvider {
	if llmProviderDto == nil {
		return nil
	}

	// 解析UUID字符串
	var applicationID uuid.UUID
	if llmProviderDto.ApplicationID != "" {
		if parsedID, err := uuid.Parse(llmProviderDto.ApplicationID); err == nil {
			applicationID = parsedID
		}
	}

	return &models.ApplicationLlmProvider{
		Name:          llmProviderDto.Name,
		Description:   llmProviderDto.Description,
		Type:          llmProviderDto.Type,
		IconUrl:       llmProviderDto.IconUrl,
		ApplicationID: applicationID,
		ApiUrl:        llmProviderDto.ApiUrl,
		ApiKey:        llmProviderDto.ApiKey,
	}
}

// LlmProviderModelListToLlmProviderDtoList 将 ApplicationLlmProvider 模型列表转换为 LlmProviderDto 列表
// 参数：llmProviders - 大语言模型提供商模型列表
// 返回：大语言模型提供商DTO列表
func LlmProviderModelListToLlmProviderDtoList(llmProviders []*models.ApplicationLlmProvider) []*dto.LlmProviderDto {
	if llmProviders == nil {
		return nil
	}

	result := make([]*dto.LlmProviderDto, len(llmProviders))
	for i, llmProvider := range llmProviders {
		result[i] = LlmProviderModelToLlmProviderDto(llmProvider)
	}
	return result
}

// LlmProviderDtoListToLlmProviderModelList 将 LlmProviderDto 列表转换为 ApplicationLlmProvider 模型列表
// 参数：llmProviderDtos - 大语言模型提供商DTO列表
// 返回：大语言模型提供商模型列表
func LlmProviderDtoListToLlmProviderModelList(llmProviderDtos []*dto.LlmProviderDto) []*models.ApplicationLlmProvider {
	if llmProviderDtos == nil {
		return nil
	}

	result := make([]*models.ApplicationLlmProvider, len(llmProviderDtos))
	for i, llmProviderDto := range llmProviderDtos {
		result[i] = LlmProviderDtoToLlmProviderModel(llmProviderDto)
	}
	return result
}

// LlmProviderSaveDtoToLlmProviderModel 将 LlmProviderSaveDto 转换为 ApplicationLlmProvider 模型
// 参数：llmProviderSaveDto - 大语言模型提供商保存DTO
// 返回：大语言模型提供商模型
func LlmProviderSaveDtoToLlmProviderModel(llmProviderSaveDto *dto.LlmProviderSaveDto) *models.ApplicationLlmProvider {
	if llmProviderSaveDto == nil {
		return nil
	}

	// 解析UUID字符串
	var id uuid.UUID
	if llmProviderSaveDto.ID != "" {
		if parsedID, err := uuid.Parse(llmProviderSaveDto.ID); err == nil {
			id = parsedID
		}
	}

	var applicationID uuid.UUID
	if llmProviderSaveDto.ApplicationID != "" {
		if parsedID, err := uuid.Parse(llmProviderSaveDto.ApplicationID); err == nil {
			applicationID = parsedID
		}
	}

	llmProvider := &models.ApplicationLlmProvider{
		Name:          llmProviderSaveDto.Name,
		Description:   llmProviderSaveDto.Description,
		Type:          llmProviderSaveDto.Type,
		IconUrl:       llmProviderSaveDto.IconUrl,
		ApplicationID: applicationID,
		ApiUrl:        llmProviderSaveDto.ApiUrl,
		ApiKey:        llmProviderSaveDto.ApiKey,
	}

	// 设置ID字段（如果存在）
	if id != uuid.Nil {
		llmProvider.ID = id
	}

	return llmProvider
}

// LlmProviderQueryDtoToLlmProviderModel 将 LlmProviderQueryDto 转换为 ApplicationLlmProvider 模型
// 参数：llmProviderQueryDto - 大语言模型提供商查询DTO
// 返回：大语言模型提供商模型（用于查询）
func LlmProviderQueryDtoToLlmProviderModel(llmProviderQueryDto *dto.LlmProviderQueryDto) *models.ApplicationLlmProvider {
	if llmProviderQueryDto == nil {
		return nil
	}

	// 解析UUID字符串
	var applicationID uuid.UUID
	if llmProviderQueryDto.ApplicationID != "" {
		if parsedID, err := uuid.Parse(llmProviderQueryDto.ApplicationID); err == nil {
			applicationID = parsedID
		}
	}

	return &models.ApplicationLlmProvider{
		Name:          llmProviderQueryDto.Name,
		Description:   llmProviderQueryDto.Description,
		Type:          llmProviderQueryDto.Type,
		IconUrl:       llmProviderQueryDto.IconUrl,
		ApplicationID: applicationID,
		ApiUrl:        llmProviderQueryDto.ApiUrl,
		ApiKey:        llmProviderQueryDto.ApiKey,
	}
}
