package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"chat-llama/pkg/cache"
	"chat-llama/pkg/db"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserStorage 用户存储实现
type UserStorage struct {
	db *gorm.DB
}

// NewUserStorage 创建用户存储
func NewUserStorage() *UserStorage {
	return &UserStorage{
		db: db.DB,
	}
}

// 用户缓存键
func userKey(id uint) string {
	return fmt.Sprintf("user:%d", id)
}

func userByUsernameKey(username string) string {
	return fmt.Sprintf("user:name:%s", username)
}

// CreateUser 创建用户
func (s *UserStorage) CreateUser(username, password string) (*User, error) {
	ctx := context.Background()

	// 检查用户名是否已存在
	var count int64
	s.db.Model(&User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &User{
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// 清除可能存在的缓存
	cache.Delete(ctx, userByUsernameKey(username))

	return user, nil
}

// GetUserByID 通过ID获取用户
func (s *UserStorage) GetUserByID(id uint) (*User, error) {
	ctx := context.Background()

	// 尝试从缓存获取
	var user User
	found, err := cache.Get(ctx, userKey(id), &user)
	if err != nil {
		return nil, err
	}

	if found {
		return &user, nil
	}

	// 从数据库获取
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 更新缓存
	cache.Set(ctx, userKey(id), user, time.Hour*24)

	return &user, nil
}

// GetUserByUsername 通过用户名获取用户
func (s *UserStorage) GetUserByUsername(username string) (*User, error) {
	ctx := context.Background()

	log.Printf("尝试获取用户: %s", username)

	// 尝试从缓存获取
	var user User
	found, err := cache.Get(ctx, userByUsernameKey(username), &user)
	if err != nil {
		log.Printf("获取缓存失败: %v", err)
		// 继续从数据库获取
	} else if found {
		log.Printf("从缓存获取用户: %s 成功", username)
		return &user, nil
	}

	// 从数据库获取
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("用户 %s 不存在", username)
			return nil, errors.New("用户不存在")
		}
		log.Printf("数据库错误: %v", err)
		return nil, err
	}

	log.Printf("从数据库获取用户: %s 成功, ID=%d", username, user.ID)

	// 更新缓存
	cache.Set(ctx, userByUsernameKey(username), user, time.Hour*24)
	cache.Set(ctx, userKey(user.ID), user, time.Hour*24)

	return &user, nil
}

// VerifyPassword 验证用户密码
func (s *UserStorage) VerifyPassword(username, password string) (*User, error) {
	log.Printf("验证用户密码: %s", username)

	// 获取用户
	user, err := s.GetUserByUsername(username)
	if err != nil {
		log.Printf("获取用户失败: %v", err)
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// 密码验证失败，清除缓存后再尝试从数据库获取新数据
		ctx := context.Background()
		cache.Delete(ctx, userByUsernameKey(username))
		cache.Delete(ctx, userKey(user.ID))

		// 重新从数据库获取用户信息
		if err := s.db.Where("username = ?", username).First(user).Error; err != nil {
			return nil, errors.New("用户不存在")
		}

		// 再次验证密码
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return nil, errors.New("密码错误")
		}
	}

	return user, nil
}
