<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getFlashSales } from '@/api/flash-sale'
import type { FlashSale } from '@/types'
import FlashSaleCard from '@/components/FlashSaleCard.vue'
import MarqueeTicker from '@/components/MarqueeTicker.vue'
import VerticalPager from '@/components/VerticalPager.vue'
import ComparisonTable from '@/components/ComparisonTable.vue'
import MapVisualization from '@/components/MapVisualization.vue'
import { Loader2, Zap, Terminal, Activity, ChevronDown, Shield, Globe, Cpu } from 'lucide-vue-next'

const router = useRouter()
const loading = ref(true)
const flashSales = ref<FlashSale[]>([])
const activeTab = ref<'active' | 'upcoming'>('active')
const currentSectionIndex = ref(0)

const sections = ['HERO', 'MARKET', 'NETWORK', 'SPECS', 'TRUST']
const sectionIds = ['hero', 'market', 'network', 'specs', 'trust']

const tickerItems = [
  'SYSTEM ONLINE', 'LOW LATENCY NETWORK DETECTED', 'ANTI-BOT PROTECTION ACTIVE',
  'HIGH-SPEED RPC READY', 'FLOW BORDER ENGINE ENGAGED', 'OBSIDIAN VELOCITY UI LOADED',
  'SECURE TRANSACTION LAYER', 'REAL-TIME STOCK SYNC', 'AI TACTICAL ANALYSIS'
]

const fetchData = async () => {
  // Keep loading true only for initial load to prevent layout shift
  if (flashSales.value.length === 0) loading.value = true
  
  try {
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

const scrollToSection = (index: number) => {
  const id = sectionIds[index]
  if (id) {
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
  }
}

const scrollToGrid = () => scrollToSection(1)

// Intersection Observer for scroll spy
let observer: IntersectionObserver | null = null

onMounted(() => {
  fetchData()

  observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const index = sectionIds.indexOf(entry.target.id)
        if (index !== -1) {
          currentSectionIndex.value = index
        }
      }
    })
  }, { threshold: 0.5 })

  sectionIds.forEach(id => {
    const el = document.getElementById(id)
    if (el) observer?.observe(el)
  })
})

onUnmounted(() => {
  observer?.disconnect()
})
</script>

