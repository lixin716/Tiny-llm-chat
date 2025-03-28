package service

// Storage 定义聊天数据的存储接口
type Storage interface {
	// 会话管理
	CreateConversation(userID uint, title string) (*Conversation, error)
	GetConversation(id string) (*Conversation, error)
	GetConversationsByUserID(userID uint) ([]*Conversation, error)
	UpdateConversationTitle(id string, title string) error
	DeleteConversation(id string) error

	// 消息管理
	AddMessage(conversationID, role, content string) (*Message, error)
	GetMessagesByConversationID(conversationID string) ([]*Message, error)
}
