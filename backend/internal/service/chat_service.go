package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"chat-llama/internal/model"
)

// ChatService 提供聊天相关功能
type ChatService struct {
	llmClient *model.LLMClient
	storage   Storage
}

// NewChatService 创建聊天服务实例
func NewChatService(llmClient *model.LLMClient, storage Storage) *ChatService {
	return &ChatService{
		llmClient: llmClient,
		storage:   storage,
	}
}

// Chat 处理聊天请求并返回模型响应
func (s *ChatService) Chat(ctx context.Context, userID uint, req *ChatRequest) (*ChatResponse, error) {
	var conversationID string
	var err error

	// 检查是新会话还是已有会话
	if req.ConversationID == "" {
		// 创建新会话
		title := createTitleFromMessage(req.Message)
		conv, err := s.storage.CreateConversation(userID, title)
		if err != nil {
			return nil, err
		}
		conversationID = conv.ID
	} else {
		// 验证会话存在且属于该用户
		conv, err := s.storage.GetConversation(req.ConversationID)
		if err != nil {
			return nil, err
		}
		if conv.UserID != userID {
			return nil, errors.New("无权访问此会话")
		}
		conversationID = req.ConversationID
	}

	// 保存用户消息
	_, err = s.storage.AddMessage(conversationID, "user", req.Message)
	if err != nil {
		return nil, err
	}

	// 构建提示词
	prompt, err := s.buildPrompt(conversationID)
	if err != nil {
		return nil, err
	}

	// 确保参数有合理默认值
	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	maxNewTokens := req.MaxNewTokens
	if maxNewTokens == 0 {
		maxNewTokens = 500
	}

	topK := req.TopK
	if topK == 0 {
		topK = 40
	}

	// 调用模型生成回复
	llmResponse, err := s.llmClient.GenerateResponse(
		ctx,
		prompt,
		temperature,
		maxNewTokens,
		topK,
	)
	if err != nil {
		log.Printf("调用LLM服务失败: %v", err)
		return nil, err
	}

	// 保存模型回复
	_, err = s.storage.AddMessage(conversationID, "assistant", llmResponse)
	if err != nil {
		return nil, err
	}

	return &ChatResponse{
		ConversationID: conversationID,
		Message:        llmResponse,
		Role:           "assistant",
	}, nil
}

// GetConversations 获取用户的所有会话
func (s *ChatService) GetConversations(userID uint) ([]*Conversation, error) {
	return s.storage.GetConversationsByUserID(userID)
}

// GetConversationHistory 获取会话的消息历史
func (s *ChatService) GetConversationHistory(userID uint, conversationID string) ([]*Message, error) {
	// 检查会话归属
	conv, err := s.storage.GetConversation(conversationID)
	if err != nil {
		return nil, err
	}

	if conv.UserID != userID {
		return nil, errors.New("无权访问此会话")
	}

	return s.storage.GetMessagesByConversationID(conversationID)
}

// DeleteConversation 删除会话
func (s *ChatService) DeleteConversation(userID uint, conversationID string) error {
	// 检查会话归属
	conv, err := s.storage.GetConversation(conversationID)
	if err != nil {
		return err
	}

	if conv.UserID != userID {
		return errors.New("无权删除此会话")
	}

	return s.storage.DeleteConversation(conversationID)
}

// buildPrompt 构建发送给LLM的提示词
func (s *ChatService) buildPrompt(conversationID string) (string, error) {
	messages, err := s.storage.GetMessagesByConversationID(conversationID)
	if err != nil {
		return "", err
	}

	// 构建提示词，包含历史对话
	var sb strings.Builder

	for _, msg := range messages {
		if msg.Role == "user" {
			sb.WriteString("用户: " + msg.Content + "\n")
		} else if msg.Role == "assistant" {
			sb.WriteString("助手: " + msg.Content + "\n")
		}
	}

	// 增加最后的提示
	sb.WriteString("助手: ")

	return sb.String(), nil
}

// createTitleFromMessage 从消息内容创建会话标题
func createTitleFromMessage(message string) string {
	if len(message) == 0 {
		return "新对话"
	}

	// 将字符串转换为[]rune以正确处理Unicode字符
	runes := []rune(message)
	if len(runes) > 20 {
		// 截取前20个字符（而不是字节）
		return string(runes[:20]) + "..."
	}

	return message
}

// GetConversation 获取会话信息
func (s *ChatService) GetConversation(conversationID string) (*Conversation, error) {
	return s.storage.GetConversation(conversationID)
}

// UpdateConversationTitle 更新会话标题
func (s *ChatService) UpdateConversationTitle(conversationID string, title string) error {
	return s.storage.UpdateConversationTitle(conversationID, title)
}
