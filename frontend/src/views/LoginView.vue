<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { checkCaptcha, getCaptcha } from '@/api/auth'
import { ArrowRight, Eye, EyeOff, ArrowLeft, RefreshCw, AlertTriangle, ShieldCheck } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const form = ref({
  username: '',
  password: '',
  captcha_id: '',
  captcha_code: ''
})
const loading = ref(false)
const error = ref('')
const showPassword = ref(false)
const needsCaptcha = ref(false)
const isLocked = ref(false)
const captchaUrl = ref('')
const captchaLoading = ref(false)

const checkNeedsCaptcha = async () => {
  if (!form.value.username) return
  try {
    const res = await checkCaptcha(form.value.username)
    needsCaptcha.value = res.needs_captcha
    isLocked.value = res.is_locked
    if (needsCaptcha.value) {
      await refreshCaptcha()
    }
  } catch (e) {}
}

const refreshCaptcha = async () => {
  captchaLoading.value = true
  try {
    const res = await getCaptcha(form.value.username)
    const captchaId = (res as any).headers?.['x-captcha-id']
    if (captchaId) {
      form.value.captcha_id = captchaId
    }
    const blob = (res as any).data || res
    captchaUrl.value = URL.createObjectURL(blob)
  } catch (e) {} finally {
    captchaLoading.value = false
  }
}

const handleSubmit = async () => {
  if (isLocked.value) {
    error.value = '账号已被锁定，请15分钟后再试'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const res = await authStore.login(form.value)
    if (res.code === 0) {
      const redirect = route.query.redirect as string || '/'
      router.push(redirect)
    } else {
      error.value = res.message
      if ((res as any).data?.needs_captcha) {
        needsCaptcha.value = true
        await refreshCaptcha()
      }
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || '用户名或密码错误'
    if (e.response?.data?.data?.needs_captcha) {
      needsCaptcha.value = true
      await refreshCaptcha()
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (form.value.username) {
    checkNeedsCaptcha()
  }
})
</script>

<template>
  <div class="min-h-screen flex flex-col md:flex-row bg-background">
    <!-- Left: Form -->
    <div class="w-full md:w-1/2 flex items-center justify-center p-8 relative z-10 bg-background border-r border-white/5">
      <button
        @click="router.push('/')"
        class="absolute top-8 left-8 text-secondary hover:text-white flex items-center gap-2 transition-colors uppercase tracking-widest text-sm font-bold group"
      >
        <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" /> Back
      </button>

      <div class="w-full max-w-md space-y-8">
        <div>
          <h1 class="text-4xl font-bold tracking-tighter text-white mb-2">WELCOME <span class="text-accent">BACK</span>.</h1>
          <p class="text-secondary text-sm tracking-widest uppercase">Initiate Session Protocol</p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <div class="space-y-2">
            <label class="text-xs text-secondary uppercase tracking-widest">Username</label>
            <input 
              v-model="form.username"
              @blur="checkNeedsCaptcha"
              type="text" 
              class="w-full bg-surface border border-white/10 px-4 py-4 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary font-mono"
              placeholder="ACCESS_ID"
              required
            />
          </div>

          <div class="space-y-2">
            <div class="flex justify-between">
              <label class="text-xs text-secondary uppercase tracking-widest">Password</label>
            </div>
            <div class="relative">
              <input 
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                class="w-full bg-surface border border-white/10 px-4 py-4 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary pr-10 font-mono"
                placeholder="••••••••"
                required
              />
              <button 
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-secondary hover:text-white transition-colors"
              >
                <Eye v-if="!showPassword" class="w-4 h-4" />
                <EyeOff v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div v-if="needsCaptcha" class="space-y-2 animate-fade-in">
            <label class="text-xs text-secondary uppercase tracking-widest">Security Check</label>
            <div class="flex gap-3">
              <input 
                v-model="form.captcha_code"
                type="text" 
                class="flex-1 bg-surface border border-white/10 px-4 py-3 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary font-mono"
                placeholder="CAPTCHA"
                maxlength="6"
                required
              />
              <div class="relative h-12 w-32 bg-surface border border-white/10 cursor-pointer overflow-hidden group" @click="refreshCaptcha">
                <img v-if="captchaUrl" :src="captchaUrl" class="h-full w-full object-cover opacity-80 group-hover:opacity-100 transition-opacity" />
                <div v-if="captchaLoading" class="absolute inset-0 flex items-center justify-center bg-surface/80">
                  <RefreshCw class="w-4 h-4 animate-spin text-accent" />
                </div>
              </div>
            </div>
          </div>

          <div v-if="error" class="flex items-center gap-3 text-red-500 text-sm bg-red-900/10 border border-red-900/20 p-4 animate-fade-in">
            <AlertTriangle class="w-5 h-5 flex-shrink-0" />
            {{ error }}
          </div>

          <div v-if="isLocked" class="flex items-center gap-3 text-orange-500 text-sm bg-orange-900/10 border border-orange-900/20 p-4 animate-fade-in">
            <ShieldCheck class="w-5 h-5 flex-shrink-0" />
            Account locked due to multiple failed attempts. Please wait 15 minutes.
          </div>

          <button 
            type="submit" 
            :disabled="loading || isLocked"
            class="w-full h-14 bg-white text-black font-bold tracking-widest uppercase hover:bg-gray-200 transition-all flex items-center justify-center gap-3 disabled:opacity-50 group relative overflow-hidden"
          >
            <div v-if="loading" class="absolute inset-0 flex items-center justify-center bg-white">
               <!-- Rotation Loader -->
               <div class="loading-rotation"><div></div></div>
            </div>
            <span v-else class="relative z-10 group-hover:translate-x-1 transition-transform">Initiate Session</span>
            <ArrowRight v-if="!loading" class="w-4 h-4 relative z-10 group-hover:translate-x-1 transition-transform" />
          </button>
        </form>

        <div class="mt-8 text-center pt-8 border-t border-white/5">
          <span class="text-secondary text-sm">NO CREDENTIALS? </span>
          <button @click="router.push('/register')" class="text-white text-sm font-bold border-b border-white hover:text-accent hover:border-accent transition-colors ml-2 uppercase tracking-wider">
            Register New ID
          </button>
        </div>
      </div>
    </div>

    <!-- Right: Visual -->
    <div class="hidden md:block w-1/2 relative overflow-hidden bg-black">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-accent/20 via-black to-black opacity-40"></div>
      <div class="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:64px_64px]"></div>
      
      <div class="absolute bottom-12 left-12">
        <div class="text-6xl font-bold text-white/10 tracking-tighter mb-4">SECURE<br/>ACCESS</div>
        <div class="flex gap-4 text-xs font-mono text-tertiary uppercase tracking-widest">
          <span>Encrypted</span>
          <span>//</span>
          <span>Distributed</span>
          <span>//</span>
          <span>Verified</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes rotation {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.loading-rotation {
  width: 24px;
  height: 24px;
  display: flex;
  justify-content: center;
  align-items: center;
  animation: rotation 1s linear infinite;
}

.loading-rotation div {
  width: 100%;
  height: 100%;
  border: 2px solid #000;
  border-top-color: transparent;
  border-radius: 50%;
}

.animate-fade-in {
  animation: fadeIn 0.3s ease-out forwards;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>