package model

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "chat-llama/internal/model/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LLMClient 封装了与 LLM 服务通信的客户端
type LLMClient struct {
	client pb.LLMServiceClient
	conn   *grpc.ClientConn
}

// NewLLMClient 创建一个新的 LLM 客户端
// serverAddr 格式为 "host:port"，例如 "localhost:50051"
func NewLLMClient(serverAddr string) (*LLMClient, error) {
	// 创建到服务器的连接
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("无法连接到 LLM 服务: %v", err)
	}

	// 创建 gRPC 客户端
	client := pb.NewLLMServiceClient(conn)
	return &LLMClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭 gRPC 连接
func (c *LLMClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GenerateResponse 调用 LLM 服务生成响应
func (c *LLMClient) GenerateResponse(ctx context.Context, prompt string, temperature float32, maxNewTokens int32, topK int32) (string, error) {
	// 创建带超时的上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 创建请求
	req := &pb.GenerateRequest{
		Prompt:       prompt,
		Temperature:  temperature,
		MaxNewTokens: maxNewTokens,
		TopK:         topK,
	}

	// 调用 gRPC 服务
	log.Printf("向 LLM 服务发送请求：prompt=%s, temperature=%.2f, maxNewTokens=%d, topK=%d",
		prompt, temperature, maxNewTokens, topK)

	resp, err := c.client.Generate(timeoutCtx, req)
	if err != nil {
		log.Printf("调用 Generate 时出错: %v", err)
		return "", err
	}

	log.Printf("收到 LLM 服务响应：%s", resp.Response)
	return resp.Response, nil
}
