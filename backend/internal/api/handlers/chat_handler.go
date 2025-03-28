package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"chat-llama/internal/service"

	"github.com/gorilla/mux"
)

// ChatHandler 处理聊天相关请求
type ChatHandler struct {
	chatService *service.ChatService
}

// NewChatHandler 创建聊天处理程序
func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// GetConversations 获取用户的所有会话
func (h *ChatHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(uint)

	// 获取会话列表
	conversations, err := h.chatService.GetConversations(userID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "获取会话失败: "+err.Error())
		return
	}

	SuccessResponse(w, conversations)
}

// GetConversationHistory 获取会话历史消息
func (h *ChatHandler) GetConversationHistory(w http.ResponseWriter, r *http.Request) {
	// 手动解析URL路径
	pathParts := strings.Split(r.URL.Path, "/")
	var conversationID string
	if len(pathParts) >= 3 {
		conversationID = pathParts[len(pathParts)-2] // 倒数第二个元素应该是ID
	}

	log.Printf("解析的会话ID: %s", conversationID)

	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(uint)

	// 获取会话历史
	messages, err := h.chatService.GetConversationHistory(userID, conversationID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "获取会话历史失败: "+err.Error())
		return
	}

	SuccessResponse(w, messages)
}

// DeleteConversation 删除会话
func (h *ChatHandler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(uint)

	// 从上下文获取会话ID (而不是从mux.Vars)
	conversationID, ok := r.Context().Value("id").(string)
	if !ok || conversationID == "" {
		ErrorResponse(w, http.StatusBadRequest, "无效的会话ID")
		return
	}

	// 删除会话
	err := h.chatService.DeleteConversation(userID, conversationID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "删除会话失败: "+err.Error())
		return
	}

	SuccessResponse(w, nil)
}

// UpdateConversationTitle 更新会话标题
func (h *ChatHandler) UpdateConversationTitle(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(uint)

	// 从URL获取会话ID
	vars := mux.Vars(r)
	conversationID := vars["id"]

	// 解析请求体
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取会话确认所有权
	conversation, err := h.chatService.GetConversation(conversationID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "获取会话失败: "+err.Error())
		return
	}

	if conversation.UserID != userID {
		ErrorResponse(w, http.StatusForbidden, "无权修改此会话")
		return
	}

	// 更新标题
	if err := h.chatService.UpdateConversationTitle(conversationID, req.Title); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "更新会话标题失败: "+err.Error())
		return
	}

	SuccessResponse(w, nil)
}

// Chat 处理聊天请求
func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(uint)

	// 解析请求
	var chatReq service.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&chatReq); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 发送聊天请求
	ctx := r.Context()
	response, err := h.chatService.Chat(ctx, userID, &chatReq)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "处理聊天请求失败: "+err.Error())
		return
	}

	SuccessResponse(w, response)
}
