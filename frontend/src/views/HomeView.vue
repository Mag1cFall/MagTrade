<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getFlashSales } from '@/api/flash-sale'
import type { FlashSale } from '@/types'
import FlashSaleCard from '@/components/FlashSaleCard.vue'
import { Loader2 } from 'lucide-vue-next'

const loading = ref(true)
const flashSales = ref<FlashSale[]>([])
const activeTab = ref<'active' | 'upcoming'>('active')

const fetchData = async () => {
  loading.value = true
  flashSales.value = [] // Clear list to avoid confusion
  try {
    // status: 1 (Active) or 0 (Pending)
    const status = activeTab.value === 'active' ? 1 : 0
    const res = await getFlashSales(1, 20, status)
    if (res.code === 0) {
      flashSales.value = res.data.flash_sales
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const switchTab = (tab: 'active' | 'upcoming') => {
  activeTab.value = tab
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="min-h-[80vh]">
    <!-- Hero Section -->
    <section class="py-16 md:py-24 border-b border-white/5">
      <div class="max-w-4xl">
        <h1 class="text-5xl md:text-7xl font-bold tracking-tighter text-white mb-6 leading-none">
          SPEED <span class="text-accent">DEFINES</span><br />WINNER.
        </h1>
        <p class="text-xl text-secondary max-w-2xl leading-relaxed">
          Experience the thrill of millisecond-level flash sales. 
          Powered by edge computing and AI analysis.
        </p>
      </div>
    </section>

    <!-- Filter Tabs -->
    <div class="flex items-center gap-8 mt-12 mb-8 border-b border-white/5 pb-4 sticky top-[73px] z-30 bg-background/95 backdrop-blur pt-4">
      <button 
        @click="switchTab('active')"
        :class="[
          'text-sm font-bold tracking-widest uppercase transition-colors duration-300 relative',
          activeTab === 'active' ? 'text-accent' : 'text-secondary hover:text-white'
        ]"
      >
        Live Now
        <span v-if="activeTab === 'active'" class="absolute -bottom-4 left-0 w-full h-0.5 bg-accent"></span>
      </button>
      <button 
        @click="switchTab('upcoming')"
        :class="[
          'text-sm font-bold tracking-widest uppercase transition-colors duration-300 relative',
          activeTab === 'upcoming' ? 'text-white' : 'text-secondary hover:text-white'
        ]"
      >
        Upcoming
        <span v-if="activeTab === 'upcoming'" class="absolute -bottom-4 left-0 w-full h-0.5 bg-white"></span>
      </button>
    </div>

    <!-- Grid -->
    <div v-if="loading" class="flex justify-center py-20">
      <Loader2 class="w-8 h-8 text-accent animate-spin" />
    </div>
    
    <div v-else-if="flashSales.length === 0" class="py-20 text-center">
      <p class="text-secondary text-lg font-mono">NO ACTIVE EVENTS_</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 pb-20">
      <FlashSaleCard 
        v-for="sale in flashSales" 
        :key="sale.id" 
        :sale="sale" 
      />
    </div>
  </div>
</template>