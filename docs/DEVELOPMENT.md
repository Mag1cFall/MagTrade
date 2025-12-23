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

# 2. 安装Go依赖
go mod tidy

# 3. 安装前端依赖
cd frontend && npm install && cd ..

# 4. 启动Docker依赖（Redis + Kafka）
docker-compose -f docker/docker-compose.dev.yml up -d

# 5. 创建数据库
$env:PGCLIENTENCODING='UTF8'
psql -U postgres -h localhost -c "CREATE DATABASE magtrade;"

# 6. 启动后端（GORM自动建表）
go run ./cmd/server

# 7. 另开终端，插入种子数据
$env:PGCLIENTENCODING='UTF8'
psql -U postgres -h localhost -d magtrade -f migrations/seed_data.sql

# 8. 启动前端
cd frontend && npm run dev
```

访问: http://localhost:5173

---

## 二、本地开发日常操作

### 干净重置数据库
```powershell
# PowerShell
$env:PGCLIENTENCODING='UTF8'
psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS magtrade;"
psql -U postgres -h localhost -c "CREATE DATABASE magtrade;"
# 重启后端让GORM建表
# 然后执行: psql -U postgres -h localhost -d magtrade -f migrations/seed_data.sql
```

### 每日启动流程
```bash
# 1. 启动Docker依赖（如果未运行）
docker-compose -f docker/docker-compose.dev.yml up -d

# 2. 启动后端
go run ./cmd/server

# 3. 启动前端（另一个终端）
cd frontend && npm run dev
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
# 在VPS上
git clone https://github.com/Mag1cFall/magtrade.git
cd magtrade

# 修改生产配置
cp config/config.prod.yaml config/config.yaml
# 编辑config.yaml修改数据库密码、JWT密钥等

# 启动全部服务
docker-compose -f docker/docker-compose.yml up -d

# 查看日志
docker-compose logs -f app
```

### 手动部署
```bash
# 1. 构建后端
go build -o magtrade ./cmd/server

# 2. 构建前端
cd frontend && npm run build

# 3. 启动（需要单独安装PostgreSQL、Redis、Kafka）
./magtrade
```

---

## 四、数据库管理最佳实践

### 架构说明
- **GORM**: 管理表结构（自动迁移）
- **SQL脚本**: 仅用于插入种子数据

### 为什么这样设计？
GORM会自动创建索引，命名格式如`uni_users_username`。如果手动SQL创建名为`idx_users_username`的索引，两者不匹配会导致迁移报错。

**解决方案（已实现）**: 在GORM模型中显式指定索引名称与SQL一致：
```go
Username string `gorm:"uniqueIndex:idx_users_username"`
```

### 迁移流程
1. 修改Go模型（如添加字段）
2. 重启后端，GORM自动迁移
3. 如需初始数据，执行seed_data.sql

---

## 五、常见问题

### 中文乱码
```powershell
$env:PGCLIENTENCODING='UTF8'
```

### Redis连接失败
```bash
docker-compose -f docker/docker-compose.dev.yml up -d
```

### GORM约束错误
```
错误: 关系 "users" 的 约束"uni_users_username" 不存在
```
原因: 数据库索引名与GORM期望不一致
解决: 删库重建（让GORM创建），或检查模型索引名

### admin登录失败
seed_data.sql中的密码hash可能过期，用项目工具重新生成：
```bash
# 创建临时文件 cmd/genhash/main.go
go run ./cmd/genhash
# 用输出的hash更新数据库
```

---

## 六、默认账户

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |

---

## 七、API端点

| 功能 | 端点 |
|------|------|
| 发送邮箱验证码 | POST /api/v1/auth/send-code |
| 注册 | POST /api/v1/auth/register |
| 登录 | POST /api/v1/auth/login |
| 商品列表 | GET /api/v1/products |
| 秒杀列表 | GET /api/v1/flash-sales |
| 参与秒杀 | POST /api/v1/flash-sales/:id/rush |
| 我的订单 | GET /api/v1/orders |
| 创建商品(admin) | POST /api/v1/admin/products |
| 创建秒杀(admin) | POST /api/v1/admin/flash-sales |
