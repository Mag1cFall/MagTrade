<template>
  <div class="live-feed-container">
    <div class="feed-header">
      <span class="live-dot"></span>
      LIVE FEED
    </div>
    <div class="feed-list" @mouseenter="paused = true" @mouseleave="paused = false">
      <TransitionGroup name="list" tag="div">
        <div
          v-for="item in visibleItems"
          :key="item.id"
          class="feed-item"
          @click="showDetail(item)"
        >
          <div class="item-avatar">{{ item.user.charAt(0) }}</div>
          <div class="item-content">
            <div class="item-row">
              <span class="username">{{ item.user }}</span>
              <span class="action-tag">{{ item.action }}</span>
            </div>
            <div class="item-product">{{ item.product }}</div>
          </div>
          <div class="item-time">{{ item.time }}</div>
        </div>
      </TransitionGroup>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface FeedItem {
  id: string
  user: string
  action: string
  product: string
  time: string
}

const rawData = [
  { user: 'Alex_K', action: 'SNIPED', product: 'iPhone 15 Pro Max', time: 'Just now' },
  { user: 'Sarah_99', action: 'BOUGHT', product: 'Sony WH-1000XM5', time: '2s ago' },
  { user: 'CryptoKing', action: 'CLEARED', product: 'RTX 4090 OC', time: '5s ago' },
  { user: 'User_8821', action: 'ORDERED', product: 'MacBook Pro M3', time: '8s ago' },
  { user: 'Trader_X', action: 'SNIPED', product: 'Switch OLED', time: '12s ago' },
  { user: 'Ghost_Protocol', action: 'LOCKED', product: 'PS5 Slim', time: '15s ago' },
]

const visibleItems = ref<FeedItem[]>([])
const paused = ref(false)
let interval: any

const generateId = () => Math.random().toString(36).substr(2, 9)

const addNextItem = () => {
  if (paused.value) return

  const randomItem = rawData[Math.floor(Math.random() * rawData.length)]
  const newItem = { ...randomItem, id: generateId(), time: 'Just now' }

  visibleItems.value.forEach((item) => {
    if (item.time === 'Just now') item.time = '2s ago'
    else if (item.time === '2s ago') item.time = '5s ago'
    else item.time = '10s+'
  })

  visibleItems.value.unshift(newItem)
  if (visibleItems.value.length > 4) {
    visibleItems.value.pop()
  }
}

const showDetail = (item: FeedItem) => {
  console.log('Viewing details for', item)
}

onMounted(() => {
  // 初始化
  visibleItems.value = rawData.slice(0, 3).map((i) => ({ ...i, id: generateId() }))
  interval = setInterval(addNextItem, 2500)
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>

<style scoped>
.live-feed-container {
  position: absolute;
  left: 40px;
  bottom: 120px; /* 下移避免遮挡标题 */
  width: 300px;
  background: rgba(10, 10, 10, 0.8);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-left: 2px solid #e33535;
  border-radius: 4px;
  padding: 16px;
  z-index: 20;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
  font-family: 'JetBrains Mono', monospace;
}

.feed-header {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #888;
  font-size: 10px;
  letter-spacing: 2px;
  margin-bottom: 12px;
  font-weight: bold;
}

.live-dot {
  width: 6px;
  height: 6px;
  background: #e33535;
  border-radius: 50%;
  box-shadow: 0 0 8px #e33535;
  animation: pulse 1s infinite;
}

.feed-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow: hidden;
}

.feed-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.feed-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.1);
  transform: translateX(5px);
}

.item-avatar {
  width: 24px;
  height: 24px;
  background: #333;
  color: #fff;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: bold;
}

.item-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.item-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.username {
  color: #fff;
  font-size: 12px;
  font-weight: bold;
}

.action-tag {
  color: #e33535;
  font-size: 10px;
  font-weight: bold;
}

.item-product {
  color: #888;
  font-size: 10px;
}

.item-time {
  color: #555;
  font-size: 9px;
  min-width: 40px;
  text-align: right;
}

/* List Transitions */
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.list-leave-active {
  position: absolute;
  width: 100%;
}

@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
  100% {
    opacity: 1;
  }
}

@media (max-width: 768px) {
  .live-feed-container {
    display: none;
  }
}
</style>
