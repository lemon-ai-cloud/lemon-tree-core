// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"fmt"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// ApplicationLlmService ApplicationLLM 业务逻辑层接口
// 定义 ApplicationLLM 相关的业务逻辑方法
type ApplicationLlmService interface {
	// FetchAndSaveModels 获取并保存所有模型
	// 从指定的 LLM 提供商获取模型列表并保存到数据库
	FetchAndSaveModels(ctx context.Context, llmProvider *models.LlmProvider) error

	// SaveApplicationLlm 保存应用模型信息
	// 如果ID为空则新增，否则更新现有记录
	SaveApplicationLlm(ctx context.Context, applicationLlm *models.ApplicationLlm) error

	// UpdateEnabledStatus 更新模型启用状态
	// 只更新 Enabled 字段
	UpdateEnabledStatus(ctx context.Context, id uuid.UUID, enabled bool) error

	// GetModelsByProviderID 根据提供商ID获取模型列表
	// 返回指定提供商下的所有模型
	GetModelsByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.ApplicationLlm, error)

	// GetModelsByApplicationID 根据应用ID获取模型列表
	// 返回指定应用下的所有模型
	GetModelsByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationLlm, error)
}

// applicationLlmService ApplicationLLM 业务逻辑层实现
// 实现 ApplicationLlmService 接口
type applicationLlmService struct {
	applicationLlmRepo repository.ApplicationLlmRepository // 数据访问层接口
}

// NewApplicationLlmService 创建 ApplicationLLM 服务实例
// 返回 ApplicationLlmService 接口的实现
// 参数：applicationLlmRepo - ApplicationLLM 数据访问层接口
func NewApplicationLlmService(applicationLlmRepo repository.ApplicationLlmRepository) ApplicationLlmService {
	return &applicationLlmService{
		applicationLlmRepo: applicationLlmRepo,
	}
}

// FetchAndSaveModels 获取并保存所有模型
// 从指定的 LLM 提供商获取模型列表并保存到数据库
func (s *applicationLlmService) FetchAndSaveModels(ctx context.Context, llmProvider *models.LlmProvider) error {
	// 检查提供商是否有必要的配置
	if llmProvider.ApiUrl == "" || llmProvider.ApiKey == "" {
		return fmt.Errorf("提供商缺少必要的配置信息")
	}

	// 根据提供商类型创建不同的客户端
	var openaiModels []openai.Model

	switch llmProvider.Type {
	case "openai_chat_completions_api", "openai_responses_api":
		// OpenAI 类型的提供商
		config := openai.DefaultConfig(llmProvider.ApiKey)
		if llmProvider.ApiUrl != "" {
			config.BaseURL = llmProvider.ApiUrl
		}
		client := openai.NewClientWithConfig(config)

		// 获取模型列表
		modelsList, err := client.ListModels(ctx)
		if err != nil {
			return fmt.Errorf("获取 OpenAI 模型列表失败: %w", err)
		}
		openaiModels = modelsList.Models

	case "ollama_api":
		// Ollama 类型的提供商
		// Ollama 通常使用本地 API，模型列表可能通过其他方式获取
		// 这里先返回空，后续可以根据 Ollama 的 API 特性进行扩展
		return fmt.Errorf("Ollama 提供商暂不支持自动获取模型列表")

	default:
		return fmt.Errorf("不支持的提供商类型: %s", llmProvider.Type)
	}

	// 获取该提供商下的现有模型记录
	existingModels, err := s.applicationLlmRepo.GetByProviderID(ctx, llmProvider.ID)
	if err != nil {
		return fmt.Errorf("获取现有模型记录失败: %w", err)
	}

	// 创建现有模型名称的映射，用于快速查找
	existingModelNames := make(map[string]bool)
	for _, existingModel := range existingModels {
		existingModelNames[existingModel.Name] = true
	}

	// 统计新增和跳过的模型数量
	addedCount := 0
	skippedCount := 0

	// 保存新的模型记录
	for _, model := range openaiModels {
		// 检查模型是否已存在
		if existingModelNames[model.ID] {
			skippedCount++
			continue // 跳过已存在的模型
		}

		applicationLlm := &models.ApplicationLlm{
			Name:          model.ID,
			Alias:         model.ID, // 默认使用模型ID作为别名
			ApplicationID: llmProvider.ApplicationID,
			LlmProviderID: llmProvider.ID,
			Enabled:       true, // 默认启用
			// 根据模型名称判断能力（这里可以根据实际需求进行调整）
			AbilityVision:         contains(model.ID, "vision") || contains(model.ID, "gpt-4-vision"),
			AbilityNetwork:        contains(model.ID, "gpt-4") || contains(model.ID, "gpt-3.5"),
			AbilityTextEmbeddings: contains(model.ID, "text-embedding") || contains(model.ID, "embedding"),
			AbilityThinking:       contains(model.ID, "gpt-4") || contains(model.ID, "gpt-3.5"),
			AbilityCallTools:      contains(model.ID, "gpt-4") || contains(model.ID, "gpt-3.5"),
			AbilityReranking:      contains(model.ID, "text-embedding-3"),
			// 设置默认计费信息
			BillingCurrency:    "USD",
			BillingPriceInput:  0.0015, // 默认价格，实际应该从配置或API获取
			BillingPriceOutput: 0.002,
		}

		if err := s.applicationLlmRepo.Create(ctx, applicationLlm); err != nil {
			return fmt.Errorf("保存模型记录失败: %w", err)
		}
		addedCount++
	}

	// 记录操作结果
	fmt.Printf("模型同步完成: 新增 %d 个模型, 跳过 %d 个已存在的模型\n", addedCount, skippedCount)

	return nil
}

