<template>
  <div 
    class="flow-border-card" 
    :style="{
      '--border-color': color || '#e33535',
      '--border-width': borderWidth || '1px',
      '--border-radius': borderRadius || '0px'
    }"
  >
    <div class="content">
      <slot></slot>
    </div>
    
    <div class="border-gradient" v-if="active"></div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  color?: string;
  active?: boolean;
  borderWidth?: string;
  borderRadius?: string;
}>();
</script>

<style scoped>
.flow-border-card {
  position: relative;
  background: var(--bg-color, #0a0a0a);
  border-radius: var(--border-radius);
  padding: var(--border-width); /* Space for the border */
  overflow: hidden;
}

.content {
  position: relative;
  z-index: 2;
  height: 100%;
  background: inherit;
  border-radius: calc(var(--border-radius) - var(--border-width));
}

.border-gradient {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 200%;
  aspect-ratio: 1;
  background: conic-gradient(
    transparent 0deg, 
    transparent 80deg, 
    var(--border-color) 120deg, 
    transparent 180deg,
    transparent 360deg
  );
  transform: translate(-50%, -50%);
  animation: rotate 4s linear infinite;
  z-index: 1;
}

@keyframes rotate {
  from { transform: translate(-50%, -50%) rotate(0deg); }
  to { transform: translate(-50%, -50%) rotate(360deg); }
}
</style>