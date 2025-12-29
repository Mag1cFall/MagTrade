<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getFlashSaleDetail, rushFlashSale, getFlashSaleStock } from '@/api/flash-sale'
import { getRecommendation } from '@/api/ai'
import type { FlashSaleDetail } from '@/types'
import CountDown from '@/components/CountDown.vue'
import SmartImage from '@/components/SmartImage.vue'
import BaseModal from '@/components/BaseModal.vue'
import { useAuthStore } from '@/stores/auth'
import {
  Loader2,
  ShieldCheck,
  Zap,
  Server,
  Activity,
  ArrowLeft,
  Bot,
  Terminal,
  TrendingUp,
  Clock,
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const loading = ref(true)
const detail = ref<FlashSaleDetail | null>(null)
const currentStock = ref(0)
const rushLoading = ref(false)
const ws = ref<WebSocket | null>(null)
const rushStatus = ref<'idle' | 'processing' | 'success' | 'failed'>('idle')

// Modal State
const showResultModal = ref(false)
const resultTitle = ref('')
const resultMessage = ref('')
const resultType = ref<'info' | 'danger'>('info')

const aiAnalysis = ref({
  difficulty: 0,
  winRate: '--%',
  advice: 'CONNECTING TO NEURAL NETWORK...',
  risk: 'CALCULATING',
})

const flashSaleId = parseInt(route.params.id as string)

const fetchData = async () => {
  try {
    const res = await getFlashSaleDetail(flashSaleId)
    if (res.code === 0) {
      detail.value = res.data
      currentStock.value = res.data.current_stock
    }

    // 无论是否登录都尝试获取 AI 分析
    getRecommendation(flashSaleId)
      .then((res) => {
        if (res.code === 0 && res.data?.analysis) {
          const { analysis } = res.data
          aiAnalysis.value = {
            difficulty: analysis.difficulty_score || 0,
            winRate: ((analysis.success_probability || 0) * 100).toFixed(1) + '%',
            advice: analysis.timing_advice || 'Analysis complete',
            risk: analysis.difficulty_reason || 'NORMAL',
          }
        } else {
          aiAnalysis.value.advice = res.message || 'AI analysis unavailable'
        }
      })
      .catch((e) => {
        console.error('AI Recommendation Error:', e)
        aiAnalysis.value.advice = 'AI SERVICE TEMPORARILY UNAVAILABLE'
      })
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

let stockInterval: number
const startStockPolling = () => {
  stockInterval = window.setInterval(async () => {
    try {
      const res = await getFlashSaleStock(flashSaleId)
      if (res.code === 0) {
        currentStock.value = res.data.stock
      }
    } catch {
      // Silent fail for stock polling
    }
  }, 3000)
}

const connectWS = () => {
  if (!authStore.token) return

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
          showModal('ACQUISITION SUCCESSFUL', msg.data.message, 'info')
          clearInterval(stockInterval)
        } else {
          rushStatus.value = 'failed'
          showModal('ACQUISITION FAILED', msg.data.message, 'danger')
        }
      }
    } catch (e) {
      console.error('WS parse error', e)
    }
  }
}

const showModal = (title: string, message: string, type: 'info' | 'danger') => {
  resultTitle.value = title
  resultMessage.value = message
  resultType.value = type
  showResultModal.value = true
}

const handleModalConfirm = () => {
  showResultModal.value = false
  if (rushStatus.value === 'success') {
    router.push('/orders')
  }
}

