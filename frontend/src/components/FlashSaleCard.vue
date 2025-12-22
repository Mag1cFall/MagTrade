<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import type { FlashSale } from '@/types'
import { ArrowRight } from 'lucide-vue-next'
import CountDown from './CountDown.vue'
import SmartImage from './SmartImage.vue'

const props = defineProps<{
  sale: FlashSale
}>()

const router = useRouter()

const discount = computed(() => {
  if (!props.sale.product) return 0
  return Math.round(((props.sale.product.original_price - props.sale.flash_price) / props.sale.product.original_price) * 100)
})

const progress = computed(() => {
  if (props.sale.total_stock === 0) return 0
  return Math.round(((props.sale.total_stock - props.sale.available_stock) / props.sale.total_stock) * 100)
})

const isStarted = computed(() => {
  const now = new Date().getTime()
  const start = new Date(props.sale.start_time).getTime()
  return now >= start
})

const isEnded = computed(() => {
  const now = new Date().getTime()
  const end = new Date(props.sale.end_time).getTime()
  return now >= end
})

const goToDetail = () => {
  router.push(`/flash-sales/${props.sale.id}`)
}
</script>

<template>
  <div 
    class="group relative bg-surface border border-white/5 overflow-hidden transition-all duration-500 hover:border-white/10 hover:shadow-2xl hover:shadow-black/50 cursor-pointer flex flex-col h-full"
    @click="goToDetail"
  >
    <!-- Image Section -->
    <div class="aspect-[4/3] overflow-hidden bg-surface-light relative">
      <SmartImage
        :src="sale.product?.image_url"
        :alt="sale.product?.name || 'Product'"
        class-name="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110 opacity-90 group-hover:opacity-100"
      />

      <!-- Status Badge -->
      <div class="absolute top-4 left-4 z-10">
        <div v-if="!isStarted" class="px-3 py-1 bg-black/80 backdrop-blur text-white text-xs font-bold tracking-widest uppercase border border-white/10">
          Upcoming
        </div>
        <div v-else-if="isEnded" class="px-3 py-1 bg-surface/80 backdrop-blur text-secondary text-xs font-bold tracking-widest uppercase border border-white/10">
          Ended
        </div>
        <div v-else class="flex items-center gap-2 px-3 py-1 bg-accent text-white text-xs font-bold tracking-widest uppercase animate-pulse-red shadow-lg shadow-accent/20">
          <span class="w-1.5 h-1.5 bg-white rounded-full animate-ping"></span>
          Live
        </div>
      </div>

      <!-- Discount Badge -->
      <div v-if="discount > 0" class="absolute top-4 right-4 bg-white text-black font-bold text-xs px-2 py-1 z-10">
        -{{ discount }}%
      </div>
      
      <!-- Overlay -->
      <div class="absolute inset-0 bg-gradient-to-t from-surface via-transparent to-transparent opacity-60"></div>
    </div>

    <!-- Content Section -->
    <div class="p-6 flex flex-col flex-grow">
      <h3 class="text-lg font-bold text-white mb-2 line-clamp-1 group-hover:text-accent transition-colors">
        {{ sale.product?.name }}
      </h3>
      
      <div class="flex items-baseline gap-3 mb-6">
        <span class="text-xl font-mono font-bold text-white">¥{{ sale.flash_price }}</span>
        <span class="text-sm text-tertiary line-through font-mono">¥{{ sale.product?.original_price }}</span>
      </div>

      <!-- Progress Bar (Only for active sales) -->
      <div v-if="isStarted && !isEnded" class="mb-6">
        <div class="flex justify-between text-xs text-secondary mb-2 uppercase tracking-wider">
          <span>Sold: {{ progress }}%</span>
          <span class="font-mono text-white">{{ sale.available_stock }}<span class="text-tertiary">/{{ sale.total_stock }}</span></span>
        </div>
        <div class="h-1 w-full bg-surface-light overflow-hidden relative">
          <div 
            class="h-full bg-accent transition-all duration-1000 ease-out absolute top-0 left-0" 
            :style="{ width: `${progress}%` }"
          ></div>
        </div>
      </div>

      <!-- Action Area -->
      <div class="flex items-center justify-between mt-auto pt-4 border-t border-white/5">
        <div v-if="!isStarted">
          <CountDown :target-time="sale.start_time" label="Starts In" size="sm" />
        </div>
        <div v-else-if="!isEnded" class="text-accent text-sm font-bold tracking-wider uppercase flex items-center gap-2 group/btn">
          Rush Now <ArrowRight class="w-4 h-4 transition-transform group-hover/btn:translate-x-1" />
        </div>
        <div v-else class="text-secondary text-sm uppercase tracking-wider">
          View Details
        </div>
      </div>
    </div>
  </div>
</template>