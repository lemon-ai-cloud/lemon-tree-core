// Package service 提供业务逻辑层功能
// 负责处理业务规则和逻辑，连接 Handler 层和 Repository 层
package service

import (
	"context"
	"fmt"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"

	"github.com/google/uuid"
)

// ApplicationService Application 业务逻辑层接口
// 定义了 Application 相关的所有业务操作接口
// 包含业务规则验证和数据处理逻辑
type ApplicationService interface {
	GetApplicationByID(ctx context.Context, id uuid.UUID) (*models.Application, error)               // 根据ID获取应用
	GetAllApplications(ctx context.Context) ([]*models.Application, error)                           // 获取所有应用
	SaveApplication(ctx context.Context, application *models.Application) error                      // 保存应用（upsert）
	DeleteApplication(ctx context.Context, id uuid.UUID) error                                       // 删除应用
	QueryApplications(ctx context.Context, query *models.Application) ([]*models.Application, error) // 动态查询应用
}

// applicationService Application 业务逻辑层实现
// 实现了 ApplicationService 接口的所有方法
// 包含业务规则验证和数据处理逻辑
type applicationService struct {
	appRepo repository.ApplicationRepository // Application 数据访问层接口
}

// NewApplicationService 创建 Application Service 实例
// 返回 ApplicationService 接口的实现
// 参数：appRepo - Application 数据访问层接口
func NewApplicationService(appRepo repository.ApplicationRepository) ApplicationService {
	return &applicationService{
		appRepo: appRepo,
	}
}

// GetApplicationByID 根据ID获取应用
// 根据 UUID 获取指定的应用信息
// 参数：ctx - 上下文，id - 应用的 UUID
// 返回：应用对象和错误信息
func (s *applicationService) GetApplicationByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	return s.appRepo.GetByID(ctx, id)
}

// GetAllApplications 获取所有应用
// 获取所有应用的信息列表
// 参数：ctx - 上下文
// 返回：应用列表和错误信息
func (s *applicationService) GetAllApplications(ctx context.Context) ([]*models.Application, error) {
	return s.appRepo.ListAll(ctx)
}

// SaveApplication 保存应用（upsert）
// 如果应用存在则更新，不存在则创建
// 参数：ctx - 上下文，application - 要保存的应用对象
// 返回：错误信息
func (s *applicationService) SaveApplication(ctx context.Context, application *models.Application) error {
	// 检查应用是否已存在
	if application.ID != uuid.Nil {
		// 更新现有应用
		existingApplication, err := s.appRepo.GetByID(ctx, application.ID)
		if err != nil {
			return fmt.Errorf("应用不存在: %w", err)
		}

		// 基于existingApplication修改值
		existingApplication.Name = application.Name
		existingApplication.Description = application.Description

		// 保存修改后的existingApplication
		return s.appRepo.Save(ctx, existingApplication)
	} else {
		// 创建新应用
		application.ID = uuid.New()
		return s.appRepo.Save(ctx, application)
	}
}

// DeleteApplication 删除应用
// 处理删除应用的业务逻辑（软删除）
// 参数：ctx - 上下文，id - 要删除的应用 UUID
// 返回：错误信息
func (s *applicationService) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	// 这里可以添加业务规则验证逻辑
	// 例如：检查应用是否存在、验证删除权限等
	return s.appRepo.DeleteByID(ctx, id)
}

// QueryApplications 动态查询应用
// 根据查询条件动态查询应用
// 参数：ctx - 上下文，query - 查询条件对象
// 返回：匹配的应用列表和错误信息
func (s *applicationService) QueryApplications(ctx context.Context, query *models.Application) ([]*models.Application, error) {
	// 这里可以添加业务规则验证逻辑
	return s.appRepo.Query(ctx, query)
}
