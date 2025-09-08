// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// ChatAgentConversationService 聊天会话 业务逻辑层接口
// 定义 聊天会话 相关的业务逻辑方法
type ChatAgentConversationService interface {
	// GetChatMessageList 获取聊天消息列表
	GetChatMessageList(ctx context.Context, chatAgentID, conversationID, lastID string, size int, sort string) ([]*models.ChatAgentMessage, error)

	// CreateConversation 创建会话
	CreateConversation(ctx context.Context, chatAgentID, serviceUserID, userMessage string) (*models.ChatAgentConversation, error)

	// GetConversationList 获取会话列表
	GetConversationList(ctx context.Context, chatAgentID, serviceUserID, lastID string, size int, sort string) ([]*models.ChatAgentConversation, error)

	// DeleteConversation 删除会话
	DeleteConversation(ctx context.Context, chatAgentID, serviceUserID, conversationID string) (*dto.DeleteConversationResponse, error)

	// UserSendMessagePredefinedAnswer 用户发送消息，返回固定答案
	UserSendMessagePredefinedAnswer(ctx context.Context, req *dto.ChatUserSendMessageRequest, streamable bool) (io.Reader, error)

	// UserSendMessage 用户发送消息
	UserSendMessage(ctx context.Context, req *dto.ChatUserSendMessageRequest, streamable bool) (io.Reader, error)

	// UploadAttachment 上传聊天附件
	UploadAttachment(ctx context.Context, chatAgentID string, file io.Reader, filename string, size int64) (*dto.UploadAttachmentResponse, error)

	// RenameConversationTitle 重命名会话标题
	RenameConversationTitle(ctx context.Context, chatAgentID, serviceUserID, conversationID, newTitle string) (*dto.RenameConversationResponse, error)
}

// chatAgentConversationService 聊天会话 业务逻辑层实现
// 实现 ChatAgentConversationService 接口
type chatAgentConversationService struct {
	db               *gorm.DB
	conversationRepo repository.ChatAgentConversationRepository
	messageRepo      repository.ChatAgentMessageRepository
	attachmentRepo   repository.ChatAgentAttachmentRepository
	chatAgentRepo    repository.ChatAgentRepository
	applicationRepo  repository.ApplicationRepository
	llmRepo          repository.ApplicationLlmRepository
	mcpConfigRepo    repository.ApplicationMcpServerConfigRepository
	mcpToolRepo      repository.ApplicationMcpServerToolRepository
}

// NewChatAgentConversationService 创建 聊天会话 服务实例
// 返回 ChatAgentConversationService 接口的实现
func NewChatAgentConversationService(
	db *gorm.DB,
	conversationRepo repository.ChatAgentConversationRepository,
	messageRepo repository.ChatAgentMessageRepository,
	attachmentRepo repository.ChatAgentAttachmentRepository,
	chatAgentRepo repository.ChatAgentRepository,
	applicationRepo repository.ApplicationRepository,
	llmRepo repository.ApplicationLlmRepository,
	mcpConfigRepo repository.ApplicationMcpServerConfigRepository,
	mcpToolRepo repository.ApplicationMcpServerToolRepository,
) ChatAgentConversationService {
	return &chatAgentConversationService{
		db:               db,
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		attachmentRepo:   attachmentRepo,
		chatAgentRepo:    chatAgentRepo,
		applicationRepo:  applicationRepo,
		llmRepo:          llmRepo,
		mcpConfigRepo:    mcpConfigRepo,
		mcpToolRepo:      mcpToolRepo,
	}
}

