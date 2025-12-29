<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getOrders, payOrder, cancelOrder } from '@/api/order'
import type { Order } from '@/types'
import SmartImage from '@/components/SmartImage.vue'
import BaseModal from '@/components/BaseModal.vue'
import PaymentModal from '@/components/PaymentModal.vue'
import { Loader2, Package, XCircle, CheckCircle, Clock, AlertCircle, Ban } from 'lucide-vue-next'

const loading = ref(true)
const orders = ref<Order[]>([])
const total = ref(0)
const actionLoading = ref<string | null>(null)

const showPayModal = ref(false)
const showCancelModal = ref(false)
const selectedOrderNo = ref('')
const selectedAmount = ref(0)
const selectedProductName = ref('')
const modalMessage = ref('')

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getOrders(1, 50)
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

const openPayModal = (order: Order) => {
  selectedOrderNo.value = order.order_no
  selectedAmount.value = order.amount
  selectedProductName.value = order.flash_sale?.product?.name || 'Unknown Product'
  showPayModal.value = true
}

const openCancelModal = (orderNo: string) => {
  selectedOrderNo.value = orderNo
  modalMessage.value = `Are you sure you want to cancel order ${orderNo}? This action cannot be undone.`
  showCancelModal.value = true
}

const confirmPay = async () => {
  showPayModal.value = false
  actionLoading.value = selectedOrderNo.value
  try {
    const res = await payOrder(selectedOrderNo.value)
    if (res.code === 0) {
      const order = orders.value.find((o: Order) => o.order_no === selectedOrderNo.value)
      if (order) {
        order.status = 1
        order.paid_at = new Date().toISOString()
      }
    }
  } catch (e) {
    console.error(e)
  } finally {
    actionLoading.value = null
  }
}

const confirmCancel = async () => {
  showCancelModal.value = false
  actionLoading.value = selectedOrderNo.value
  try {
    const res = await cancelOrder(selectedOrderNo.value)
    if (res.code === 0) {
      const order = orders.value.find((o: Order) => o.order_no === selectedOrderNo.value)
      if (order) order.status = 2
    }
  } catch (e) {
    console.error(e)
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
      <h1 class="text-4xl font-bold text-white tracking-tighter mb-2">
        ORDER <span class="text-accent">HISTORY</span>.
      </h1>
      <p class="text-secondary text-sm tracking-widest uppercase">Track your acquisitions</p>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <Loader2 class="w-8 h-8 text-accent animate-spin" />
    </div>

    <div
      v-else-if="orders.length === 0"
      class="py-20 text-center border border-dashed border-white/5 bg-surface/50 rounded-lg"
    >
      <Package class="w-12 h-12 text-secondary mx-auto mb-4 opacity-50" />
      <p class="text-secondary text-lg font-mono">NO ORDERS FOUND_</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="order in orders"
        :key="order.id"
        class="bg-surface border border-white/5 p-6 transition-all hover:border-white/20 group relative overflow-hidden"
      >
        <div
          class="absolute top-0 left-0 w-1 h-full bg-white/5 group-hover:bg-accent transition-colors duration-300"
        ></div>

        <div class="flex flex-col md:flex-row md:items-center justify-between gap-6 pl-4">
          <div class="flex-grow flex gap-6">
            <div
              class="w-24 h-24 flex-shrink-0 bg-surface-light border border-white/5 overflow-hidden"
            >
              <SmartImage
                :src="order.flash_sale?.product?.image_url"
                :alt="order.flash_sale?.product?.name || 'Product'"
                class-name="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
              />
            </div>

            <div>
              <div class="flex items-center gap-4 mb-2">
                <span class="text-xs text-tertiary font-mono tracking-widest">{{
                  order.order_no
                }}</span>
                <span
                  :class="[
                    'text-xs font-bold px-2 py-0.5 border flex items-center gap-1 uppercase tracking-wider',
                    getStatusConfig(order.status).color,
                    `border-current`,
                  ]"
                >
                  <component :is="getStatusConfig(order.status).icon" class="w-3 h-3" />
                  {{ getStatusConfig(order.status).label }}
                </span>
              </div>
              <h3 class="text-xl font-bold text-white mb-2">
                {{ order.flash_sale?.product?.name || 'Unknown Product' }}
              </h3>
              <div class="text-sm text-secondary flex gap-4">
                <span
                  >Quantity: <span class="text-white font-mono">{{ order.quantity }}</span></span
                >
                <span class="text-tertiary">|</span>
                <span
                  >Time:
                  <span class="text-white font-mono">{{
                    new Date(order.created_at).toLocaleString()
                  }}</span></span
                >
              </div>
            </div>
          </div>

          <div class="flex items-center gap-8 md:min-w-[300px] justify-between md:justify-end">
            <div class="text-right">
              <div class="text-xs text-secondary mb-1 uppercase tracking-wider">Total Amount</div>
              <div class="text-2xl font-mono font-bold text-white">
                ¥{{ order.amount.toFixed(2) }}
              </div>
            </div>

            <div v-if="order.status === 0" class="flex gap-3">
              <button
                :disabled="!!actionLoading"
                class="p-3 border border-white/10 text-secondary hover:text-red-500 hover:border-red-500 hover:bg-red-500/10 transition-colors disabled:opacity-50"
                title="Cancel Order"
                @click="openCancelModal(order.order_no)"
              >
                <Ban class="w-5 h-5" />
              </button>
              <button
                :disabled="!!actionLoading"
                class="px-6 py-3 bg-white text-black text-sm font-bold hover:bg-accent hover:text-white transition-all duration-300 disabled:opacity-50 flex items-center gap-2 uppercase tracking-widest"
                @click="openPayModal(order)"
              >
                <Loader2 v-if="actionLoading === order.order_no" class="w-4 h-4 animate-spin" />
                <span v-else>Pay</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <PaymentModal
      :show="showPayModal"
      :order-no="selectedOrderNo"
      :amount="selectedAmount"
      :product-name="selectedProductName"
      @confirm="confirmPay"
      @cancel="showPayModal = false"
    />

    <BaseModal
      :show="showCancelModal"
      title="Cancel Order"
      :message="modalMessage"
      type="danger"
      confirm-text="Yes, Cancel"
      @confirm="confirmCancel"
      @cancel="showCancelModal = false"
    />
  </div>
</template>
