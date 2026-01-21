# MagTrade 项目优化任务清单

## 项目现状

代码规模: Go 9,143行 + 前端58文件
测试覆盖: 6个测试文件 (仅 model/utils/validator)
部署状态: 6容器运行中 (backend/frontend/postgres/redis/kafka/zookeeper)

## 已实现功能

[x] 秒杀核心 - Redis Lua原子操作 + 分布式锁 + Kafka异步下单
[x] gRPC库存服务 - Server(50051) + Client完整实现
[x] AI Agent - 智能客服/策略推荐/异常检测/LLM客户端
[x] Kafka消息队列 - 3个Topic (flash-sale-orders/order-status-change/ai-analysis-tasks)
[x] WebSocket推送 - 秒杀结果实时通知
[x] JWT双Token - Access + Refresh Token机制
[x] 邮箱验证 - 验证码发送 + 验证逻辑
[x] 定时任务 - 活动状态更新 + 过期订单取消(15分钟)
[x] 审计日志 - 中间件记录

## 待优化任务

### 高优先级 (简历亮点)

[ ] Cloudflare 配置
    - 已添加 magtrade.yukiyuki.cfd A记录
    - 需申请新SSL证书 (acme.sh 添加域名)
    - 更新 Nginx 配置 (已完成)
    - SSL模式设为 Full (Strict)
    - 开启 Rate Limiting / WAF

[ ] 补充核心业务测试
    - internal/service/ 业务逻辑测试
    - internal/handler/ API集成测试
    - internal/cache/ Redis操作测试
    目标覆盖率: 60%+

[ ] 添加 Prometheus 监控
    - /metrics 端点输出 Prometheus 格式
    - 添加 prometheus + grafana 容器
    - 创建秒杀面板 (QPS/库存/订单)

[ ] gRPC服务拆分示例
    - 库存服务独立 docker-compose service
    - 通过环境变量切换直连/gRPC模式

### 中优先级 (进阶展示)

[ ] Kafka 多分区配置
    - flash-sale-orders: 6 partitions
    - 消费者组多实例

[ ] PostgreSQL 读写分离
    - 添加只读副本
    - GORM 配置读写分离

[ ] OpenTelemetry 链路追踪
    - 集成 otel-go
    - 部署 Jaeger

[ ] 集成 Stripe 支付 (测试模式)
    - 创建 Checkout Session
    - Webhook 处理支付成功

### 低优先级 (锦上添花)

[ ] Elasticsearch 商品搜索
[ ] 购物车功能
[ ] 优惠券/促销系统
[ ] 分库分表 (sharding)

## 变成真正电商的差距

缺失功能:
  - 支付集成 (支付宝/微信/Stripe)
  - 物流系统 (发货/追踪)
  - 完整库存管理 (SKU/仓库)
  - 购物车
  - 优惠券/促销
  - 评价系统
  - 全文搜索

安全合规:
  - CSRF保护 (部分)
  - 敏感数据加密 (未实现)
  - GDPR合规 (未实现)

## Cloudflare 配置说明

当前状态:
  域名: magtrade.yukiyuki.cfd
  A记录: 35.194.138.153
  代理: 已开启 (橙色云)

需要操作:
  1. 申请新SSL证书:
     acme.sh --issue -d gcptw.yukiyuki.cfd -d magtrade.yukiyuki.cfd --webroot /var/www/html

  2. 或使用 Cloudflare Origin Certificate:
     Cloudflare Dashboard -> SSL/TLS -> Origin Server -> Create Certificate

  3. 更新服务器Nginx:
     sudo cp /home/mgf/MagTrade/docker/nginx-magtrade.conf /etc/nginx/sites-available/magtrade
     sudo nginx -t && sudo systemctl reload nginx

  4. Cloudflare SSL设置:
     SSL/TLS -> 加密模式 -> Full (Strict)

副作用:
  - 服务器日志记录的是CF IP (已通过 set_real_ip_from 修复)
  - WebSocket 需要 CF Pro 或开启 WebSocket 支持 (免费版默认开启)
  - 上传大文件可能需要调整 CF 限制 (免费版100MB)

## 数据库优化方向

当前: 单PostgreSQL实例

优化路径:
  1. 连接池优化 (max_connections)
  2. 索引分析 (EXPLAIN ANALYZE)
  3. 读写分离 (主从复制)
  4. 分库分表 (按用户ID, 大规模时)
  5. TimescaleDB (时序数据)

## 常用命令

make build      零停机更新后端
make logs       查看日志
make backup     备份数据库
make db         进入数据库CLI
make health     检查服务状态

## 运维信息

项目路径: /home/mgf/MagTrade
Nginx配置: /etc/nginx/sites-enabled/magtrade
Docker配置: /home/mgf/MagTrade/docker/docker-compose.yml
环境变量: /home/mgf/MagTrade/docker/.env (敏感,不提交)

访问地址:
  https://gcptw.yukiyuki.cfd (原)
  https://magtrade.yukiyuki.cfd (新,通过CF)
