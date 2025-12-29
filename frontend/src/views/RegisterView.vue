<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { sendEmailCode } from '@/api/auth'
import { ArrowRight, ArrowLeft, Check, AlertTriangle } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  email: '',
  password: '',
  email_code: '',
})
const loading = ref(false)
const sendingCode = ref(false)
const error = ref('')
const codeSent = ref(false)
const countdown = ref(0)
let countdownTimer: ReturnType<typeof setInterval> | null = null

const canSendCode = computed(() => {
  return form.value.email && !sendingCode.value && countdown.value === 0
})

const startCountdown = () => {
  countdown.value = 60
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      if (countdownTimer) clearInterval(countdownTimer)
    }
  }, 1000)
}

const handleSendCode = async () => {
  if (!canSendCode.value) return
  sendingCode.value = true
  error.value = ''
  try {
    const res = await sendEmailCode(form.value.email)
    if (res.code === 0) {
      codeSent.value = true
      startCountdown()
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || '发送失败，请稍后重试'
  } finally {
    sendingCode.value = false
  }
}

const handleSubmit = async () => {
  if (!form.value.email_code) {
    error.value = '请输入验证码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const res = await authStore.register(form.value)
    if (res.code === 0) {
      router.push('/')
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || e.message || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex flex-col md:flex-row bg-background">
    <!-- Left: Form -->
    <div
      class="w-full md:w-1/2 flex items-center justify-center p-8 relative z-10 bg-background border-r border-white/5"
    >
      <button
        class="absolute top-8 left-8 text-secondary hover:text-white flex items-center gap-2 transition-colors uppercase tracking-widest text-sm font-bold group"
        @click="router.push('/')"
      >
        <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" /> Back
      </button>

      <div class="w-full max-w-md space-y-8">
        <div>
          <h1 class="text-4xl font-bold tracking-tighter text-white mb-2">
            JOIN <span class="text-accent">NOW</span>.
          </h1>
          <p class="text-secondary text-sm tracking-widest uppercase">
            Create Your Digital Identity
          </p>
        </div>

        <form class="space-y-6" @submit.prevent="handleSubmit">
          <div class="space-y-2">
            <label class="text-xs text-secondary uppercase tracking-widest">Username</label>
            <input
              v-model="form.username"
              type="text"
              class="w-full bg-surface border border-white/10 px-4 py-4 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary font-mono"
              placeholder="CREATE_ID"
              required
            />
          </div>

          <div class="space-y-2">
            <label class="text-xs text-secondary uppercase tracking-widest"
              >Email Verification</label
            >
            <div class="flex gap-3">
              <input
                v-model="form.email"
                type="email"
                class="flex-1 bg-surface border border-white/10 px-4 py-4 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary font-mono"
                placeholder="MAIL@ADDRESS"
                required
              />
              <button
                type="button"
                :disabled="!canSendCode"
                class="px-6 bg-surface border border-white/10 text-white font-bold text-sm uppercase tracking-wider hover:bg-white/5 hover:border-accent transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 min-w-[120px]"
                @click="handleSendCode"
              >
                <div v-if="sendingCode" class="loading-rotation"><div></div></div>
                <template v-else>
                  <span v-if="countdown > 0">{{ countdown }}s</span>
                  <span v-else>Send</span>
                </template>
              </button>
            </div>
            <p v-if="codeSent" class="text-xs text-green-500 mt-2 flex items-center gap-1">
              <Check class="w-3 h-3" /> Code sent to inbox
            </p>
          </div>

          <div class="space-y-2">
            <label class="text-xs text-secondary uppercase tracking-widest"
              >Verification Code</label
            >
            <input
              v-model="form.email_code"
              type="text"
              class="w-full bg-surface border border-white/10 px-4 py-4 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary tracking-[0.5em] text-center font-mono text-lg"
              placeholder="000000"
              maxlength="6"
              required
            />
          </div>

          <div class="space-y-2">
            <label class="text-xs text-secondary uppercase tracking-widest">Password</label>
            <input
              v-model="form.password"
              type="password"
              class="w-full bg-surface border border-white/10 px-4 py-4 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary font-mono"
              placeholder="••••••••"
              required
            />
          </div>

          <div
            v-if="error"
            class="flex items-center gap-3 text-red-500 text-sm bg-red-900/10 border border-red-900/20 p-4 animate-fade-in"
          >
            <AlertTriangle class="w-5 h-5 flex-shrink-0" />
            {{ error }}
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full h-14 bg-white text-black font-bold tracking-widest uppercase hover:bg-gray-200 transition-all flex items-center justify-center gap-3 disabled:opacity-50 group relative overflow-hidden"
          >
            <div v-if="loading" class="absolute inset-0 flex items-center justify-center bg-white">
              <div class="loading-rotation"><div></div></div>
            </div>
            <span v-else class="relative z-10 group-hover:translate-x-1 transition-transform"
              >Initialize Account</span
            >
            <ArrowRight
              v-if="!loading"
              class="w-4 h-4 relative z-10 group-hover:translate-x-1 transition-transform"
            />
          </button>
        </form>

        <div class="mt-8 text-center pt-8 border-t border-white/5">
          <span class="text-secondary text-sm">ALREADY REGISTERED? </span>
          <button
            class="text-white text-sm font-bold border-b border-white hover:text-accent hover:border-accent transition-colors ml-2 uppercase tracking-wider"
            @click="router.push('/login')"
          >
            Login Here
          </button>
        </div>
      </div>
    </div>

    <!-- Right: Visual -->
    <div class="hidden md:block w-1/2 relative overflow-hidden bg-black">
      <div
        class="absolute inset-0 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-accent/20 via-black to-black opacity-40"
      ></div>
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:64px_64px]"
      ></div>

      <div class="absolute bottom-12 left-12">
        <div class="text-6xl font-bold text-white/10 tracking-tighter mb-4">NEW<br />ENTITY</div>
        <div class="flex gap-4 text-xs font-mono text-tertiary uppercase tracking-widest">
          <span>Registration Protocol</span>
          <span>//</span>
          <span>v2.0.4</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes rotation {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.loading-rotation {
  width: 20px;
  height: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  animation: rotation 1s linear infinite;
}

.loading-rotation div {
  width: 100%;
  height: 100%;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
}

.animate-fade-in {
  animation: fadeIn 0.3s ease-out forwards;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
