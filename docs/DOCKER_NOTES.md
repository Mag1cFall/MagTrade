# MagTrade Docker 完整操作笔记

## 一、Docker Compose 文件区分

| 文件 | 用途 | 启动内容 |
|------|------|----------|
| `docker-compose.dev.yml` | 开发环境 | 仅基础设施（PostgreSQL、Redis、Kafka） |
| `docker-compose.yml` | 生产环境 | 全部服务（前后端 + 基础设施） |

**为什么开发环境不启动前后端？**
- 热更新：本地 `go run` 和 `npm run dev` 代码改动即时生效
- 调试方便：可以断点调试
- 构建耗时：每次改动都重新构建Docker镜像太慢

---

## 二、开发环境启动流程

```bash
# 1. 启动基础设施
docker-compose -f docker/docker-compose.dev.yml up -d

# 2. 验证服务健康
docker-compose -f docker/docker-compose.dev.yml ps

# 3. 本地启动后端
go run ./cmd/server

# 4. 本地启动前端（另一个终端）
cd frontend && npm run dev
```

**停止服务：**
```bash
docker-compose -f docker/docker-compose.dev.yml down      # 停止
docker-compose -f docker/docker-compose.dev.yml down -v   # 停止+删除数据
```

---

## 三、生产环境一键部署

```bash
# 设置环境变量
export DB_PASSWORD=your_secure_password
export JWT_SECRET=your_jwt_secret
export AI_API_KEY=your_api_key

# 一键启动全部
docker-compose -f docker/docker-compose.yml up -d --build

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f backend

# 热更新前端
docker-compose -f docker/docker-compose.yml build --no-cache frontend
docker-compose -f docker/docker-compose.yml up -d frontend
```

---

## 四、热更新操作速查

### 仅更新前端
```bash
docker-compose build --no-cache frontend
docker-compose up -d frontend
```

### 仅更新后端
```bash
docker-compose build --no-cache backend
docker-compose up -d backend
```

### 全量更新
```bash
docker-compose up --build -d
```

### 清理未使用镜像
```bash
docker image prune -f
docker builder prune -f
```

---

## 五、数据库管理

### 备份
```bash
docker exec mt-postgres pg_dump -U postgres magtrade > backup.sql
```

### 恢复
```bash
cat backup.sql | docker exec -i mt-postgres psql -U postgres magtrade
```

### 进入数据库CLI
```bash
docker exec -it mt-postgres psql -U postgres -d magtrade
```

---

## 六、问题排查记录

### 问题1：后端容器无法连接PostgreSQL

**错误信息：**
```
failed to connect to `user=postgres database=magtrade`:
127.0.0.1:5432 (localhost): dial error: dial tcp 127.0.0.1:5432: connect: connection refused
```

**原因分析：**
后端容器内`localhost`指向容器自身，不是宿主机。Docker网络中需要用服务名`postgres`。

**根本原因：**
Go配置文件读取的是 `config.dev.yaml` 中的 `host: localhost`，但环境变量 `DB_HOST=postgres` 没有生效。

**代码问题（config.go）：**
```go
// 原代码只识别 ${VAR} 格式
func expandEnv(s string) string {
    if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
        // ...
    }
    return s  // 直接返回原值，忽略了环境变量
}
```

**解决方案：**
修改 `internal/config/config.go` 的 `expandEnvVars` 函数，优先读取环境变量：
```go
func expandEnvVars(c *Config) {
    if v := os.Getenv("DB_HOST"); v != "" {
        c.Database.Host = v
    } else {
        c.Database.Host = expandEnv(c.Database.Host)
    }
    // ... 其他配置项同理
}
```

---

### 问题2：Seed数据未导入

**错误信息：**
PostgreSQL日志显示：
```
psql:/docker-entrypoint-initdb.d/02_seed.sql:7: ERROR: relation "users" does not exist
```

**原因分析：**
PostgreSQL的 `docker-entrypoint-initdb.d` 脚本在**容器首次启动时**执行，但此时：
- 后端还没启动
- GORM还没创建表
- 所以seed脚本找不到users表

