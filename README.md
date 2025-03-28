# Tiny-llm-chat

# Chat Llama

一个简单的AI聊天应用，提供类似ChatGPT的用户体验，支持多会话管理、实时响应和美观的蓝白主题界面。

## 功能特性

- 🤖 自训练并微调的参数量较小的大语言模型（中文），基于llama（参考baby-llama项目）
- 💬 多会话创建与管理
- 🔄 历史对话查看与继续
- 👤 用户注册与登录系统
- 🎨 清新蓝白主题界面
- 📱 响应式设计，可以支持PC和移动设备

## 技术栈

### 模型
- llama
- chatGLM分词器
- 语料包含通用语料（预训练）和医疗领域语料（微调）
- 已封装，支持后续扩展，替换成自己训练的模型

### 后端
- Golang
- Gin and Gorm框架
- gRPC (与LLM服务通信)
- JWT认证
- redis and MySQL数据库

### 前端
- React 
- CSS3

## 快速开始
'cd backend'

