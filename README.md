# MagTrade 高并发分布式秒杀系统

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24-blue.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/Vue-3.5-42b883.svg" alt="Vue Version">
  <img src="https://img.shields.io/badge/Gin-1.11-green.svg" alt="Gin Version">
  <img src="https://img.shields.io/badge/PostgreSQL-16-blue.svg" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Redis-7-red.svg" alt="Redis">
  <img src="https://img.shields.io/badge/Kafka-7.5-orange.svg" alt="Kafka">
  <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License">
</p>

## 📖 项目简介

MagTrade 是一个面向高并发场景的分布式秒杀系统，支持百万级并发请求下的商品抢购。系统集成 **AI Agent** 模块，提供智能客服、秒杀策略推荐、异常检测等增值能力。

### ✨ 核心特性

- 🚀 **高并发秒杀**：Redis 预扣库存 + Lua 原子操作 + Kafka 异步下单
- 🔒 **分布式锁**：基于 Redis 的分布式锁防止超卖
- 📊 **流量削峰**：Kafka 消息队列解耦请求与订单处理
- 🛡️ **安全防护**：登录验证码、失败锁定、邮箱验证、IP 限流
- 🤖 **AI Agent**：智能客服、策略推荐、异常检测、流式对话
- 📡 **实时通知**：WebSocket 推送秒杀结果和订单状态
- 🐳 **容器化部署**：多阶段 Dockerfile + Docker Compose 一键部署
- 🎨 **现代前端**：Vue 3 + TypeScript + Tailwind CSS 赛博朋克风格

## 🏗️ 技术栈

### 后端
| 层级 | 技术 | 说明 |
|------|------|------|
| Web Framework | Gin 1.11 | 高性能 HTTP 框架 |
| ORM | GORM 1.31 | PostgreSQL 操作 |
| 缓存 | Redis 7 | 库存缓存、分布式锁、验证码 |
| 消息队列 | Kafka 7.5 | 异步下单、流量削峰 |
| 数据库 | PostgreSQL 16 | 持久化存储 |
| 认证 | JWT | 双 Token 无状态认证 |
| AI | SiliconFlow + DeepSeek(默认) | 智能对话与决策 |
| 日志 | Zap | 结构化高性能日志 |
| ID 生成 | Snowflake | 分布式唯一 ID |

### 前端
| 层级 | 技术 | 说明 |
|------|------|------|
| 框架 | Vue 3.5 | Composition API |
| 语言 | TypeScript | 类型安全 |
| 状态管理 | Pinia | 轻量状态管理 |
| 样式 | Tailwind CSS 4 | 原子化 CSS |
| 路由 | Vue Router 4 | SPA 路由 |
| 构建 | Vite 7 | 极速开发服务器 |

## 📁 项目结构

```
magtrade/
├── cmd/
│   ├── server/main.go              # 应用入口
│   └── benchmark/main.go           # 压测工具
├── internal/
│   ├── cache/                      # Redis 操作
│   │   ├── redis.go                # 连接管理
│   │   ├── scripts.go              # Lua 脚本
│   │   └── stock.go                # 库存服务 + 分布式锁
│   ├── config/                     # 配置管理
│   ├── database/                   # PostgreSQL 连接池
│   ├── handler/                    # HTTP 处理器
│   │   ├── auth_handler.go         # 认证
│   │   ├── flash_sale_handler.go   # 秒杀
│   │   ├── order_handler.go        # 订单
│   │   ├── product_handler.go      # 商品
│   │   ├── ai_handler.go           # AI 对话
│   │   ├── ws_handler.go           # WebSocket
│   │   └── ...                     # 验证码、上传、指标
│   ├── logger/                     # Zap 日志
│   ├── middleware/                 # 中间件
│   │   ├── auth.go                 # JWT 认证
│   │   ├── ratelimit.go            # 令牌桶限流
│   │   ├── security.go             # 安全防护
│   │   ├── idempotency.go          # 幂等性检查
│   │   └── ...                     # CORS、审计、日志
│   ├── model/                      # 数据模型 (User, Product, Order...)
│   ├── mq/                         # Kafka 生产者/消费者
│   ├── pkg/                        # 工具包
│   │   ├── utils/                  # JWT、密码、Snowflake
│   │   ├── response/               # 统一响应格式
│   │   └── validator/              # 数据验证
│   ├── repository/                 # 数据访问层
│   ├── router/                     # 路由配置
│   ├── service/                    # 业务逻辑层
│   │   ├── flash_sale_service.go   # 秒杀核心
│   │   ├── order_service.go        # 订单处理
│   │   ├── auth_service.go         # 认证服务
│   │   └── ai/                     # AI Agent 模块
│   │       ├── llm_client.go       # LLM 客户端
│   │       ├── customer_service.go # 智能客服
│   │       ├── strategy_advisor.go # 策略推荐
│   │       └── anomaly_detector.go # 异常检测
│   └── worker/                     # 后台定时任务
├── frontend/                       # Vue 3 前端
│   ├── src/
│   │   ├── api/                    # API 封装
│   │   ├── components/             # 30+ 组件
│   │   ├── layouts/                # 布局组件
│   │   ├── router/                 # 路由配置
│   │   ├── stores/                 # Pinia 状态
│   │   ├── views/                  # 10 个页面视图
│   │   └── style.css               # 全局样式
│   └── public/fonts/               # 自定义字体 (Git LFS)
├── config/                         # YAML 配置文件
├── migrations/                     # SQL 迁移脚本
│   └── seeds/                      # 种子数据
├── docker/                         # Docker Compose
├── docs/                           # 文档 (Swagger, 开发指南)
└── .github/workflows/              # CI/CD
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 16+
- Redis 7+
- Docker & Docker Compose

### 本地开发

1. **克隆项目**
```bash
git clone https://github.com/Mag1cFall/magtrade.git
cd magtrade
```

2. **启动依赖服务 (Redis + Kafka)**
```bash
docker-compose -f docker/docker-compose.dev.yml up -d
```

3. **创建数据库**
```bash
# 使用 psql 连接本地 PostgreSQL
psql -U postgres -f migrations/001_init_schema.sql
```

4. **安装依赖并编译**
```bash
go mod tidy
go build -o bin/magtrade.exe ./cmd/server
```

5. **启动服务**
```bash
# Windows
.\bin\magtrade.exe

