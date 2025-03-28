package service

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryStorage 提供基于内存的存储实现，适用于开发和测试
type MemoryStorage struct {
	conversations map[string]*Conversation
	messages      map[string][]*Message
	mutex         sync.RWMutex
}

// NewMemoryStorage 创建内存存储实例
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		conversations: make(map[string]*Conversation),
		messages:      make(map[string][]*Message),
	}
}

// CreateConversation 创建新会话
func (s *MemoryStorage) CreateConversation(userID uint, title string) (*Conversation, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	conv := &Conversation{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.conversations[conv.ID] = conv
	s.messages[conv.ID] = []*Message{}

	return conv, nil
}

// GetConversation 获取会话
func (s *MemoryStorage) GetConversation(id string) (*Conversation, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	conv, exists := s.conversations[id]
	if !exists {
		return nil, errors.New("会话不存在")
	}

	return conv, nil
}

// GetConversationsByUserID 获取用户的所有会话
func (s *MemoryStorage) GetConversationsByUserID(userID uint) ([]*Conversation, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*Conversation
	for _, conv := range s.conversations {
		if conv.UserID == userID {
			result = append(result, conv)
		}
	}

	return result, nil
}

// UpdateConversationTitle 更新会话标题
func (s *MemoryStorage) UpdateConversationTitle(id string, title string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	conv, exists := s.conversations[id]
	if !exists {
		return errors.New("会话不存在")
	}

	conv.Title = title
	conv.UpdatedAt = time.Now()

	return nil
}

// DeleteConversation 删除会话
func (s *MemoryStorage) DeleteConversation(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.conversations[id]; !exists {
		return errors.New("会话不存在")
	}

	delete(s.conversations, id)
	delete(s.messages, id)

	return nil
}

// AddMessage 添加消息
func (s *MemoryStorage) AddMessage(conversationID, role, content string) (*Message, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.conversations[conversationID]; !exists {
		return nil, errors.New("会话不存在")
	}

	msg := &Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	s.messages[conversationID] = append(s.messages[conversationID], msg)

	// 更新会话的最后更新时间
	s.conversations[conversationID].UpdatedAt = time.Now()

	return msg, nil
}

// GetMessagesByConversationID 获取会话的所有消息
func (s *MemoryStorage) GetMessagesByConversationID(conversationID string) ([]*Message, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	msgs, exists := s.messages[conversationID]
	if !exists {
		return nil, errors.New("会话不存在")
	}

	return msgs, nil
}
