// Package converter 提供数据转换功能
// 负责在不同层之间转换数据格式，如模型到DTO的转换
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
)

// ApplicationMcpServerConfigModelToApplicationMcpServerConfigDto 将模型转换为DTO
// 将数据库模型转换为用于API响应的DTO
// 参数：model - 数据库模型
// 返回：DTO对象
func ApplicationMcpServerConfigModelToApplicationMcpServerConfigDto(model *models.ApplicationMcpServerConfig) dto.ApplicationMcpServerConfigDto {
	return dto.ApplicationMcpServerConfigDto{
		ID:                   model.ID.String(),
		ApplicationID:        model.ApplicationID.String(),
		Name:                 model.Name,
		Description:          model.Description,
		Version:              model.Version,
		McpServerConnectType: model.McpServerConnectType,
		McpServerTimeout:     model.McpServerTimeout,
		McpServerUrl:         model.McpServerUrl,
		McpServerHeader:      model.McpServerHeader,
		McpServerCommand:     model.McpServerCommand,
		McpServerArgs:        model.McpServerArgs,
		McpServerEnv:         model.McpServerEnv,
		CreatedAt:            model.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:            model.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ApplicationMcpServerConfigModelListToApplicationMcpServerConfigDtoList 将模型列表转换为DTO列表
// 将数据库模型列表转换为用于API响应的DTO列表
// 参数：models - 数据库模型列表
// 返回：DTO列表
func ApplicationMcpServerConfigModelListToApplicationMcpServerConfigDtoList(models []*models.ApplicationMcpServerConfig) []dto.ApplicationMcpServerConfigDto {
	dtoList := make([]dto.ApplicationMcpServerConfigDto, len(models))
	for i, model := range models {
		dtoList[i] = ApplicationMcpServerConfigModelToApplicationMcpServerConfigDto(model)
	}
	return dtoList
}

// SaveApplicationMcpServerConfigRequestToApplicationMcpServerConfigModel 将保存请求转换为模型
// 将前端保存请求转换为数据库模型
// 参数：request - 保存请求
// 返回：数据库模型
func SaveApplicationMcpServerConfigRequestToApplicationMcpServerConfigModel(request *dto.SaveApplicationMcpServerConfigRequest) *models.ApplicationMcpServerConfig {
	model := &models.ApplicationMcpServerConfig{
		Name:                 request.Name,
		Description:          request.Description,
		Version:              request.Version,
		McpServerConnectType: request.McpServerConnectType,
		McpServerTimeout:     request.McpServerTimeout,
		McpServerUrl:         request.McpServerUrl,
		McpServerHeader:      request.McpServerHeader,
		McpServerCommand:     request.McpServerCommand,
		McpServerArgs:        request.McpServerArgs,
		McpServerEnv:         request.McpServerEnv,
	}

	// 解析应用ID
	if applicationID, err := uuid.Parse(request.ApplicationID); err == nil {
		model.ApplicationID = applicationID
	}

	// 如果有ID，则解析ID（用于更新操作）
	if request.ID != nil && *request.ID != "" {
		if id, err := uuid.Parse(*request.ID); err == nil {
			model.ID = id
		}
	}

	return model
}
