# MagTrade 项目结构文档

## 根目录

```
MagTrade/
├── cmd/                    # 可执行程序入口
├── config/                 # 配置文件
├── docker/                 # Docker相关文件
├── docs/                   # 文档
├── frontend/               # Vue前端
├── internal/               # Go后端核心代码
├── migrations/             # 数据库迁移脚本
├── go.mod / go.sum         # Go依赖
├── README.md               # 项目说明
└── LICENSE                 # MIT许可证
```

---

## cmd/ - 程序入口

| 目录 | 功能 |
|------|------|
| `server/main.go` | 主服务入口，启动HTTP服务器、数据库、Redis、Kafka |
| `benchmark/main.go` | 压测工具，测试API性能 |

---

## config/ - 配置文件

| 文件 | 用途 |
|------|------|
| `config.dev.yaml` | 本地开发配置 |
| `config.prod.yaml` | 生产环境配置 |

**配置热更新**: 不支持，修改后需重启后端

---

## docker/ - 容器化

| 文件 | 用途 |
|------|------|
| `Dockerfile` | 后端镜像构建 |
| `docker-compose.dev.yml` | 开发环境依赖（Redis + Kafka） |
| `docker-compose.yml` | 生产环境全栈部署 |

---

## docs/ - 文档

| 文件 | 内容 |
|------|------|
| `DEVELOPMENT.md` | 开发指南 |
| `TEST_REPORT.md` | 测试报告 |
| `swagger.yaml` | API文档（OpenAPI 3.0） |
| `swagger.html` | Swagger UI |
| `PROJECT_STRUCTURE.md` | 本文档 |

---

## migrations/ - 数据库

| 文件 | 用途 |
|------|------|
| `001_init_schema.sql` | 初始表结构（备用，推荐用GORM自动建表） |
| `002_add_email_verified_audit_logs.sql` | 增量迁移 |
| `seed_data.sql` | 种子数据（admin用户、测试商品） |

**最佳实践**: GORM管理表结构，SQL仅用于种子数据

---

## internal/ - 后端核心

### 架构层次

```
internal/
├── handler/      # HTTP处理器（Controller层）
├── service/      # 业务逻辑（Service层）
├── repository/   # 数据访问（DAO层）
├── model/        # 数据模型（Entity）
├── middleware/   # 中间件
├── router/       # 路由配置
├── config/       # 配置加载
├── database/     # 数据库连接
├── cache/        # Redis缓存
├── mq/           # Kafka消息队列
├── logger/       # 日志
├── pkg/          # 公共工具包
└── worker/       # 后台任务
```

### 各模块详解

#### handler/ - HTTP处理器
| 文件 | 功能 |
|------|------|
| `auth_handler.go` | 登录、注册、邮箱验证码 |
| `product_handler.go` | 商品CRUD |
| `flash_sale_handler.go` | 秒杀活动、抢购 |
| `order_handler.go` | 订单管理 |
| `ai_handler.go` | AI客服、推荐 |
| `captcha_handler.go` | 图片验证码 |
| `ws_handler.go` | WebSocket通知 |
| `metrics_handler.go` | 监控指标 |

#### service/ - 业务逻辑
| 文件 | 功能 |
|------|------|
| `auth_service.go` | 认证逻辑（含验证码校验） |
| `email_service.go` | 邮件发送（验证码） |
| `captcha_service.go` | 验证码生成、验证、登录失败计数 |
| `flash_sale_service.go` | 秒杀核心逻辑 |
| `order_service.go` | 订单处理 |
| `ai/` | AI模块（智能推荐、异常检测） |

#### model/ - 数据模型
| 文件 | 对应表 |
|------|--------|
| `user.go` | users |
| `product.go` | products |
| `flash_sale.go` | flash_sales |
| `order.go` | orders |
| `chat_history.go` | chat_histories |
| `ai_recommendation.go` | ai_recommendations |

#### middleware/ - 中间件
| 文件 | 功能 |
|------|------|
| `auth.go` | JWT认证 |
| `rate_limit.go` | 限流 |
| `cors.go` | 跨域 |
| `logger.go` | 请求日志 |
| `admin.go` | 管理员权限 |

#### pkg/ - 工具包
| 目录 | 功能 |
|------|------|
| `utils/` | JWT、密码Hash、雪花ID |
| `validator/` | 输入校验 |
| `response/` | 统一响应格式 |

#### mq/ - 消息队列
| 文件 | 功能 |
|------|------|
| `producer.go` | Kafka生产者 |
| `consumer.go` | Kafka消费者 |

#### worker/ - 后台任务
| 文件 | 功能 |
|------|------|
| `order_worker.go` | 异步订单处理 |
| `scheduler_worker.go` | 定时任务 |

---

## frontend/ - Vue前端

### 技术栈
- Vue 3 + TypeScript + Vite
- Pinia（状态管理）
- TailwindCSS（样式）

### 目录结构
```
frontend/src/
├── api/          # API调用
├── components/   # 可复用组件
├── layouts/      # 布局组件
├── router/       # 路由配置
├── stores/       # Pinia状态
├── types/        # TypeScript类型
├── utils/        # 工具函数
├── views/        # 页面组件
├── App.vue       # 根组件
└── main.ts       # 入口
```

### 主要页面
| 文件 | 页面 |
|------|------|
| `HomeView.vue` | 首页（秒杀列表） |
| `LoginView.vue` | 登录（含验证码） |
| `RegisterView.vue` | 注册（含邮箱验证码） |
| `FlashSaleDetailView.vue` | 秒杀详情 |
| `OrdersView.vue` | 我的订单 |
| `AdminView.vue` | 管理后台 |

---

## 快速命令

```bash
# 重启后端（配置修改后）
Ctrl+C  # 停止
go run ./cmd/server  # 启动

# 重置数据库
$env:PGCLIENTENCODING='UTF8'
psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS magtrade;"
psql -U postgres -h localhost -c "CREATE DATABASE magtrade;"
# 然后重启后端

# 压测
go run ./cmd/benchmark -c 50 -n 500

# 单元测试
go test ./... -v
```
