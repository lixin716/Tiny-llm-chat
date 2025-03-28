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

首先把项目clone到本地,然后在backend/config/config.yaml中更改自己的配置

### 模型服务器启动
进入model_api文件夹,启动模型服务（确保需要的依赖都已安装）
>cd model/model_api

>python grpc_server.py

启动成功会显示模型加载成功，正在监听端口xxx

### 后端服务启动
首先确保已经更改配置为自己的，然后启动redis和mysql，最后启动后端服务
>redis-server
>cd backend
>go mod tidy
>go run main.go
启动成功会显示
>服务加载成功，http://localhost:xxxx

### 前端启动
>cd site
>npm install
>npm start
之后会自动打开浏览器相应端口，服务正常启动
