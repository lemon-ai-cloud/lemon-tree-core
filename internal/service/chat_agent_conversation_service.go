// Package service 提供业务逻辑层功能
// 负责处理业务逻辑、数据验证、调用数据访问层和返回业务结果
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lemon-tree-core/internal/al_client"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/manager"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"
	"lemon-tree-core/internal/utils"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"gorm.io/gorm"
)

// ChatAgentConversationService 聊天会话 业务逻辑层接口
// 定义 聊天会话 相关的业务逻辑方法
type ChatAgentConversationService interface {
	// GetChatMessageList 获取聊天消息列表
	GetChatMessageList(ctx context.Context, conversationID, lastID string, size int) ([]*models.ChatAgentMessage, error)

	// CreateConversation 创建会话
	CreateConversation(ctx context.Context, serviceUserID, userMessage string) (*models.ChatAgentConversation, error)

	// GetConversationList 获取会话列表
	GetConversationList(ctx context.Context, serviceUserID, lastID string, size int) ([]*models.ChatAgentConversation, error)

	// DeleteConversation 删除会话
	DeleteConversation(ctx context.Context, serviceUserID, conversationID string) (*dto.DeleteConversationResponse, error)

	// UserSendMessagePredefinedAnswer 用户发送消息，返回固定答案
	UserSendMessagePredefinedAnswer(ctx context.Context, req *dto.ChatUserSendMessageRequest, streamable bool) (io.Reader, error)

	// UserSendMessage 用户发送消息
	UserSendMessage(ctx context.Context, req *dto.ChatUserSendMessageRequest, streamable bool) (io.Reader, error)

	// UploadAttachment 上传聊天附件
	UploadAttachment(ctx context.Context, file io.Reader, filename string, size int64) (*dto.UploadAttachmentResponse, error)

	// RenameConversationTitle 重命名会话标题
	RenameConversationTitle(ctx context.Context, serviceUserID, conversationID, newTitle string) (*dto.RenameConversationResponse, error)

	// GetChatAgentMcpServerTools 获取聊天智能体启用的MCP工具列表
	// 根据chatAgentID查询启用的工具，并从MCP服务器获取最新的工具信息
	GetChatAgentMcpServerTools(ctx context.Context) ([]al_client.Tool, error)
}

// chatAgentConversationService 聊天会话 业务逻辑层实现
// 实现 ChatAgentConversationService 接口
type chatAgentConversationService struct {
	db                         *gorm.DB
	conversationRepo           repository.ChatAgentConversationRepository
	messageRepo                repository.ChatAgentMessageRepository
	attachmentRepo             repository.ChatAgentAttachmentRepository
	chatAgentRepo              repository.ChatAgentRepository
	applicationRepo            repository.ApplicationRepository
	llmRepo                    repository.ApplicationLlmRepository
	mcpConfigRepo              repository.ApplicationMcpServerConfigRepository
	mcpToolRepo                repository.ApplicationMcpServerToolRepository
	chatAgentMcpServerToolRepo repository.ChatAgentMcpServerToolRepository
	llmProviderRepo            repository.LlmProviderRepository
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
	chatAgentMcpServerToolRepo repository.ChatAgentMcpServerToolRepository,
	llmProviderRepo repository.LlmProviderRepository,
) ChatAgentConversationService {
	return &chatAgentConversationService{
		db:                         db,
		conversationRepo:           conversationRepo,
		messageRepo:                messageRepo,
		attachmentRepo:             attachmentRepo,
		chatAgentRepo:              chatAgentRepo,
		applicationRepo:            applicationRepo,
		llmRepo:                    llmRepo,
		mcpConfigRepo:              mcpConfigRepo,
		mcpToolRepo:                mcpToolRepo,
		chatAgentMcpServerToolRepo: chatAgentMcpServerToolRepo,
		llmProviderRepo:            llmProviderRepo,
	}
}

// GetChatMessageList 获取聊天消息列表
func (s *chatAgentConversationService) GetChatMessageList(ctx context.Context, conversationID, lastID string, size int) ([]*models.ChatAgentMessage, error) {
	_, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("无效的智能体ID: %w", err)
	}

	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, fmt.Errorf("无效的会话ID: %w", err)
	}

	// 构建查询条件
	query := s.db.Where("chat_agent_id = ? AND conversation_id = ? AND deleted_at IS NULL", chatAgent.ID, convID)

	// 处理游标分页
	if lastID != "" {
		lastMsgID, err := uuid.Parse(lastID)
		if err != nil {
			return nil, fmt.Errorf("无效的last_id: %w", err)
		}

		// 获取lastID对应消息的创建时间
		var lastMessage models.ChatAgentMessage
		if err := s.db.Where("id = ?", lastMsgID).First(&lastMessage).Error; err == nil {
			// 从该消息的创建时间往前获取更早的消息
			query = query.Where("created_at < ?", lastMessage.CreatedAt)
		}
	}

	// 按创建时间倒序排列
	query = query.Order("created_at DESC")

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
func (s *chatAgentConversationService) CreateConversation(ctx context.Context, serviceUserID, userMessage string) (*models.ChatAgentConversation, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	application, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}

	// 过滤消息，去除消息中无需保存的数据，比如qwen中的/no_think
	filteredUserMessage := userMessage
	if strings.HasPrefix(userMessage, "/no_think") {
		filteredUserMessage = userMessage[9:]
	}

	conversation := &models.ChatAgentConversation{
		Title:         filteredUserMessage,
		ApplicationID: application.ID,
		ChatAgentID:   chatAgent.ID,
		ServiceUserID: serviceUserID,
	}

	if err := s.conversationRepo.Create(ctx, conversation); err != nil {
		return nil, fmt.Errorf("创建会话失败: %w", err)
	}

	return conversation, nil
}

