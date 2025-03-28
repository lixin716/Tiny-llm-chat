package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// Config Redis配置
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// InitRedis 初始化Redis连接
func InitRedis(config Config) error {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("连接Redis失败: %v", err)
	}

	return nil
}

// Set 设置缓存，默认过期时间1小时
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = time.Hour
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, jsonData, expiration).Err()
}

// Get 获取缓存并解析到目标结构
func Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 缓存未命中
		}
		return false, err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Delete 删除缓存
func Delete(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// ClearPattern 清除匹配模式的所有键
func ClearPattern(ctx context.Context, pattern string) error {
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		err := RedisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
