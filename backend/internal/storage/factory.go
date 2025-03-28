package storage

import (
	"chat-llama/internal/service"
	"chat-llama/pkg/db"
)

// 创建存储接口的工厂方法
func NewStorage() service.Storage {
	// 检查数据库连接
	if db.DB == nil {
		panic("数据库未初始化")
	}

	// 初始化数据库表
	if err := InitTables(db.DB); err != nil {
		panic("初始化数据库表失败: " + err.Error())
	}

	// 返回MySQL存储实现
	return NewMySQLStorage()
}
