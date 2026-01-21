MagTrade 项目上下文

部署地址: https://gcptw.yukiyuki.cfd
项目路径: /home/mgf/MagTrade

SSL证书:
  证书: /root/ygkkkca/cert.crt
  私钥: /root/ygkkkca/private.key
  域名: gcptw.yukiyuki.cfd

Nginx配置: /etc/nginx/sites-enabled/magtrade
Docker配置: /home/mgf/MagTrade/docker/docker-compose.yml
环境变量: /home/mgf/MagTrade/docker/.env (敏感,不提交git)

端口映射:
  80/443 -> Nginx
  9080 -> Frontend (mt-frontend)
  9081 -> Backend (mt-backend)
  50051 -> gRPC (内部)

Docker容器:
  mt-backend    Go API服务
  mt-frontend   Vue静态资源
  mt-postgres   PostgreSQL 16
  mt-redis      Redis 7
  mt-kafka      Kafka 7.5
  mt-zookeeper  Zookeeper

项目结构:
  cmd/server/main.go              应用入口
  internal/config/                配置管理
  internal/handler/               HTTP处理器
  internal/service/               业务逻辑
  internal/service/ai/            AI Agent模块
  internal/repository/            数据访问层
  internal/cache/                 Redis操作+分布式锁
  internal/mq/                    Kafka生产者/消费者
  internal/grpc/                  gRPC库存服务
  internal/middleware/            中间件(认证/限流/日志)
  internal/worker/                后台任务(订单处理/定时任务)
  frontend/                       Vue 3前端

API路由:
  /health                         健康检查(含DB/Redis状态)
  /metrics                        监控指标
  /api/v1/auth/*                  认证(登录/注册/刷新)
  /api/v1/products                商品列表
  /api/v1/flash-sales             秒杀活动
  /api/v1/flash-sales/:id/rush    秒杀抢购(核心)
  /api/v1/orders                  订单管理
  /api/v1/ai/*                    AI客服/策略
  /api/v1/admin/*                 管理后台
  /ws/notifications               WebSocket通知

核心技术:
  后端: Go 1.24 + Gin + GORM + JWT双Token
  前端: Vue 3 + TypeScript + Pinia + Tailwind
  数据库: PostgreSQL 16
  缓存: Redis 7 + Lua脚本原子操作
  消息队列: Kafka (异步下单)
  AI: SiliconFlow + DeepSeek-V3.2

常用命令:
  make build    零停机更新后端
  make logs     查看日志
  make backup   备份数据库
  make db       进入数据库CLI
