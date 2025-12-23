<template>
  <FlowBorderCard 
    :active="isStarted && !isEnded" 
    :color="isStarted && !isEnded ? '#e33535' : '#333'"
    border-radius="12px"
    border-width="2px"
    class="h-full group cursor-pointer"
    @click="goToDetail"
  >
    <div class="relative h-full flex flex-col bg-surface overflow-hidden rounded-[10px]">
      <!-- Image Area -->
      <div class="relative h-64 overflow-hidden bg-gray-900">
        <img :src="item?.image || '/placeholder.png'" :alt="item?.name" class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110" />
        
        <!-- Status Badge -->
        <div class="absolute top-2 right-2 px-2 py-1 bg-black/70 backdrop-blur-md rounded text-xs font-mono text-white border border-white/10 z-10">
          {{ statusText }}
        </div>
        
        <!-- Progress Bar (Overlay at bottom of image) -->
        <div v-if="item?.status === 1" class="absolute bottom-0 inset-x-0 h-1 bg-gray-800 z-10">
          <div class="h-full bg-accent animate-pulse" :style="{ width: progress + '%' }"></div>
        </div>

         <!-- Hover Overlay -->
        <div class="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
        </div>
      </div>
      
      <!-- Content Area -->
      <div class="p-4 flex flex-col flex-grow">
        <div class="flex justify-between items-start mb-2">
          <h3 class="text-lg font-bold text-white line-clamp-1 group-hover:text-accent transition-colors">{{ item?.name }}</h3>
          <span v-if="item?.discount" class="text-xs font-bold text-accent bg-accent/10 px-2 py-1 rounded">-{{ item.discount }}%</span>
        </div>
        
        <div class="flex items-baseline gap-2 mb-4">
          <span class="text-2xl font-bold text-white">${{ item?.price }}</span>
          <span v-if="item?.original_price" class="text-sm text-gray-500 line-through decoration-white/20">${{ item.original_price }}</span>
          <!-- Fallback calc if original_price not present but discount is -->
          <span v-else-if="item?.price && item?.discount" class="text-sm text-gray-500 line-through decoration-white/20">
             ${{ (item.price * 100 / (100 - item.discount)).toFixed(0) }}
          </span>
        </div>

        <div class="space-y-2 mt-auto">
          <div class="flex justify-between text-xs text-gray-400 font-mono">
            <span>{{ item?.available_stock }} left</span>
            <span>{{ progress }}% claimed</span>
          </div>
          <div class="w-full bg-white/5 rounded-full h-1.5 overflow-hidden">
            <div class="bg-gradient-to-r from-accent to-orange-500 h-full rounded-full transition-all duration-1000" :style="{ width: progress + '%' }"></div>
          </div>
        </div>

        <div class="mt-4 pt-4 border-t border-white/5 flex justify-between items-center text-xs text-gray-500 group-hover:text-white transition-colors duration-300">
           <span class="flex items-center gap-1">
             <Clock class="w-3 h-3" /> 
             {{ timeText }}
           </span>
           <span class="group-hover:translate-x-1 transition-transform flex items-center gap-1">
             View Details <ArrowRight class="w-3 h-3" />
           </span>
        </div>
      </div>
    </div>
  </FlowBorderCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Clock, ArrowRight } from 'lucide-vue-next'
import FlowBorderCard from './FlowBorderCard.vue'

// Define props to accept 'item' instead of 'sale' to match HomeView usage
// Or adapt HomeView to pass 'sale'. Let's stick to 'item' as used in recent updates.
const props = defineProps<{
  item: any
}>()

const router = useRouter()

const goToDetail = () => {
  if (props.item?.id) {
    router.push(`/sale/${props.item.id}`)
  }
}

const statusText = computed(() => {
  if (props.item?.status === 1) return 'LIVE NOW'
  if (props.item?.status === 0) return 'UPCOMING'
  return 'ENDED'
})

const isStarted = computed(() => props.item?.status === 1 || props.item?.status === 2)
const isEnded = computed(() => props.item?.status === 2)

const timeText = computed(() => {
  if (props.item?.status === 1) return 'Ends soon'
  return 'Coming soon'
})

const progress = computed(() => {
  if (!props.item?.total_stock) return 0
  const sold = props.item.total_stock - props.item.available_stock
  return Math.round((sold / props.item.total_stock) * 100)
})
</script>