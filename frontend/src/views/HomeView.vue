<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getFlashSales } from '@/api/flash-sale'
import type { FlashSale } from '@/types'
import FlashSaleCard from '@/components/FlashSaleCard.vue'
import MarqueeTicker from '@/components/MarqueeTicker.vue'
import VerticalPager from '@/components/VerticalPager.vue'
import ComparisonTable from '@/components/ComparisonTable.vue'
import MapVisualization from '@/components/MapVisualization.vue'
import ParticleBackground from '@/components/ParticleBackground.vue'
import LiveFeed from '@/components/LiveFeed.vue'
import StatsCounter from '@/components/StatsCounter.vue'
import BrandWall from '@/components/BrandWall.vue'
import FilterToolbar from '@/components/FilterToolbar.vue'
import TagCloud from '@/components/TagCloud.vue'
import NodeLatencyList from '@/components/NodeLatencyList.vue'
import StatusRing from '@/components/StatusRing.vue'
import TechStack from '@/components/TechStack.vue'
import TestimonialCarousel from '@/components/TestimonialCarousel.vue'
import HowItWorks from '@/components/HowItWorks.vue'
import PlatformStats from '@/components/PlatformStats.vue'
import FeatureCards from '@/components/FeatureCards.vue'
import PartnersCarousel from '@/components/PartnersCarousel.vue'
import FAQAccordion from '@/components/FAQAccordion.vue'
import NewsletterCTA from '@/components/NewsletterCTA.vue'
import { Loader2, Zap, Activity, ChevronDown, Shield, Globe, CheckCircle, ArrowRight } from 'lucide-vue-next'

const router = useRouter()
const loading = ref(true)
const flashSales = ref<FlashSale[]>([])
const activeTab = ref<'active' | 'upcoming'>('active')
const currentSectionIndex = ref(0)

const sections = ['HOME', 'SYSTEM', 'FEATURES', 'MARKET', 'NETWORK', 'SPECS', 'SECURITY', 'FAQ']
const sectionIds = ['hero', 'stats', 'features', 'market', 'network', 'specs', 'trust', 'faq']

const tickerItems = [
  'SYSTEM ONLINE', 'LOW LATENCY DETECTED', 'ANTI-BOT PROTECTION ACTIVE',
  'HIGH-SPEED RPC READY', 'OBSIDIAN ENGINE ENGAGED', 'REAL-TIME SYNC',
  'SECURE TRANSACTION LAYER', 'AI TACTICAL ANALYSIS', 'GLOBAL NODES: 35+'
]

const fetchData = async () => {
  if (flashSales.value.length === 0) loading.value = true
  try {
    const status = activeTab.value === 'active' ? 1 : 0
    const res = await getFlashSales(1, 6, status)
    if (res.code === 0) {
      flashSales.value = res.data.flash_sales
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

watch(activeTab, () => fetchData())

const scrollToSection = (index: number) => {
  const id = sectionIds[index]
  if (id) {
    const el = document.getElementById(id)
    if (el) {
      el.scrollIntoView({ behavior: 'smooth' })
    }
  }
}

const goToShop = () => router.push('/shop')
const goToRegister = () => router.push('/register')

let observer: IntersectionObserver | null = null

onMounted(() => {
  fetchData()
  observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const index = sectionIds.indexOf(entry.target.id)
        if (index !== -1) currentSectionIndex.value = index
      }
    })
  }, { threshold: 0.5 })
  sectionIds.forEach(id => {
    const el = document.getElementById(id)
    if (el) observer?.observe(el)
  })
})

onUnmounted(() => observer?.disconnect())
</script>