// GetChatMessageList 获取聊天消息列表
func (s *chatAgentConversationService) GetChatMessageList(ctx context.Context, chatAgentID, conversationID, lastID string, size int, sort string) ([]*models.ChatAgentMessage, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	appID, agentID, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}

	// 验证传入的chatAgentID是否与上下文中的一致
	requestedAgentID, err := uuid.Parse(chatAgentID)
	if err != nil {
		return nil, fmt.Errorf("无效的智能体ID: %w", err)
	}
	if requestedAgentID != agentID {
		return nil, fmt.Errorf("智能体ID不匹配")
	}

	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, fmt.Errorf("无效的会话ID: %w", err)
	}

	// 构建查询条件
	query := s.db.Where("application_id = ? AND chat_agent_id = ? AND conversation_id = ? AND deleted_at IS NULL", appID, agentID, convID)

	// 处理游标分页
	if lastID != "" {
		lastMsgID, err := uuid.Parse(lastID)
		if err != nil {
			return nil, fmt.Errorf("无效的last_id: %w", err)
		}

		// 根据排序方式确定游标条件
		switch sort {
		case "-created_at":
			// 倒序：查询创建时间小于last_id对应记录的时间
			var lastMessage models.ChatAgentMessage
			if err := s.db.Where("id = ?", lastMsgID).First(&lastMessage).Error; err == nil {
				query = query.Where("created_at < ?", lastMessage.CreatedAt)
			}
		case "created_at":
			// 正序：查询创建时间大于last_id对应记录的时间
			var lastMessage models.ChatAgentMessage
			if err := s.db.Where("id = ?", lastMsgID).First(&lastMessage).Error; err == nil {
				query = query.Where("created_at > ?", lastMessage.CreatedAt)
			}
		case "-id":
			// 按ID倒序
			query = query.Where("id < ?", lastMsgID)
		case "id":
			// 按ID正序
			query = query.Where("id > ?", lastMsgID)
		}
	}

	// 处理排序
	switch sort {
	case "-created_at":
		query = query.Order("created_at DESC")
	case "created_at":
		query = query.Order("created_at ASC")
	case "-id":
		query = query.Order("id DESC")
	case "id":
		query = query.Order("id ASC")
	default:
		// 默认按创建时间倒序
		query = query.Order("created_at DESC")
	}

	// 限制返回数量
	query = query.Limit(size)

	// 执行查询
	var messages []*models.ChatAgentMessage
	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %w", err)
	}

	return messages, nil
}

// CreateConversation 创建会话
func (s *chatAgentConversationService) CreateConversation(ctx context.Context, chatAgentID, serviceUserID, userMessage string) (*models.ChatAgentConversation, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	appID, agentID, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}

	// 验证传入的chatAgentID是否与上下文中的一致
	requestedAgentID, err := uuid.Parse(chatAgentID)
	if err != nil {
		return nil, fmt.Errorf("无效的智能体ID: %w", err)
	}
	if requestedAgentID != agentID {
		return nil, fmt.Errorf("智能体ID不匹配")
	}

	// 过滤消息，去除消息中无需保存的数据，比如qwen中的/no_think
	filteredUserMessage := userMessage
	if strings.HasPrefix(userMessage, "/no_think") {
		filteredUserMessage = userMessage[9:]
	}

	conversation := &models.ChatAgentConversation{
		Title:         filteredUserMessage,
		ApplicationID: appID,
		ChatAgentID:   agentID,
		ServiceUserID: serviceUserID,
	}

	if err := s.conversationRepo.Create(ctx, conversation); err != nil {
		return nil, fmt.Errorf("创建会话失败: %w", err)
	}

	return conversation, nil
}

