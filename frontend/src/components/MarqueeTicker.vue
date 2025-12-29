<template>
  <div class="marquee-wrapper" :style="wrapperStyle">
    <div class="marquee-track" :style="trackStyle">
      <div ref="contentRef" class="marquee-content">
        <slot></slot>
      </div>
      <div class="marquee-content">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  duration?: number
  direction?: 'normal' | 'reverse'
  pauseOnHover?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  duration: 20,
  direction: 'normal',
  pauseOnHover: true,
})

const wrapperStyle = computed(() => ({
  '--pause-state': props.pauseOnHover ? 'paused' : 'running',
}))

const trackStyle = computed(() => ({
  '--animation-duration': `${props.duration}s`,
  '--animation-direction': props.direction,
}))
</script>

<style scoped>
.marquee-wrapper {
  width: 100%;
  overflow: hidden;
  position: relative;
  background: transparent;
  mask-image: linear-gradient(to right, transparent, black 10%, black 90%, transparent);
  -webkit-mask-image: linear-gradient(to right, transparent, black 10%, black 90%, transparent);
}

.marquee-track {
  display: flex;
  width: max-content;
  animation: scroll var(--animation-duration) linear infinite;
  animation-direction: var(--animation-direction);
}

.marquee-wrapper:hover .marquee-track {
  animation-play-state: var(--pause-state);
}

.marquee-content {
  display: flex;
  flex-shrink: 0;
  gap: 2rem;
  padding-right: 2rem;
}

@keyframes scroll {
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(-50%);
  }
}
</style>
