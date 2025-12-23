<template>
  <div class="stats-counter-container">
    <div v-for="(stat, index) in stats" :key="index" class="stat-box">
      <div class="stat-value">
        <span class="value-text">{{ stat.displayValue }}</span>
        <span class="unit">{{ stat.unit }}</span>
      </div>
      <div class="stat-label">{{ stat.label }}</div>
      <div class="stat-bar">
        <div class="stat-fill" :style="{ width: stat.percent + '%' }"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

interface StatItem {
  label: string;
  value: number;
  displayValue: string;
  unit: string;
  percent: number;
  increment: number;
}

const stats = ref<StatItem[]>([
  { label: 'TOTAL VOLUME', value: 482000, displayValue: '482,000', unit: '$', percent: 75, increment: 120 },
  { label: 'ACTIVE USERS', value: 2450, displayValue: '2,450', unit: '', percent: 60, increment: 5 },
  { label: 'SYS LOAD', value: 34, displayValue: '34.2', unit: '%', percent: 34, increment: 0.1 },
]);

// 简单的数字加法模拟真实感
let timer: any;

const updateStats = () => {
  stats.value.forEach(stat => {
    // 随机增加
    if (Math.random() > 0.5) {
      stat.value += stat.increment * Math.random();
      
      // 格式化显示
      if (stat.label === 'SYS LOAD') {
        // 负载波动
        stat.value = 30 + Math.random() * 15;
        stat.displayValue = stat.value.toFixed(1);
        stat.percent = stat.value;
      } else {
        stat.displayValue = Math.floor(stat.value).toLocaleString();
        // 简单计算百分比用于进度条动画（假定一个最大值）
        if (stat.label.includes('VOLUME')) stat.percent = (stat.value % 1000000) / 10000;
        if (stat.label.includes('USERS')) stat.percent = (stat.value % 5000) / 50;
      }
    }
  });
};

onMounted(() => {
  timer = setInterval(updateStats, 1000);
});

onUnmounted(() => {
  clearInterval(timer);
});
</script>

<style scoped>
.stats-counter-container {
  position: absolute;
  right: 40px;
  bottom: 120px; /* 在跑马灯上方 */
  display: flex;
  gap: 30px;
  z-index: 20;
}

.stat-box {
  display: flex;
  flex-direction: column;
  align-items: flex-end; /* 右对齐 */
  min-width: 120px;
}

.stat-value {
  font-family: 'JetBrains Mono', monospace;
  font-size: 28px;
  font-weight: 800;
  color: #fff;
  display: flex;
  align-items: baseline;
  gap: 4px;
  text-shadow: 0 0 20px rgba(255, 255, 255, 0.3);
}

.unit {
  font-size: 14px;
  color: #e33535;
}

.stat-label {
  font-size: 10px;
  color: #666;
  letter-spacing: 2px;
  margin-top: 4px;
  text-transform: uppercase;
}

.stat-bar {
  width: 100%;
  height: 2px;
  background: rgba(255, 255, 255, 0.1);
  margin-top: 8px;
  position: relative;
  overflow: hidden;
}

.stat-fill {
  height: 100%;
  background: #e33535;
  transition: width 0.5s ease;
  box-shadow: 0 0 8px #e33535;
}

/* 移动端隐藏 */
@media (max-width: 768px) {
  .stats-counter-container {
    display: none;
  }
}
</style>
