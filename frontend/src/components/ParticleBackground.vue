<template>
  <div ref="container" class="canvas-container">
    <canvas ref="canvas"></canvas>
    <div class="overlay">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, onUnmounted } from 'vue'

const canvas = ref<HTMLCanvasElement | null>(null)
const container = ref<HTMLElement | null>(null)
let animationFrame: number

const initCanvas = () => {
  const ctx = canvas.value?.getContext('2d')
  if (!ctx || !canvas.value || !container.value) return

  const resize = () => {
    if (container.value && canvas.value) {
      canvas.value.width = container.value.clientWidth
      canvas.value.height = container.value.clientHeight
    }
  }
  resize()
  window.addEventListener('resize', resize)

  // 增加粒子数量到 80，提高连接距离
  const particleCount = 80
  const particles: { x: number; y: number; vx: number; vy: number; size: number }[] = []

  for (let i = 0; i < particleCount; i++) {
    particles.push({
      x: Math.random() * canvas.value.width,
      y: Math.random() * canvas.value.height,
      vx: (Math.random() - 0.5) * 0.3, // 速度稍慢，更优雅
      vy: (Math.random() - 0.5) * 0.3,
      size: Math.random() * 2 + 0.5, // 随机大小
    })
  }

  const animate = () => {
    if (!ctx || !canvas.value) return
    ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)

    // 粒子颜色：Shyft 风格的金色/橙色
    const particleColor = 'rgba(251, 185, 1, 0.4)'
    const lineColor = 'rgba(251, 185, 1, 0.15)'

    particles.forEach((p, i) => {
      p.x += p.vx
      p.y += p.vy

      // 边界反弹
      if (p.x < 0 || p.x > canvas.value!.width) p.vx *= -1
      if (p.y < 0 || p.y > canvas.value!.height) p.vy *= -1

      // 绘制粒子
      ctx.beginPath()
      ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2)
      ctx.fillStyle = particleColor
      ctx.fill()

      // 绘制连接线
      for (let j = i + 1; j < particles.length; j++) {
        const p2 = particles[j]
        const dist = Math.hypot(p.x - p2.x, p.y - p2.y)
        // 连接距离增加到 150
        if (dist < 150) {
          ctx.beginPath()
          ctx.moveTo(p.x, p.y)
          ctx.lineTo(p2.x, p2.y)
          ctx.strokeStyle = lineColor
          ctx.lineWidth = 1 - dist / 150 // 距离越远越细
          ctx.stroke()
        }
      }
    })
    animationFrame = requestAnimationFrame(animate)
  }
  animate()
}

onMounted(() => initCanvas())
onUnmounted(() => {
  cancelAnimationFrame(animationFrame)
})
</script>

<style scoped>
.canvas-container {
  position: absolute; /* 改为 absolute 以覆盖背景 */
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 0; /* 在内容之下 */
  pointer-events: none; /* 不干扰交互 */
}

canvas {
  display: block;
}

.overlay {
  position: relative;
  z-index: 1;
  width: 100%;
  height: 100%;
}
</style>