// GetConversationList 获取会话列表
func (s *chatAgentConversationService) GetConversationList(ctx context.Context, serviceUserID, lastID string, size int) ([]*models.ChatAgentConversation, error) {
	_, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("无效的智能体ID: %w", err)
	}

	// 构建查询条件
	query := s.db.Where("chat_agent_id = ? AND service_user_id = ? AND deleted_at IS NULL", chatAgent.ID, serviceUserID)

	// 处理游标分页
	if lastID != "" {
		lastConvID, err := uuid.Parse(lastID)
		if err != nil {
			return nil, fmt.Errorf("无效的last_id: %w", err)
		}

		// 获取lastID对应会话的创建时间
		var lastConversation models.ChatAgentConversation
		if err := s.db.Where("id = ?", lastConvID).First(&lastConversation).Error; err == nil {
			// 从该会话的创建时间往前获取更早的会话
			query = query.Where("created_at < ?", lastConversation.CreatedAt)
		}
	}

	// 按创建时间倒序排列
	query = query.Order("created_at DESC")

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
func (s *chatAgentConversationService) DeleteConversation(ctx context.Context, serviceUserID, conversationID string) (*dto.DeleteConversationResponse, error) {
	_, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr("无效的智能体ID"),
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

	if conversation.ChatAgentID != chatAgent.ID || conversation.ServiceUserID != serviceUserID {
		return &dto.DeleteConversationResponse{
			Success: false,
			Error:   stringPtr("无权删除此会话"),
		}, nil
	}

	// 2. 删除会话相关的所有消息
	var messages []*models.ChatAgentMessage
	if err := s.db.Where("chat_agent_id = ? AND conversation_id = ?", chatAgent.ID, convID).Find(&messages).Error; err != nil {
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
		if err := s.db.Where("chat_agent_id = ? AND message_id IN ?", chatAgent.ID, messageIDs).Find(&attachments).Error; err != nil {
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
	application, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}

	// 如果conversation_id为空，则认为是新的会话
	var conversation *models.ChatAgentConversation
	var conversationIDStr string
	var historyMessages []al_client.ChatMessage

	if req.ConversationID == nil {
		// 创建新会话
		conversation, err = s.CreateConversation(ctx, req.ServiceUserID, req.UserMessage)
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
			conversation, err = s.CreateConversation(ctx, req.ServiceUserID, req.UserMessage)
			if err != nil {
				return nil, fmt.Errorf("创建会话失败: %w", err)
			}
			conversationIDStr = conversation.ID.String()
		} else {
			conversationIDStr = *req.ConversationID
			// 是历史会话，查询历史消息
			messageList, err := s.GetChatMessageList(ctx, conversationIDStr, "", 100)
			if err != nil {
				return nil, fmt.Errorf("获取历史消息失败: %w", err)
			}

			// 构建历史消息列表（为将来的AI对话功能预留）
			historyMessages = make([]al_client.ChatMessage, 0, len(messageList))
			for _, messageItem := range messageList {
				// 只处理普通消息类型，跳过函数调用相关消息
				if messageItem.Type == "message" && messageItem.Role != "" {
					historyMessages = append(historyMessages, al_client.ChatMessage{
						Role:    messageItem.Role,
						Content: messageItem.Content,
					})
				}
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
		ApplicationID:  application.ID,
		ChatAgentID:    chatAgent.ID,
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
		ApplicationID:  application.ID,
		ChatAgentID:    chatAgent.ID,
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
	// 从上下文中获取ApplicationID和ChatAgentID
	application, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}

	// 如果conversation_id为空，则认为是新的会话
	var conversation *models.ChatAgentConversation
	var conversationIDStr string
	var historyMessages []al_client.ChatMessage

	if req.ConversationID == nil || *req.ConversationID == "" {
		// 创建新会话
		conversation, err = s.CreateConversation(ctx, req.ServiceUserID, req.UserMessage)
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
			conversation, err = s.CreateConversation(ctx, req.ServiceUserID, req.UserMessage)
			if err != nil {
				return nil, fmt.Errorf("创建会话失败: %w", err)
			}
			conversationIDStr = conversation.ID.String()
		} else {
			conversationIDStr = *req.ConversationID
			// 是历史会话，查询历史消息
			messageList, err := s.GetChatMessageList(ctx, conversationIDStr, "", 100)
			if err != nil {
				return nil, fmt.Errorf("获取历史消息失败: %w", err)
			}

			// 构建历史消息列表
			historyMessages = make([]al_client.ChatMessage, 0, len(messageList))
			for _, messageItem := range messageList {
				// 只处理普通消息类型，跳过函数调用相关消息
				if messageItem.Type == "message" && messageItem.Role != "" {
					historyMessages = append(historyMessages, al_client.ChatMessage{
						Role:    messageItem.Role,
						Content: messageItem.Content,
					})
				}
			}
		}
	}

	if conversation == nil {
		return nil, fmt.Errorf("会话创建失败")
	}

	// 准备工具列表
	openaiToolsList, err := s.prepareToolsList(ctx, req.UsedMcpToolList, req.UsedInternalToolList)
	if err != nil {
		log.Printf("获取工具列表失败: %v", err)
		// 工具获取失败不影响主流程，使用空工具列表
		openaiToolsList = []al_client.Tool{}
	}

	// 生成请求id
	requestID := uuid.New().String()

	// 将用户的消息存储到数据库
	userMessageObj := &models.ChatAgentMessage{
		ApplicationID:  application.ID,
		ChatAgentID:    chatAgent.ID,
		ConversationID: conversation.ID,
		RequestID:      requestID,
		Type:           "message",
		Role:           "user",
		Content:        req.UserMessage,
	}
	if err := s.messageRepo.Create(ctx, userMessageObj); err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %w", err)
	}

	// 处理附件
	attachmentsPrompt := ""
	if len(req.Attachments) > 0 {
		attachmentInfoList, err := s.processMessageAttachments(ctx, req.Attachments, userMessageObj)
		if err != nil {
			log.Printf("处理附件失败: %v", err)
		} else {
			// 更新消息的附件信息
			attachmentInfoJSON, _ := json.Marshal(attachmentInfoList)
			userMessageObj.AttachmentsInfo = string(attachmentInfoJSON)
			if err := s.messageRepo.Update(ctx, userMessageObj); err != nil {
				log.Printf("更新消息附件信息失败: %v", err)
			}
			attachmentsPrompt = "用户上传的附件文件ID数组：" + string(attachmentInfoJSON) + "\n\n"
		}
	}

	// 构建完整的消息列表
	messages := make([]al_client.ChatMessage, 0, len(historyMessages)+2)

	// 添加系统提示词
	messages = append(messages, al_client.ChatMessage{
		Role:    "system",
		Content: req.SystemPrompt,
	})

	// 添加历史消息
	messages = append(messages, historyMessages...)

	// 添加当前用户消息（包含附件信息）
	messages = append(messages, al_client.ChatMessage{
		Role:    "user",
		Content: attachmentsPrompt + req.UserMessage,
	})

	log.Printf("开始处理消息，请求id:%s, 请求工具：%v, 工具列表数量：%d", requestID, req.UsedMcpToolList, len(openaiToolsList))

	// 交给AI处理消息
	if streamable {
		return s.aiProcessStreamable(ctx, conversationIDStr, requestID, messages, openaiToolsList) //openaiToolsList)
	} else {
		return s.aiProcess(ctx, conversationIDStr, requestID, messages, openaiToolsList)
	}
}

