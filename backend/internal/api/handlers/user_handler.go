package handlers

import (
	"net/http"
	"strings"
	"time"

	"chat-llama/internal/storage"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// UserHandler 处理用户相关请求
type UserHandler struct {
	userStorage *storage.UserStorage
	jwtSecret   []byte
}

// NewUserHandler 创建用户处理程序
func NewUserHandler(userStorage *storage.UserStorage, jwtSecret string) *UserHandler {
	return &UserHandler{
		userStorage: userStorage,
		jwtSecret:   []byte(jwtSecret),
	}
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWT Claims
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// Register 用户注册
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := ParseJSON(r, &req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 验证输入
	if req.Username == "" || req.Password == "" {
		ErrorResponse(w, http.StatusBadRequest, "用户名和密码不能为空")
		return
	}

	// 创建用户
	user, err := h.userStorage.CreateUser(req.Username, req.Password)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "注册失败: "+err.Error())
		return
	}

	// 返回成功响应
	SuccessResponse(w, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Login 用户登录
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := ParseJSON(r, &req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 记录输入的密码(调试用)
	logger.Printf("登录尝试 - 用户名: %s, 密码: %s", req.Username, maskPassword(req.Password))

	// 验证用户名和密码
	user, err := h.userStorage.VerifyPassword(req.Username, req.Password)
	if err != nil {
		ErrorResponse(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 生成JWT令牌
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	// 返回令牌和用户信息
	SuccessResponse(w, LoginResponse{
		Token: tokenString,
		User: map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// GetProfile 获取用户信息
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(uint)

	// 查询用户信息
	user, err := h.userStorage.GetUserByID(userID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	// 返回用户信息
	SuccessResponse(w, map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	})
}

// 遮盖密码，只显示长度和前两个字符(如果存在)
func maskPassword(password string) string {
	if len(password) <= 2 {
		return strings.Repeat("*", len(password))
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}