# Linux/macOS
./bin/magtrade
```

服务启动后访问：http://localhost:8080/health

### Docker 部署

```bash
# 一键启动所有服务
docker-compose -f docker/docker-compose.yml up -d

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f app
```

## 📡 API 文档

### 验证码与安全模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/v1/captcha` | 获取验证码图片 | ❌ |
| GET | `/api/v1/captcha/check` | 检查是否需要验证码 | ❌ |
| GET | `/metrics` | 数据库连接池监控 | ❌ |

### 认证模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/send-code` | 发送邮箱验证码（60秒冷却，15分钟有效） | ❌ |
| POST | `/api/v1/auth/register` | 用户注册（需邮箱验证码） | ❌ |
| POST | `/api/v1/auth/login` | 用户登录（失败3次需验证码，5次锁定15分钟） | ❌ |
| POST | `/api/v1/auth/refresh` | 刷新Token | ✅ |
| GET | `/api/v1/auth/me` | 获取当前用户 | ✅ |

### 商品模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/v1/products` | 商品列表 | ❌ |
| GET | `/api/v1/products/:id` | 商品详情 | ❌ |

### 秒杀模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/v1/flash-sales` | 秒杀活动列表 | ❌ |
| GET | `/api/v1/flash-sales/:id` | 活动详情 | ❌ |
| GET | `/api/v1/flash-sales/:id/stock` | 实时库存 | ❌ |
| POST | `/api/v1/flash-sales/:id/rush` | 🔥 秒杀抢购 | ✅ |

### 订单模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/v1/orders` | 我的订单 | ✅ |
| GET | `/api/v1/orders/:order_no` | 订单详情 | ✅ |
| POST | `/api/v1/orders/:order_no/pay` | 支付订单 | ✅ |
| POST | `/api/v1/orders/:order_no/cancel` | 取消订单 | ✅ |

### AI Agent 模块 🤖

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/v1/ai/chat` | 智能客服对话 | ✅ |
| POST | `/api/v1/ai/chat/stream` | 流式对话 (SSE) | ✅ |
| GET | `/api/v1/ai/chat/history` | 对话历史 | ✅ |
| DELETE | `/api/v1/ai/chat/history` | 清除对话历史 | ✅ |
| GET | `/api/v1/ai/recommendations/:id` | 策略推荐 | ✅ |

### 管理员模块 🔐

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/v1/admin/products` | 创建商品 | ✅ Admin |
| PUT | `/api/v1/admin/products/:id` | 更新商品 | ✅ Admin |
| DELETE | `/api/v1/admin/products/:id` | 删除商品 | ✅ Admin |
| POST | `/api/v1/admin/flash-sales` | 创建秒杀活动 | ✅ Admin |
| POST | `/api/v1/admin/upload` | 上传图片 | ✅ Admin |
| POST | `/api/v1/admin/ai/analyze/:id` | 触发 AI 分析 | ✅ Admin |

### WebSocket

| 路径 | 描述 |
|------|------|
| `/ws/notifications?token=xxx` | 实时通知 |

## 🔧 配置说明

配置文件位于 `config/` 目录：

- `config.dev.yaml` - 开发环境
- `config.prod.yaml` - 生产环境

通过环境变量 `APP_ENV` 切换环境（默认 `dev`）。

### 关键配置项

```yaml
# AI 配置
ai:
  provider: "siliconflow"
  base_url: "https://api.siliconflow.cn/v1"
  api_key: "your-api-key"
  model: "deepseek-ai/DeepSeek-V3"

# JWT 配置
jwt:
  secret: "your-secret-key"
  access_token_expire: "2h"
  refresh_token_expire: "168h"
```

## 🧪 测试

```bash
# 运行单元测试
go test ./... -v

# 运行测试并生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📊 性能指标

| 指标 | 目标值 |
|------|--------|
| 秒杀接口 QPS | ≥ 5000 |
| 平均响应时间 | ≤ 50ms |
| 错误率 | ≤ 0.1% |
| 库存准确性 | 100% |

## 🤖 AI Agent 功能

### 1. 智能客服
- 回答秒杀活动、商品信息问题
- 实时查询库存和订单状态
- 多轮对话上下文理解
- **流式响应** (SSE) 实时输出

### 2. 策略推荐
- 分析活动热度和抢购难度
- 提供最佳抢购时机建议
- 预测成功概率
- 个性化策略分析

### 3. 异常检测
- IP 频率限制 (令牌桶算法)
- 机器人行为识别
- 多设备异常检测
- 秒杀请求模式分析

### 4. 后续扩展方向
- 价格预测 (基于历史数据)
- 个性化商品推荐
- 自动补货建议
- 风控规则自学习

## 📝 License

MIT License - 详见 [LICENSE](LICENSE) 文件
