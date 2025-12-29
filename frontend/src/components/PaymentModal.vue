<template>
  <teleport to="body">
    <transition name="modal">
      <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/80 backdrop-blur-sm" @click="$emit('cancel')"></div>

        <div
          class="relative w-full max-w-md bg-[#1a1a1a] border border-white/10 shadow-2xl overflow-hidden animate-in"
        >
          <!-- 顶部渐变条 -->
          <div class="h-1 bg-gradient-to-r from-accent via-orange-500 to-accent"></div>

          <div class="p-6">
            <!-- 订单详情 -->
            <div class="mb-6">
              <h2 class="text-xl font-bold text-white mb-4">Order Details</h2>
              <div class="bg-black/40 border border-white/5 rounded-lg p-4 space-y-3">
                <div class="flex justify-between items-center">
                  <span class="text-secondary text-sm">{{ productName }}</span>
                  <span class="text-white font-mono">¥{{ amount.toFixed(2) }}</span>
                </div>
                <div class="border-t border-white/5 pt-3 flex justify-between items-center">
                  <span class="text-secondary text-sm">Subtotal</span>
                  <span class="text-white font-mono">¥{{ amount.toFixed(2) }}</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="text-cyan-400 font-semibold">Total due today</span>
                  <span class="text-cyan-400 font-mono font-bold text-lg"
                    >¥{{ amount.toFixed(2) }}</span
                  >
                </div>
              </div>
            </div>

            <!-- 支付方式 -->
            <div class="mb-6">
              <h3 class="text-white font-semibold mb-3">Payment method</h3>
              <div class="grid grid-cols-5 gap-2 mb-4">
                <!-- Card -->
                <button
                  :class="[
                    'py-2.5 px-2 border rounded-lg flex flex-col items-center gap-1 transition-all',
                    paymentMethod === 'card'
                      ? 'border-accent bg-accent/10 text-accent'
                      : 'border-white/10 text-secondary hover:border-white/30',
                  ]"
                  @click="paymentMethod = 'card'"
                >
                  <CreditCard class="w-5 h-5" />
                  <span class="text-[10px] font-medium">Card</span>
                </button>
                <!-- Alipay -->
                <button
                  :class="[
                    'py-2.5 px-2 border rounded-lg flex flex-col items-center gap-1 transition-all',
                    paymentMethod === 'alipay'
                      ? 'border-blue-500 bg-blue-500/10 text-blue-500'
                      : 'border-white/10 text-secondary hover:border-white/30',
                  ]"
                  @click="paymentMethod = 'alipay'"
                >
                  <Wallet class="w-5 h-5" />
                  <span class="text-[10px] font-medium">Alipay</span>
                </button>
                <!-- WeChat -->
                <button
                  :class="[
                    'py-2.5 px-2 border rounded-lg flex flex-col items-center gap-1 transition-all',
                    paymentMethod === 'wechat'
                      ? 'border-green-500 bg-green-500/10 text-green-500'
                      : 'border-white/10 text-secondary hover:border-white/30',
                  ]"
                  @click="paymentMethod = 'wechat'"
                >
                  <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                    <path
                      d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.472c-.078.295-.003.617.213.828a.72.72 0 0 0 .88.121l1.994-.984a.587.587 0 0 1 .504-.023c.87.351 1.82.555 2.809.596-.156-.484-.239-.996-.239-1.524 0-3.61 3.494-6.546 7.79-6.546.291 0 .578.013.86.039C16.704 5.103 13.032 2.188 8.691 2.188zm-2.51 4.66a1.073 1.073 0 1 1 0-2.145 1.073 1.073 0 0 1 0 2.146zm5.082 0a1.073 1.073 0 1 1 0-2.145 1.073 1.073 0 0 1 0 2.146z"
                    />
                    <path
                      d="M23.99 16.234c0-3.252-3.262-5.891-7.283-5.891-4.02 0-7.282 2.64-7.282 5.891 0 3.252 3.262 5.891 7.282 5.891.88 0 1.721-.13 2.498-.368a.49.49 0 0 1 .421.019l1.664.82a.6.6 0 0 0 .734-.101.59.59 0 0 0 .177-.69l-.325-1.227a.49.49 0 0 1 .177-.553c1.527-1.12 2.537-2.828 2.537-4.791zm-9.633-.902a.894.894 0 1 1 0-1.789.894.894 0 0 1 0 1.79zm4.7 0a.894.894 0 1 1 0-1.789.894.894 0 0 1 0 1.79z"
                    />
                  </svg>
                  <span class="text-[10px] font-medium">WeChat</span>
                </button>
                <!-- Apple Pay -->
                <button
                  :class="[
                    'py-2.5 px-2 border rounded-lg flex flex-col items-center gap-1 transition-all',
                    paymentMethod === 'apple'
                      ? 'border-white bg-white/10 text-white'
                      : 'border-white/10 text-secondary hover:border-white/30',
                  ]"
                  @click="paymentMethod = 'apple'"
                >
                  <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                    <path
                      d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"
                    />
                  </svg>
                  <span class="text-[10px] font-medium">Apple</span>
                </button>
                <!-- Google Pay -->
                <button
                  :class="[
                    'py-2.5 px-2 border rounded-lg flex flex-col items-center gap-1 transition-all',
                    paymentMethod === 'google'
                      ? 'border-blue-400 bg-blue-400/10 text-blue-400'
                      : 'border-white/10 text-secondary hover:border-white/30',
                  ]"
                  @click="paymentMethod = 'google'"
                >
                  <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                    <path
                      d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                      fill="#4285F4"
                    />
                    <path
                      d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                      fill="#34A853"
                    />
                    <path
                      d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                      fill="#FBBC05"
                    />
                    <path
                      d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                      fill="#EA4335"
                    />
                  </svg>
                  <span class="text-[10px] font-medium">Google</span>
                </button>
              </div>

              <!-- 信用卡表单 -->
              <div v-if="paymentMethod === 'card'" class="space-y-4">
                <div>
                  <label class="text-secondary text-xs uppercase tracking-wider block mb-1.5"
                    >Full name</label
                  >
                  <input
                    v-model="cardForm.name"
                    type="text"
                    placeholder="Cardholder Name"
                    class="w-full bg-black border border-white/10 rounded-lg px-4 py-3 text-white placeholder-gray-600 focus:border-accent outline-none transition-colors"
                  />
                </div>

                <div>
                  <label class="text-secondary text-xs uppercase tracking-wider block mb-1.5"
                    >Card number</label
                  >
                  <div class="relative">
                    <input
                      v-model="cardForm.number"
                      type="text"
                      placeholder="1234 1234 1234 1234"
                      maxlength="19"
                      class="w-full bg-black border border-white/10 rounded-lg px-4 py-3 pr-32 text-white placeholder-gray-600 focus:border-accent outline-none transition-colors font-mono"
                      @input="formatCardNumber"
                    />
                    <div class="absolute right-3 top-1/2 -translate-y-1/2 flex gap-1 items-center">
                      <div
                        class="w-8 h-5 bg-[#1A1F71] rounded flex items-center justify-center text-[6px] text-white font-bold"
                      >
                        VISA
                      </div>
                      <div
                        class="w-8 h-5 bg-gradient-to-br from-[#EB001B] to-[#F79E1B] rounded flex items-center justify-center"
                      >
                        <div class="w-2 h-2 rounded-full bg-[#EB001B] -mr-1"></div>
                        <div class="w-2 h-2 rounded-full bg-[#F79E1B]"></div>
                      </div>
                      <div
                        class="w-8 h-5 bg-white rounded flex items-center justify-center text-[5px] text-blue-600 font-bold"
                      >
                        JCB
                      </div>
                      <div
                        class="w-8 h-5 bg-[#006FCF] rounded flex items-center justify-center text-[5px] text-white font-bold"
                      >
                        AMEX
                      </div>
                    </div>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="text-secondary text-xs uppercase tracking-wider block mb-1.5"
                      >Expiration</label
                    >
                    <input
                      v-model="cardForm.expiry"
                      type="text"
                      placeholder="MM / YY"
                      maxlength="7"
                      class="w-full bg-black border border-white/10 rounded-lg px-4 py-3 text-white placeholder-gray-600 focus:border-accent outline-none transition-colors font-mono"
                      @input="formatExpiry"
                    />
                  </div>
                  <div>
                    <label class="text-secondary text-xs uppercase tracking-wider block mb-1.5"
                      >CVC</label
                    >
                    <div class="relative">
                      <input
                        v-model="cardForm.cvc"
                        type="text"
                        placeholder="CVC"
                        maxlength="4"
                        class="w-full bg-black border border-white/10 rounded-lg px-4 py-3 text-white placeholder-gray-600 focus:border-accent outline-none transition-colors font-mono"
                      />
                      <CreditCard
                        class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-600"
                      />
                    </div>
                  </div>
                </div>
              </div>

              <!-- 支付宝/微信/Apple/Google 提示 -->
              <div
                v-else-if="paymentMethod === 'alipay'"
                class="bg-blue-500/10 border border-blue-500/30 rounded-lg p-4 text-center"
              >
                <Wallet class="w-8 h-8 text-blue-500 mx-auto mb-2" />
                <p class="text-blue-400 text-sm">You will be redirected to Alipay</p>
              </div>
              <div
                v-else-if="paymentMethod === 'wechat'"
                class="bg-green-500/10 border border-green-500/30 rounded-lg p-4 text-center"
              >
                <svg
                  class="w-8 h-8 text-green-500 mx-auto mb-2"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.472c-.078.295-.003.617.213.828a.72.72 0 0 0 .88.121l1.994-.984a.587.587 0 0 1 .504-.023c.87.351 1.82.555 2.809.596-.156-.484-.239-.996-.239-1.524 0-3.61 3.494-6.546 7.79-6.546.291 0 .578.013.86.039C16.704 5.103 13.032 2.188 8.691 2.188zm-2.51 4.66a1.073 1.073 0 1 1 0-2.145 1.073 1.073 0 0 1 0 2.146zm5.082 0a1.073 1.073 0 1 1 0-2.145 1.073 1.073 0 0 1 0 2.146z"
                  />
                  <path
                    d="M23.99 16.234c0-3.252-3.262-5.891-7.283-5.891-4.02 0-7.282 2.64-7.282 5.891 0 3.252 3.262 5.891 7.282 5.891.88 0 1.721-.13 2.498-.368a.49.49 0 0 1 .421.019l1.664.82a.6.6 0 0 0 .734-.101.59.59 0 0 0 .177-.69l-.325-1.227a.49.49 0 0 1 .177-.553c1.527-1.12 2.537-2.828 2.537-4.791zm-9.633-.902a.894.894 0 1 1 0-1.789.894.894 0 0 1 0 1.79zm4.7 0a.894.894 0 1 1 0-1.789.894.894 0 0 1 0 1.79z"
                  />
                </svg>
                <p class="text-green-400 text-sm">Scan QR code with WeChat</p>
              </div>
              <div
                v-else-if="paymentMethod === 'apple'"
                class="bg-white/5 border border-white/20 rounded-lg p-4 text-center"
              >
                <svg
                  class="w-8 h-8 text-white mx-auto mb-2"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"
                  />
                </svg>
                <p class="text-white/70 text-sm">Continue with Apple Pay</p>
              </div>
              <div
                v-else-if="paymentMethod === 'google'"
                class="bg-blue-400/10 border border-blue-400/30 rounded-lg p-4 text-center"
              >
                <svg class="w-8 h-8 mx-auto mb-2" viewBox="0 0 24 24">
                  <path
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                    fill="#4285F4"
                  />
                  <path
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                    fill="#34A853"
                  />
                  <path
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                    fill="#FBBC05"
                  />
                  <path
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                    fill="#EA4335"
                  />
                </svg>
                <p class="text-blue-300 text-sm">Continue with Google Pay</p>
              </div>
            </div>

            <!-- 协议 -->
            <div class="flex items-start gap-3 mb-6">
              <input
                v-model="agreed"
                type="checkbox"
                class="mt-1 w-4 h-4 rounded border-white/20 bg-black text-accent focus:ring-accent"
              />
              <p class="text-xs text-secondary leading-relaxed">
                You agree that MagTrade will charge your payment method in the amount above.
                <router-link to="/terms" class="text-accent hover:underline">Terms</router-link>
                apply.
              </p>
            </div>

            <!-- 按钮 -->
            <div class="flex gap-3">
              <button
                class="flex-1 py-3 border border-white/10 text-secondary hover:text-white hover:border-white/30 transition-colors rounded-lg font-medium"
                @click="$emit('cancel')"
              >
                Cancel
              </button>
              <button
                :disabled="!isFormValid || processing"
                class="flex-1 py-3 bg-accent text-white font-bold rounded-lg hover:bg-accent/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                @click="handlePay"
              >
                <Loader2 v-if="processing" class="w-4 h-4 animate-spin" />
                <span>{{ processing ? 'Processing...' : `Pay ¥${amount.toFixed(2)}` }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { CreditCard, Wallet, Loader2 } from 'lucide-vue-next'

defineProps<{
  show: boolean
  orderNo: string
  amount: number
  productName: string
}>()

const emit = defineEmits<{
  (e: 'cancel'): void
  (e: 'confirm'): void
}>()

const paymentMethod = ref<'card' | 'alipay' | 'wechat' | 'apple' | 'google'>('card')
const agreed = ref(false)
const processing = ref(false)

const cardForm = ref({
  name: '',
  number: '',
  expiry: '',
  cvc: '',
})

const formatCardNumber = (e: Event) => {
  const input = e.target as HTMLInputElement
  let value = input.value.replace(/\s/g, '').replace(/\D/g, '')
  value = value.match(/.{1,4}/g)?.join(' ') || value
  cardForm.value.number = value
}

const formatExpiry = (e: Event) => {
  const input = e.target as HTMLInputElement
  let value = input.value.replace(/\D/g, '')
  if (value.length >= 2) {
    value = value.slice(0, 2) + ' / ' + value.slice(2, 4)
  }
  cardForm.value.expiry = value
}

const isFormValid = computed(() => {
  if (!agreed.value) return false
  if (paymentMethod.value === 'alipay') return true
  return (
    cardForm.value.name.length > 0 &&
    cardForm.value.number.replace(/\s/g, '').length >= 16 &&
    cardForm.value.expiry.length >= 7 &&
    cardForm.value.cvc.length >= 3
  )
})

const handlePay = async () => {
  if (!isFormValid.value) return
  processing.value = true
  // 模拟支付处理
  await new Promise((resolve) => setTimeout(resolve, 1500))
  processing.value = false
  emit('confirm')
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from .animate-in,
.modal-leave-to .animate-in {
  transform: scale(0.95) translateY(10px);
}

@keyframes animate-in {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.animate-in {
  animation: animate-in 0.3s ease-out forwards;
}
</style>
