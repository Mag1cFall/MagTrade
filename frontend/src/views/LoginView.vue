<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Loader2, ArrowRight, Eye, EyeOff, ArrowLeft } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const form = ref({
  username: '',
  password: ''
})
const loading = ref(false)
const error = ref('')
const showPassword = ref(false)

const handleSubmit = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await authStore.login(form.value)
    if (res.code === 0) {
      const redirect = route.query.redirect as string || '/'
      router.push(redirect)
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    // 优先显示后端返回的具体错误信息
    error.value = e.response?.data?.message || '用户名或密码错误'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background p-6 relative">
    <!-- Back Button -->
    <button
      @click="router.push('/')"
      class="absolute top-6 left-6 text-secondary hover:text-white flex items-center gap-2 transition-colors uppercase tracking-widest text-sm font-bold"
    >
      <ArrowLeft class="w-4 h-4" /> Back
    </button>

    <div class="w-full max-w-md">
      <div class="mb-12 text-center">
        <h1 class="text-4xl font-bold tracking-tighter text-white mb-2">WELCOME <span class="text-accent">BACK</span>.</h1>
        <p class="text-secondary text-sm tracking-widest uppercase">Enter the void</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="space-y-2">
          <label class="text-xs text-secondary uppercase tracking-widest">Username</label>
          <input 
            v-model="form.username"
            type="text" 
            class="w-full bg-surface border border-white/10 px-4 py-3 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary"
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
              class="w-full bg-surface border border-white/10 px-4 py-3 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary pr-10"
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

        <div v-if="error" class="text-red-500 text-sm bg-red-900/10 border border-red-900/20 p-3 text-center animate-fade-in">
          {{ error }}
        </div>

        <button 
          type="submit" 
          :disabled="loading"
          class="w-full h-14 bg-white text-black font-bold tracking-widest uppercase hover:bg-gray-200 transition-colors flex items-center justify-center gap-2"
        >
          <Loader2 v-if="loading" class="w-5 h-5 animate-spin" />
          <span v-else>Initiate Session</span>
          <ArrowRight v-if="!loading" class="w-4 h-4" />
        </button>
      </form>

      <div class="mt-8 text-center">
        <span class="text-secondary text-sm">NO CREDENTIALS? </span>
        <button @click="router.push('/register')" class="text-white text-sm font-bold border-b border-white hover:text-accent hover:border-accent transition-colors ml-2">
          REGISTER NEW ID
        </button>
      </div>
    </div>
  </div>
</template>