// UploadAttachment 上传聊天附件
func (s *chatAgentConversationService) UploadAttachment(ctx context.Context, file io.Reader, filename string, size int64) (*dto.UploadAttachmentResponse, error) {
	// 从上下文中获取ApplicationID和ChatAgentID
	application, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return &dto.UploadAttachmentResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("获取上下文信息失败: %v", err)),
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
		ApplicationID:    application.ID,
		ChatAgentID:      chatAgent.ID,
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
func (s *chatAgentConversationService) RenameConversationTitle(ctx context.Context, serviceUserID, conversationID, newTitle string) (*dto.RenameConversationResponse, error) {
	_, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return &dto.RenameConversationResponse{
			Success: false,
			Error:   stringPtr("无效的智能体ID"),
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

	if conversation.ChatAgentID != chatAgent.ID || conversation.ServiceUserID != serviceUserID {
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
func getContextInfo(ctx context.Context) (*models.Application, *models.ChatAgent, error) {
	// 从上下文中获取ChatAgent
	chatAgentValue := ctx.Value(define.AppContextKeyCurrentChatAgent)
	if chatAgentValue == nil {
		return nil, nil, fmt.Errorf("上下文中未找到ChatAgent信息")
	}

	chatAgent, ok := chatAgentValue.(*models.ChatAgent)
	if !ok {
		return nil, nil, fmt.Errorf("ChatAgent类型转换失败")
	}

	// 从上下文中获取Application
	applicationValue := ctx.Value(define.AppContextKeyCurrentApplication)
	if applicationValue == nil {
		return nil, nil, fmt.Errorf("上下文中未找到Application信息")
	}

	application, ok := applicationValue.(*models.Application)
	if !ok {
		return nil, nil, fmt.Errorf("Application类型转换失败")
	}

	return application, chatAgent, nil
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

// prepareToolsList 准备工具列表
func (s *chatAgentConversationService) prepareToolsList(ctx context.Context, usedMcpToolList []dto.ChatMessageUseToolDto, usedInternalToolList []string) ([]al_client.Tool, error) {
	//_, chatAgent, err := getContextInfo(ctx)
	var openaiToolsList []al_client.Tool

	// 处理MCP工具
	mcpTools, err := s.GetChatAgentMcpServerTools(ctx)
	if err != nil {
		return nil, err
	}
	openaiToolsList = append(openaiToolsList, mcpTools...)

	// 处理内部工具
	for _, internalToolName := range usedInternalToolList {
		// 这里需要根据实际的内部工具实现来构建工具定义
		// 暂时使用简单的实现
		openaiTool := al_client.Tool{
			Type: "function",
			Function: &al_client.FunctionDefinition{
				Name:        fmt.Sprintf("__lai__%s", internalToolName),
				Description: fmt.Sprintf("内部工具: %s", internalToolName),
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "查询参数",
						},
					},
					"required": []string{"query"},
				},
			},
		}
		openaiToolsList = append(openaiToolsList, openaiTool)
	}

	return openaiToolsList, nil
}

// processMessageAttachments 处理消息附件
func (s *chatAgentConversationService) processMessageAttachments(ctx context.Context, attachments []string, message *models.ChatAgentMessage) ([]dto.ChatMessageAttachmentInfoDto, error) {
	var attachmentInfoList []dto.ChatMessageAttachmentInfoDto

	for _, attachmentID := range attachments {
		attachmentUUID, err := uuid.Parse(attachmentID)
		if err != nil {
			log.Printf("无效的附件ID: %s", attachmentID)
			continue
		}

		attachment, err := s.attachmentRepo.GetByID(ctx, attachmentUUID)
		if err != nil {
			log.Printf("获取附件失败: %v", err)
			continue
		}

		if attachment != nil {
			attachmentInfoList = append(attachmentInfoList, dto.ChatMessageAttachmentInfoDto{
				ID:   attachmentID,
				Name: attachment.OriginalFileName,
			})

			// 更新附件的消息ID和会话ID
			attachment.MessageID = message.ID
			attachment.ConversationID = message.ConversationID
			if err := s.attachmentRepo.Update(ctx, attachment); err != nil {
				log.Printf("更新附件关联信息失败: %v", err)
			}
		}
	}

	return attachmentInfoList, nil
}

// aiProcessStreamable 处理AI消息 - 流式调用AI
func (s *chatAgentConversationService) aiProcessStreamable(ctx context.Context, conversationID, requestID string, messages []al_client.ChatMessage, aiTools []al_client.Tool) (io.Reader, error) {
	// 从上下文中获取ChatAgentID
	application, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}
	// 创建流式响应管道
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		// 获取应用配置
		llmProvider, llm, err := s.getChatAgentChatLlmConfig(ctx)
		if err != nil {
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "error",
				Content:        fmt.Sprintf("获取应用配置失败: %v", err),
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
			return
		}

		// 动态创建AI客户端
		aiClient, err := s.createAIClient(llmProvider)
		if err != nil {
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "error",
				Content:        fmt.Sprintf("创建AI客户端失败: %v", err),
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
			return
		}

		// 构建请求
		req := al_client.SendMessageRequest{
			Model:       llm.Name,
			Messages:    messages,
			Stream:      true,
			Tools:       aiTools,
			Temperature: chatAgent.ModelParamTemperature,
			TopP:        chatAgent.ModelParamTopP,
			ToolChoice:  "auto",
		}

		// 创建流式请求
		stream, err := aiClient.SendMessageStream(ctx, req)
		if err != nil {
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "error",
				Content:        fmt.Sprintf("AI Process error: %v", err),
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
			return
		}
		defer stream.Close()

		isNeedAiProcessContinue := false
		finalToolCalls := make(map[string]al_client.ToolCall)
		currentToolCall := al_client.ToolCall{}
		currentToolCallID := ""
		answerFullContent := ""

		// 处理AI的返回数据流
		for {
			chunk, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("处理流式数据时出错: %v", err)
				continue
			}

			if len(chunk.Choices) == 0 {
				continue
			}

			choice := chunk.Choices[0]

			// 处理工具调用
			if choice.Delta.ToolCalls != nil {
				for _, toolCall := range choice.Delta.ToolCalls {
					if toolCall.ID != "" {
						currentToolCallID = toolCall.ID
						currentToolCall = al_client.ToolCall{
							ID:   toolCall.ID,
							Type: toolCall.Type,
							Function: al_client.FunctionCall{
								Name:      toolCall.Function.Name,
								Arguments: toolCall.Function.Arguments,
							},
						}
						finalToolCalls[toolCall.ID] = currentToolCall
					} else {
						// 继续构建工具调用
						if currentToolCallID != "" && currentToolCall.ID == currentToolCallID {
							if toolCall.Function.Name != "" {
								currentToolCall.Function.Name = toolCall.Function.Name
							}
							if toolCall.Function.Arguments != "" {
								currentToolCall.Function.Arguments += toolCall.Function.Arguments
							}
							finalToolCalls[currentToolCallID] = currentToolCall
						}
					}
				}
			}

			// 处理正常消息内容
			if choice.Delta.Content != "" {
				answerFullContent += choice.Delta.Content
				event := dto.ChatMessageResponseEventDto{
					ConversationID: conversationID,
					RequestID:      requestID,
					MessageType:    "answer_delta",
					Content:        choice.Delta.Content,
				}
				eventJSON, _ := json.Marshal(event)
				pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
			}

			// 处理完成原因
			if choice.FinishReason != "" {
				if choice.FinishReason == "stop" {
					// 生成最终消息并保存到数据库
					finalAssistantMessageObj := &models.ChatAgentMessage{
						ApplicationID:  application.ID,
						ChatAgentID:    chatAgent.ID,
						ConversationID: uuid.MustParse(conversationID),
						RequestID:      requestID,
						Type:           "message",
						Role:           "assistant",
						Content:        answerFullContent,
					}

					// 保存消息
					if err := s.messageRepo.Create(ctx, finalAssistantMessageObj); err != nil {
						log.Printf("保存助手消息失败: %v", err)
					}

					// 返回最终答案
					event := dto.ChatMessageResponseEventDto{
						ConversationID: conversationID,
						RequestID:      requestID,
						MessageType:    "answer",
						Content:        answerFullContent,
					}
					eventJSON, _ := json.Marshal(event)
					pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
					break
				}
			}
		}

		// 处理工具调用
		for _, toolCall := range finalToolCalls {
			isNeedAiProcessContinue = true

			// 告诉调用者，有工具调用
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "tool_call",
				Content:        toolCall.Function.Name,
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))

			// 保存工具调用消息到数据库
			functionCallMessageObj := &models.ChatAgentMessage{
				ApplicationID:         application.ID,
				ChatAgentID:           chatAgent.ID,
				ConversationID:        uuid.MustParse(conversationID),
				RequestID:             requestID,
				Type:                  "function_call",
				FunctionCallID:        toolCall.ID,
				FunctionCallName:      toolCall.Function.Name,
				FunctionCallArguments: toolCall.Function.Arguments,
			}
			if err := s.messageRepo.Create(ctx, functionCallMessageObj); err != nil {
				log.Printf("保存工具调用消息失败: %v", err)
			}

			// 告诉调用者，工具调用处理中
			event = dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "tool_call_processing",
				Content:        toolCall.Function.Name,
			}
			eventJSON, _ = json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))

			// 调用工具
			toolResult, err := s.callTool(ctx, chatAgent.ID, toolCall)
			if err != nil {
				log.Printf("调用工具失败: %v", err)
				toolResult = "调用工具失败"
			}

			// 保存工具调用结果到数据库
			functionCallOutputMessageObj := &models.ChatAgentMessage{
				ApplicationID:      application.ID,
				ChatAgentID:        chatAgent.ID,
				ConversationID:     uuid.MustParse(conversationID),
				RequestID:          requestID,
				Type:               "function_call_output",
				FunctionCallID:     toolCall.ID,
				FunctionCallName:   toolCall.Function.Name,
				FunctionCallOutput: toolResult,
			}
			if err := s.messageRepo.Create(ctx, functionCallOutputMessageObj); err != nil {
				log.Printf("保存工具调用结果失败: %v", err)
			}

			// 告诉调用者，工具调用结束
			event = dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "tool_call_end",
				Content:        toolCall.Function.Name,
			}
			eventJSON, _ = json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))

			// 更新消息列表，添加工具调用和结果
			messages = append(messages, al_client.ChatMessage{
				Role:      "assistant",
				Content:   "",
				ToolCalls: []al_client.ToolCall{toolCall},
			})
			messages = append(messages, al_client.ChatMessage{
				Role:       "tool",
				Content:    toolResult,
				ToolCallID: toolCall.ID,
			})
		}

		// 如果需要继续AI处理
		if isNeedAiProcessContinue {
			// 递归调用AI处理
			recursiveReader, err := s.aiProcessStreamable(ctx, conversationID, requestID, messages, aiTools)
			if err != nil {
				event := dto.ChatMessageResponseEventDto{
					ConversationID: conversationID,
					RequestID:      requestID,
					MessageType:    "error",
					Content:        fmt.Sprintf("递归AI处理出错: %v", err),
				}
				eventJSON, _ := json.Marshal(event)
				pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
				return
			}

			// 复制递归结果到当前流
			io.Copy(pw, recursiveReader)
		}
	}()

	return pr, nil
}

