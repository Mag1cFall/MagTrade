# MagTrade 前端设计方案 - "Neon Velocity"

## 1. 设计愿景
打造一个**沉浸式、高性能、未来感**的秒杀交易终端。摒弃传统电商的拥挤布局，采用类似金融交易终端或科幻控制台的 **Bento Grid（便当盒）** 布局，强调信息的层级和实时性。

- **关键词**：极速 (Velocity)、深邃 (Deep Space)、智能 (Intelligence)。
- **视觉风格**：Cyber-Minimalism（赛博极简主义）。
  - **背景**：深空灰 (`#0f172a`) 搭配极细的磨砂玻璃 (`Backdrop Blur`)。
  - **强调色**：电光紫 (`#8b5cf6`) 用于 AI，信号绿 (`#10b981`) 用于成功/库存，警示红 (`#ef4444`) 用于倒计时/售罄。
  - **字体**：Inter / Roboto Mono (数字/代码相关)。

## 2. 技术栈架构 (Vue 3 + TS 生态)

| 模块 | 技术选型 | 选择理由 |
| :--- | :--- | :--- |
| **核心框架** | **Vue 3 (Script Setup)** | 性能优异，Composition API 逻辑复用能力强。 |
| **语言** | **TypeScript** | 强类型约束，与后端 Go 的结构体对齐。 |
| **构建工具** | **Vite** | 极速冷启动，秒级热更新 (HMR)。 |
| **状态管理** | **Pinia** | 轻量级，Type safe，支持持久化插件。 |
| **UI 引擎** | **Tailwind CSS** | 原子化 CSS，实现极致的定制设计，避免“组件库味”。 |
| **组件库** | **Naive UI** | 尤雨溪推荐，全 TS 编写，主题定制能力强，自带暗黑模式。 |
| **动画引擎** | **Motion One / GSAP** | 实现高性能的复杂过渡动画 (60fps)。 |
| **实时通信** | **Native WebSocket** | 无额外依赖，轻量对接后端的 Gorilla WebSocket。 |
| **倒计时** | **Web Worker** | 将倒计时逻辑放入 Worker 线程，防止主线程阻塞导致计时偏差。 |

## 3. 核心功能与酷炫实现

### 3.1 首页 / 驾驶舱 (Dashboard)
- **布局**：Bento Grid 布局，左侧为 AI 助手状态，中间为热门秒杀（大卡片），右侧为实时战报（滚动列表）。
- **动效**：
  - 卡片悬停时产生光晕流动效果 (Glow Effect)。
  - 实时战报采用无缝滚动 (Infinite Scroll)。

### 3.2 秒杀详情与倒计时
- **高性能倒计时组件**：
  - 使用 `requestAnimationFrame` 替代 `setInterval` 确保毫秒级渲染流畅。
  - 数字变化采用类似“老虎机”的滚动动画。
- **抢购按钮**：
  - 状态流转：`等待中` (灰色) -> `开抢` (流光呼吸) -> `排队中` (Loading 粒子) -> `成功` (全屏庆祝)。
  - 点击反馈：点击瞬间产生波纹扩散效果 (Ripple)。

### 3.3 AI 智能客服 (DeepSeek Driver)
- **形态**：页面右下角悬浮的“全息球体”或“光环”。
- **交互**：
  - 支持流式输出 (Streaming Typing Effect)，模拟真实打字感。
  - 上下文感知：自动读取当前浏览的商品信息，无需用户重复输入。
- **Markdown 渲染**：支持渲染富文本建议（如秒杀策略表格）。

### 3.4 实时库存系统
- **动态库存条**：使用 SVG 路径动画模拟液态流动，颜色随库存比例变化（绿 -> 黄 -> 红）。
- **WS 推送**：秒级更新库存，并在页面角落弹出类似“游戏成就”的 Toast 通知 ("用户 xxx 刚刚抢到了!").

## 4. 目录结构规划

```
frontend/
├── src/
│   ├── assets/             # 静态资源 (Lottie JSON, Images)
│   ├── components/         # 通用组件
│   │   ├── common/         # 按钮, 卡片, 弹窗
│   │   ├── business/       # 倒计时, 库存条, AI对话框
│   ├── composables/        # 组合式函数 (useWebSocket, useCountDown)
│   ├── stores/             # Pinia 状态 (user, flashSale, order)
│   ├── views/              # 页面视图 (Home, Product, Order)
│   ├── router/             # 路由配置
│   ├── api/                # Axios 封装与接口定义
│   ├── utils/              # 工具函数
│   ├── types/              # TS 类型定义
│   ├── workers/            # Web Workers (计时器)
│   ├── App.vue
│   └── main.ts
├── index.html
├── tailwind.config.js
├── tsconfig.json
└── vite.config.ts
```

## 5. 开发路线图

1.  **初始化**：搭建 Vite + Vue3 + TS + Tailwind 环境。
2.  **基础建设**：封装 Axios (拦截器处理 JWT)，封装 WebSocket (心跳保活)。
3.  **UI 开发**：构建核心布局组件，配置暗黑主题。
4.  **业务对接**：
    - 认证模块 (登录/注册/JWT存储)。
    - 秒杀列表与详情。
    - 订单流程 (下单/支付/WebSocket通知)。
5.  **AI 集成**：开发 Chat 组件，对接流式接口。
6.  **优化与打磨**：添加动效，性能调优 (Lighthouse)。