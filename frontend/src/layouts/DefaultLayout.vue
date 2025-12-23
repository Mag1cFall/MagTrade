<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { LogOut, Menu, X, Zap } from 'lucide-vue-next'
import AIChat from '@/components/AIChat.vue'

const route = useRoute()
const authStore = useAuthStore()
const isScrolled = ref(false)
const isMobileMenuOpen = ref(false)

const handleScroll = () => {
  isScrolled.value = window.scrollY > 20
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

const navClasses = computed(() => {
  return `absolute top-0 w-full z-50 transition-all duration-500 ease-out border-b ${
    isScrolled.value
      ? 'fixed bg-background/80 backdrop-blur-xl border-white/5 py-4'
      : 'bg-transparent border-transparent py-6'
  }`
})

const linkClasses = (path: string) => {
  const isActive = route.path === path
  return `relative text-sm font-medium tracking-wide transition-colors duration-300 ${
    isActive ? 'text-white' : 'text-secondary hover:text-white'
  } ${isActive ? 'after:content-[""] after:absolute after:-bottom-1 after:left-0 after:w-full after:h-[1px] after:bg-accent' : ''}`
}

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}
</script>

<template>
  <div class="min-h-screen bg-background text-primary selection:bg-accent selection:text-white flex flex-col font-sans">
    <!-- Navigation -->
    <nav :class="navClasses">
      <div class="container mx-auto px-6 md:px-12 flex justify-between items-center">
        <!-- Logo -->
        <RouterLink to="/" class="group flex items-center gap-2 z-50">
          <div class="relative w-8 h-8 bg-white flex items-center justify-center rounded-sm overflow-hidden group-hover:scale-105 transition-transform duration-300">
            <Zap class="w-5 h-5 text-black fill-current" />
            <div class="absolute inset-0 bg-accent/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300 mix-blend-multiply"></div>
          </div>
          <span class="text-xl font-bold tracking-tighter">MAG<span class="text-accent">TRADE</span>.</span>
        </RouterLink>

        <!-- Desktop Menu -->
        <div class="hidden md:flex items-center gap-8">
          <RouterLink to="/" :class="linkClasses('/')">首页</RouterLink>
          <RouterLink to="/orders" :class="linkClasses('/orders')" v-if="authStore.isAuthenticated">我的订单</RouterLink>
          <RouterLink to="/admin" :class="linkClasses('/admin')" v-if="authStore.user?.role === 'admin'" class="text-accent font-bold">控制台</RouterLink>
          
          <div class="w-px h-4 bg-white/10 mx-2"></div>

          <template v-if="authStore.isAuthenticated">
            <div class="flex items-center gap-4">
              <span class="text-sm text-secondary font-mono">{{ authStore.user?.username }}</span>
              <button @click="authStore.logout" class="text-secondary hover:text-accent transition-colors" title="退出登录">
                <LogOut class="w-5 h-5" />
              </button>
            </div>
          </template>
          <template v-else>
            <RouterLink to="/login" class="text-sm font-medium text-white hover:text-accent transition-colors">登录</RouterLink>
            <RouterLink to="/register" class="px-5 py-2 bg-white text-black text-sm font-bold rounded hover:bg-gray-200 transition-colors">注册</RouterLink>
          </template>
        </div>

        <!-- Mobile Menu Button -->
        <button @click="toggleMobileMenu" class="md:hidden z-50 text-white">
          <Menu v-if="!isMobileMenuOpen" class="w-6 h-6" />
          <X v-else class="w-6 h-6" />
        </button>
      </div>

      <!-- Mobile Menu Overlay -->
      <div v-show="isMobileMenuOpen" class="fixed inset-0 bg-black/95 backdrop-blur-xl z-40 flex flex-col items-center justify-center gap-8 md:hidden transition-opacity duration-300">
        <RouterLink to="/" @click="isMobileMenuOpen = false" class="text-2xl font-bold hover:text-accent">首页</RouterLink>
        <RouterLink to="/orders" @click="isMobileMenuOpen = false" v-if="authStore.isAuthenticated" class="text-2xl font-bold hover:text-accent">我的订单</RouterLink>
        <RouterLink to="/admin" @click="isMobileMenuOpen = false" v-if="authStore.user?.role === 'admin'" class="text-2xl font-bold text-accent">商家后台</RouterLink>
        
        <div class="w-12 h-px bg-white/10"></div>

        <template v-if="authStore.isAuthenticated">
          <div class="text-secondary">{{ authStore.user?.username }}</div>
          <button @click="authStore.logout(); isMobileMenuOpen = false" class="text-xl text-red-500">退出登录</button>
        </template>
        <template v-else>
          <RouterLink to="/login" @click="isMobileMenuOpen = false" class="text-xl">登录</RouterLink>
          <RouterLink to="/register" @click="isMobileMenuOpen = false" class="text-xl text-accent">注册账号</RouterLink>
        </template>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="flex-grow">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Footer -->
    <footer class="border-t border-white/5 mt-20 py-12">
      <div class="container mx-auto px-6 md:px-12 flex flex-col md:flex-row justify-between items-center gap-6">
        <div class="text-secondary text-sm">
          © 2025 MagTrade Inc. All rights reserved.
        </div>
        <div class="flex gap-6 text-secondary text-sm">
          <RouterLink to="/privacy" class="hover:text-white transition-colors">Privacy</RouterLink>
          <RouterLink to="/terms" class="hover:text-white transition-colors">Terms</RouterLink>
          <RouterLink to="/contact" class="hover:text-white transition-colors">Support</RouterLink>
        </div>
      </div>
    </footer>

    <!-- AI Chat Floating Button -->
    <AIChat />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>