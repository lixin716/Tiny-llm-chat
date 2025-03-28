package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"chat-llama/internal/service"

	"github.com/gorilla/websocket"
)

const (
	// 写入超时
	writeWait = 10 * time.Second

	// 读取超时
	pongWait = 60 * time.Second

	// Ping 间隔
	pingPeriod = (pongWait * 9) / 10

	// 消息大小限制
	maxMessageSize = 512 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源，生产环境应该检查来源
		return true
	},
}

// 消息类型
const (
	TypeChat    = "chat"
	TypeHistory = "history"
	TypeError   = "error"
)

// WebSocketMessage WebSocket消息结构
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

// WebSocketClient WebSocket客户端
type WebSocketClient struct {
	conn        *websocket.Conn
	userID      uint
	send        chan []byte
	chatService *service.ChatService
}

// NewWebSocketClient 创建新的WebSocket客户端
func NewWebSocketClient(conn *websocket.Conn, userID uint, chatService *service.ChatService) *WebSocketClient {
	return &WebSocketClient{
		conn:        conn,
		userID:      userID,
		send:        make(chan []byte, 256),
		chatService: chatService,
	}
}

// WebSocketHandler 处理WebSocket连接
type WebSocketHandler struct {
	chatService *service.ChatService
}

// NewWebSocketHandler 创建新的WebSocket处理程序
func NewWebSocketHandler(chatService *service.ChatService) *WebSocketHandler {
	return &WebSocketHandler{
		chatService: chatService,
	}
}

// HandleWebSocket 处理WebSocket连接请求
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(uint)

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 创建客户端
	client := NewWebSocketClient(conn, userID, h.chatService)

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

// readPump 从WebSocket连接读取消息
func (c *WebSocketClient) readPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		// 处理接收到的消息
		c.handleMessage(message)
	}
}

// writePump 向WebSocket连接发送消息
func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 通道已关闭
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加队列中的消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理接收到的消息
func (c *WebSocketClient) handleMessage(data []byte) {
	// 解析消息
	var msg WebSocketMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.sendError("无效的消息格式")
		return
	}

	// 根据消息类型处理
	switch msg.Type {
	case TypeChat:
		var chatReq service.ChatRequest
		if err := json.Unmarshal(msg.Content, &chatReq); err != nil {
			c.sendError("无效的聊天请求")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// 处理聊天请求
		resp, err := c.chatService.Chat(ctx, c.userID, &chatReq)
		if err != nil {
			c.sendError("处理聊天请求失败: " + err.Error())
			return
		}

		// 发送响应
		c.sendResponse(TypeChat, resp)

	case TypeHistory:
		var historyReq struct {
			ConversationID string `json:"conversation_id"`
		}
		if err := json.Unmarshal(msg.Content, &historyReq); err != nil {
			c.sendError("无效的历史请求")
			return
		}

		// 获取会话历史
		messages, err := c.chatService.GetConversationHistory(c.userID, historyReq.ConversationID)
		if err != nil {
			c.sendError("获取会话历史失败: " + err.Error())
			return
		}

		// 发送响应
		c.sendResponse(TypeHistory, messages)

	default:
		c.sendError("不支持的消息类型: " + msg.Type)
	}
}

// sendResponse 发送响应
func (c *WebSocketClient) sendResponse(msgType string, data interface{}) {
	resp := WebSocketMessage{
		Type: msgType,
	}

	// 序列化内容
	content, err := json.Marshal(data)
	if err != nil {
		c.sendError("序列化响应失败")
		return
	}
	resp.Content = content

	// 序列化完整消息
	message, err := json.Marshal(resp)
	if err != nil {
		c.sendError("序列化消息失败")
		return
	}

	select {
	case c.send <- message:
	default:
		// 发送缓冲区已满，关闭连接
		c.conn.Close()
	}
}

// sendError 发送错误消息
func (c *WebSocketClient) sendError(errMsg string) {
	resp := WebSocketMessage{
		Type: TypeError,
	}

	// 序列化错误内容
	content, err := json.Marshal(map[string]string{"message": errMsg})
	if err != nil {
		log.Printf("序列化错误消息失败: %v", err)
		return
	}
	resp.Content = content

	// 序列化完整消息
	message, err := json.Marshal(resp)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}

	select {
	case c.send <- message:
	default:
		// 发送缓冲区已满，关闭连接
		c.conn.Close()
	}
}
