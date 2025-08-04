// Package converter 提供模型与DTO之间的转换功能
// 负责将内部模型转换为前端DTO，以及将DTO转换为内部模型
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
)

// ApplicationModelToApplicationDto 将 Application 模型转换为 ApplicationDto
// 参数：application - 应用模型
// 返回：ApplicationDto - 应用DTO
func ApplicationModelToApplicationDto(application *models.Application) *dto.ApplicationDto {
	if application == nil {
		return nil
	}

	var deletedAt *int64
	if !application.DeletedAt.Time.IsZero() {
		timestamp := application.DeletedAt.Time.UnixMilli()
		deletedAt = &timestamp
	}

	return &dto.ApplicationDto{
		BaseModelDto: dto.BaseModelDto{
			ID:        application.ID,
			CreatedAt: application.CreatedAt.UnixMilli(),
			UpdatedAt: application.UpdatedAt.UnixMilli(),
			DeletedAt: deletedAt,
		},
		Name:        application.Name,
		Description: application.Description,
	}
}

// ApplicationModelListToApplicationDtoList 将 Application 模型列表转换为 ApplicationDto 列表
// 参数：applications - 应用模型列表
// 返回：[]*ApplicationDto - 应用DTO列表
func ApplicationModelListToApplicationDtoList(applications []*models.Application) []*dto.ApplicationDto {
	if applications == nil {
		return nil
	}

	result := make([]*dto.ApplicationDto, len(applications))
	for i, application := range applications {
		result[i] = ApplicationModelToApplicationDto(application)
	}
	return result
}

// ApplicationSaveDtoToApplicationModel 将 ApplicationSaveDto 转换为 Application 模型
// 参数：applicationDto - 应用保存DTO
// 返回：*Application - 应用模型
func ApplicationSaveDtoToApplicationModel(applicationDto *dto.ApplicationSaveDto) *models.Application {
	if applicationDto == nil {
		return nil
	}

	application := &models.Application{
		Name:        applicationDto.Name,
		Description: applicationDto.Description,
	}

	// 如果提供了ID，则解析UUID
	if applicationDto.ID != "" {
		if id, err := uuid.Parse(applicationDto.ID); err == nil {
			application.ID = id
		}
	}

	return application
}

// ApplicationQueryDtoToApplicationModel 将 ApplicationQueryDto 转换为 Application 模型（用于查询）
// 参数：applicationDto - 应用查询DTO
// 返回：*Application - 应用模型
func ApplicationQueryDtoToApplicationModel(applicationDto *dto.ApplicationQueryDto) *models.Application {
	if applicationDto == nil {
		return nil
	}

	return &models.Application{
		Name:        applicationDto.Name,
		Description: applicationDto.Description,
	}
}
