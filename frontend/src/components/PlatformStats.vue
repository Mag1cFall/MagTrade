<template>
  <div class="py-20 bg-black relative overflow-hidden border-b border-white/5 platform-stats">
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-accent/5 via-transparent to-transparent"></div>
    
    <div class="max-w-6xl mx-auto px-4 relative z-10">
      <div class="text-center mb-16">
        <h2 class="text-3xl md:text-5xl font-black text-white mb-4 uppercase tracking-tighter italic">Platform Statistics</h2>
        <p class="text-secondary font-mono tracking-widest text-xs uppercase opacity-60">Real-time throughput metrics // Network Load</p>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 md:gap-8">
        <div 
          v-for="stat in stats" 
          :key="stat.label"
          class="text-center group p-8 border border-white/5 rounded-none bg-white/5 backdrop-blur-sm hover:border-accent transition-all duration-500 hover:bg-accent/5"
        >
          <div class="flex flex-col items-center justify-center space-y-2">
            <div class="text-3xl sm:text-4xl lg:text-5xl font-black text-white font-mono group-hover:text-accent transition-colors leading-none tracking-tighter">
              <span>{{ stat.displayValue }}</span><span class="text-accent text-xl ml-1">{{ stat.suffix }}</span>
            </div>
            <div class="text-secondary text-[10px] sm:text-xs uppercase tracking-[0.3em] font-bold h-4 flex items-center justify-center opacity-50 group-hover:opacity-100 transition-opacity">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const stats = ref([
  { value: 2500000, displayValue: '0', suffix: '+', label: 'Active Users' },
  { value: 15000000, displayValue: '0', suffix: '+', label: 'Transactions' },
  { value: 98.7, displayValue: '0', suffix: '%', label: 'Success Rate' },
  { value: 850, displayValue: '0', suffix: 'M', label: 'Volume (¥)' }
])

const animateNumbers = () => {
  stats.value.forEach((stat, index) => {
    const duration = 2000
    const startTime = performance.now()
    const end = stat.value
    
    const animate = (currentTime: number) => {
      const elapsed = currentTime - startTime
      const progress = Math.min(elapsed / duration, 1)
      const eased = 1 - Math.pow(1 - progress, 3)
      const current = end * eased
      
      if (stat.suffix === '%') {
        stats.value[index].displayValue = current.toFixed(1)
      } else if (stat.suffix === 'M') {
        stats.value[index].displayValue = current.toFixed(0)
      } else {
        stats.value[index].displayValue = Math.floor(current).toLocaleString()
      }
      
      if (progress < 1) {
        requestAnimationFrame(animate)
      }
    }
    
    setTimeout(() => requestAnimationFrame(animate), index * 200)
  })
}

onMounted(() => {
  const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) {
      animateNumbers()
      observer.disconnect()
    }
  }, { threshold: 0.3 })
  
  const el = document.querySelector('.platform-stats')
  if (el) observer.observe(el)
})
</script>