// SaveApplicationLlm 保存应用模型信息
// 如果ID为空则新增，否则更新现有记录
func (s *applicationLlmService) SaveApplicationLlm(ctx context.Context, applicationLlm *models.ApplicationLlm) error {
	// 数据验证
	if err := s.validateApplicationLlm(applicationLlm); err != nil {
		return err
	}

	if applicationLlm.ID == uuid.Nil {
		// 新增：生成新的UUID
		applicationLlm.ID = uuid.New()
		return s.applicationLlmRepo.Create(ctx, applicationLlm)
	} else {
		// 更新：检查记录是否存在
		existing, err := s.applicationLlmRepo.GetByID(ctx, applicationLlm.ID)
		if err != nil {
			return fmt.Errorf("模型不存在: %w", err)
		}
		if existing == nil {
			return fmt.Errorf("模型不存在")
		}
		return s.applicationLlmRepo.Update(ctx, applicationLlm)
	}
}

// UpdateEnabledStatus 更新模型启用状态
// 只更新 Enabled 字段
func (s *applicationLlmService) UpdateEnabledStatus(ctx context.Context, id uuid.UUID, enabled bool) error {
	// 检查记录是否存在
	existing, err := s.applicationLlmRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("模型不存在: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("模型不存在")
	}

	// 只更新 Enabled 字段
	existing.Enabled = enabled
	return s.applicationLlmRepo.Update(ctx, existing)
}

// GetModelsByProviderID 根据提供商ID获取模型列表
// 返回指定提供商下的所有模型
func (s *applicationLlmService) GetModelsByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.ApplicationLlm, error) {
	return s.applicationLlmRepo.GetByProviderID(ctx, providerID)
}

// GetModelsByApplicationID 根据应用ID获取模型列表
// 返回指定应用下的所有模型
func (s *applicationLlmService) GetModelsByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*models.ApplicationLlm, error) {
	return s.applicationLlmRepo.GetByApplicationID(ctx, applicationID)
}

// validateApplicationLlm 验证应用模型数据
// 检查必填字段是否为空
func (s *applicationLlmService) validateApplicationLlm(applicationLlm *models.ApplicationLlm) error {
	if applicationLlm == nil {
		return fmt.Errorf("模型不能为空")
	}

	if applicationLlm.Name == "" {
		return fmt.Errorf("模型名称不能为空")
	}

	if applicationLlm.Alias == "" {
		return fmt.Errorf("模型别名不能为空")
	}

	if applicationLlm.ApplicationID == uuid.Nil {
		return fmt.Errorf("所属应用ID不能为空")
	}

	if applicationLlm.LlmProviderID == uuid.Nil {
		return fmt.Errorf("所属模型提供商ID不能为空")
	}

	return nil
}

// contains 检查字符串是否包含指定的子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

// containsSubstring 检查字符串中间是否包含子字符串
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