// GetConversationList 获取会话列表
func (s *chatAgentConversationService) GetConversationList(ctx context.Context, chatAgentID, serviceUserID, lastID string, size int, sort string) ([]*models.ChatAgentConversation, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	appID, agentID, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}

	// 验证传入的chatAgentID是否与上下文中的一致
	requestedAgentID, err := uuid.Parse(chatAgentID)
	if err != nil {
		return nil, fmt.Errorf("无效的智能体ID: %w", err)
	}
	if requestedAgentID != agentID {
		return nil, fmt.Errorf("智能体ID不匹配")
	}

	// 构建查询条件
	query := s.db.Where("application_id = ? AND chat_agent_id = ? AND service_user_id = ? AND deleted_at IS NULL", appID, agentID, serviceUserID)

	// 处理游标分页
	if lastID != "" {
		lastConvID, err := uuid.Parse(lastID)
		if err != nil {
			return nil, fmt.Errorf("无效的last_id: %w", err)
		}

		// 根据排序方式确定游标条件
		switch sort {
		case "-created_at":
			// 倒序：查询创建时间小于last_id对应记录的时间
			var lastConversation models.ChatAgentConversation
			if err := s.db.Where("id = ?", lastConvID).First(&lastConversation).Error; err == nil {
				query = query.Where("created_at < ?", lastConversation.CreatedAt)
			}
		case "created_at":
			// 正序：查询创建时间大于last_id对应记录的时间
			var lastConversation models.ChatAgentConversation
			if err := s.db.Where("id = ?", lastConvID).First(&lastConversation).Error; err == nil {
				query = query.Where("created_at > ?", lastConversation.CreatedAt)
			}
		case "-id":
			// 按ID倒序
			query = query.Where("id < ?", lastConvID)
		case "id":
			// 按ID正序
			query = query.Where("id > ?", lastConvID)
		}
	}

	// 处理排序
	switch sort {
	case "-created_at":
		query = query.Order("created_at DESC")
	case "created_at":
		query = query.Order("created_at ASC")
	case "-id":
		query = query.Order("id DESC")
	case "id":
		query = query.Order("id ASC")
	default:
		// 默认按创建时间倒序
		query = query.Order("created_at DESC")
	}

	// 限制返回数量
	query = query.Limit(size)

	// 执行查询
	var conversations []*models.ChatAgentConversation
	if err := query.Find(&conversations).Error; err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}

	return conversations, nil
}

// DeleteConversation 删除会话
func (s *chatAgentConversationService) DeleteConversation(ctx context.Context, chatAgentID, serviceUserID, conversationID string) (*dto.DeleteConversationResponse, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	appID, agentID, err := getContextInfo(ctx)
	if err != nil {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("获取上下文信息失败: %v", err)),
		}, nil
	}

	// 验证传入的chatAgentID是否与上下文中的一致
	requestedAgentID, err := uuid.Parse(chatAgentID)
	if err != nil {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr("无效的智能体ID"),
		}, nil
	}
	if requestedAgentID != agentID {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr("智能体ID不匹配"),
		}, nil
	}

	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr("无效的会话ID"),
		}, nil
	}

	// 1. 验证会话是否存在且属于该用户和智能体
	conversation, err := s.conversationRepo.GetByID(ctx, convID)
	if err != nil {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr("会话不存在"),
		}, nil
	}

	if conversation.ApplicationID != appID || conversation.ChatAgentID != agentID || conversation.ServiceUserID != serviceUserID {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr("无权删除此会话"),
		}, nil
	}

	// 2. 删除会话相关的所有消息
	var messages []*models.ChatAgentMessage
	if err := s.db.Where("application_id = ? AND chat_agent_id = ? AND conversation_id = ?", appID, agentID, convID).Find(&messages).Error; err != nil {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("查询消息失败: %v", err)),
		}, nil
	}

	// 收集需要删除附件的消息ID
	var messageIDs []uuid.UUID
	for _, msg := range messages {
		messageIDs = append(messageIDs, msg.ID)
	}

	// 删除消息
	for _, message := range messages {
		if err := s.messageRepo.DeleteByID(ctx, message.ID); err != nil {
			log.Printf("删除消息失败: %v", err)
		}
	}

	// 3. 删除会话相关的所有附件
	var attachments []*models.ChatAgentAttachment
	if len(messageIDs) > 0 {
		if err := s.db.Where("application_id = ? AND chat_agent_id = ? AND message_id IN ?", appID, agentID, messageIDs).Find(&attachments).Error; err != nil {
			log.Printf("查询附件失败: %v", err)
		} else {
			// 删除附件文件和目录
			for _, attachment := range attachments {
				// 删除附件文件目录
				if attachment.FilePath != "" {
					attachmentDir := filepath.Dir(attachment.FilePath)
					if _, err := os.Stat(attachmentDir); err == nil {
						os.RemoveAll(attachmentDir)
					}
				}

				// 删除数据库记录
				if err := s.attachmentRepo.DeleteByID(ctx, attachment.ID); err != nil {
					log.Printf("删除附件记录失败: %v", err)
				}
			}
		}
	}

	// 4. 删除会话本身
	if err := s.conversationRepo.DeleteByID(ctx, convID); err != nil {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("删除会话失败: %v", err)),
		}, nil
	}

	return &dto.DeleteConversationResponse{
		Success:                 true,
		Message:                 stringPtr("会话删除成功"),
		DeletedMessagesCount:    intPtr(len(messages)),
		DeletedAttachmentsCount: intPtr(len(attachments)),
	}, nil
}

