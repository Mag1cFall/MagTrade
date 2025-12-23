<template>
  <div class="filter-toolbar">
    <div class="filter-group">
      <div 
        v-for="filter in filters" 
        :key="filter.id"
        class="filter-btn"
        :class="{ active: modelValue === filter.id }"
        @click="$emit('update:modelValue', filter.id)"
      >
        <span>{{ filter.label }}</span>
        <!-- 滑块背景 -->
        <div v-if="modelValue === filter.id" class="active-bg" layoutId="highlight"></div>
      </div>
    </div>

    <!-- 右侧小工具 (装饰性) -->
    <div class="tools-placeholder">
      <div class="tool-icon">Filter</div>
      <div class="tool-icon">Sort</div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Filter {
  id: string;
  label: string;
}

defineProps<{
  modelValue: string;
}>();

defineEmits(['update:modelValue']);

const filters: Filter[] = [
  { id: 'all', label: '全部' },
  { id: 'digital', label: '数码' },
  { id: 'electronics', label: '电子' },
  { id: 'home', label: '生活' },
  { id: 'luxury', label: '奢品' }
];
</script>

<style scoped>
.filter-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: #111;
  padding: 6px;
  border-radius: 8px;
  border: 1px solid #222;
}

.filter-group {
  display: flex;
  gap: 4px;
  background: #0a0a0a;
  padding: 4px;
  border-radius: 6px;
}

.filter-btn {
  position: relative;
  padding: 8px 16px;
  color: #666;
  font-size: 12px;
  font-weight: bold;
  cursor: pointer;
  transition: color 0.3s;
  z-index: 1;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.filter-btn.active {
  color: #fff;
}

.filter-btn:hover:not(.active) {
  color: #999;
}

.active-bg {
  position: absolute;
  inset: 0;
  background: #222;
  border: 1px solid #333;
  border-radius: 4px;
  z-index: -1;
  /* 实际项目中如果引入了 Framer Motion for Vue 可以做 Layout 动画，这里 CSS 简单模拟 */
}

.tools-placeholder {
  display: flex;
  gap: 12px;
  padding-right: 12px;
}

.tool-icon {
  font-size: 10px;
  color: #444;
  border: 1px solid #222;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: not-allowed;
}
</style>