<template>
  <div class="min-h-screen bg-background text-primary selection:bg-accent selection:text-white">
    <VerticalPager 
      :sections="sections" 
      :active-index="currentSectionIndex" 
      @change="scrollToSection"
    />

    <!-- SECTION 1: HERO -->
    <section id="hero" class="h-screen flex flex-col relative overflow-hidden border-b border-white/5">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-surface-light/20 via-background to-background opacity-50 pointer-events-none"></div>
      
      <div class="flex-grow flex flex-col justify-center items-center text-center px-4 z-10">
        <div class="mb-8 inline-flex items-center gap-2 px-4 py-1.5 border border-accent/30 rounded-full bg-accent/5 text-accent text-xs font-mono tracking-widest uppercase animate-pulse-fast">
          <span class="w-2 h-2 bg-accent rounded-full animate-ping"></span>
          MagTrade v2.0 Live
        </div>
        
        <h1 class="text-6xl md:text-8xl lg:text-9xl font-bold tracking-tighter text-white mb-8 leading-none mix-blend-difference animate-slide-up">
          SPEED <span class="text-transparent bg-clip-text bg-gradient-to-r from-white to-gray-500">KILLS</span><br />
          <span class="text-accent">WINNER</span> TAKES ALL.
        </h1>
        
        <p class="text-lg md:text-xl text-secondary max-w-2xl leading-relaxed mb-12 animate-fade-in" style="animation-delay: 0.2s">
          Dominate the market with millisecond precision.
          <br class="hidden md:block" />
          Powered by cutting-edge AI and distributed ledger technology.
        </p>

        <button 
          @click="scrollToGrid"
          class="group relative overflow-hidden flex items-center gap-3 px-10 py-5 bg-accent text-white font-bold tracking-widest uppercase transition-all duration-300 animate-fade-in shadow-[0_0_30px_rgba(227,53,53,0.3)] hover:shadow-[0_0_50px_rgba(227,53,53,0.6)] hover:scale-105"
          style="animation-delay: 0.4s"
        >
          <div class="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300"></div>
          <Zap class="w-5 h-5 fill-current relative z-10" />
          <span class="relative z-10">Start Trading</span>
          <ChevronDown class="w-4 h-4 relative z-10 opacity-0 group-hover:opacity-100 transition-all duration-300 translate-y-1 group-hover:translate-y-0" />
        </button>
      </div>

      <div class="w-full py-3 border-t border-white/5 bg-surface/30 backdrop-blur-md z-10">
        <MarqueeTicker :duration="40">
          <div v-for="(item, i) in tickerItems" :key="i" class="flex items-center gap-2 text-xs font-mono text-tertiary uppercase tracking-widest px-4">
            <Activity class="w-3 h-3 text-accent" />
            {{ item }}
          </div>
        </MarqueeTicker>
      </div>
    </section>

    <!-- SECTION 2: MARKET GRID -->
    <section id="market" class="min-h-screen py-24 px-6 md:px-12 bg-background relative flex flex-col border-b border-white/5">
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-8 mb-12">
        <div>
          <h2 class="text-5xl font-bold text-white tracking-tighter flex items-center gap-4">
            <Terminal class="w-10 h-10 text-accent" />
            MARKET <span class="text-secondary">FEED</span>
          </h2>
          <p class="text-secondary text-sm mt-3 font-mono uppercase tracking-widest pl-1">Real-time inventory status</p>
        </div>

        <!-- Styled Tabs -->
        <div class="flex p-1 bg-surface border border-white/10 rounded-lg">
          <button 
            @click="switchTab('active')"
            :class="[
              'px-8 py-3 text-xs font-bold uppercase tracking-widest transition-all duration-300 rounded-md relative overflow-hidden',
              activeTab === 'active' ? 'text-white' : 'text-secondary hover:text-white'
            ]"
          >
            <div v-if="activeTab === 'active'" class="absolute inset-0 bg-accent opacity-100"></div>
            <span class="relative z-10">Live Auctions</span>
          </button>
          <button 
            @click="switchTab('upcoming')"
            :class="[
              'px-8 py-3 text-xs font-bold uppercase tracking-widest transition-all duration-300 rounded-md relative overflow-hidden',
              activeTab === 'upcoming' ? 'text-black' : 'text-secondary hover:text-white'
            ]"
          >
            <div v-if="activeTab === 'upcoming'" class="absolute inset-0 bg-white opacity-100"></div>
            <span class="relative z-10">Upcoming</span>
          </button>
        </div>
      </div>

      <div class="flex-grow min-h-[600px] relative">
        <transition name="fade" mode="out-in">
          <div v-if="loading" class="absolute inset-0 flex items-center justify-center">
            <Loader2 class="w-16 h-16 text-accent animate-spin" />
          </div>
          
          <div v-else-if="flashSales.length === 0" class="absolute inset-0 flex flex-col items-center justify-center border border-dashed border-white/10 rounded-lg">
            <div class="p-4 bg-surface rounded-full mb-4">
              <Activity class="w-8 h-8 text-secondary" />
            </div>
            <p class="text-secondary text-xl font-mono">NO SIGNAL DETECTED</p>
            <p class="text-tertiary text-sm mt-2">Check back later for new drops</p>
          </div>

          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            <FlashSaleCard 
              v-for="sale in flashSales" 
              :key="sale.id" 
              :sale="sale" 
            />
          </div>
        </transition>
      </div>
    </section>

    <!-- SECTION 3: NETWORK MAP -->
    <section id="network" class="min-h-screen flex flex-col justify-center py-24 px-6 md:px-12 bg-surface/30 border-b border-white/5 relative">
      <div class="max-w-7xl mx-auto w-full">
        <div class="mb-16">
          <div class="flex items-center gap-3 mb-4">
            <Globe class="w-6 h-6 text-accent" />
            <span class="text-accent text-sm font-mono tracking-widest uppercase">Global Infrastructure</span>
          </div>
          <h2 class="text-5xl md:text-7xl font-bold text-white tracking-tighter mb-6">
            LOW LATENCY.<br />ANYWHERE.
          </h2>
          <p class="text-secondary text-lg max-w-2xl leading-relaxed">
            Our distributed edge nodes ensure you have the fastest connection to the matching engine, no matter where you are located on the planet.
          </p>
        </div>

        <MapVisualization />
      </div>
    </section>

    <!-- SECTION 4: SPECS -->
    <section id="specs" class="min-h-screen flex flex-col justify-center py-24 px-6 md:px-12 border-b border-white/5 relative bg-background">
      <div class="max-w-7xl mx-auto w-full">
        <div class="mb-16 flex flex-col md:flex-row md:items-end justify-between gap-8">
          <div>
            <div class="flex items-center gap-3 mb-4">
              <Cpu class="w-6 h-6 text-accent" />
              <span class="text-accent text-sm font-mono tracking-widest uppercase">System Architecture</span>
            </div>
            <h2 class="text-5xl md:text-6xl font-bold text-white tracking-tighter">
              BUILT FOR SPEED.
            </h2>
          </div>
          <p class="text-secondary text-right max-w-md">
            See how MagTrade's Obsidian Velocity protocol outperforms traditional monolithic e-commerce architectures.
          </p>
        </div>

        <ComparisonTable />
      </div>
    </section>

    <!-- SECTION 5: TRUST & FOOTER -->
    <section id="trust" class="min-h-[80vh] flex flex-col justify-between pt-24 px-6 md:px-12 bg-surface relative">
      <div class="max-w-7xl mx-auto w-full grid grid-cols-1 md:grid-cols-3 gap-12">
        <div class="p-8 border border-white/10 bg-background/50 backdrop-blur hover:border-accent/50 transition-colors group">
          <Shield class="w-10 h-10 text-white mb-6 group-hover:text-accent transition-colors" />
          <h3 class="text-2xl font-bold text-white mb-4">Bank-Grade Security</h3>
          <p class="text-secondary leading-relaxed">
            All transactions are encrypted with AES-256 and processed through PCI-DSS compliant gateways. Your assets are safe.
          </p>
        </div>
        <div class="p-8 border border-white/10 bg-background/50 backdrop-blur hover:border-accent/50 transition-colors group">
          <Zap class="w-10 h-10 text-white mb-6 group-hover:text-accent transition-colors" />
          <h3 class="text-2xl font-bold text-white mb-4">Instant Settlement</h3>
          <p class="text-secondary leading-relaxed">
            No more waiting. Our atomic swap engine ensures that ownership transfer happens in the same block as payment.
          </p>
        </div>
        <div class="p-8 border border-white/10 bg-background/50 backdrop-blur hover:border-accent/50 transition-colors group">
          <Terminal class="w-10 h-10 text-white mb-6 group-hover:text-accent transition-colors" />
          <h3 class="text-2xl font-bold text-white mb-4">Developer API</h3>
          <p class="text-secondary leading-relaxed">
            Full programmatic access to our matching engine via gRPC and WebSocket. Build your own trading bots.
          </p>
        </div>
      </div>

      <footer class="mt-24 py-12 border-t border-white/10 flex flex-col md:flex-row justify-between items-center gap-6">
        <div class="text-secondary text-sm">
          © 2025 MagTrade Inc. All rights reserved.
        </div>
        <div class="flex gap-8 text-sm font-bold uppercase tracking-widest">
          <button @click="router.push('/legal')" class="text-secondary hover:text-white transition-colors">Privacy</button>
          <button @click="router.push('/legal')" class="text-secondary hover:text-white transition-colors">Terms</button>
          <button @click="router.push('/legal')" class="text-secondary hover:text-white transition-colors">Support</button>
        </div>
      </footer>
    </section>
  </div>
</template>

<style scoped>
/* Ensure smooth fade transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>