// UserSendMessagePredefinedAnswer 用户发送消息，返回固定答案
func (s *chatAgentConversationService) UserSendMessagePredefinedAnswer(ctx context.Context, req *dto.ChatUserSendMessageRequest, streamable bool) (io.Reader, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	appID, agentID, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}

	// 如果conversation_id为空，则认为是新的会话
	var conversation *models.ChatAgentConversation
	var conversationIDStr string
	var historyMessages []openai.ChatCompletionMessage

	if req.ConversationID == nil {
		// 创建新会话
		conversation, err = s.CreateConversation(ctx, agentID.String(), req.ServiceUserID, req.UserMessage)
		if err != nil {
			return nil, fmt.Errorf("创建会话失败: %w", err)
		}
		conversationIDStr = conversation.ID.String()
	} else {
		// 根据会话id进行查询，如果找不到那么就创建一个新的会话
		convID, err := uuid.Parse(*req.ConversationID)
		if err != nil {
			return nil, fmt.Errorf("无效的会话ID: %w", err)
		}

		conversation, err = s.conversationRepo.GetByID(ctx, convID)
		if err != nil || conversation == nil {
			// 创建新会话
			conversation, err = s.CreateConversation(ctx, agentID.String(), req.ServiceUserID, req.UserMessage)
			if err != nil {
				return nil, fmt.Errorf("创建会话失败: %w", err)
			}
			conversationIDStr = conversation.ID.String()
		} else {
			conversationIDStr = *req.ConversationID
			// 是历史会话，查询历史消息
			messageList, err := s.GetChatMessageList(ctx, agentID.String(), conversationIDStr, "", 100, "created_at")
			if err != nil {
				return nil, fmt.Errorf("获取历史消息失败: %w", err)
			}

			// 构建历史消息列表（为将来的AI对话功能预留）
			historyMessages = make([]openai.ChatCompletionMessage, 0, len(messageList))
			for _, messageItem := range messageList {
				historyMessages = append(historyMessages, openai.ChatCompletionMessage{
					Role:    messageItem.Role,
					Content: messageItem.Content,
				})
			}
			// 暂时使用historyMessages变量以避免编译器警告（将来AI对话功能会使用）
			_ = len(historyMessages)
		}
	}

	if conversation == nil {
		return nil, fmt.Errorf("会话创建失败")
	}

	// 生成请求id
	requestID := uuid.New().String()

	// 将用户的消息存储到数据库
	userMessageObj := &models.ChatAgentMessage{
		ApplicationID:  appID,
		ChatAgentID:    agentID,
		ConversationID: conversation.ID,
		RequestID:      requestID,
		Type:           "message",
		Role:           "user",
		Content:        req.UserMessage,
	}
	if err := s.messageRepo.Create(ctx, userMessageObj); err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %w", err)
	}

	// 将预制的系统回复答案消息存储到数据库
	assistantMessageObj := &models.ChatAgentMessage{
		ApplicationID:  appID,
		ChatAgentID:    agentID,
		ConversationID: conversation.ID,
		RequestID:      requestID,
		Type:           "message",
		Role:           "assistant",
		Content:        *req.PredefinedAnswer,
	}
	if err := s.messageRepo.Create(ctx, assistantMessageObj); err != nil {
		return nil, fmt.Errorf("保存助手消息失败: %w", err)
	}

	// 创建流式响应
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()

		// 将固定的答案用streamable流的形式逐字返回给用户
		if streamable {
			// 流式返回，逐字返回
			for _, char := range *req.PredefinedAnswer {
				event := dto.ChatMessageResponseEventDto{
					ConversationID: conversationIDStr,
					RequestID:      requestID,
					MessageType:    "answer_delta",
					Content:        string(char),
				}
				eventJSON, _ := json.Marshal(event)
				pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
			}
			// 最后返回完整答案
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationIDStr,
				RequestID:      requestID,
				MessageType:    "answer",
				Content:        *req.PredefinedAnswer,
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
		} else {
			// 非流式返回，直接返回完整答案
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationIDStr,
				RequestID:      requestID,
				MessageType:    "answer",
				Content:        *req.PredefinedAnswer,
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
		}
	}()

	return pr, nil
}