// aiProcess 处理AI消息 - 非流式调用AI
func (s *chatAgentConversationService) aiProcess(ctx context.Context, conversationID, requestID string, messages []al_client.ChatMessage, aiTools []al_client.Tool) (io.Reader, error) {
	// 从上下文中获取ChatAgentID
	application, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取上下文信息失败: %w", err)
	}
	// 创建响应管道
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		// 获取应用配置
		llmProvider, llm, err := s.getChatAgentChatLlmConfig(ctx)
		if err != nil {
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "error",
				Content:        fmt.Sprintf("获取应用配置失败: %v", err),
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
			return
		}

		// 动态创建AI客户端
		aiClient, err := s.createAIClient(llmProvider)
		if err != nil {
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "error",
				Content:        fmt.Sprintf("创建AI客户端失败: %v", err),
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
			return
		}

		// 构建请求
		req := al_client.SendMessageRequest{
			Model:       llm.Name,
			Messages:    messages,
			Stream:      false,
			Tools:       aiTools,
			Temperature: chatAgent.ModelParamTemperature,
			TopP:        chatAgent.ModelParamTopP,
			ToolChoice:  "auto",
			MaxTokens:   chatAgent.MaxOutputTokenCountLimit,
		}

		// 发送请求
		response, err := aiClient.SendMessage(ctx, req)
		if err != nil {
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "error",
				Content:        fmt.Sprintf("AI处理出错: %v", err),
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
			return
		}

		isNeedAiProcessContinue := false

		// 处理工具调用
		if response.Choices[0].Message.ToolCalls != nil {
			for _, toolCall := range response.Choices[0].Message.ToolCalls {
				isNeedAiProcessContinue = true

				// 告诉调用者，有工具调用
				event := dto.ChatMessageResponseEventDto{
					ConversationID: conversationID,
					RequestID:      requestID,
					MessageType:    "tool_call",
					Content:        fmt.Sprintf("调用工具: %s", toolCall.Function.Name),
					ToolCall: &dto.ToolCallDto{
						ID:   toolCall.ID,
						Type: toolCall.Type,
						Function: dto.FunctionCallDto{
							Name:      toolCall.Function.Name,
							Arguments: toolCall.Function.Arguments,
						},
					},
				}
				eventJSON, _ := json.Marshal(event)
				pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))

				// 调用工具
				toolResult, err := s.callTool(ctx, chatAgent.ID, toolCall)
				if err != nil {
					log.Printf("调用工具失败: %v", err)
					toolResult = fmt.Sprintf("工具调用失败: %v", err)
				}

				// 告诉调用者，工具调用完成
				event = dto.ChatMessageResponseEventDto{
					ConversationID: conversationID,
					RequestID:      requestID,
					MessageType:    "tool_result",
					Content:        toolResult,
					ToolCall: &dto.ToolCallDto{
						ID:   toolCall.ID,
						Type: toolCall.Type,
						Function: dto.FunctionCallDto{
							Name:      toolCall.Function.Name,
							Arguments: toolCall.Function.Arguments,
						},
					},
				}
				eventJSON, _ = json.Marshal(event)
				pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))

				// 将工具调用和结果添加到消息历史
				messages = append(messages, al_client.ChatMessage{
					Role:      "assistant",
					Content:   "",
					ToolCalls: []al_client.ToolCall{toolCall},
				})

				messages = append(messages, al_client.ChatMessage{
					Role:       "tool",
					Content:    toolResult,
					ToolCallID: toolCall.ID,
				})
			}
		}

		// 如果需要继续AI处理
		if isNeedAiProcessContinue {
			// 递归调用AI处理
			recursiveReader, err := s.aiProcess(ctx, conversationID, requestID, messages, aiTools)
			if err != nil {
				event := dto.ChatMessageResponseEventDto{
					ConversationID: conversationID,
					RequestID:      requestID,
					MessageType:    "error",
					Content:        fmt.Sprintf("递归AI处理出错: %v", err),
				}
				eventJSON, _ := json.Marshal(event)
				pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
				return
			}

			// 复制递归结果到当前流
			io.Copy(pw, recursiveReader)
		} else {
			// 有最终消息，无需调用工具
			assistantMessageObj := &models.ChatAgentMessage{
				ApplicationID:  application.ID,
				ChatAgentID:    chatAgent.ID,
				ConversationID: uuid.MustParse(conversationID),
				RequestID:      requestID,
				Type:           "message",
				Role:           "assistant",
				Content:        response.Choices[0].Message.Content,
			}

			if err := s.messageRepo.Create(ctx, assistantMessageObj); err != nil {
				log.Printf("保存助手消息失败: %v", err)
			}

			// 返回最终答案
			event := dto.ChatMessageResponseEventDto{
				ConversationID: conversationID,
				RequestID:      requestID,
				MessageType:    "answer",
				Content:        response.Choices[0].Message.Content,
			}
			eventJSON, _ := json.Marshal(event)
			pw.Write([]byte(fmt.Sprintf("data: %s\n\n", eventJSON)))
		}
	}()

	return pr, nil
}

