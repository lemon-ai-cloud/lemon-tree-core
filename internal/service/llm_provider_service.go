// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// LlmProviderService 大语言模型提供商业务逻辑层接口
// 定义大语言模型提供商相关的业务逻辑方法
type LlmProviderService interface {
	// GetLlmProviderByID 根据ID获取大语言模型提供商
	// 根据UUID获取指定的提供商信息
	GetLlmProviderByID(ctx context.Context, id uuid.UUID) (*models.ApplicationLlmProvider, error)

	// SaveLlmProvider 保存大语言模型提供商（新增或更新）
	// 如果ID为空则新增，否则更新现有记录
	SaveLlmProvider(ctx context.Context, llmProvider *models.ApplicationLlmProvider) error

	// DeleteLlmProvider 删除大语言模型提供商
	// 根据ID删除指定的提供商
	DeleteLlmProvider(ctx context.Context, id uuid.UUID) error

	// GetAllLlmProviders 获取所有大语言模型提供商
	// 返回所有提供商的列表
	GetAllLlmProviders(ctx context.Context) ([]*models.ApplicationLlmProvider, error)

	// QueryLlmProviders 动态查询大语言模型提供商
	// 根据查询条件动态查询提供商
	QueryLlmProviders(ctx context.Context, query *models.ApplicationLlmProvider) ([]*models.ApplicationLlmProvider, error)

	// GetLlmProvidersByApplicationID 根据应用ID获取大语言模型提供商列表
	// 返回指定应用下的所有提供商
	GetLlmProvidersByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationLlmProvider, error)
}

// llmProviderService 大语言模型提供商业务逻辑层实现
// 实现 LlmProviderService 接口
type llmProviderService struct {
	llmProviderRepo       repository.LlmProviderRepository // 数据访问层接口
	applicationLlmService ApplicationLlmService            // 应用模型服务接口
}

// NewLlmProviderService 创建大语言模型提供商服务实例
// 返回 LlmProviderService 接口的实现
// 参数：llmProviderRepo - 大语言模型提供商数据访问层接口
// 参数：applicationLlmService - 应用模型服务接口
func NewLlmProviderService(llmProviderRepo repository.LlmProviderRepository, applicationLlmService ApplicationLlmService) LlmProviderService {
	return &llmProviderService{
		llmProviderRepo:       llmProviderRepo,
		applicationLlmService: applicationLlmService,
	}
}

// GetLlmProviderByID 根据ID获取大语言模型提供商
// 根据UUID获取指定的提供商信息
func (s *llmProviderService) GetLlmProviderByID(ctx context.Context, id uuid.UUID) (*models.ApplicationLlmProvider, error) {
	return s.llmProviderRepo.GetByID(ctx, id)
}

// QueryLlmProviders 动态查询大语言模型提供商
// 根据查询条件动态查询提供商
func (s *llmProviderService) QueryLlmProviders(ctx context.Context, query *models.ApplicationLlmProvider) ([]*models.ApplicationLlmProvider, error) {
	return s.llmProviderRepo.Query(ctx, query)
}

// SaveLlmProvider 保存大语言模型提供商（新增或更新）
// 如果ID为空则新增，否则更新现有记录
func (s *llmProviderService) SaveLlmProvider(ctx context.Context, llmProvider *models.ApplicationLlmProvider) error {
	// 数据验证
	if err := s.validateLlmProvider(llmProvider); err != nil {
		return err
	}

	// 处理图片保存
	if err := s.handleIconSave(llmProvider); err != nil {
		return fmt.Errorf("保存图片失败: %w", err)
	}

	// 根据ID判断是新增还是更新
	var err error
	if llmProvider.ID == uuid.Nil {
		// 新增：生成新的UUID
		llmProvider.ID = uuid.New()
		err = s.llmProviderRepo.Create(ctx, llmProvider)
	} else {
		// 更新：检查记录是否存在
		existing, err := s.llmProviderRepo.GetByID(ctx, llmProvider.ID)
		if err != nil {
			return fmt.Errorf("提供商不存在: %w", err)
		}
		if existing == nil {
			return fmt.Errorf("提供商不存在")
		}
		err = s.llmProviderRepo.Update(ctx, llmProvider)
	}

	// 如果保存成功且有API配置，尝试获取并保存模型列表
	if err == nil && llmProvider.ApiUrl != "" && llmProvider.ApiKey != "" {
		// 异步获取模型列表，不阻塞主流程
		go func() {
			if fetchErr := s.applicationLlmService.FetchAndSaveModels(context.Background(), llmProvider); fetchErr != nil {
				// 记录错误日志，但不影响主流程
				fmt.Printf("自动获取模型列表失败: %v\n", fetchErr)
			}
		}()
	}

	return err
}

