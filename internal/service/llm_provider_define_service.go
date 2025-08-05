// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"
	"os"
	"path/filepath"
	"strings"
)

// LlmProviderDefineService 大语言模型提供商定义业务逻辑层接口
// 定义大语言模型提供商定义相关的业务逻辑方法
type LlmProviderDefineService interface {
	// GetLlmProviderDefineByID 根据ID获取大语言模型提供商定义
	// 根据UUID获取指定的提供商定义信息
	GetLlmProviderDefineByID(ctx context.Context, id uuid.UUID) (*models.LlmProviderDefine, error)

	// SaveLlmProviderDefine 保存大语言模型提供商定义（新增或更新）
	// 如果ID为空则新增，否则更新现有记录
	SaveLlmProviderDefine(ctx context.Context, llmProviderDefine *models.LlmProviderDefine) error

	// DeleteLlmProviderDefine 删除大语言模型提供商定义
	// 根据ID删除指定的提供商定义
	DeleteLlmProviderDefine(ctx context.Context, id uuid.UUID) error

	// GetAllLlmProviderDefines 获取所有大语言模型提供商定义
	// 返回所有提供商定义的列表
	GetAllLlmProviderDefines(ctx context.Context) ([]*models.LlmProviderDefine, error)

	// QueryLlmProviderDefines 动态查询大语言模型提供商定义
	// 根据查询条件动态查询提供商定义
	QueryLlmProviderDefines(ctx context.Context, query *models.LlmProviderDefine) ([]*models.LlmProviderDefine, error)
}

// llmProviderDefineService 大语言模型提供商定义业务逻辑层实现
// 实现 LlmProviderDefineService 接口
type llmProviderDefineService struct {
	llmProviderDefineRepo repository.LlmProviderDefineRepository // 数据访问层接口
}

// NewLlmProviderDefineService 创建大语言模型提供商定义服务实例
// 返回 LlmProviderDefineService 接口的实现
// 参数：llmProviderDefineRepo - 大语言模型提供商定义数据访问层接口
func NewLlmProviderDefineService(llmProviderDefineRepo repository.LlmProviderDefineRepository) LlmProviderDefineService {
	return &llmProviderDefineService{
		llmProviderDefineRepo: llmProviderDefineRepo,
	}
}

// GetLlmProviderDefineByID 根据ID获取大语言模型提供商定义
// 根据UUID获取指定的提供商定义信息
func (s *llmProviderDefineService) GetLlmProviderDefineByID(ctx context.Context, id uuid.UUID) (*models.LlmProviderDefine, error) {
	return s.llmProviderDefineRepo.GetByID(ctx, id)
}

// QueryLlmProviderDefines 动态查询大语言模型提供商定义
// 根据查询条件动态查询提供商定义
func (s *llmProviderDefineService) QueryLlmProviderDefines(ctx context.Context, query *models.LlmProviderDefine) ([]*models.LlmProviderDefine, error) {
	return s.llmProviderDefineRepo.Query(ctx, query)
}

// SaveLlmProviderDefine 保存大语言模型提供商定义（新增或更新）
// 如果ID为空则新增，否则更新现有记录
func (s *llmProviderDefineService) SaveLlmProviderDefine(ctx context.Context, llmProviderDefine *models.LlmProviderDefine) error {
	// 数据验证
	if err := s.validateLlmProviderDefine(llmProviderDefine); err != nil {
		return err
	}

	// 处理图片保存
	if err := s.handleIconSave(llmProviderDefine); err != nil {
		return fmt.Errorf("保存图片失败: %w", err)
	}

	// 根据ID判断是新增还是更新
	if llmProviderDefine.ID == uuid.Nil {
		// 新增：生成新的UUID
		llmProviderDefine.ID = uuid.New()
		return s.llmProviderDefineRepo.Create(ctx, llmProviderDefine)
	} else {
		// 更新：检查记录是否存在
		existing, err := s.llmProviderDefineRepo.GetByID(ctx, llmProviderDefine.ID)
		if err != nil {
			return fmt.Errorf("提供商定义不存在: %w", err)
		}
		if existing == nil {
			return fmt.Errorf("提供商定义不存在")
		}
		return s.llmProviderDefineRepo.Update(ctx, llmProviderDefine)
	}
}

// DeleteLlmProviderDefine 删除大语言模型提供商定义
// 根据ID删除指定的提供商定义
func (s *llmProviderDefineService) DeleteLlmProviderDefine(ctx context.Context, id uuid.UUID) error {
	// 检查记录是否存在
	existing, err := s.llmProviderDefineRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("提供商定义不存在: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("提供商定义不存在")
	}

	// 执行删除
	return s.llmProviderDefineRepo.DeleteByID(ctx, id)
}

// GetAllLlmProviderDefines 获取所有大语言模型提供商定义
// 返回所有提供商定义的列表
func (s *llmProviderDefineService) GetAllLlmProviderDefines(ctx context.Context) ([]*models.LlmProviderDefine, error) {
	return s.llmProviderDefineRepo.ListAll(ctx)
}

// validateLlmProviderDefine 验证大语言模型提供商定义数据
// 检查必填字段是否为空
func (s *llmProviderDefineService) validateLlmProviderDefine(llmProviderDefine *models.LlmProviderDefine) error {
	if llmProviderDefine == nil {
		return fmt.Errorf("提供商定义不能为空")
	}

	if llmProviderDefine.Name == "" {
		return fmt.Errorf("提供商名称不能为空")
	}

	if llmProviderDefine.Description == "" {
		return fmt.Errorf("提供商描述不能为空")
	}

	if llmProviderDefine.IconUrl == "" {
		return fmt.Errorf("提供商图标URL不能为空")
	}

	if llmProviderDefine.Type == "" {
		return fmt.Errorf("提供商类型不能为空")
	}

	return nil
}

// handleIconSave 处理图标保存
// 如果 IconUrl 是 base64 格式，则保存为本地文件并更新为相对路径
// 如果 IconUrl 不是 base64 格式，则保持原内容不变
func (s *llmProviderDefineService) handleIconSave(llmProviderDefine *models.LlmProviderDefine) error {
	// 检查是否是 base64 格式
	if !strings.HasPrefix(llmProviderDefine.IconUrl, "data:image/") {
		return nil // 不是 base64 格式，保持原内容不变
	}

	// 解析 base64 数据
	parts := strings.Split(llmProviderDefine.IconUrl, ",")
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
	saveDir := filepath.Join(workspacePath, define.WorkspaceDirNameLlmProviderDefineIcon)

	// 确保目录存在
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 如果是更新操作，尝试删除旧文件
	if llmProviderDefine.ID != uuid.Nil {
		// 获取现有记录以获取旧的图片路径
		existing, err := s.llmProviderDefineRepo.GetByID(context.Background(), llmProviderDefine.ID)
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

	// 更新 IconUrl 为相对路径（以 WorkspaceDirNameLlmProviderDefineIcon 开头）
	llmProviderDefine.IconUrl = define.WorkspaceDirNameLlmProviderDefineIcon + newFileName

	return nil
}