// callTool 调用工具
func (s *chatAgentConversationService) callTool(ctx context.Context, agentID uuid.UUID, toolCall al_client.ToolCall) (string, error) {
	toolName := toolCall.Function.Name
	toolArgs := toolCall.Function.Arguments

	var toolCallParams map[string]interface{}

	// 解析 JSON
	err := json.Unmarshal([]byte(toolArgs), &toolCallParams)
	if err != nil {
		fmt.Println("解析 JSON 出错:", err)
		return "", err
	}
	var callToolResult any
	var callToolErr error
	// 判断是否为内部工具
	if strings.HasPrefix(toolName, "__lai__") {
		// 调用内部工具
		callToolResult, callToolErr = s.callInternalTool(ctx, agentID, toolName, toolCallParams)
	} else {
		// 调用MCP工具
		callToolResult, callToolErr = s.callMcpTool(ctx, agentID, toolName, toolCallParams)
	}
	if callToolErr != nil {
		return "", callToolErr
	}
	jsonBytes, marshalJsonErr := json.Marshal(callToolResult)
	if marshalJsonErr != nil {
		fmt.Println("转换 JSON 出错:", err)
		return "", marshalJsonErr
	}
	return string(jsonBytes), nil
}

// callInternalTool 调用内部工具
func (s *chatAgentConversationService) callInternalTool(ctx context.Context, agentID uuid.UUID, toolName string, toolArgs map[string]interface{}) (string, error) {
	// 这里需要根据实际的内部工具实现来调用
	// 暂时返回一个简单的实现
	log.Printf("调用内部工具: %s, 参数: %s", toolName, toolArgs)
	return "内部工具调用结果", nil
}

