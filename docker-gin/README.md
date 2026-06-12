# AI 智能单词本

## 一、项目基本信息

| 项目 | 内容 |
|------|------|
| 姓名 | 陈慧 |
| 学校 | 华中师范大学 |
| 学号 | 2023214479 |
| 项目名称 | AI 智能单词本 |
| 完成时间 | 2026年5月16日 |

## 二、开发任务索引

| 序号 | 任务 | 完成情况 |
|------|------|----------|
| 1 | 前后端分离架构 | ✅ |
| 2 | 用户注册（密码bcrypt加密，严禁明文） | ✅ |
| 3 | 用户登录（JWT Token） | ✅ |
| 4 | 前端Token存储（localStorage） | ✅ |
| 5 | AI单词查询（通义千问） | ✅ |
| 6 | 查询缓存（已保存单词直接返回） | ✅ |
| 7 | 手动保存单词 | ✅ |
| 8 | 分页单词列表 | ✅ |
| 9 | 软删除单词 | ✅ |
| 10 | 后端严禁CORS（Vite proxy + Nginx反向代理） | ✅ |
| 11 | Docker Compose三服务编排（db/backend/frontend） | ✅ |
| 12 | 手动SQL初始化（严禁AutoMigrate） | ✅ |
| 13 | API接口文档 | ✅ |
| 14 | 数据库设计文档 | ✅ |

## 三、项目简介

AI智能单词本是一个前后端分离的英语学习Web应用。用户通过前端页面查询单词，后端调用通义千问AI生成单词的精准释义和3条例句。查询结果返回前端后，用户可手动点击“保存”按钮，将单词持久化到个人单词本中，方便日后复习。

### 架构说明

本项目采用前后端分离架构，通过Docker Compose编排三个服务：

1. **前端服务（frontend）**：基于Vue3 + Vite构建，使用Nginx作为Web服务器。生产环境中Nginx同时提供静态文件服务和反向代理功能，将`/api/`请求转发至后端容器。

2. **后端服务（backend）**：基于Go + Gin框架，提供RESTful API。包含JWT认证中间件、用户注册登录、AI单词查询、单词保存、单词列表分页、软删除等功能。

3. **数据库服务（db）**：MySQL 8.0，存储用户信息和单词本数据。通过GORM连接，表结构使用init.sql手动初始化。

**跨域处理**：后端代码严禁配置CORS中间件。开发环境通过Vite proxy解决跨域，生产环境通过Nginx反向代理实现同源访问。

**数据流向**：用户 → Nginx:80 → 静态页面或代理转发 → Backend:8080 → MySQL:3306

## 四、运行指南

### 4.1 前置依赖

| 软件 | 版本要求 | 检查命令 |
|------|----------|----------|
| Docker Desktop | 20.10+ | `docker --version` |
| Docker Compose | 2.0+ | `docker-compose --version` |

### 4.2 配置 AI 的 API Key

本项目使用**阿里云百炼（通义千问）**服务。

**步骤1：获取 API Key**

1. 访问 https://bailian.console.aliyun.com/
2. 登录阿里云账号
3. 开通百炼服务
4. 在 API Key 管理页面创建密钥（格式：`sk-xxx`）

**步骤2：创建 .env 文件**

在项目根目录 `week05/homework/docker-gin/` 下创建 `.env` 文件：

```bash
cd week05/homework/docker-gin
notepad .env
```

写入以下内容：

```
QWEN_API_KEY=sk-你的百炼API密钥
```

### 4.3一键启动

```
cd week05/homework/docker-gin
docker-compose up -d
```

首次启动会自动完成：

1. 拉取 MySQL、Nginx、Golang 镜像
2. 构建后端 Go 应用
3. 构建前端 Vue 应用
4. 初始化数据库（执行 `docs/init.sql`）
5. 启动三个容器

### 4.4验证服务状态

```
docker-compose ps
```

预期输出三个容器状态均为 Up。

### 4.5访问服务

| 服务     | 地址                                  | 说明                |
| :------- | :------------------------------------ | :------------------ |
| 前端页面 | [http://localhost](http://localhost/) | 主要访问入口        |
| 后端 API | http://localhost/api                  | 通过 Nginx 代理访问 |