// DeleteLlmProvider 删除大语言模型提供商
// 根据ID删除指定的提供商
func (s *llmProviderService) DeleteLlmProvider(ctx context.Context, id uuid.UUID) error {
	// 检查记录是否存在
	existing, err := s.llmProviderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("提供商不存在: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("提供商不存在")
	}

	// 执行删除
	return s.llmProviderRepo.DeleteByID(ctx, id)
}

// GetAllLlmProviders 获取所有大语言模型提供商
// 返回所有提供商的列表
func (s *llmProviderService) GetAllLlmProviders(ctx context.Context) ([]*models.ApplicationLlmProvider, error) {
	return s.llmProviderRepo.ListAll(ctx)
}

// GetLlmProvidersByApplicationID 根据应用ID获取大语言模型提供商列表
// 返回指定应用下的所有提供商
func (s *llmProviderService) GetLlmProvidersByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationLlmProvider, error) {
	return s.llmProviderRepo.GetByApplicationID(ctx, applicationID)
}

// validateLlmProvider 验证大语言模型提供商数据
// 检查必填字段是否为空
func (s *llmProviderService) validateLlmProvider(llmProvider *models.ApplicationLlmProvider) error {
	if llmProvider == nil {
		return fmt.Errorf("提供商不能为空")
	}

	if llmProvider.Name == "" {
		return fmt.Errorf("提供商名称不能为空")
	}

	if llmProvider.Description == "" {
		return fmt.Errorf("提供商描述不能为空")
	}

	if llmProvider.IconUrl == "" {
		return fmt.Errorf("提供商图标URL不能为空")
	}

	if llmProvider.Type == "" {
		return fmt.Errorf("提供商类型不能为空")
	}

	if llmProvider.ApplicationID == uuid.Nil {
		return fmt.Errorf("所属应用ID不能为空")
	}

	if llmProvider.ApiUrl == "" {
		return fmt.Errorf("API URL不能为空")
	}

	if llmProvider.ApiKey == "" {
		return fmt.Errorf("API Key不能为空")
	}

	return nil
}

// handleIconSave 处理图标保存
// 如果 IconUrl 是 base64 格式，则保存为本地文件并更新为相对路径
// 如果 IconUrl 不是 base64 格式，则保持原内容不变
func (s *llmProviderService) handleIconSave(llmProvider *models.ApplicationLlmProvider) error {
	// 检查是否是 base64 格式
	if !strings.HasPrefix(llmProvider.IconUrl, "data:image/") {
		return nil // 不是 base64 格式，保持原内容不变
	}

	// 解析 base64 数据
	parts := strings.Split(llmProvider.IconUrl, ",")
	if len(parts) != 2 {
		return fmt.Errorf("无效的 base64 格式")
	}

	// 获取图片类型和 base64 数据
	header := parts[0]
	data := parts[1]

	// 确定文件扩展名
	var ext string
	if strings.Contains(header, "image/jpeg") || strings.Contains(header, "image/jpg") {
		ext = ".jpg"
	} else if strings.Contains(header, "image/png") {
		ext = ".png"
	} else if strings.Contains(header, "image/gif") {
		ext = ".gif"
	} else if strings.Contains(header, "image/webp") {
		ext = ".webp"
	} else {
		return fmt.Errorf("不支持的图片格式")
	}

	// 解码 base64 数据
	imageData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("解码 base64 失败: %w", err)
	}

	// 获取工作区路径
	workspacePath := os.Getenv("WORKSPACE_PUBLIC_PATH")
	if workspacePath == "" {
		return fmt.Errorf("环境变量 WORKSPACE_PUBLIC_PATH 未设置")
	}

	// 构建保存路径
	saveDir := filepath.Join(workspacePath, define.WorkspaceDirNameLlmProviderIcon)

	// 确保目录存在
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 如果是更新操作，尝试删除旧文件
	if llmProvider.ID != uuid.Nil {
		// 获取现有记录以获取旧的图片路径
		existing, err := s.llmProviderRepo.GetByID(context.Background(), llmProvider.ID)
		if err == nil && existing != nil && existing.IconUrl != "" {
			// 检查旧文件是否存在并删除
			oldFilePath := filepath.Join(workspacePath, existing.IconUrl)
			if _, err := os.Stat(oldFilePath); err == nil {
				// 文件存在，删除它
				_ = os.Remove(oldFilePath)
			}
		}
	}

	// 生成新的随机文件名
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(saveDir, newFileName)

	// 保存文件
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}

	// 更新 IconUrl 为相对路径（以 WorkspaceDirNameLlmProviderIcon 开头）
	llmProvider.IconUrl = define.WorkspaceDirNameLlmProviderIcon + newFileName

	return nil
}
