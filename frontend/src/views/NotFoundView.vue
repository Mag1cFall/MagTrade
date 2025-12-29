<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'

const router = useRouter()
</script>

<template>
  <div class="min-h-[80vh] flex flex-col items-center justify-center text-center px-4">
    <h1 class="text-9xl font-mono font-bold text-white/5 mb-4">404</h1>
    <h2 class="text-2xl font-bold text-white mb-4">VOID DETECTED</h2>
    <p class="text-secondary max-w-md mb-8">
      The coordinates you are attempting to reach do not exist in this dimension.
    </p>
    <button
      class="px-8 py-3 bg-white text-black font-bold tracking-widest uppercase hover:bg-gray-200 transition-colors flex items-center gap-2"
      @click="router.push('/')"
    >
      <ArrowLeft class="w-4 h-4" /> Return to Base
    </button>
  </div>
</template>