const handleRush = async () => {
  if (!authStore.isAuthenticated) {
    router.push(`/login?redirect=${route.fullPath}`)
    return
  }

  rushLoading.value = true
  rushStatus.value = 'processing'

  try {
    const res = await rushFlashSale(flashSaleId)
    if (res.code === 0) {
      // Waiting for WS
    } else {
      rushLoading.value = false
      rushStatus.value = 'failed'
      showModal('REQUEST REJECTED', res.message, 'danger')
    }
  } catch (e: any) {
    rushLoading.value = false
    rushStatus.value = 'failed'
    // 解析错误消息
    const status = e.response?.status
    const serverMsg = e.response?.data?.message || e.message || 'Network Error'

    if (status === 409) {
      showModal('已参与过活动', '每人限购一次，您已经抢购过此商品', 'danger')
    } else if (status === 400) {
      showModal('抢购失败', serverMsg, 'danger')
    } else {
      showModal('系统错误', serverMsg, 'danger')
    }
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

const discount = computed(() => {
  if (!detail.value?.flash_sale.product) return 0
  const p = detail.value.flash_sale.product
  return Math.round(
    ((p.original_price - detail.value.flash_sale.flash_price) / p.original_price) * 100
  )
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

const handleCountdownFinish = () => {
  // 倒计时结束，刷新数据获取最新状态
  fetchData()
}
</script>

<template>
  <div v-if="loading" class="min-h-screen flex items-center justify-center">
    <Loader2 class="w-10 h-10 text-accent animate-spin" />
  </div>

  <div v-else-if="!detail" class="min-h-screen flex flex-col items-center justify-center gap-4">
    <div class="text-2xl font-bold">EVENT NOT FOUND</div>
    <button
      class="text-secondary hover:text-white flex items-center gap-2"
      @click="router.push('/')"
    >
      <ArrowLeft class="w-4 h-4" /> Return to Base
    </button>
  </div>

  <div v-else class="min-h-screen pb-32 lg:pb-20">
    <!-- Breadcrumb -->
    <div class="container mx-auto px-4 md:px-8 py-4 mt-20 relative z-20">
      <button
        class="text-sm text-secondary hover:text-white flex items-center gap-2 uppercase tracking-widest font-bold cursor-pointer"
        @click="router.push('/')"
      >
        <ArrowLeft class="w-4 h-4" /> Back to Lobby
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 mt-8">
      <!-- Left: Visuals (7 Cols) -->
      <div class="lg:col-span-7 space-y-8">
        <div
          class="aspect-square lg:aspect-[16/10] bg-surface-light border border-white/5 relative overflow-hidden group"
        >
          <SmartImage
            :src="detail.flash_sale.product?.image_url"
            :alt="detail.flash_sale.product?.name || 'Product'"
            class-name="w-full h-full object-cover transition-transform duration-1000 group-hover:scale-105"
          />

          <div class="absolute top-0 right-0 bg-white text-black text-xl font-bold px-4 py-2">
            -{{ discount }}%
          </div>
        </div>

        <!-- AI Analysis Panel -->
        <div class="bg-surface border border-white/10 p-6 relative overflow-hidden hidden lg:block">
          <div class="absolute top-0 right-0 p-4 opacity-5">
            <Bot class="w-32 h-32" />
          </div>
          <h3
            class="text-sm font-bold text-accent uppercase tracking-widest mb-6 flex items-center gap-2"
          >
            <Terminal class="w-4 h-4" /> AI Tactical Analysis
          </h3>
          <div class="grid grid-cols-3 gap-8 mb-6">
            <div>
              <div class="text-xs text-secondary mb-2 tracking-widest">DIFFICULTY</div>
              <div class="flex items-end gap-2">
                <div class="text-3xl font-mono font-bold text-white">
                  {{ aiAnalysis.difficulty }}
                </div>
                <div class="text-sm text-tertiary mb-1">/ 10</div>
              </div>
              <div class="h-1 w-full bg-surface-light mt-2 rounded-full overflow-hidden">
                <div
                  class="h-full bg-accent"
                  :style="{ width: `${aiAnalysis.difficulty * 10}%` }"
                ></div>
              </div>
            </div>
            <div>
              <div class="text-xs text-secondary mb-2 tracking-widest">WIN PROBABILITY</div>
              <div class="text-3xl font-mono font-bold text-white">{{ aiAnalysis.winRate }}</div>
            </div>
            <div>
              <div class="text-xs text-secondary mb-2 tracking-widest">RISK FACTOR</div>
              <div class="text-xl font-mono font-bold text-white uppercase">
                {{ aiAnalysis.risk }}
              </div>
            </div>
          </div>
          <div class="flex gap-3 items-start border-t border-white/10 pt-4 mt-4">
            <TrendingUp class="w-5 h-5 text-accent mt-0.5" />
            <p class="text-sm text-secondary leading-relaxed">
              <span class="text-white font-bold">ADVISOR:</span> {{ aiAnalysis.advice }}
            </p>
          </div>
        </div>
      </div>

      <!-- Right: Info & Action (5 Cols) -->
      <div class="lg:col-span-5 relative">
        <div class="lg:sticky lg:top-32 space-y-8">
          <div>
            <h1 class="text-3xl md:text-5xl font-bold text-white mb-4 leading-tight">
              {{ detail.flash_sale.product?.name }}
            </h1>
            <p class="text-secondary text-lg leading-relaxed line-clamp-3">
              {{ detail.flash_sale.product?.description }}
            </p>
          </div>

          <!-- Price Block -->
          <div class="p-6 bg-surface-light/50 border border-white/5 backdrop-blur-sm">
            <div class="flex items-baseline gap-4 mb-2">
              <span class="text-5xl font-mono font-bold text-white"
                >¥{{ detail.flash_sale.flash_price }}</span
              >
              <span class="text-xl text-tertiary line-through font-mono"
                >¥{{ detail.flash_sale.product?.original_price }}</span
              >
            </div>
            <div
              class="flex items-center gap-2 text-xs text-accent uppercase tracking-widest font-bold"
            >
              <Activity class="w-3 h-3" /> Real-time Pricing
            </div>
          </div>

          <!-- Status & Countdown -->
          <div class="space-y-4">
            <div class="flex justify-between items-end border-b border-white/10 pb-2">
              <span class="text-xs text-secondary uppercase tracking-widest">Inventory Status</span>
              <span class="font-mono text-xl text-white"
                >{{ currentStock }}
                <span class="text-tertiary text-sm"
                  >/ {{ detail.flash_sale.total_stock }}</span
                ></span
              >
            </div>

            <div
              v-if="!isStarted"
              class="p-4 bg-surface border border-white/10 flex items-center justify-between"
            >
              <span class="text-sm font-bold text-white uppercase tracking-widest"
                >Dropping In</span
              >
              <CountDown
                :target-time="detail.flash_sale.start_time"
                size="md"
                @finish="handleCountdownFinish"
              />
            </div>
          </div>

          <!-- Desktop Action Button -->
          <div class="hidden lg:block pt-4">
            <button
              v-if="!isStarted"
              disabled
              class="w-full h-16 bg-surface-light border border-white/10 text-secondary font-bold tracking-widest uppercase cursor-not-allowed"
            >
              Wait for Drop
            </button>
            <button
              v-else-if="!isEnded"
              :disabled="rushLoading || currentStock <= 0 || rushStatus === 'success'"
              :class="[
                'w-full h-16 text-xl font-bold tracking-widest uppercase transition-all duration-300 flex items-center justify-center gap-3 hover:shadow-[0_0_30px_rgba(227,53,53,0.3)]',
                rushStatus === 'success'
                  ? 'bg-green-600 text-white'
                  : currentStock <= 0
                    ? 'bg-surface-light text-tertiary'
                    : 'bg-accent text-white hover:bg-accent-hover',
              ]"
              @click="handleRush"
            >
              <template v-if="rushLoading"
                ><Loader2 class="w-6 h-6 animate-spin" /> PROCESSING</template
              >
              <template v-else-if="currentStock <= 0">SOLD OUT</template>
              <template v-else><Zap class="w-6 h-6 fill-current" /> RUSH NOW</template>
            </button>
            <button
              v-else
              disabled
              class="w-full h-16 bg-surface-light border border-white/10 text-tertiary font-bold tracking-widest uppercase cursor-not-allowed"
            >
              Event Ended
            </button>
          </div>

          <!-- Features -->
          <div class="grid grid-cols-2 gap-4 pt-8 border-t border-white/5">
            <div class="flex items-center gap-3 text-secondary text-sm">
              <ShieldCheck class="w-5 h-5 text-accent" />
              <span>Anti-Bot Verified</span>
            </div>
            <div class="flex items-center gap-3 text-secondary text-sm">
              <Server class="w-5 h-5 text-accent" />
              <span>Edge Network</span>
            </div>
            <div class="flex items-center gap-3 text-secondary text-sm">
              <Clock class="w-5 h-5 text-accent" />
              <span>Instant Settlement</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile Fixed Bottom Bar -->
    <div
      class="lg:hidden fixed bottom-0 left-0 w-full z-50 bg-background/90 backdrop-blur-xl border-t border-white/10 p-4 pb-8"
    >
      <div class="flex gap-4 items-center">
        <div class="flex-1">
          <div class="text-xs text-secondary uppercase tracking-wider mb-1">Flash Price</div>
          <div class="text-2xl font-mono font-bold text-white">
            ¥{{ detail.flash_sale.flash_price }}
          </div>
        </div>
        <button
          v-if="isStarted && !isEnded"
          :disabled="rushLoading || currentStock <= 0"
          class="flex-1 h-12 bg-accent text-white font-bold tracking-widest uppercase flex items-center justify-center gap-2 shadow-lg shadow-accent/20"
          @click="handleRush"
        >
          <template v-if="rushLoading"><Loader2 class="w-5 h-5 animate-spin" /></template>
          <template v-else-if="currentStock <= 0">Sold Out</template>
          <template v-else>Rush Now</template>
        </button>
        <button
          v-else
          disabled
          class="flex-1 h-12 bg-surface-light text-secondary font-bold tracking-widest uppercase"
        >
          {{ !isStarted ? 'Wait' : 'Ended' }}
        </button>
      </div>
    </div>

    <BaseModal
      :show="showResultModal"
      :title="resultTitle"
      :message="resultMessage"
      :type="resultType"
      confirm-text="View Orders"
      cancel-text="Close"
      @confirm="handleModalConfirm"
      @cancel="showResultModal = false"
    />
  </div>
</template>
