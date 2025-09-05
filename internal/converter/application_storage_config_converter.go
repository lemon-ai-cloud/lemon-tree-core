// Package converter 提供数据转换功能
// 负责在不同层之间转换数据格式，如模型到DTO的转换
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
)

// ApplicationStorageConfigModelToApplicationStorageConfigDto 将模型转换为DTO
// 将数据库模型转换为用于API响应的DTO
// 参数：model - 数据库模型
// 返回：DTO对象
func ApplicationStorageConfigModelToApplicationStorageConfigDto(model *models.ApplicationStorageConfig) dto.ApplicationStorageConfigDto {
	return dto.ApplicationStorageConfigDto{
		ID:            model.ID.String(),
		Type:          model.Type,
		ApplicationID: model.ApplicationID.String(),
		RootPath:      model.RootPath,
		Endpoint:      model.Endpoint,
		Region:        model.Region,
		BucketName:    model.BucketName,
		SecretId:      model.SecretId,
		SecretKey:     model.SecretKey,
		KeyPrefix:     model.KeyPrefix,
		CreatedAt:     model.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     model.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// SaveApplicationStorageConfigRequestToApplicationStorageConfigModel 将保存请求转换为模型
// 将前端保存请求转换为数据库模型
// 参数：request - 保存请求
// 返回：数据库模型
func SaveApplicationStorageConfigRequestToApplicationStorageConfigModel(request *dto.SaveApplicationStorageConfigRequest) *models.ApplicationStorageConfig {
	model := &models.ApplicationStorageConfig{
		Type:       request.Type,
		RootPath:   request.RootPath,
		Endpoint:   request.Endpoint,
		Region:     request.Region,
		BucketName: request.BucketName,
		SecretId:   request.SecretId,
		SecretKey:  request.SecretKey,
		KeyPrefix:  request.KeyPrefix,
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
