<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Loader2, ArrowRight, ArrowLeft } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  email: '',
  password: ''
})
const loading = ref(false)
const error = ref('')

const handleSubmit = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await authStore.register(form.value)
    if (res.code === 0) {
      // 注册成功直接跳转首页（后端逻辑为自动激活）
      router.push('/')
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    // 优先显示后端返回的具体错误信息
    error.value = e.response?.data?.message || e.message || 'Registration failed'
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
        <h1 class="text-4xl font-bold tracking-tighter text-white mb-2">JOIN <span class="text-accent">NOW</span>.</h1>
        <p class="text-secondary text-sm tracking-widest uppercase">Create your identity</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="space-y-2">
          <label class="text-xs text-secondary uppercase tracking-widest">Username</label>
          <input 
            v-model="form.username"
            type="text" 
            class="w-full bg-surface border border-white/10 px-4 py-3 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary"
            placeholder="CREATE_ID"
            required
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-secondary uppercase tracking-widest">Email</label>
          <input 
            v-model="form.email"
            type="email" 
            class="w-full bg-surface border border-white/10 px-4 py-3 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary"
            placeholder="MAIL@ADDRESS"
            required
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-secondary uppercase tracking-widest">Password</label>
          <input 
            v-model="form.password"
            type="password" 
            class="w-full bg-surface border border-white/10 px-4 py-3 text-white focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary"
            placeholder="••••••••"
            required
          />
        </div>

        <div v-if="error" class="text-red-500 text-sm bg-red-900/10 border border-red-900/20 p-3 text-center">
          {{ error }}
        </div>

        <button 
          type="submit" 
          :disabled="loading"
          class="w-full h-14 bg-white text-black font-bold tracking-widest uppercase hover:bg-gray-200 transition-colors flex items-center justify-center gap-2"
        >
          <Loader2 v-if="loading" class="w-5 h-5 animate-spin" />
          <span v-else>Register ID</span>
          <ArrowRight v-if="!loading" class="w-4 h-4" />
        </button>
      </form>

      <div class="mt-8 text-center">
        <span class="text-secondary text-sm">ALREADY REGISTERED? </span>
        <button @click="router.push('/login')" class="text-white text-sm font-bold border-b border-white hover:text-accent hover:border-accent transition-colors ml-2">
          LOGIN HERE
        </button>
      </div>
    </div>
  </div>
</template>