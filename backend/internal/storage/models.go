package storage

import (
	"chat-llama/internal/service"
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;not null;unique" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"` // 密码不输出到JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Conversation 会话模型
type Conversation struct {
	ID        string         `gorm:"primarykey;type:varchar(36)" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Title     string         `gorm:"type:text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Message 消息模型
type Message struct {
	ID             string         `gorm:"primarykey;type:varchar(36)" json:"id"`
	ConversationID string         `gorm:"index;type:varchar(36);not null" json:"conversation_id"`
	Role           string         `gorm:"size:20;not null" json:"role"` // "user" 或 "assistant"
	Content        string         `gorm:"type:text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null" json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名设置
func (User) TableName() string {
	return "users"
}

func (Conversation) TableName() string {
	return "conversations"
}

func (Message) TableName() string {
	return "messages"
}

// 初始化数据库表
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Conversation{}, &Message{})
}

// 数据库模型转换为服务层模型
func (c *Conversation) ToServiceModel() *service.Conversation {
	return &service.Conversation{
		ID:        c.ID,
		UserID:    c.UserID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (m *Message) ToServiceModel() *service.Message {
	return &service.Message{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		Role:           m.Role,
		Content:        m.Content,
		CreatedAt:      m.CreatedAt,
	}
}
