package api

import (
	"context"
	"net/http"

	"chat-llama/config"
	"chat-llama/internal/api/handlers"
	"chat-llama/internal/api/middleware"
	"chat-llama/internal/service"
	"chat-llama/internal/storage"

	"github.com/gin-gonic/gin"
)

// Router API路由器
type Router struct {
	engine      *gin.Engine
	chatService *service.ChatService
	userStorage *storage.UserStorage
}

// NewRouter 创建新路由器
func NewRouter(chatService *service.ChatService, userStorage *storage.UserStorage) *Router {
	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	return &Router{
		engine:      engine,
		chatService: chatService,
		userStorage: userStorage,
	}
}

// Setup 设置路由
func (r *Router) Setup() http.Handler {
	// 获取配置
	cfg := config.GetConfig()

	// 创建处理程序
	userHandler := handlers.NewUserHandler(r.userStorage, cfg.Server.JWTSecret)
	chatHandler := handlers.NewChatHandler(r.chatService)
	wsHandler := handlers.NewWebSocketHandler(r.chatService)

	// 创建中间件包装器
	jwtMiddleware := func(c *gin.Context) {
		middleware.NewJWTMiddleware(cfg.Server.JWTSecret).Middleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.Request = r
				c.Next()
			}),
		).ServeHTTP(c.Writer, c.Request)

		if c.IsAborted() {
			return
		}
	}

	loggerMiddleware := func(c *gin.Context) {
		middleware.NewLoggerMiddleware().Middleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.Next()
			}),
		).ServeHTTP(c.Writer, c.Request)
	}

	corsMiddleware := func(c *gin.Context) {
		middleware.NewCORSMiddleware().Middleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if c.Request.Method == "OPTIONS" {
					c.AbortWithStatus(http.StatusOK)
					return
				}
				c.Next()
			}),
		).ServeHTTP(c.Writer, c.Request)
	}

	// 应用全局中间件
	r.engine.Use(loggerMiddleware)
	r.engine.Use(corsMiddleware)
	r.engine.Use(gin.Recovery()) // 添加Gin的Recovery中间件处理panic

	// API路由组
	api := r.engine.Group("/api")

	// 无需认证的路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", gin.WrapF(userHandler.Register))
		auth.POST("/login", gin.WrapF(userHandler.Login))
	}

	// 需要认证的API路由
	protected := api.Group("")
	protected.Use(jwtMiddleware)
	{
		// 用户相关路由
		user := protected.Group("/user")
		{
			user.GET("/profile", gin.WrapF(userHandler.GetProfile))
		}

		// 聊天相关路由
		protected.GET("/conversations", gin.WrapF(chatHandler.GetConversations))
		protected.GET("/conversations/:id/messages", func(c *gin.Context) {
			// 提取参数并设置到请求上下文
			id := c.Param("id")
			r := c.Request
			ctx := context.WithValue(r.Context(), "conversationID", id)
			c.Request = r.WithContext(ctx)

			// 调用原始处理程序
			chatHandler.GetConversationHistory(c.Writer, c.Request)
		})
		protected.DELETE("/conversations/:id", func(c *gin.Context) {
			// 提取参数并设置到请求上下文
			id := c.Param("id")
			r := c.Request
			ctx := context.WithValue(r.Context(), "id", id)
			c.Request = r.WithContext(ctx)

			// 调用原始处理程序
			chatHandler.DeleteConversation(c.Writer, c.Request)
		})
		protected.PUT("/conversations/:id/title", gin.WrapF(chatHandler.UpdateConversationTitle))
		protected.POST("/chat", gin.WrapF(chatHandler.Chat))

		// WebSocket路由
		protected.GET("/ws", gin.WrapF(wsHandler.HandleWebSocket))
	}

	return r.engine
}

// ServeHTTP 实现http.Handler接口
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}
