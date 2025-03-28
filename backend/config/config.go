package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LLMConfig LLM服务配置
type LLMConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

var cfg *Config

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 创建新的viper实例
	v := viper.New()

	// 设置配置文件名和路径
	v.SetConfigName("config")

	if configPath != "" {
		// 使用指定的配置路径
		v.AddConfigPath(configPath)
	} else {
		// 默认配置路径
		v.AddConfigPath("./config")
		v.AddConfigPath(".")
	}

	// 设置配置文件类型
	v.SetConfigType("yaml")

	// 从环境变量读取
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析到配置结构
	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 设置全局配置
	cfg = config

	return config, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if cfg == nil {
		var err error
		cfg, err = LoadConfig("")
		if err != nil {
			log.Fatalf("无法加载配置: %v", err)
		}
	}
	return cfg
}

// 获取数据库DSN
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// 获取Redis地址
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// 获取LLM服务地址
func (c *LLMConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// 初始化日志目录
func (c *LogConfig) InitLogDir() error {
	if c.Path == "" {
		c.Path = "./logs"
	}

	// 创建日志目录
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(c.Path, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %w", err)
		}
	}

	return nil
}

// 获取日志文件路径
func (c *LogConfig) GetLogFilePath() string {
	return filepath.Join(c.Path, "app.log")
}
