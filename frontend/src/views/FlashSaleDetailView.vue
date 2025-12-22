<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getFlashSaleDetail, rushFlashSale, getFlashSaleStock } from '@/api/flash-sale'
import { getRecommendation } from '@/api/ai'
import type { FlashSaleDetail } from '@/types'
import CountDown from '@/components/CountDown.vue'
import SmartImage from '@/components/SmartImage.vue'
import { useAuthStore } from '@/stores/auth'
import { Loader2, ShieldCheck, Zap, Server, Activity, ArrowLeft, Bot, Terminal } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const loading = ref(true)
const detail = ref<FlashSaleDetail | null>(null)
const currentStock = ref(0)
const rushLoading = ref(false)
const ws = ref<WebSocket | null>(null)
const rushStatus = ref<'idle' | 'processing' | 'success' | 'failed'>('idle')
const rushMessage = ref('')

const aiAnalysis = ref({
  difficulty: 0,
  winRate: '--%',
  advice: '正在接入神经元网络分析...',
  risk: '计算中'
})

const flashSaleId = parseInt(route.params.id as string)

const fetchData = async () => {
  try {
    const res = await getFlashSaleDetail(flashSaleId)
    if (res.code === 0) {
      detail.value = res.data
      currentStock.value = res.data.current_stock
    }

    if (authStore.isAuthenticated) {
      getRecommendation(flashSaleId).then(res => {
        if (res.code === 0) {
          const { analysis } = res.data
          aiAnalysis.value = {
            difficulty: analysis.difficulty_score,
            winRate: (analysis.success_probability * 100).toFixed(1) + '%',
            advice: analysis.timing_advice,
            risk: analysis.difficulty_reason
          }
        }
      }).catch(() => {
        aiAnalysis.value.advice = 'AI 分析服务暂时不可用'
      })
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

// 轮询库存 (简单实现)
let stockInterval: number
const startStockPolling = () => {
  stockInterval = window.setInterval(async () => {
    try {
      const res = await getFlashSaleStock(flashSaleId)
      if (res.code === 0) {
        currentStock.value = res.data.stock
      }
    } catch (e) {
      // ignore error
    }
  }, 3000)
}

// WebSocket 连接
const connectWS = () => {
  if (!authStore.token) return
  
  // 使用 Vite 代理配置的 /ws 路径
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const wsUrl = `${protocol}//${host}/ws/notifications?token=${authStore.token}`
  
  ws.value = new WebSocket(wsUrl)
  
  ws.value.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      if (msg.type === 'flash_sale_result' && msg.data.flash_sale_id === flashSaleId) {
        rushLoading.value = false
        if (msg.data.success) {
          rushStatus.value = 'success'
          rushMessage.value = msg.data.message
          // 成功后停止库存轮询，避免视觉干扰
          clearInterval(stockInterval)
        } else {
          rushStatus.value = 'failed'
          rushMessage.value = msg.data.message
        }
      }
    } catch (e) {
      console.error('WS parse error', e)
    }
  }
}

const handleRush = async () => {
  if (!authStore.isAuthenticated) {
    router.push(`/login?redirect=${route.fullPath}`)
    return
  }

  rushLoading.value = true
  rushStatus.value = 'processing'
  rushMessage.value = '请求已发送，正在排队...'

  try {
    const res = await rushFlashSale(flashSaleId)
    if (res.code === 0) {
      // 请求提交成功，等待 WS 通知结果
      rushMessage.value = res.data.message || '排队中...'
    } else {
      rushLoading.value = false
      rushStatus.value = 'failed'
      rushMessage.value = res.message
    }
  } catch (e: any) {
    rushLoading.value = false
    rushStatus.value = 'failed'
    rushMessage.value = e.message || '请求失败'
  }
}

onMounted(() => {
  fetchData()
  startStockPolling()
  if (authStore.isAuthenticated) {
    connectWS()
  }
})

onUnmounted(() => {
  clearInterval(stockInterval)
  if (ws.value) ws.value.close()
})

// 计算属性
const discount = computed(() => {
  if (!detail.value?.flash_sale.product) return 0
  const p = detail.value.flash_sale.product
  return Math.round(((p.original_price - detail.value.flash_sale.flash_price) / p.original_price) * 100)
})

const isStarted = computed(() => {
  if (!detail.value) return false
  const now = new Date().getTime()
  const start = new Date(detail.value.flash_sale.start_time).getTime()
  return now >= start
})

const isEnded = computed(() => {
  if (!detail.value) return false
  const now = new Date().getTime()
  const end = new Date(detail.value.flash_sale.end_time).getTime()
  return now >= end
})
</script>

<template>
  <div v-if="loading" class="min-h-screen flex items-center justify-center">
    <Loader2 class="w-10 h-10 text-accent animate-spin" />
  </div>

  <div v-else-if="!detail" class="min-h-screen flex flex-col items-center justify-center gap-4">
    <div class="text-2xl font-bold">EVENT NOT FOUND</div>
    <button @click="router.push('/')" class="text-secondary hover:text-white flex items-center gap-2">
      <ArrowLeft class="w-4 h-4" /> 返回大厅
    </button>
  </div>

  <div v-else class="min-h-screen pb-20">
    <!-- Breadcrumb -->
    <div class="py-6 border-b border-white/5">
      <button @click="router.push('/')" class="text-sm text-secondary hover:text-white flex items-center gap-2 uppercase tracking-widest">
        <ArrowLeft class="w-4 h-4" /> Back to Lobby
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 mt-12">
      <!-- Left: Visuals -->
      <div class="space-y-8">
        <div class="aspect-square bg-surface-light border border-white/5 relative overflow-hidden group">
          <SmartImage
            :src="detail.flash_sale.product?.image_url"
            :alt="detail.flash_sale.product?.name || 'Product'"
            class-name="w-full h-full object-cover transition-transform duration-1000 group-hover:scale-105"
          />
          
          <!-- Discount Tag -->
          <div class="absolute top-0 right-0 bg-white text-black text-xl font-bold px-4 py-2">
            -{{ discount }}%
          </div>
        </div>

        <!-- AI Analysis Panel -->
        <div class="bg-surface border border-white/10 p-6 relative overflow-hidden">
          <div class="absolute top-0 right-0 p-4 opacity-10">
            <Bot class="w-24 h-24" />
          </div>
          <h3 class="text-sm font-bold text-accent uppercase tracking-widest mb-6 flex items-center gap-2">
            <Terminal class="w-4 h-4" /> AI Tactical Analysis
          </h3>
          <div class="grid grid-cols-3 gap-4 mb-6">
            <div>
              <div class="text-xs text-secondary mb-1">DIFFICULTY</div>
              <div class="text-2xl font-mono text-white">{{ aiAnalysis.difficulty }}<span class="text-sm text-tertiary">/10</span></div>
            </div>
            <div>
              <div class="text-xs text-secondary mb-1">EST. WIN RATE</div>
              <div class="text-2xl font-mono text-white">{{ aiAnalysis.winRate }}</div>
            </div>
            <div>
              <div class="text-xs text-secondary mb-1">RISK LEVEL</div>
              <div class="text-2xl font-mono text-white tracking-tighter">{{ aiAnalysis.risk }}</div>
            </div>
          </div>
          <div class="text-sm text-secondary leading-relaxed border-t border-white/10 pt-4">
            <span class="text-accent font-bold">ADVICE: </span> {{ aiAnalysis.advice }}
          </div>
        </div>
      </div>

      <!-- Right: Info & Action -->
      <div class="flex flex-col h-full">
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4 leading-tight">
          {{ detail.flash_sale.product?.name }}
        </h1>
        <p class="text-secondary text-lg mb-8 leading-relaxed">
          {{ detail.flash_sale.product?.description }}
        </p>

        <div class="flex items-baseline gap-4 mb-8">
          <span class="text-5xl font-mono font-bold text-white">¥{{ detail.flash_sale.flash_price }}</span>
          <span class="text-xl text-tertiary line-through font-mono">¥{{ detail.flash_sale.product?.original_price }}</span>
        </div>

        <!-- Status Bar -->
        <div class="bg-surface-light border-l-2 border-accent p-4 mb-8 flex justify-between items-center">
          <div class="flex flex-col">
            <span class="text-xs text-secondary uppercase tracking-widest mb-1">Status</span>
            <span v-if="!isStarted" class="text-white font-bold">UPCOMING</span>
            <span v-else-if="isEnded" class="text-secondary font-bold">ENDED</span>
            <span v-else class="text-accent font-bold animate-pulse">LIVE NOW</span>
          </div>
          <div class="flex flex-col items-end">
             <span class="text-xs text-secondary uppercase tracking-widest mb-1">Stock</span>
             <span class="font-mono text-xl text-white">{{ currentStock }}<span class="text-sm text-tertiary">/{{ detail.flash_sale.total_stock }}</span></span>
          </div>
        </div>

        <!-- Action Area -->
        <div class="mt-auto space-y-6">
          <div v-if="!isStarted" class="p-6 border border-white/10 bg-surface">
            <div class="flex items-center justify-between">
              <span class="text-sm text-secondary uppercase tracking-widest">Dropping In</span>
              <CountDown :target-time="detail.flash_sale.start_time" size="lg" />
            </div>
          </div>

          <div v-else-if="!isEnded">
            <!-- Rush Button State Machine -->
            <button 
              @click="handleRush"
              :disabled="rushLoading || currentStock <= 0 || rushStatus === 'success'"
              :class="[
                'w-full h-16 text-xl font-bold tracking-widest uppercase transition-all duration-300 flex items-center justify-center gap-3',
                rushStatus === 'success' ? 'bg-green-600 text-white cursor-default' :
                rushStatus === 'failed' ? 'bg-red-900/50 text-red-500 border border-red-500 hover:bg-red-900/80' :
                currentStock <= 0 ? 'bg-surface-light text-tertiary cursor-not-allowed' :
                rushLoading ? 'bg-surface border border-accent text-accent cursor-wait' :
                'bg-accent text-white hover:bg-accent-hover hover:shadow-[0_0_20px_rgba(255,59,48,0.4)]'
              ]"
            >
              <template v-if="rushStatus === 'success'">
                <ShieldCheck class="w-6 h-6" /> SECURED
              </template>
              <template v-else-if="rushLoading">
                <Loader2 class="w-6 h-6 animate-spin" /> PROCESSING
              </template>
              <template v-else-if="currentStock <= 0">
                SOLD OUT
              </template>
              <template v-else>
                <Zap class="w-6 h-6 fill-current" /> RUSH NOW
              </template>
            </button>
            
            <div v-if="rushMessage" :class="[
              'mt-4 text-center text-sm font-bold tracking-wide p-3 border',
              rushStatus === 'success' ? 'text-green-500 border-green-900 bg-green-900/10' : 
              rushStatus === 'failed' ? 'text-red-500 border-red-900 bg-red-900/10' :
              'text-accent border-accent/20 bg-accent/5'
            ]">
              > {{ rushMessage }}
            </div>
          </div>

          <div v-else class="w-full h-16 bg-surface-light border border-white/5 flex items-center justify-center text-secondary font-bold tracking-widest uppercase">
            Event Ended
          </div>
          
          <div class="flex justify-between text-xs text-tertiary uppercase tracking-wider">
            <span class="flex items-center gap-1"><Server class="w-3 h-3" /> Edge Network</span>
            <span class="flex items-center gap-1"><Activity class="w-3 h-3" /> Low Latency</span>
            <span class="flex items-center gap-1"><ShieldCheck class="w-3 h-3" /> Anti-Bot</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>