<template>
  <div class="min-h-screen bg-background text-primary selection:bg-accent selection:text-white">
    <VerticalPager :sections="sections" :active-index="currentSectionIndex" @change="scrollToSection" />

    <section id="hero" class="h-[100vh] w-full relative flex flex-col items-center justify-center overflow-hidden border-b border-white/5">
      <ParticleBackground />
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-surface-light/10 via-background/80 to-background opacity-80 pointer-events-none z-0"></div>

      <div class="flex flex-col justify-center items-center text-center px-6 z-10 relative">
        <div class="mb-8 inline-flex items-center gap-2 px-4 py-1.5 border border-accent/30 rounded-full bg-accent/5 text-accent text-xs font-mono tracking-widest uppercase animate-pulse-fast">
          <span class="w-2 h-2 bg-accent rounded-full animate-ping"></span>
          MagTrade v2.0 Live
        </div>
        
        <h1 class="font-bold tracking-tighter text-white mb-8 leading-none hero-title scale-90 md:scale-100 lg:scale-110">
          <span class="stagger-text" style="--i:1">SPEED</span>
          <span class="mx-2 md:mx-4"></span>
          <span class="text-transparent bg-clip-text bg-gradient-to-r from-white to-gray-500 stagger-text" style="--i:2">KILLS.</span>
          <br />
          <span class="text-accent stagger-text" style="--i:3">WINNER</span>
          <span class="mx-2 md:mx-4"></span>
          <span class="stagger-text" style="--i:4">TAKES</span>
          <span class="mx-2 md:mx-4"></span>
          <span class="stagger-text" style="--i:5">ALL.</span>
        </h1>
        
        <p class="text-lg md:text-xl text-secondary max-w-2xl leading-relaxed mb-12 animate-fade-in" style="animation-delay: 0.8s">
          Dominate the market with millisecond precision.
          <br class="hidden md:block" />
          Powered by cutting-edge AI and distributed ledger technology.
        </p>

        <div class="flex flex-col sm:flex-row gap-4">
          <button 
            @click="goToShop"
            class="group relative overflow-hidden flex items-center gap-3 px-10 py-5 bg-accent text-white font-bold tracking-widest uppercase transition-all duration-300 shadow-[0_0_30px_rgba(227,53,53,0.3)] hover:shadow-[0_0_50px_rgba(227,53,53,0.8)] hover:scale-105"
          >
            <div class="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300"></div>
            <Zap class="w-5 h-5 fill-current relative z-10" />
            <span class="relative z-10">Start Trading</span>
            <ArrowRight class="w-4 h-4 relative z-10 opacity-0 group-hover:opacity-100 transition-all" />
          </button>
          
          <button 
            @click="scrollToSection(1)"
            class="px-10 py-5 border border-white/20 text-white font-bold tracking-widest uppercase hover:bg-white/5 transition-colors"
          >
            Learn More
          </button>
        </div>
      </div>

      <div class="absolute bottom-0 left-0 right-0 z-20 flex flex-col pointer-events-none">
         <div class="w-full py-3 border-t border-white/5 bg-surface/40 backdrop-blur-xl pointer-events-auto">
            <MarqueeTicker :duration="40">
              <div v-for="(item, i) in tickerItems" :key="i" class="flex items-center gap-2 text-xs font-mono text-tertiary uppercase tracking-widest px-4">
                <Activity class="w-3 h-3 text-accent" />
                {{ item }}
              </div>
            </MarqueeTicker>
          </div>
         <StatsCounter class="pointer-events-auto" />
         <LiveFeed class="pointer-events-auto" />
      </div>
    </section>

    <section id="stats" class="bg-background pt-24 pb-12">
      <div class="container mx-auto px-6 md:px-12">
        <PlatformStats />
        <HowItWorks />
      </div>
    </section>

    <section id="features" class="bg-background py-12">
      <div class="container mx-auto px-6 md:px-12">
        <FeatureCards />
        <PartnersCarousel />
      </div>
    </section>

    <section id="market" class="min-h-screen bg-background border-b border-white/5 relative z-10 py-20">
      <div class="mb-16">
        <BrandWall />
      </div>

      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-5xl font-bold text-white mb-4 tracking-tight">FEATURED DROPS</h2>
          <p class="text-secondary">The hottest flash sales happening right now</p>
          <div class="h-1 w-20 bg-accent mx-auto mt-4"></div>
        </div>

        <div class="flex flex-col lg:flex-row gap-8">
          <div class="lg:w-1/4 hidden lg:block space-y-6">
            <TagCloud />
            <div class="p-4 border border-accent/30 rounded-xl bg-accent/5">
              <p class="text-accent text-sm font-semibold mb-2">🔥 Pro Tip</p>
              <p class="text-secondary text-xs">Use our AI advisor to get success probability before each drop!</p>
            </div>
          </div>

          <div class="lg:w-3/4">
            <FilterToolbar v-model="activeTab" />

            <div v-if="loading" class="flex justify-center py-20">
              <Loader2 class="w-10 h-10 animate-spin text-accent" />
            </div>
            
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <FlashSaleCard v-for="item in flashSales" :key="item.id" :item="item" class="h-full" />
            </div>
            
            <div class="mt-12 text-center">
              <button @click="goToShop" class="px-8 py-3 border border-white/20 text-white hover:bg-white/5 transition-colors rounded">
                View All Products <ArrowRight class="w-4 h-4 inline ml-2" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="network" class="min-h-screen relative bg-black flex items-center overflow-hidden border-b border-white/5 py-20">
      <div class="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:40px_40px]"></div>

      <div class="max-w-7xl mx-auto px-4 w-full relative z-10 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
        <div>
          <div class="inline-flex items-center gap-2 px-3 py-1 border border-accent/30 rounded-full bg-accent/5 text-accent text-xs font-mono mb-6">
            <Globe class="w-3 h-3" />
            GLOBAL INFRASTRUCTURE
          </div>
          
          <h2 class="text-4xl md:text-6xl font-bold text-white mb-6 tracking-tighter">
            ZERO LATENCY.<br/>EVERYWHERE.
          </h2>
          
          <p class="text-secondary text-lg mb-10 leading-relaxed max-w-md">
            Our distributed edge network ensures your orders hit the matching engine faster than the blink of an eye.
          </p>

          <div class="flex gap-8 mb-10">
            <StatusRing :value="99" label="Uptime" />
            <StatusRing :value="35" label="Nodes" />
          </div>

          <NodeLatencyList />
        </div>

        <div class="h-[500px] w-full bg-surface/20 rounded-2xl border border-white/5 overflow-hidden relative">
          <MapVisualization />
          <div class="absolute inset-0 pointer-events-none bg-gradient-to-t from-black/80 via-transparent to-transparent"></div>
          <div class="absolute bottom-6 left-6 text-xs font-mono text-gray-500">
            LIVE MAP VISUALIZATION // SYSTEM ACTIVE
          </div>
        </div>
      </div>
    </section>

    <section id="specs" class="min-h-screen bg-background border-b border-white/5 py-20">
      <div class="max-w-7xl mx-auto px-4">
        <div class="text-center mb-20">
          <h2 class="text-3xl md:text-5xl font-bold text-white mb-6">BUILT FOR SPEED</h2>
          <p class="text-secondary max-w-2xl mx-auto">
            Engineered with a modern stack designed for millions of concurrent connections.
          </p>
        </div>

        <div class="mb-24">
          <TechStack />
        </div>

        <div class="bg-surface/30 border border-white/10 rounded-2xl overflow-hidden p-8 md:p-12 relative">
          <div class="absolute top-0 right-0 p-4 text-xs font-mono text-gray-600 uppercase tracking-widest border-l border-b border-white/5">
            Architecture Specs
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
            <div>
              <h3 class="text-2xl font-bold text-white mb-6">Why We Win</h3>
              <ul class="space-y-4">
                <li v-for="i in ['Event-Driven Architecture', 'In-Memory Matching Engine', 'Geo-Distributed Consistency', 'AI Fraud Detection']" :key="i" class="flex items-center gap-3 text-secondary">
                  <div class="w-6 h-6 rounded-full bg-accent/20 flex items-center justify-center text-accent">
                    <CheckCircle class="w-4 h-4" />
                  </div>
                  {{ i }}
                </li>
              </ul>
            </div>
            
            <div class="bg-black/50 rounded-xl p-4 border border-white/5">
              <ComparisonTable />
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="trust" class="min-h-screen flex flex-col justify-center py-20 bg-black relative">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_top,_var(--tw-gradient-stops))] from-gray-900 via-black to-black opacity-50"></div>

      <div class="max-w-7xl w-full mx-auto px-4 relative z-10">
        <div class="mb-20">
          <h2 class="text-8xl font-bold text-white/5 absolute top-0 w-full text-center transform -translate-y-1/2 select-none pointer-events-none">TRUSTED</h2>
          <div class="text-center relative">
            <Shield class="w-12 h-12 text-accent mx-auto mb-6" />
            <h2 class="text-4xl font-bold text-white mb-4">BANK-GRADE SECURITY</h2>
            <p class="text-secondary">Audited by top security firms. Your assets are SAFU.</p>
          </div>
        </div>

        <div class="mb-20">
          <TestimonialCarousel />
        </div>
      </div>
    </section>

    <section id="faq">
      <FAQAccordion />
      <NewsletterCTA />
    </section>

    <footer class="bg-black border-t border-white/10 py-12">
      <div class="max-w-6xl mx-auto px-4">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
          <div>
            <h3 class="text-white font-bold text-lg mb-4">MagTrade</h3>
            <p class="text-secondary text-sm">The world's fastest flash sale platform. Powered by AI.</p>
          </div>
          <div>
            <h4 class="text-white font-semibold mb-4">Product</h4>
            <ul class="space-y-2 text-sm text-secondary">
              <li><router-link to="/shop" class="hover:text-white">Shop</router-link></li>
              <li><a href="#features" class="hover:text-white">Features</a></li>
              <li><a href="#specs" class="hover:text-white">Technology</a></li>
            </ul>
          </div>
          <div>
            <h4 class="text-white font-semibold mb-4">Company</h4>
            <ul class="space-y-2 text-sm text-secondary">
              <li><router-link to="/privacy" class="hover:text-white">Privacy Policy</router-link></li>
              <li><router-link to="/terms" class="hover:text-white">Terms of Service</router-link></li>
              <li><router-link to="/contact" class="hover:text-white">Contact Support</router-link></li>
            </ul>
          </div>
          <div>
            <h4 class="text-white font-semibold mb-4">Connect</h4>
            <ul class="space-y-2 text-sm text-secondary">
              <li><a href="#" class="hover:text-white">Twitter</a></li>
              <li><a href="#" class="hover:text-white">Discord</a></li>
              <li><a href="#" class="hover:text-white">GitHub</a></li>
            </ul>
          </div>
        </div>
        
        <div class="border-t border-white/10 pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
          <p class="text-xs text-gray-600">© 2025 MagTrade System. All rights reserved.</p>
          <button @click="goToRegister" class="px-6 py-2 bg-accent text-white text-sm font-semibold rounded hover:bg-accent/90 transition-colors">
            Join the Revolution
          </button>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.hero-title {
  font-size: clamp(2.5rem, 10vw, 8.5rem);
  line-height: 0.9;
  letter-spacing: -0.05em;
  width: 100%;
  padding: 0 1rem;
}

.stagger-text {
  display: inline-block;
  opacity: 0;
  transform: translateY(20px);
  animation: slideUpFade 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  animation-delay: calc(var(--i) * 0.1s);
}

@keyframes slideUpFade {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>