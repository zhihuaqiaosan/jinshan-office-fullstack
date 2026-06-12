# Gin-Vue-Admin 二次开发项目

## 📋 项目基本信息

| 项目信息 | 内容 |
|---------|------|
| **学校名称** | 华中师范大学 |
| **姓名** | 陈慧 |
| **学号** | 2023214479 |
| **开发日期** | 2026年4月30日 |
| **项目名称** | Gin-Vue-Admin 全栈后台管理系统 |
| **项目地址** | `week06/homework/gin-fullstack` |

---

## ✅ 开发任务清单

### 任务1：环境搭建与初始化 - 数据库迁移到SQLite

- [x] **1.1 数据库切换配置**
  - 修改 `server/config.yaml` 配置文件
  - 将数据库类型从 MySQL 切换为 SQLite
  - 配置 SQLite 数据库文件存放路径

- [x] **1.2 SQLite驱动安装**
  - 安装 `gorm.io/driver/sqlite` 驱动
  - 更新 `go.mod` 依赖文件

- [x] **1.3 系统启动验证**
  - 后端服务启动成功（使用 `go run main.go`）
  - 前端服务启动成功（使用 `npm run dev`）
  - 浏览器访问初始化页面，选择 SQLite 完成初始化

### 任务2：用户行为追踪功能开发

- [x] **2.1 数据库字段扩展**
  - 在 `SysUser` 模型中添加 `LastLoginIP` 字段（最后登录IP）
  - 在 `SysUser` 模型中添加 `LastLoginTime` 字段（最后登录时间）
  - 使用 GORM AutoMigrate 自动迁移数据库

- [x] **2.2 后端登录逻辑改造**
  - 修改登录服务方法，获取客户端真实IP
  - 登录成功后更新用户的最后登录IP和时间
  - 添加错误处理，确保登录记录失败不影响正常登录

- [x] **2.3 用户列表接口调整**
  - 修改用户列表查询逻辑，返回新增的字段
  - 确保数据能正确传递给前端

- [x] **2.4 前端用户管理页面改造**
  - 在用户列表表格中新增"登录IP"列
  - 在用户列表表格中新增"登录时间"列
  - 使用 `dayjs` 库格式化时间：`YYYY-MM-DD HH:mm`

- [x] **2.5 功能验证**
  - 管理员登录创建测试用户a
  - 另一浏览器登录用户a
  - 管理员页面查看用户a的登录信息

---

## 🔧 核心技术实现

### 任务1：SQLite数据库迁移

#### 技术选型
- **ORM框架**：GORM v1.25.2
- **SQLite驱动**：gorm.io/driver/sqlite
- **配置管理**：viper + yaml

#### 实现思路

1. **配置文件修改**
```yaml
# server/config.yaml
db-list:
  - disable: false
    type: "sqlite"
    path: "./db"
    db-name: "gva.db"
    max-idle-conns: 10
    max-open-conns: 100