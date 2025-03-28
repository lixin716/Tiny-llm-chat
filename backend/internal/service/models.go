package service

import "time"

// Conversation 表示一个聊天会话
type Conversation struct {
	ID        string    `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Message 表示聊天消息
type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	Role           string    `json:"role"` // "user" 或 "assistant"
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// ChatRequest 表示聊天请求
type ChatRequest struct {
	ConversationID string  `json:"conversation_id,omitempty"`
	Message        string  `json:"message"`
	Temperature    float32 `json:"temperature,omitempty"`
	MaxNewTokens   int32   `json:"max_new_tokens,omitempty"`
	TopK           int32   `json:"top_k,omitempty"`
}

// ChatResponse 表示聊天响应
type ChatResponse struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
	Role           string `json:"role"`
}
