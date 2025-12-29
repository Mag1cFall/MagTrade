<template>
  <div
    class="min-h-screen bg-background text-primary selection:bg-accent selection:text-white overflow-x-hidden"
  >
    <div
      class="pt-20 pb-8 px-4 md:px-8 border-b border-white/5 bg-black/80 backdrop-blur-3xl sticky top-0 z-30"
    >
      <div class="container mx-auto flex flex-col lg:flex-row justify-between items-end gap-6">
        <div class="space-y-2">
          <div class="flex items-center gap-2">
            <div class="w-1 h-4 bg-accent"></div>
            <span class="text-[10px] font-mono tracking-[0.3em] text-accent uppercase"
              >秒杀商店</span
            >
          </div>
          <h1 class="text-3xl md:text-5xl font-black tracking-tighter text-white leading-none">
            MAG<span class="text-accent">STORE</span>
          </h1>
          <p class="text-secondary tracking-[0.2em] uppercase text-[10px] font-medium">
            限时抢购 // 先到先得
          </p>
        </div>

        <div class="flex flex-col items-end gap-4 w-full lg:w-auto">
          <!-- 搜索框 -->
          <div class="relative group w-full lg:w-64">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索商品..."
              class="w-full bg-transparent border-b border-white/10 py-2 px-0 text-white placeholder-gray-700 focus:border-accent focus:outline-none transition-all text-sm"
            />
            <Search
              class="absolute right-0 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-700 group-focus-within:text-accent transition-colors"
            />
          </div>
          <!-- 过滤和排序 -->
          <div class="flex items-center gap-3">
            <select
              v-model="statusFilter"
              class="bg-black border border-white/10 text-white text-xs px-3 py-2 rounded focus:border-accent outline-none"
            >
              <option value="all">全部状态</option>
              <option value="active">进行中</option>
              <option value="upcoming">即将开始</option>
            </select>
            <select
              v-model="sortBy"
              class="bg-black border border-white/10 text-white text-xs px-3 py-2 rounded focus:border-accent outline-none"
            >
              <option value="default">默认排序</option>
              <option value="price-asc">价格从低到高</option>
              <option value="price-desc">价格从高到低</option>
              <option value="discount">折扣力度</option>
              <option value="stock">剩余库存</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <div class="container mx-auto px-4 md:px-8 py-12">
      <div v-if="loading" class="flex flex-col items-center justify-center py-40 gap-6">
        <Loader2 class="w-12 h-12 animate-spin text-accent" />
        <span class="text-xs font-mono tracking-widest text-secondary animate-pulse"
          >加载中...</span
        >
      </div>

      <div v-else-if="filteredFlashSales.length === 0" class="py-40 text-center">
        <Ban class="w-24 h-24 mx-auto mb-4 opacity-10" />
        <h3 class="text-2xl font-black text-white mb-2">暂无秒杀活动</h3>
        <p class="text-secondary text-sm">请稍后再来</p>
      </div>

      <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <div
          v-for="item in filteredFlashSales"
          :key="item.id"
          class="group relative cursor-pointer"
          @click="goToDetail(item.id)"
        >
          <FlowBorderCard
            :active="item.status === 1"
            :color="item.status === 1 ? '#FF3B30' : '#333'"
            border-radius="0px"
            border-width="1px"
            duration="3"
            class="transition-all duration-700"
          >
            <div class="relative aspect-square bg-surface-light overflow-hidden flex flex-col">
              <div
                class="absolute inset-0 bg-gradient-to-br from-accent/5 to-transparent pointer-events-none z-10"
              ></div>

              <div
                class="w-full h-full p-4 flex items-center justify-center relative overflow-hidden"
              >
                <div
                  v-if="item.product?.image_url"
                  class="w-full h-full relative z-10 transition-transform duration-1000 ease-out group-hover:scale-110"
                >
                  <SmartImage
                    :src="item.product.image_url"
                    :alt="item.product.name"
                    class-name="w-full h-full object-contain"
                  />
                </div>
                <div
                  v-else
                  class="text-5xl opacity-20 transition-all duration-500 group-hover:opacity-40 relative z-10"
                >
                  {{ getProductIcon(item.product?.name) }}
                </div>
              </div>

              <div
                class="absolute inset-x-0 bottom-0 p-4 translate-y-full group-hover:translate-y-0 transition-transform duration-300 ease-out z-20 bg-black/70 backdrop-blur-xl border-t border-white/10"
              >
                <button
                  :disabled="rushing || item.status !== 1"
                  class="w-full h-10 bg-accent text-white text-xs font-bold uppercase hover:bg-white hover:text-black transition-all flex items-center justify-center gap-2 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
                  @click.stop="handleRush(item)"
                >
                  <Zap class="w-3 h-3 fill-current" :class="{ 'animate-pulse': rushing }" />
                  {{ rushing ? '抢购中...' : `立即抢购 ¥${item.flash_price}` }}
                </button>
              </div>

              <div class="absolute top-2 left-2 z-20 flex flex-col gap-1">
                <div
                  v-if="item.status === 1"
                  class="flex items-center gap-1 px-2 py-0.5 bg-accent text-white text-[8px] font-bold uppercase animate-pulse"
                >
                  <span class="w-1 h-1 bg-white rounded-full"></span>
                  限时
                </div>
                <div class="px-2 py-0.5 bg-white text-black text-[8px] font-bold uppercase">
                  -{{ calculateDiscount(item) }}%
                </div>
              </div>

              <div class="absolute top-2 right-2 z-20 font-mono text-[8px] text-secondary/50">
                #{{ item.id }} | {{ item.available_stock }}库存
              </div>
            </div>
          </FlowBorderCard>

          <div class="mt-3 px-1">
            <h3
              class="text-sm font-bold text-white truncate group-hover:text-accent transition-colors"
            >
              {{ item.product?.name || '未知商品' }}
            </h3>
            <p class="text-secondary text-[10px] line-clamp-1 mt-1 opacity-60">
              {{ item.product?.description || 'No description' }}
            </p>
            <div class="flex items-center justify-between mt-2">
              <div class="text-[10px] text-tertiary line-through">
                ¥{{ item.product?.original_price }}
              </div>
              <div class="text-lg font-black text-accent font-mono">¥{{ item.flash_price }}</div>
            </div>
          </div>
        </div>
      </div>

      <ConfirmModal
        :show="showRushResult"
        :title="rushResultTitle"
        :message="rushResultMessage"
        :type="rushSuccess ? 'success' : 'warning'"
        :confirm-text="rushSuccess ? '查看订单' : '知道了'"
        @confirm="onRushResultConfirm"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getFlashSales, rushFlashSale } from '@/api/flash-sale'