// callMcpTool 调用MCP工具
func (s *chatAgentConversationService) callMcpTool(ctx context.Context, agentID uuid.UUID, toolName string, toolArgs map[string]interface{}) (any, error) {
	// 解析工具名称，格式为: configID_____toolName
	toolNameItems := strings.Split(toolName, "_____")
	if len(toolNameItems) != 2 {
		return "", fmt.Errorf("无效的工具名称格式: %s", toolName)
	}

	configID := toolNameItems[0]

	// 这里需要根据实际的MCP客户端实现来调用
	// 暂时返回一个简单的实现
	mcpServerConfig, getMcpServerConfigErr := s.mcpConfigRepo.GetByConfigID(ctx, configID)
	if getMcpServerConfigErr != nil {
		return "", fmt.Errorf("获取MCP配置失败: %w", getMcpServerConfigErr)
	}
	mcpClient, getMcpClientError := manager.GetMcpClient(ctx, mcpServerConfig)
	if getMcpClientError != nil {
		return "", fmt.Errorf("创建MCP客户端失败: %w", getMcpClientError)
	}
	callToolResult, callToolErr := mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolNameItems[1],
			Arguments: toolArgs,
		},
	})
	if callToolErr != nil {
		return "", fmt.Errorf("调用MCP工具失败: %w", callToolErr)
	}
	log.Printf("调用MCP工具: %s, 配置ID: %s, 参数: %v, 结果: %v", toolNameItems[1], configID, toolArgs, callToolResult)
	return callToolResult.Content, nil
}

