// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"fmt"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"

	"github.com/google/uuid"
)

// ApplicationStorageConfigService 应用存储配置 业务逻辑层接口
// 定义 应用存储配置 相关的业务逻辑方法
type ApplicationStorageConfigService interface {
	// SaveApplicationStorageConfig 保存应用存储配置
	// 根据ApplicationID保存配置，如果存在则覆盖，不存在则创建
	// 保证一个ApplicationID只有一条配置记录
	SaveApplicationStorageConfig(ctx context.Context, config *models.ApplicationStorageConfig) error

	// GetApplicationStorageConfigByApplicationID 根据应用ID获取存储配置
	// 返回指定应用的存储配置
	GetApplicationStorageConfigByApplicationID(ctx context.Context, applicationID uuid.UUID) (*models.ApplicationStorageConfig, error)
}

// applicationStorageConfigService 应用存储配置 业务逻辑层实现
// 实现 ApplicationStorageConfigService 接口
type applicationStorageConfigService struct {
	applicationStorageConfigRepo repository.ApplicationStorageConfigRepository // 数据访问层接口
}

// NewApplicationStorageConfigService 创建 应用存储配置 服务实例
// 返回 ApplicationStorageConfigService 接口的实现
// 参数：applicationStorageConfigRepo - 应用存储配置 数据访问层接口
func NewApplicationStorageConfigService(applicationStorageConfigRepo repository.ApplicationStorageConfigRepository) ApplicationStorageConfigService {
	return &applicationStorageConfigService{
		applicationStorageConfigRepo: applicationStorageConfigRepo,
	}
}

// SaveApplicationStorageConfig 保存应用存储配置
// 根据ApplicationID保存配置，如果存在则覆盖，不存在则创建
// 保证一个ApplicationID只有一条配置记录
func (s *applicationStorageConfigService) SaveApplicationStorageConfig(ctx context.Context, config *models.ApplicationStorageConfig) error {
	// 数据验证
	if err := s.validateApplicationStorageConfig(config); err != nil {
		return err
	}

	// 检查是否已存在该应用的配置
	existingConfig, err := s.applicationStorageConfigRepo.GetByApplicationID(ctx, config.ApplicationID)
	if err != nil {
		// 如果查询出错且不是"未找到"错误，则返回错误
		return fmt.Errorf("查询现有配置失败: %w", err)
	}

	if existingConfig != nil {
		// 存在则更新：保留ID和创建时间，更新其他字段
		existingConfig.Type = config.Type
		existingConfig.RootPath = config.RootPath
		existingConfig.Endpoint = config.Endpoint
		existingConfig.Region = config.Region
		existingConfig.BucketName = config.BucketName
		existingConfig.SecretId = config.SecretId
		existingConfig.SecretKey = config.SecretKey
		existingConfig.KeyPrefix = config.KeyPrefix

		return s.applicationStorageConfigRepo.Update(ctx, existingConfig)
	} else {
		// 不存在则创建：生成新的UUID
		config.ID = uuid.New()
		return s.applicationStorageConfigRepo.Create(ctx, config)
	}
}

// GetApplicationStorageConfigByApplicationID 根据应用ID获取存储配置
// 返回指定应用的存储配置
func (s *applicationStorageConfigService) GetApplicationStorageConfigByApplicationID(ctx context.Context, applicationID uuid.UUID) (*models.ApplicationStorageConfig, error) {
	return s.applicationStorageConfigRepo.GetByApplicationID(ctx, applicationID)
}

// validateApplicationStorageConfig 验证应用存储配置数据
// 检查必填字段是否为空
func (s *applicationStorageConfigService) validateApplicationStorageConfig(config *models.ApplicationStorageConfig) error {
	if config == nil {
		return fmt.Errorf("存储配置不能为空")
	}

	if config.Type == "" {
		return fmt.Errorf("存储类型不能为空")
	}

	if config.ApplicationID == uuid.Nil {
		return fmt.Errorf("所属应用ID不能为空")
	}

	// 根据存储类型验证必填字段
	switch config.Type {
	case "file_system":
		if config.RootPath == "" {
			return fmt.Errorf("文件系统存储类型下，根路径不能为空")
		}
	case "s3":
		if config.Endpoint == "" {
			return fmt.Errorf("S3存储类型下，Endpoint不能为空")
		}
		if config.Region == "" {
			return fmt.Errorf("S3存储类型下，Region不能为空")
		}
		if config.BucketName == "" {
			return fmt.Errorf("S3存储类型下，存储桶名称不能为空")
		}
		if config.SecretId == "" {
			return fmt.Errorf("S3存储类型下，SecretId不能为空")
		}
		if config.SecretKey == "" {
			return fmt.Errorf("S3存储类型下，SecretKey不能为空")
		}
	default:
		return fmt.Errorf("不支持的存储类型: %s", config.Type)
	}

	return nil
}
