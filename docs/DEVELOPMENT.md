# MagTrade 完整操作指南

## 项目概述

MagTrade是**高并发分布式秒杀系统**，面向电商场景的限时抢购平台。

### 角色说明
- **普通用户**: 浏览商品、参与秒杀、管理订单
- **管理员admin**: 发布商品、创建秒杀活动、查看数据

---

## 一、从零开始（GitHub克隆后）

```bash
# 1. 克隆项目
git clone https://github.com/Mag1cFall/magtrade.git
cd magtrade

# 2. 创建本地配置文件并填写密钥
cp config/config.dev.yaml.example config/config.dev.yaml
# 编辑 config.dev.yaml 填入 API Key、邮箱密码等

# 3. 安装依赖
go mod tidy
cd frontend && npm install && cd ..

# 4. 启动Docker基础设施（PostgreSQL + Redis + Kafka）
docker-compose -f docker/docker-compose.dev.yml up -d

# 5. 启动后端（GORM自动建表，SeedData自动创建admin用户）
go run ./cmd/server

# 6. 另开终端，启动前端
cd frontend && npm run dev
```

访问: http://localhost:5173

---

## 二、本地开发日常操作

### 每日启动流程
```bash
# 1. 启动Docker服务
docker-compose -f docker/docker-compose.dev.yml up -d

# 2. 启动后端
go run ./cmd/server

# 3. 启动前端（另一个终端）
cd frontend && npm run dev
```

### 完全重置数据库
```bash
# 删除Docker数据卷重建
docker-compose -f docker/docker-compose.dev.yml down -v
docker-compose -f docker/docker-compose.dev.yml up -d
# 重启后端让GORM重建表
```

### 运行测试
```bash
go test ./... -v        # 全部测试
go test ./... -short    # 快速测试
```

### 运行压测
```bash
go run ./cmd/benchmark -c 50 -n 500
```

---

## 三、部署到VPS

### Docker一键部署（推荐）
```bash
git clone https://github.com/Mag1cFall/magtrade.git
cd magtrade

# 设置生产环境变量（必须）
export DB_PASSWORD=your_secure_db_password
export JWT_SECRET=your_jwt_secret_key
export AI_API_KEY=your_ai_api_key
export ADMIN_INIT_PASSWORD=your_admin_password  # 首次启动必需

# 启动全部服务
docker-compose -f docker/docker-compose.yml up -d

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f backend
```

> **注意**: `ADMIN_INIT_PASSWORD` 仅在首次启动（users表为空）时使用，之后可移除。

---

## 四、数据库管理

### 架构说明
- **Docker PostgreSQL**: 自动创建数据库 + 导入种子数据
- **GORM**: 管理表结构（自动迁移）
- **seed_data.sql**: 初始数据（admin用户、测试商品）

### 迁移流程
1. 修改Go模型（如添加字段）
2. 重启后端，GORM自动迁移

---

## 五、常见问题

### Docker服务未启动
```bash
docker-compose -f docker/docker-compose.dev.yml up -d
docker-compose -f docker/docker-compose.dev.yml ps  # 检查状态
```

### 重置一切
```bash
docker-compose -f docker/docker-compose.dev.yml down -v
docker-compose -f docker/docker-compose.dev.yml up -d
```

---

## 六、默认账户

### 开发环境
| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |

### 生产环境
首次启动时通过环境变量 `ADMIN_INIT_PASSWORD` 设置admin密码。

---

## 七、Docker服务

| 服务 | 端口 | 说明 |
|------|------|------|
| PostgreSQL | 5432 | 数据库（自动创建） |
| Redis | 6379 | 缓存 |
| Kafka | 9092 | 消息队列 |

---

## 八、API端点

| 功能 | 端点 |
|------|------|
| 发送邮箱验证码 | POST /api/v1/auth/send-code |
| 注册 | POST /api/v1/auth/register |
| 登录 | POST /api/v1/auth/login |
| AI聊天(流式) | POST /api/v1/ai/chat/stream |
| 商品列表 | GET /api/v1/products |
| 秒杀列表 | GET /api/v1/flash-sales |
| 参与秒杀 | POST /api/v1/flash-sales/:id/rush |
| 我的订单 | GET /api/v1/orders |
