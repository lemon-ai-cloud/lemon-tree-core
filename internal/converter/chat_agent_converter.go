// Package converter 提供数据转换功能
// 负责在不同层之间转换数据格式，如模型到DTO的转换
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
)

// ChatAgentModelToChatAgentDto 将模型转换为DTO
// 将数据库模型转换为用于API响应的DTO
// 参数：model - 数据库模型
// 返回：DTO对象
func ChatAgentModelToChatAgentDto(model *models.ChatAgent) dto.ChatAgentDto {
	return dto.ChatAgentDto{
		ID:                             model.ID.String(),
		Name:                           model.Name,
		Description:                    model.Description,
		ApplicationID:                  model.ApplicationID.String(),
		AvatarUrl:                      model.AvatarUrl,
		ChatSystemPrompt:               model.ChatSystemPrompt,
		ChatModelID:                    model.ChatModelID.String(),
		ConversationNamingPrompt:       model.ConversationNamingPrompt,
		ConversationNamingModelID:      model.ConversationNamingModelID.String(),
		ModelParamTemperature:          model.ModelParamTemperature,
		ModelParamTopP:                 model.ModelParamTopP,
		EnableContextLengthLimit:       model.EnableContextLengthLimit,
		ContextLengthLimit:             model.ContextLengthLimit,
		EnableMaxOutputTokenCountLimit: model.EnableMaxOutputTokenCountLimit,
		MaxOutputTokenCountLimit:       model.MaxOutputTokenCountLimit,
		DefaultStreamable:              model.DefaultStreamable,
		CreatedAt:                      model.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:                      model.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ChatAgentModelListToChatAgentDtoList 将模型列表转换为DTO列表
// 将数据库模型列表转换为用于API响应的DTO列表
// 参数：models - 数据库模型列表
// 返回：DTO列表
func ChatAgentModelListToChatAgentDtoList(models []*models.ChatAgent) []dto.ChatAgentDto {
	dtoList := make([]dto.ChatAgentDto, len(models))
	for i, model := range models {
		dtoList[i] = ChatAgentModelToChatAgentDto(model)
	}
	return dtoList
}

// SaveChatAgentRequestToChatAgentModel 将保存请求转换为模型
// 将前端保存请求转换为数据库模型
// 参数：request - 保存请求
// 返回：数据库模型
func SaveChatAgentRequestToChatAgentModel(request *dto.SaveChatAgentRequest) *models.ChatAgent {
	model := &models.ChatAgent{
		Name:                           request.Name,
		Description:                    request.Description,
		AvatarUrl:                      request.AvatarUrl,
		ChatSystemPrompt:               request.ChatSystemPrompt,
		ConversationNamingPrompt:       request.ConversationNamingPrompt,
		ModelParamTemperature:          request.ModelParamTemperature,
		ModelParamTopP:                 request.ModelParamTopP,
		EnableContextLengthLimit:       request.EnableContextLengthLimit,
		ContextLengthLimit:             request.ContextLengthLimit,
		EnableMaxOutputTokenCountLimit: request.EnableMaxOutputTokenCountLimit,
		MaxOutputTokenCountLimit:       request.MaxOutputTokenCountLimit,
		DefaultStreamable:              request.DefaultStreamable,
	}

	// 解析应用ID
	if applicationID, err := uuid.Parse(request.ApplicationID); err == nil {
		model.ApplicationID = applicationID
	}

	// 解析聊天模型ID
	if chatModelID, err := uuid.Parse(request.ChatModelID); err == nil {
		model.ChatModelID = chatModelID
	}

	// 解析会话命名模型ID
	if conversationNamingModelID, err := uuid.Parse(request.ConversationNamingModelID); err == nil {
		model.ConversationNamingModelID = conversationNamingModelID
	}

	// 如果有ID，则解析ID（用于更新操作）
	if request.ID != nil && *request.ID != "" {
		if id, err := uuid.Parse(*request.ID); err == nil {
			model.ID = id
		}
	}

	return model
}
