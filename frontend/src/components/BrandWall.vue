<template>
  <div class="brand-wall-container">
    <div class="brand-track" :style="trackStyle">
      <!-- 第一组 -->
      <div v-for="(brand, i) in brands" :key="`a-${i}`" class="brand-item group">
        <component :is="brand.icon" class="brand-icon" />
        <span class="brand-name">{{ brand.name }}</span>
        
        <!-- Hover Popover -->
        <div class="brand-popover">
          <div class="pop-title">{{ brand.name }}</div>
          <div class="pop-desc">Official Partner</div>
        </div>
      </div>
      
      <!-- 克隆组实现无缝滚动 -->
      <div v-for="(brand, i) in brands" :key="`b-${i}`" class="brand-item group">
        <component :is="brand.icon" class="brand-icon" />
        <span class="brand-name">{{ brand.name }}</span>
        
         <!-- Hover Popover -->
         <div class="brand-popover">
          <div class="pop-title">{{ brand.name }}</div>
          <div class="pop-desc">Official Partner</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h } from 'vue';
import { Apple, Smartphone, Laptop, Headphones, Watch, Gamepad2, Camera, Speaker } from 'lucide-vue-next';

const brands = [
  { name: 'Apple', icon: Apple },
  { name: 'Samsung', icon: Smartphone },
  { name: 'Dell', icon: Laptop },
  { name: 'Sony', icon: Headphones },
  { name: 'Rolex', icon: Watch },
  { name: 'Nintendo', icon: Gamepad2 },
  { name: 'Canon', icon: Camera },
  { name: 'JBL', icon: Speaker },
  { name: 'Xiaomi', icon: Smartphone }, // 复用图标
  { name: 'Asus', icon: Laptop }
];

const trackStyle = computed(() => ({
  animationDuration: '30s'
}));
</script>

<style scoped>
.brand-wall-container {
  width: 100%;
  overflow: hidden;
  padding: 24px 0;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  position: relative;
}

.brand-wall-container::before,
.brand-wall-container::after {
  content: "";
  position: absolute;
  top: 0;
  width: 150px;
  height: 100%;
  z-index: 2;
  pointer-events: none;
}

.brand-wall-container::before {
  left: 0;
  background: linear-gradient(to right, #0a0a0a, transparent);
}

.brand-wall-container::after {
  right: 0;
  background: linear-gradient(to left, #0a0a0a, transparent);
}

.brand-track {
  display: flex;
  width: max-content;
  gap: 60px;
  animation: scroll linear infinite;
}

.brand-wall-container:hover .brand-track {
  animation-play-state: paused;
}

.brand-item {
  display: flex;
  align-items: center;
  gap: 12px;
  opacity: 0.4;
  transition: all 0.3s;
  cursor: pointer;
  position: relative;
}

.brand-item:hover {
  opacity: 1;
  transform: scale(1.05);
  filter: drop-shadow(0 0 8px rgba(255, 255, 255, 0.5));
}

.brand-icon {
  width: 24px;
  height: 24px;
}

.brand-name {
  font-family: 'Kanit', sans-serif;
  font-weight: 500;
  font-size: 18px;
  color: #fff;
}

/* Popover */
.brand-popover {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(-10px);
  background: #1a1a1a;
  border: 1px solid #333;
  padding: 10px 16px;
  border-radius: 8px;
  min-width: 120px;
  text-align: center;
  opacity: 0;
  visibility: hidden;
  transition: all 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  z-index: 10;
  pointer-events: none;
}

.brand-item:hover .brand-popover {
  opacity: 1;
  visibility: visible;
  transform: translateX(-50%) translateY(-15px);
}

.pop-title {
  color: #fff;
  font-weight: bold;
  font-size: 14px;
}

.pop-desc {
  color: #666;
  font-size: 10px;
  margin-top: 2px;
}

@keyframes scroll {
  0% { transform: translateX(0); }
  100% { transform: translateX(-50%); }
}
</style>