func (s *chatAgentConversationService) getChatAgentChatLlmConfig(ctx context.Context) (llmProvider *models.ApplicationLlmProvider, chatModel *models.ApplicationLlm, err error) {
	// 获取应用的LLM配置
	_, chatAgent, err := getContextInfo(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("获取应用Agent LLM配置失败: %w", err)
	}

	chatLlm, getChatLlmErr := s.llmRepo.GetByID(ctx, chatAgent.ChatModelID)
	if getChatLlmErr != nil {
		return nil, nil, fmt.Errorf("ChatAgent未配置Chat LLM")
	}

	chatLlmProvider, getChatLlmProviderErr := s.llmProviderRepo.GetByID(ctx, chatLlm.LlmProviderID)
	if getChatLlmProviderErr != nil {
		return nil, nil, fmt.Errorf("ChatAgent未配置Chat LLM提供商")
	}

	return chatLlmProvider, chatLlm, nil
}
func (s *chatAgentConversationService) getChatAgentNamingLlmConfig(ctx context.Context) (llmProvider *models.ApplicationLlmProvider, chatModel *models.ApplicationLlm, err error) {
	_, chatAgent, err := getContextInfo(ctx)
	// 获取应用的LLM配置
	if err != nil {
		return nil, nil, fmt.Errorf("获取应用Agent LLM配置失败: %w", err)
	}

	namingLlm, getNamingLlmErr := s.llmRepo.GetByID(ctx, chatAgent.ConversationNamingModelID)
	if getNamingLlmErr != nil {
		return nil, nil, fmt.Errorf("ChatAgent未配置Conversation Naming LLM")
	}

	namingLlmProvider, getNamingLlmProviderErr := s.llmProviderRepo.GetByID(ctx, namingLlm.LlmProviderID)
	if getNamingLlmProviderErr != nil {
		return nil, nil, fmt.Errorf("ChatAgent未配置Conversation Naming LLM提供商")
	}

	return namingLlmProvider, namingLlm, nil
}

// createAIClient 根据LLM提供商配置创建AI客户端
func (s *chatAgentConversationService) createAIClient(llmProvider *models.ApplicationLlmProvider) (al_client.LemonAiClient, error) {
	// 根据LLM提供商类型创建相应的AI客户端
	switch llmProvider.Type {
	case "openai_chat_completions_api":
		return al_client.NewOpenAIChatCompletionsClient(llmProvider.ApiKey), nil
	case "ollama":
		// TODO: 实现Ollama客户端
		return nil, fmt.Errorf("Ollama客户端尚未实现")
	case "volcano_engine":
		// TODO: 实现火山引擎客户端
		return nil, fmt.Errorf("火山引擎客户端尚未实现")
	default:
		// 默认使用OpenAI
		return al_client.NewOpenAIChatCompletionsClient(llmProvider.ApiKey), nil
	}
}

