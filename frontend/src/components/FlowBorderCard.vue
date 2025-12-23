<template>
  <div class="flow-border-wrapper" :style="wrapperStyle">
    <div class="flow-border-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  color?: string;
  borderWidth?: string;
  borderRadius?: string;
  duration?: number;
  active?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  color: '#e33535',
  borderWidth: '1px',
  borderRadius: '0px',
  duration: 4,
  active: true
});

const wrapperStyle = computed(() => ({
  '--border-color': props.color,
  '--border-width': props.borderWidth,
  '--border-radius': props.borderRadius,
  '--animation-duration': `${props.duration}s`,
  '--opacity': props.active ? 1 : 0
}));
</script>

<style scoped>
.flow-border-wrapper {
  position: relative;
  padding: var(--border-width);
  border-radius: var(--border-radius);
  overflow: hidden;
  background: rgba(255, 255, 255, 0.02);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.flow-border-wrapper::before {
  content: "";
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: conic-gradient(
    from 0deg,
    transparent 0%,
    transparent 40%,
    var(--border-color) 50%,
    transparent 60%,
    transparent 100%
  );
  animation: rotate-border var(--animation-duration) linear infinite;
  opacity: var(--opacity);
  transition: opacity 0.3s ease;
}

.flow-border-content {
  position: relative;
  z-index: 1;
  background: #050505;
  border-radius: calc(var(--border-radius) - var(--border-width));
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

@keyframes rotate-border {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>