import type { FlashSale } from '@/types'
import FlowBorderCard from '@/components/FlowBorderCard.vue'
import SmartImage from '@/components/SmartImage.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import { Loader2, Search, Zap, Ban } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(true)
const flashSales = ref<FlashSale[]>([])
const searchQuery = ref('')
const statusFilter = ref<'all' | 'active' | 'upcoming'>('all')
const sortBy = ref<'default' | 'price-asc' | 'price-desc' | 'discount' | 'stock'>('default')
const rushing = ref(false)
const showRushResult = ref(false)
const rushSuccess = ref(false)
const rushResultTitle = ref('')
const rushResultMessage = ref('')

const getProductIcon = (name: string = '') => {
  const n = name.toLowerCase()
  if (n.includes('phone') || n.includes('手机')) return '📱'
  if (n.includes('laptop') || n.includes('电脑')) return '💻'
  if (n.includes('watch') || n.includes('手表')) return '⌚'
  if (n.includes('box') || n.includes('盒')) return '📦'
  if (n.includes('audio') || n.includes('音')) return '🎧'
  return '⚡'
}

const calculateDiscount = (item: FlashSale) => {
  if (!item.product?.original_price) return 0
  return Math.round((1 - item.flash_price / item.product.original_price) * 100)
}

const goToDetail = (id: number) => {
  router.push(`/flash-sales/${id}`)
}

const filteredFlashSales = computed(() => {
  let result = flashSales.value

  // 状态过滤
  if (statusFilter.value === 'active') {
    result = result.filter((item) => item.status === 1)
  } else if (statusFilter.value === 'upcoming') {
    result = result.filter((item) => item.status === 0)
  }

  // 搜索过滤
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(
      (item) =>
        item.product?.name?.toLowerCase().includes(q) ||
        item.product?.description?.toLowerCase().includes(q)
    )
  }

  // 排序
  if (sortBy.value === 'price-asc') {
    result = [...result].sort((a, b) => a.flash_price - b.flash_price)
  } else if (sortBy.value === 'price-desc') {
    result = [...result].sort((a, b) => b.flash_price - a.flash_price)
  } else if (sortBy.value === 'discount') {
    result = [...result].sort((a, b) => calculateDiscount(b) - calculateDiscount(a))
  } else if (sortBy.value === 'stock') {
    result = [...result].sort((a, b) => b.available_stock - a.available_stock)
  }

  return result
})

const fetchData = async () => {
  loading.value = true
  try {
    // 获取进行中和即将开始的活动
    const [activeRes, upcomingRes] = await Promise.all([
      getFlashSales(1, 100, 1),
      getFlashSales(1, 100, 0),
    ])
    const active = activeRes.code === 0 ? activeRes.data.flash_sales || [] : []
    const upcoming = upcomingRes.code === 0 ? upcomingRes.data.flash_sales || [] : []
    // 合并：进行中的排前面
    flashSales.value = [...active, ...upcoming]
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleRush = async (item: FlashSale) => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  if (rushing.value) return

  rushing.value = true
  try {
    const res = await rushFlashSale(item.id, 1)
    if (res.code === 0) {
      rushSuccess.value = true
      rushResultTitle.value = '抢购成功！'
      rushResultMessage.value = `恭喜抢到 ${item.product?.name}，请在订单页面完成支付`
    } else {
      rushSuccess.value = false
      rushResultTitle.value = '抢购失败'
      rushResultMessage.value = res.message || '库存不足或系统繁忙'
    }
  } catch (e: any) {
    rushSuccess.value = false
    rushResultTitle.value = '抢购失败'
    rushResultMessage.value = e.response?.data?.message || '网络错误'
  } finally {
    rushing.value = false
    showRushResult.value = true
  }
}

const onRushResultConfirm = () => {
  showRushResult.value = false
  if (rushSuccess.value) {
    router.push('/orders')
  }
}

onMounted(() => fetchData())
</script>

<style scoped>
.shop-card-entry {
  animation: fadeInUp 0.6s ease-out forwards;
  animation-delay: var(--delay);
  opacity: 0;
}

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
  from {
    opacity: 0;
    transform: translateY(20px);
  }
}
</style>
