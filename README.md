# 金山实习全栈项目集

> 华中师范大学 · 计算机科学与技术 · 陈慧 · 2023214479

本仓库收录了我在**金山实习期间**完成的 7 个全栈开发项目，涵盖后端服务、前端应用、AI 集成、容器化部署、系统监控等多个技术领域。

---

## 📊 项目总览

| 序号 | 项目名称                                                     | 技术栈                     | 核心亮点                              |
| :--: | ------------------------------------------------------------ | -------------------------- | ------------------------------------- |
|  1   | [智能阶梯计费系统](#1-智能阶梯计费系统-billing)              | Go                         | 阶梯电价 + 峰谷时段调节 + 单元测试    |
|  2   | [在线学习管理平台](#2-在线学习管理平台-course)               | React + Koa + SQLite       | JWT 认证 + ECharts 可视化 + 完整 CRUD |
|  3   | [AI 智能单词本](#3-ai-智能单词本-docker-gin)                 | Vue3 + Go + MySQL + Docker | 通义千问 AI + 三服务容器编排          |
|  4   | [Gin-Vue-Admin 二次开发](#4-gin-vue-admin-二次开发-gin-fullstack) | Gin + Vue3 + SQLite        | 数据库迁移 + 用户行为追踪             |
|  5   | [灵犀 AI 对话助手](#5-灵犀-ai-对话助手-lingxi)               | 原生 JS + 通义千问         | 流式响应 + Markdown 渲染 + 主题切换   |
|  6   | [服务健康探测器](#6-服务健康探测器-monitor)                  | Go                         | 高并发探测 + 超时控制 + 报表生成      |
|  7   | [Gin 全栈应用](#7-gin-全栈应用-gin-fullstack-示例)           | Gin                        | RESTful API + 静态文件服务            |

---

## 🚀 项目详情

### 1. 智能阶梯计费系统 (`billing`)

**功能特性**：
- 阶梯电价计算（3档：0.5/0.8/1.2元）
- 峰谷时段调节（高峰+10% / 低谷-20%）
- 异常处理（时段格式错误默认高峰）
- 单元测试 27 个用例全部通过

**技术栈**：Go + Go testing

**快速运行**：

```bash
cd billing
go run billing.go
```

### 2. 在线学习管理平台 (`course`)

**功能特性**：

- 用户认证：JWT Token + 路由守卫
- 工作台：ECharts 数据可视化（柱状图/折线图/饼图）
- 课程管理：分页/搜索/筛选/增删改查/发布下架
- 学生管理：分页/搜索/筛选/增删改查/多选课程关联
- 学习总结：Markdown 渲染（表格 + 代码高亮）

**技术栈**：React 18 + TypeScript + Vite + Ant Design + ECharts + Koa + SQLite + JWT

**快速运行**：

bash

```
cd course/client && npm install && npm run dev
cd course/server && npm install && npm run dev
```



------

### 3. AI 智能单词本 (`docker-gin`)

**功能特性**：

- 用户注册（bcrypt 加密）/ 登录（JWT）
- AI 单词查询（通义千问 qwen-vl-plus）
- 查询缓存（已查单词直接返回）
- 单词本管理（保存/列表/软删除）
- Docker Compose 三服务编排（MySQL + Backend + Nginx）

**技术栈**：Vue3 + Vite + Go + Gin + GORM + MySQL + Docker + 阿里云百炼

**快速运行**：

bash

```
cd docker-gin
# 配置 .env 文件（QWEN_API_KEY）
docker-compose up -d
```



**访问地址**：[http://localhost](http://localhost/)

------

### 4. Gin-Vue-Admin 二次开发 (`gin-fullstack`)

**改造内容**：

- 数据库迁移：MySQL → SQLite
- 用户行为追踪：新增 `LastLoginIP` + `LastLoginTime` 字段
- 登录时记录真实 IP 和时间
- 前端用户列表新增登录信息列

**技术栈**：Gin + GORM + SQLite + Vue3 + Element Plus + Gin-Vue-Admin

**快速运行**：

bash

```
cd gin-fullstack/server && go run main.go
cd gin-fullstack/web && npm install && npm run dev
```



------

### 5. 灵犀 AI 对话助手 (`lingxi`)

**功能特性**：

- 多模态对话（文本 + 图片）
- 流式响应 + 打字机效果
- Markdown 渲染（marked.js）
- 代码语法高亮 + 一键复制（Prism.js）
- 深色/浅色主题切换（localStorage 持久化）
- 图片上传预览 + 大图弹窗
- 停止生成 / 清除对话
- 响应式适配（PC/平板/手机）

**技术栈**：原生 HTML/CSS/JS + 阿里云百炼（通义千问 qwen-vl-plus）

**快速运行**：

```
cd lingxi
# 使用本地服务器打开 index.html
npx serve .
```



------

### 6. 服务健康探测器 (`monitor`)

**功能特性**：

- 支持 HTTP/TCP 协议探测
- 可配置重试机制（≤3次）
- 高并发探测引擎（goroutine + channel）
- 超时控制（context.WithTimeout）
- 命令行参数支持（`--config` / `--timeout` / `-v`）
- 报表生成（成功/失败比例 + 响应时间分布 + 最慢 TOP3）
- 日志文件输出

**技术栈**：Go 1.21+

**快速运行**：

```
cd monitor
go run main.go --config config.json -v
go test -v  # 运行测试
```



------

### 7. Gin 全栈应用 (`gin-fullstack` 示例)

**功能特性**：

- RESTful API 设计
- 前端静态文件服务
- 完整项目结构示例

**技术栈**：Gin

------

## 🛠️ 技术栈总览

| 类别          | 技术                                |
| :------------ | :---------------------------------- |
| **后端语言**  | Go、Node.js (Koa)                   |
| **前端框架**  | React 18、Vue3、原生 JavaScript     |
| **UI 组件库** | Ant Design、Element Plus            |
| **构建工具**  | Vite                                |
| **数据库**    | MySQL、SQLite                       |
| **ORM**       | GORM                                |
| **AI 服务**   | 阿里云百炼（通义千问 qwen-vl-plus） |
| **容器化**    | Docker、Docker Compose              |
| **认证**      | JWT、bcrypt                         |
| **可视化**    | ECharts                             |
| **Markdown**  | marked.js、react-markdown           |
| **代码高亮**  | Prism.js                            |

------

## 📁 目录结构

```
jinshan-office-fullstack/
├── billing/          # 智能阶梯计费系统
├── course/           # 在线学习管理平台
├── docker-gin/       # AI 智能单词本（Docker 编排）
├── gin-fullstack/    # Gin-Vue-Admin 二次开发
├── lingxi/           # 灵犀 AI 对话助手
├── monitor/          # 服务健康探测器
└── README.md         # 本文件
```



------

## 👩‍💻 作者信息

| 项目     | 内容                |
| :------- | :------------------ |
| **姓名** | 陈慧                |
| **学校** | 华中师范大学（211） |
| **学号** | 2023214479          |
| **专业** | 计算机科学与技术    |

------

## ⚠️ 免责声明

- 本仓库仅用于**个人作品集展示**
- 不包含任何公司敏感信息
- 所有 API Key 均已移除，使用时请自行配置

*最后更新：2026年6月*