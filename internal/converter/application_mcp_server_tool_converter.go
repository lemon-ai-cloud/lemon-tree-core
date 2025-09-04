// Package converter 提供数据转换功能
// 负责在不同层之间转换数据格式，如模型到DTO的转换
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"
)

// ApplicationMcpServerToolModelToApplicationMcpServerToolDto 将模型转换为DTO
// 将数据库模型转换为用于API响应的DTO
// 参数：model - 数据库模型
// 返回：DTO对象
func ApplicationMcpServerToolModelToApplicationMcpServerToolDto(model *models.ApplicationMcpServerTool) dto.ApplicationMcpServerToolDto {
	return dto.ApplicationMcpServerToolDto{
		ID:                           model.ID.String(),
		ApplicationID:                model.ApplicationID.String(),
		ApplicationMcpServerConfigID: model.ApplicationMcpServerConfigID.String(),
		Name:                         model.Name,
		Title:                        model.Title,
		Description:                  model.Description,
		CreatedAt:                    model.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:                    model.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ApplicationMcpServerToolModelListToApplicationMcpServerToolDtoList 将模型列表转换为DTO列表
// 将数据库模型列表转换为用于API响应的DTO列表
// 参数：models - 数据库模型列表
// 返回：DTO列表
func ApplicationMcpServerToolModelListToApplicationMcpServerToolDtoList(models []*models.ApplicationMcpServerTool) []dto.ApplicationMcpServerToolDto {
	dtoList := make([]dto.ApplicationMcpServerToolDto, len(models))
	for i, model := range models {
		dtoList[i] = ApplicationMcpServerToolModelToApplicationMcpServerToolDto(model)
	}
	return dtoList
}
