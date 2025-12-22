<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getOrders, payOrder, cancelOrder } from '@/api/order'
import type { Order } from '@/types'
import SmartImage from '@/components/SmartImage.vue'
import { Loader2, Package, XCircle, CheckCircle, Clock, AlertCircle } from 'lucide-vue-next'

const loading = ref(true)
const orders = ref<Order[]>([])
const total = ref(0)
const actionLoading = ref<string | null>(null)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getOrders(1, 50) // Simple pagination for now
    if (res.code === 0) {
      orders.value = res.data.orders
      total.value = res.data.total
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handlePay = async (orderNo: string) => {
  if (!confirm('确认支付该订单吗？')) return
  actionLoading.value = orderNo
  try {
    const res = await payOrder(orderNo)
    if (res.code === 0) {
      // Update local state
      const order = orders.value.find(o => o.order_no === orderNo)
      if (order) {
        order.status = 1 // Paid
        order.paid_at = new Date().toISOString()
      }
    }
  } catch (e) {
    alert('支付失败')
  } finally {
    actionLoading.value = null
  }
}

const handleCancel = async (orderNo: string) => {
  if (!confirm('确认取消该订单吗？此操作不可恢复。')) return
  actionLoading.value = orderNo
  try {
    const res = await cancelOrder(orderNo)
    if (res.code === 0) {
      const order = orders.value.find(o => o.order_no === orderNo)
      if (order) order.status = 2 // Cancelled
    }
  } catch (e) {
    alert('取消失败')
  } finally {
    actionLoading.value = null
  }
}

const statusMap: Record<number, { label: string; color: string; icon: any }> = {
  0: { label: 'PENDING', color: 'text-yellow-500', icon: Clock },
  1: { label: 'PAID', color: 'text-green-500', icon: CheckCircle },
  2: { label: 'CANCELLED', color: 'text-red-500', icon: XCircle },
  3: { label: 'REFUNDED', color: 'text-secondary', icon: AlertCircle },
}

const getStatusConfig = (status: number) => {
  return statusMap[status] || { label: 'UNKNOWN', color: 'text-secondary', icon: AlertCircle }
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="min-h-screen pb-20">
    <div class="py-12 border-b border-white/5 mb-8">
      <h1 class="text-4xl font-bold text-white tracking-tighter mb-2">ORDER <span class="text-accent">HISTORY</span>.</h1>
      <p class="text-secondary text-sm tracking-widest uppercase">Track your acquisitions</p>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <Loader2 class="w-8 h-8 text-accent animate-spin" />
    </div>

    <div v-else-if="orders.length === 0" class="py-20 text-center border border-white/5 bg-surface/50">
      <Package class="w-12 h-12 text-secondary mx-auto mb-4 opacity-50" />
      <p class="text-secondary text-lg font-mono">NO ORDERS FOUND_</p>
    </div>

    <div v-else class="space-y-4">
      <div 
        v-for="order in orders" 
        :key="order.id"
        class="bg-surface border border-white/5 p-6 transition-all hover:border-white/10"
      >
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <!-- Order Info -->
          <div class="flex-grow flex gap-6">
            <!-- Smart Image Thumbnail -->
            <div class="w-20 h-20 flex-shrink-0 bg-surface-light border border-white/5 rounded-sm overflow-hidden">
              <SmartImage
                :src="order.flash_sale?.product?.image_url"
                :alt="order.flash_sale?.product?.name || 'Product'"
                class-name="w-full h-full object-cover"
              />
            </div>

            <div>
              <div class="flex items-center gap-4 mb-2">
                <span class="text-xs text-secondary font-mono tracking-widest">{{ order.order_no }}</span>
                <span :class="['text-xs font-bold px-2 py-0.5 border flex items-center gap-1', getStatusConfig(order.status).color, `border-current`]">
                  <component :is="getStatusConfig(order.status).icon" class="w-3 h-3" />
                  {{ getStatusConfig(order.status).label }}
                </span>
              </div>
              <h3 class="text-lg font-bold text-white mb-1">{{ order.flash_sale?.product?.name || 'Unknown Product' }}</h3>
              <div class="text-sm text-secondary">
                Quantity: <span class="text-white">{{ order.quantity }}</span> |
                Time: <span class="text-white">{{ new Date(order.created_at).toLocaleString() }}</span>
              </div>
            </div>
          </div>

          <!-- Price & Actions -->
          <div class="flex items-center gap-8 md:min-w-[300px] justify-between md:justify-end">
            <div class="text-right">
              <div class="text-xs text-secondary mb-1">TOTAL AMOUNT</div>
              <div class="text-xl font-mono font-bold text-white">¥{{ order.amount.toFixed(2) }}</div>
            </div>

            <div class="flex gap-3" v-if="order.status === 0">
              <button 
                @click="handleCancel(order.order_no)"
                :disabled="actionLoading === order.order_no"
                class="px-4 py-2 border border-white/10 text-secondary text-sm font-bold hover:bg-white/5 transition-colors disabled:opacity-50"
              >
                CANCEL
              </button>
              <button 
                @click="handlePay(order.order_no)"
                :disabled="actionLoading === order.order_no"
                class="px-4 py-2 bg-white text-black text-sm font-bold hover:bg-gray-200 transition-colors disabled:opacity-50 flex items-center gap-2"
              >
                <Loader2 v-if="actionLoading === order.order_no" class="w-4 h-4 animate-spin" />
                <span v-else>PAY NOW</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>