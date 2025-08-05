// Package converter 提供数据转换功能
// 负责在模型和DTO之间进行数据转换
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
)

// LlmProviderDefineModelToLlmProviderDefineDto 将模型转换为DTO
// 将 LlmProviderDefine 模型转换为 LlmProviderDefineDto
// 参数：model - LlmProviderDefine 模型
// 返回：LlmProviderDefineDto
func LlmProviderDefineModelToLlmProviderDefineDto(model *models.LlmProviderDefine) *dto.LlmProviderDefineDto {
	if model == nil {
		return nil
	}

	return &dto.LlmProviderDefineDto{
		BaseModelDto: dto.BaseModelDto{
			ID:        model.ID,
			CreatedAt: model.CreatedAt.UnixMilli(),
			UpdatedAt: model.UpdatedAt.UnixMilli(),
		},
		Name:        model.Name,
		Description: model.Description,
		IconUrl:     model.IconUrl,
		Type:        model.Type,
	}
}

// LlmProviderDefineModelListToLlmProviderDefineDtoList 将模型列表转换为DTO列表
// 将 LlmProviderDefine 模型列表转换为 LlmProviderDefineDto 列表
// 参数：models - LlmProviderDefine 模型列表
// 返回：LlmProviderDefineDto 列表
func LlmProviderDefineModelListToLlmProviderDefineDtoList(models []*models.LlmProviderDefine) []*dto.LlmProviderDefineDto {
	if models == nil {
		return nil
	}

	dtos := make([]*dto.LlmProviderDefineDto, len(models))
	for i, model := range models {
		dtos[i] = LlmProviderDefineModelToLlmProviderDefineDto(model)
	}
	return dtos
}

// LlmProviderDefineSaveDtoToLlmProviderDefineModel 将保存DTO转换为模型
// 将 LlmProviderDefineSaveDto 转换为 LlmProviderDefine 模型
// 参数：dto - LlmProviderDefineSaveDto
// 返回：LlmProviderDefine 模型
func LlmProviderDefineSaveDtoToLlmProviderDefineModel(dto *dto.LlmProviderDefineSaveDto) *models.LlmProviderDefine {
	if dto == nil {
		return nil
	}

	model := &models.LlmProviderDefine{
		Name:        dto.Name,
		Description: dto.Description,
		IconUrl:     dto.IconUrl,
		Type:        dto.Type,
	}

	// 如果ID不为空，则解析UUID
	if dto.ID != "" {
		if id, err := uuid.Parse(dto.ID); err == nil {
			model.ID = id
		}
	}

	return model
}

// LlmProviderDefineQueryDtoToLlmProviderDefineModel 将查询DTO转换为模型
// 将 LlmProviderDefineQueryDto 转换为 LlmProviderDefine 模型（用于查询）
// 参数：dto - LlmProviderDefineQueryDto
// 返回：LlmProviderDefine 模型
func LlmProviderDefineQueryDtoToLlmProviderDefineModel(dto *dto.LlmProviderDefineQueryDto) *models.LlmProviderDefine {
	if dto == nil {
		return nil
	}

	return &models.LlmProviderDefine{
		Name:        dto.Name,
		Description: dto.Description,
		Type:        dto.Type,
	}
}
