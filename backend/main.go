package main

import (
	"log"

	"chat-llama/config"
	"chat-llama/internal/api"
	"chat-llama/internal/model"
	"chat-llama/internal/service"
	"chat-llama/internal/storage"
	"chat-llama/pkg/cache"
	"chat-llama/pkg/db"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志目录
	if err := cfg.Log.InitLogDir(); err != nil {
		log.Fatalf("初始化日志目录失败: %v", err)
	}

	// 初始化数据库连接
	dbConfig := db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	}

	if err := db.InitMySQL(dbConfig); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close()

	// 初始化Redis连接
	redisConfig := cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	if err := cache.InitRedis(redisConfig); err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}

	// 初始化LLM客户端
	llmClient, err := model.NewLLMClient(cfg.LLM.GetAddr())
	if err != nil {
		log.Fatalf("初始化LLM客户端失败: %v", err)
	}
	defer llmClient.Close()

	// 初始化存储
	store := storage.NewStorage()
	userStorage := storage.NewUserStorage()

	// 初始化服务
	chatService := service.NewChatService(llmClient, store)

	// 初始化路由
	router := api.NewRouter(chatService, userStorage)
	handler := router.Setup()

	// 创建并启动服务器
	server := api.NewServer(handler)
	if err := server.Start(); err != nil {
		log.Fatalf("服务器错误: %v", err)
	}
}
