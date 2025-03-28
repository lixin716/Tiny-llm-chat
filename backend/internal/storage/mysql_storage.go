package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"chat-llama/internal/service"
	"chat-llama/pkg/cache"
	"chat-llama/pkg/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MySQLStorage 提供MySQL数据库实现
type MySQLStorage struct {
	db *gorm.DB
}

// NewMySQLStorage 创建新的MySQL存储
func NewMySQLStorage() *MySQLStorage {
	return &MySQLStorage{
		db: db.DB,
	}
}

// 缓存键生成函数
func conversationKey(id string) string {
	return fmt.Sprintf("conversation:%s", id)
}

func userConversationsKey(userID uint) string {
	return fmt.Sprintf("user:%d:conversations", userID)
}

func conversationMessagesKey(conversationID string) string {
	return fmt.Sprintf("conversation:%s:messages", conversationID)
}

// CreateConversation 创建新会话
func (s *MySQLStorage) CreateConversation(userID uint, title string) (*service.Conversation, error) {
	ctx := context.Background()

	// 生成UUID
	id := uuid.New().String()

	// 创建会话记录
	conversation := &Conversation{
		ID:        id,
		UserID:    userID,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	if err := s.db.Create(conversation).Error; err != nil {
		return nil, err
	}

	// 缓存会话
	serviceConv := conversation.ToServiceModel()
	cache.Set(ctx, conversationKey(id), serviceConv, time.Hour)

	// 清除用户会话列表缓存
	cache.Delete(ctx, userConversationsKey(userID))

	return serviceConv, nil
}

// GetConversation 获取会话
func (s *MySQLStorage) GetConversation(id string) (*service.Conversation, error) {
	ctx := context.Background()

	// 尝试从缓存获取
	var serviceConv service.Conversation
	found, err := cache.Get(ctx, conversationKey(id), &serviceConv)
	if err != nil {
		return nil, err
	}

	if found {
		return &serviceConv, nil
	}

	// 缓存未命中，从数据库获取
	var conversation Conversation
	if err := s.db.Where("id = ?", id).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("会话不存在")
		}
		return nil, err
	}

	// 更新缓存
	serviceConv = *conversation.ToServiceModel()
	cache.Set(ctx, conversationKey(id), serviceConv, time.Hour)

	return &serviceConv, nil
}

// GetConversationsByUserID 获取用户的所有会话
func (s *MySQLStorage) GetConversationsByUserID(userID uint) ([]*service.Conversation, error) {
	ctx := context.Background()

	// 尝试从缓存获取
	var serviceConvs []*service.Conversation
	found, err := cache.Get(ctx, userConversationsKey(userID), &serviceConvs)
	if err != nil {
		return nil, err
	}

	if found {
		return serviceConvs, nil
	}

	// 缓存未命中，从数据库获取
	var conversations []Conversation
	if err := s.db.Where("user_id = ?", userID).Order("updated_at DESC").Find(&conversations).Error; err != nil {
		return nil, err
	}

	// 转换为服务层模型
	serviceConvs = make([]*service.Conversation, len(conversations))
	for i, conv := range conversations {
		serviceConvs[i] = conv.ToServiceModel()
	}

	// 更新缓存
	cache.Set(ctx, userConversationsKey(userID), serviceConvs, time.Hour)

	return serviceConvs, nil
}

// UpdateConversationTitle 更新会话标题
func (s *MySQLStorage) UpdateConversationTitle(id string, title string) error {
	ctx := context.Background()

	// 获取会话以检查存在性
	var conversation Conversation
	if err := s.db.Where("id = ?", id).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("会话不存在")
		}
		return err
	}

	// 更新会话
	conversation.Title = title
	conversation.UpdatedAt = time.Now()

	if err := s.db.Save(&conversation).Error; err != nil {
		return err
	}

	// 清除缓存
	cache.Delete(ctx, conversationKey(id))
	cache.Delete(ctx, userConversationsKey(conversation.UserID))

	return nil
}

// DeleteConversation 删除会话
func (s *MySQLStorage) DeleteConversation(id string) error {
	ctx := context.Background()

	// 先查询会话以获取用户ID
	var conversation Conversation
	if err := s.db.Where("id = ?", id).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("会话不存在")
		}
		return err
	}

	// 开始事务
	tx := s.db.Begin()

	// 删除会话中的所有消息
	if err := tx.Where("conversation_id = ?", id).Delete(&Message{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除会话
	if err := tx.Delete(&conversation).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 清除缓存
	cache.Delete(ctx, conversationKey(id))
	cache.Delete(ctx, userConversationsKey(conversation.UserID))
	cache.Delete(ctx, conversationMessagesKey(id))

	return nil
}

// AddMessage 添加消息
func (s *MySQLStorage) AddMessage(conversationID, role, content string) (*service.Message, error) {
	ctx := context.Background()

	// 验证会话是否存在
	var conversation Conversation
	if err := s.db.Where("id = ?", conversationID).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("会话不存在")
		}
		return nil, err
	}

	// 生成UUID
	id := uuid.New().String()

	// 创建消息
	message := &Message{
		ID:             id,
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 保存消息
	if err := s.db.Create(message).Error; err != nil {
		return nil, err
	}

	// 更新会话的最后更新时间
	conversation.UpdatedAt = time.Now()
	if err := s.db.Save(&conversation).Error; err != nil {
		return nil, err
	}

	// 清除缓存
	cache.Delete(ctx, conversationMessagesKey(conversationID))
	cache.Delete(ctx, conversationKey(conversationID))
	cache.Delete(ctx, userConversationsKey(conversation.UserID))

	// 返回服务层消息模型
	return message.ToServiceModel(), nil
}

// GetMessagesByConversationID 获取会话的所有消息
func (s *MySQLStorage) GetMessagesByConversationID(conversationID string) ([]*service.Message, error) {
	ctx := context.Background()

	// 尝试从缓存获取
	var serviceMessages []*service.Message
	found, err := cache.Get(ctx, conversationMessagesKey(conversationID), &serviceMessages)
	if err != nil {
		return nil, err
	}

	if found {
		return serviceMessages, nil
	}

	// 缓存未命中，从数据库获取
	var messages []Message
	if err := s.db.Where("conversation_id = ?", conversationID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}

	// 转换为服务层模型
	serviceMessages = make([]*service.Message, len(messages))
	for i, msg := range messages {
		serviceMessages[i] = msg.ToServiceModel()
	}

	// 更新缓存
	cache.Set(ctx, conversationMessagesKey(conversationID), serviceMessages, time.Hour)

	return serviceMessages, nil
}