// UserSendMessage 用户发送消息
func (s *chatAgentConversationService) UserSendMessage(ctx context.Context, req *dto.ChatUserSendMessageRequest, streamable bool) (io.Reader, error) {
	// 这里需要实现完整的AI对话逻辑
	// 包括工具调用、流式响应等
	// 由于代码较长，这里先返回一个简单的实现
	return s.UserSendMessagePredefinedAnswer(ctx, req, streamable)
}

// UploadAttachment 上传聊天附件
func (s *chatAgentConversationService) UploadAttachment(ctx context.Context, chatAgentID string, file io.Reader, filename string, size int64) (*dto.UploadAttachmentResponse, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	appID, agentID, err := getContextInfo(ctx)
	if err != nil {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("获取上下文信息失败: %v", err)),
		}, nil
	}

	// 验证传入的chatAgentID是否与上下文中的一致
	requestedAgentID, err := uuid.Parse(chatAgentID)
	if err != nil {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr("无效的智能体ID"),
		}, nil
	}
	if requestedAgentID != agentID {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr("智能体ID不匹配"),
		}, nil
	}

	// 检查文件大小（限制为50MB）
	maxFileSize := int64(50 * 1024 * 1024) // 50MB
	if size > maxFileSize {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr("文件大小超过限制（最大50MB）"),
		}, nil
	}

	// 获取文件扩展名
	fileExtension := strings.ToLower(filepath.Ext(filename))
	if fileExtension == "" {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr("文件名不能为空"),
		}, nil
	}

	// 生成附件ID
	attachmentID := uuid.New()

	// 创建存储目录
	baseStorageDir := "chat_attachment_files"
	attachmentDir := filepath.Join(baseStorageDir, attachmentID.String())
	if err := os.MkdirAll(attachmentDir, 0755); err != nil {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("创建存储目录失败: %v", err)),
		}, nil
	}

	// 保存原始文件
	filePath := filepath.Join(attachmentDir, fmt.Sprintf("file%s", fileExtension))
	dst, err := os.Create(filePath)
	if err != nil {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("创建文件失败: %v", err)),
		}, nil
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("保存文件失败: %v", err)),
		}, nil
	}

	// 确定附件类型
	attachmentType := "other"
	if isDocumentFile(fileExtension) {
		attachmentType = "document"
	} else if isImageFile(fileExtension) {
		attachmentType = "image"
	}

	// 创建附件记录
	attachment := &models.ChatAgentAttachment{
		ApplicationID:    appID,
		ChatAgentID:      agentID,
		OriginalFileName: filename,
		FileExtension:    fileExtension,
		FileSize:         size,
		MimeType:         getMimeType(fileExtension),
		FilePath:         filePath,
		AttachmentType:   attachmentType,
		IsProcessed:      false,
		ProcessingError:  "",
	}

	// 保存到数据库
	if err := s.attachmentRepo.Create(ctx, attachment); err != nil {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("保存附件记录失败: %v", err)),
		}, nil
	}

	return &dto.UploadAttachmentResponse{
		Success:          true,
		AttachmentID:     stringPtr(attachment.ID.String()),
		OriginalFileName: stringPtr(filename),
		FileSize:         &size,
		AttachmentType:   stringPtr(attachmentType),
		IsProcessed:      boolPtr(false),
	}, nil
}