**解决方案：**
1. 移除docker-compose中的seed_data.sql挂载
2. 在Go代码中添加 `SeedData()` 函数
3. 在 `main.go` 的 `AutoMigrate()` 之后调用

**代码修改（postgres.go）：**
```go
func SeedData() error {
    var count int64
    db.Model(&model.User{}).Count(&count)
    if count > 0 {
        return nil  // 已有数据，跳过
    }
    
    admin := model.User{
        Username:     "admin",
        Email:        "admin@magtrade.com",
        PasswordHash: "$2a$10$bKmO1qYASj5loeB5oT4nBO2NXISkf2E7vaOPQQ5/pBXy0FSQpQO0m",
        Role:         "admin",
        Status:       1,
        EmailVerified: true,
    }
    return db.Create(&admin).Error
}
```

**调用位置（main.go）：**
```go
if err := database.AutoMigrate(); err != nil {
    log.Fatal("failed to auto migrate", zap.Error(err))
}
if err := database.SeedData(); err != nil {
    log.Warn("failed to seed data", zap.Error(err))
}
```

---

### 问题3：密码Hash不匹配

**错误信息：**
```
{"code":401,"message":"用户名或密码错误"}
```

**原因分析：**
硬编码的密码hash是 `password` 的hash，不是 `admin123`。

**解决方案：**
使用 `seed_data.sql` 中正确的hash值：
```go
PasswordHash: "$2a$10$bKmO1qYASj5loeB5oT4nBO2NXISkf2E7vaOPQQ5/pBXy0FSQpQO0m"  // admin123
```

---

### 问题4：PostgreSQL容器不健康

**错误信息：**
```
container mt-postgres is unhealthy
```

**原因分析：**
seed_data.sql在init阶段执行失败（找不到表），导致PostgreSQL初始化中断。

**解决方案：**
移除seed_data.sql的挂载，完全依赖Go代码的SeedData()函数。

---

## 七、八股文简记

### Docker网络
- 容器间通信用**服务名**（如`postgres`），不是`localhost`
- `localhost`在容器内指向容器自身
- Docker Compose自动创建bridge网络，服务名自动解析

### Docker Volume
- **命名卷**：`volume_name:/path`，数据持久化
- **挂载卷**：`./local:/container`，开发调试用
- 删除容器不删数据，需要 `docker-compose down -v`

### docker-entrypoint-initdb.d
- PostgreSQL容器**首次启动**时执行该目录下的脚本
- 按文件名字母序执行（01_xxx.sql → 02_xxx.sql）
- Volume已存在数据则**跳过**初始化
- 不适合依赖应用层表结构的数据（GORM表）

### 环境变量覆盖
- `docker-compose.yml`中用 `${VAR:-default}` 语法
- 应用代码需要显式读取 `os.Getenv("VAR")`
- 不能假设YAML配置自动读取环境变量

### 健康检查
- `healthcheck`定义容器健康探测命令
- `depends_on.condition: service_healthy` 等待依赖健康
- 常见问题：依赖服务已启动但未**就绪**

### Docker构建缓存
- 层缓存：只有变化的层重新构建
- `--no-cache`：强制全量重建
- `.dockerignore`：排除不需要的文件（node_modules等）

---

## 八、容器命名规范

| 服务 | 容器名 | 说明 |
|------|--------|------|
| PostgreSQL | mt-postgres | 数据库 |
| Redis | mt-redis | 缓存 |
| Kafka | mt-kafka | 消息队列 |
| Zookeeper | mt-zookeeper | Kafka协调器 |
| Backend | mt-backend | Go后端 |
| Frontend | mt-frontend | Vue前端 |

---

## 九、常用命令速查

```bash
# 查看容器状态
docker-compose -f docker/docker-compose.dev.yml ps

# 查看实时日志
docker logs -f mt-backend

# 进入容器Shell
docker exec -it mt-backend sh

# 查看网络
docker network ls
docker network inspect docker_default

# 查看卷
docker volume ls
docker volume inspect docker_mt_postgres_data
```