// GetChatAgentMcpServerTools 获取聊天智能体启用的MCP工具列表
// 根据chatAgentID查询启用的工具，并从MCP服务器获取最新的工具信息
func (s *chatAgentConversationService) GetChatAgentMcpServerTools(ctx context.Context) ([]al_client.Tool, error) {
	_, chatAgent, getContextErr := getContextInfo(ctx)
	if getContextErr != nil {
		return nil, fmt.Errorf("获取上下文数据失败: %w", getContextErr)
	}
	// 1. 查询该聊天智能体启用的MCP工具配置
	enabledTools, err := s.chatAgentMcpServerToolRepo.GetByChatAgentID(ctx, chatAgent.ID)
	if err != nil {
		return nil, fmt.Errorf("获取聊天智能体MCP工具配置失败: %w", err)
	}

	// 过滤出启用的工具
	var enabledToolIDs []uuid.UUID
	for _, tool := range enabledTools {
		if tool.Enabled {
			enabledToolIDs = append(enabledToolIDs, tool.ApplicationMcpServerToolID)
		}
	}

	if len(enabledToolIDs) == 0 {
		return []al_client.Tool{}, nil
	}

	// 2. 根据工具ID逐个获取工具信息，并按MCP配置分组
	configToolsMap := make(map[uuid.UUID][]*models.ApplicationMcpServerTool)
	for _, toolID := range enabledToolIDs {
		tool, err := s.mcpToolRepo.GetByID(ctx, toolID)
		if err != nil {
			log.Printf("获取MCP工具信息失败: %v", err)
			continue
		}
		if tool != nil {
			configToolsMap[tool.ApplicationMcpServerConfigID] = append(configToolsMap[tool.ApplicationMcpServerConfigID], tool)
		}
	}

	var openaiTools []al_client.Tool

	// 4. 为每个MCP配置获取最新的工具信息
	for configID, configTools := range configToolsMap {
		// 获取MCP配置信息
		config, err := s.mcpConfigRepo.GetByID(ctx, configID)
		if err != nil {
			log.Printf("获取MCP配置失败: %v", err)
			continue
		}

		// 从MCP服务器获取最新的工具信息
		mcpTools, err := s.getToolsFromMcpServer(ctx, config)
		if err != nil {
			log.Printf("从MCP服务器获取工具失败: %v", err)
			continue
		}

		// 创建工具名称映射
		mcpToolsMap := make(map[string]mcp.Tool)
		for _, mcpTool := range mcpTools {
			mcpToolsMap[mcpTool.Name] = mcpTool
		}

		// 为每个启用的工具创建OpenAI工具格式
		for _, tool := range configTools {
			if mcpTool, exists := mcpToolsMap[tool.Name]; exists {
				configShortID, _ := utils.ShortUUID(configID.String())
				openaiTool := s.convertMcpToolToOpenAI(mcpTool, configShortID)
				openaiTools = append(openaiTools, openaiTool)
			}
		}
	}

	return openaiTools, nil
}

// getToolsFromMcpServer 从MCP服务器获取工具列表
func (s *chatAgentConversationService) getToolsFromMcpServer(ctx context.Context, config *models.ApplicationMcpServerConfig) ([]mcp.Tool, error) {
	// 使用公共的MCP客户端管理器创建客户端
	c, err := manager.GetMcpClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("创建MCP客户端失败: %w", err)
	}

	// 获取工具列表
	toolsRequest := mcp.ListToolsRequest{}
	toolsResult, err := c.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, fmt.Errorf("获取MCP工具列表失败: %w", err)
	}

	log.Printf("MCP服务器有 %d 个可用工具", len(toolsResult.Tools))
	return toolsResult.Tools, nil
}

// convertMcpToolToOpenAI 将MCP工具转换为AI工具格式
func (s *chatAgentConversationService) convertMcpToolToOpenAI(mcpTool mcp.Tool, configID string) al_client.Tool {
	// 构建工具名称，格式为: configID_____toolName
	toolName := fmt.Sprintf("%s_____%s", configID, mcpTool.Name)

	// 转换InputSchema为map[string]interface{}
	parameters := make(map[string]interface{})
	if mcpTool.InputSchema.Type != "" {
		parameters["type"] = mcpTool.InputSchema.Type
	}
	if mcpTool.InputSchema.Properties != nil {
		parameters["properties"] = mcpTool.InputSchema.Properties
	}
	if len(mcpTool.InputSchema.Required) > 0 {
		parameters["required"] = mcpTool.InputSchema.Required
	}

	return al_client.Tool{
		Type: "function",
		Function: &al_client.FunctionDefinition{
			Name:        toolName,
			Description: mcpTool.Description,
			Parameters:  parameters,
		},
	}
}