// RenameConversationTitle 重命名会话标题
func (s *chatAgentConversationService) RenameConversationTitle(ctx context.Context, chatAgentID, serviceUserID, conversationID, newTitle string) (*dto.RenameConversationResponse, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	appID, agentID, err := getContextInfo(ctx)
	if err != nil {
		return &dto.RenameConversationResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("获取上下文信息失败: %v", err)),
		}, nil
	}

	// 验证传入的chatAgentID是否与上下文中的一致
	requestedAgentID, err := uuid.Parse(chatAgentID)
	if err != nil {
		return &dto.RenameConversationResponse{
			Success: false,
			Error:   stringPtr("无效的智能体ID"),
		}, nil
	}
	if requestedAgentID != agentID {
		return &dto.RenameConversationResponse{
			Success: false,
			Error:   stringPtr("智能体ID不匹配"),
		}, nil
	}

	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return &dto.RenameConversationResponse{
			Success: false,
			Error:   stringPtr("无效的会话ID"),
		}, nil
	}

	// 1. 验证会话是否存在且属于该用户和智能体
	conversation, err := s.conversationRepo.GetByID(ctx, convID)
	if err != nil {
		return &dto.RenameConversationResponse{
			Success: false,
			Error:   stringPtr("会话不存在"),
		}, nil
	}

	if conversation.ApplicationID != appID || conversation.ChatAgentID != agentID || conversation.ServiceUserID != serviceUserID {
		return &dto.RenameConversationResponse{
			Success: false,
			Error:   stringPtr("无权重命名此会话"),
		}, nil
	}

	// 2. 更新会话标题
	conversation.Title = newTitle
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return &dto.RenameConversationResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("重命名会话失败: %v", err)),
		}, nil
	}

	return &dto.RenameConversationResponse{
		Success:        true,
		Message:        stringPtr("会话重命名成功"),
		ConversationID: stringPtr(conversationID),
		NewTitle:       stringPtr(newTitle),
	}, nil
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

// getContextInfo 从上下文中获取ApplicationID和ChatAgentID
func getContextInfo(ctx context.Context) (applicationID, chatAgentID uuid.UUID, err error) {
	// 从上下文中获取ChatAgent
	chatAgentValue := ctx.Value(define.AppContextKeyCurrentChatAgent)
	if chatAgentValue == nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("上下文中未找到ChatAgent信息")
	}

	chatAgent, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		return uuid.Nil, uuid.Nil, fmt.Errorf("ChatAgent类型转换失败")
	}

	// 从上下文中获取Application
	applicationValue := ctx.Value(define.AppContextKeyCurrentApplication)
	if applicationValue == nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("上下文中未找到Application信息")
	}

	application, ok := applicationValue.(*models.Application)
	if !ok {
		return uuid.Nil, uuid.Nil, fmt.Errorf("Application类型转换失败")
	}

	return application.ID, chatAgent.ID, nil
}

func isDocumentFile(ext string) bool {
	documentExts := []string{".doc", ".docx", ".pdf", ".txt", ".md", ".xls", ".xlsx", ".ppt", ".pptx"}
	for _, docExt := range documentExts {
		if ext == docExt {
			return true
		}
	}
	return false
}

func isImageFile(ext string) bool {
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

func getMimeType(ext string) string {
	mimeTypes := map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".md":   "text/markdown",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".webp": "image/webp",
	}
	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	return "application/octet-stream"
}
