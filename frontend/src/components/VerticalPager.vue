<template>
  <aside class="sidebar">
    <div class="sidebar__line"></div>
    <ul class="sidebar__pagers">
      <li 
        v-for="(item, index) in sections" 
        :key="index"
        :class="['sidebar__num', { 'sidebar__num--active': activeIndex === index }]"
        @click="$emit('change', index)"
      >
        {{ formatNumber(index + 1) }}
      </li>
    </ul>
  </aside>
</template>

<script setup lang="ts">
// defineProps and defineEmits are compiler macros and do not need to be imported

const props = defineProps<{
  sections: string[]
  activeIndex: number
}>()

const emit = defineEmits(['change'])

const formatNumber = (num: number) => num < 10 ? `0${num}` : `${num}`
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 50%;
  right: 2rem;
  transform: translateY(-50%);
  z-index: 50;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.sidebar__line {
  width: 1px;
  height: 100px;
  background: linear-gradient(to bottom, transparent, #333, transparent);
  margin-bottom: 1rem;
}

.sidebar__pagers {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.sidebar__num {
  font-family: 'Impact', sans-serif;
  font-size: 1.2rem;
  color: #333;
  cursor: pointer;
  writing-mode: vertical-rl;
  user-select: none;
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  position: relative;
}

.sidebar__num:hover {
  color: #666;
  transform: translateX(-5px);
}

.sidebar__num--active {
  color: #d8fa00; /* ZZZ Style Green */
  font-size: 1.5rem;
  text-shadow: 0 0 10px rgba(216, 250, 0, 0.5);
}

.sidebar__num--active::after {
  content: '';
  position: absolute;
  right: -10px;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 4px;
  background-color: #d8fa00;
  box-shadow: 0 0 5px #d8fa00;
}
</style>