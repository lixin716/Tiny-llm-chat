package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"chat-llama/config"
)

// Server API服务器
type Server struct {
	server *http.Server
	router http.Handler
}

// NewServer 创建新的API服务器
func NewServer(router http.Handler) *Server {
	// 获取配置
	cfg := config.GetConfig()
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: server,
		router: router,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 优雅关闭的通道
	idleConnsClosed := make(chan struct{})

	// 监听系统信号
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("正在关闭服务器...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("服务器关闭错误: %v", err)
		}

		close(idleConnsClosed)
	}()

	// 启动服务器
	log.Printf("服务器启动在 %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("服务器启动失败: %v", err)
	}

	<-idleConnsClosed
	log.Println("服务器已关闭")
	return nil
}
