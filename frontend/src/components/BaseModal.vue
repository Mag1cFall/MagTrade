<script setup lang="ts">
import { X, AlertTriangle } from 'lucide-vue-next'

defineProps<{
  show: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  type?: 'danger' | 'info'
}>()

const emit = defineEmits(['confirm', 'cancel'])
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/80 backdrop-blur-sm" @click="emit('cancel')"></div>
      
      <div class="relative w-full max-w-md bg-surface border border-white/10 shadow-2xl shadow-black overflow-hidden slide-up">
        <div :class="['h-1 w-full', type === 'danger' ? 'bg-accent' : 'bg-primary']"></div>
        
        <div class="p-6">
          <div class="flex items-start gap-4">
            <div v-if="type === 'danger'" class="p-2 bg-accent/10 rounded-full">
              <AlertTriangle class="w-6 h-6 text-accent" />
            </div>
            <div class="flex-1">
              <h3 class="text-xl font-bold text-white tracking-tight mb-2">{{ title }}</h3>
              <p class="text-secondary text-sm leading-relaxed">{{ message }}</p>
            </div>
            <button @click="emit('cancel')" class="text-tertiary hover:text-white transition-colors">
              <X class="w-5 h-5" />
            </button>
          </div>

          <div class="mt-8 flex gap-3">
            <button 
              @click="emit('cancel')"
              class="flex-1 px-4 py-3 border border-white/5 bg-white/5 text-white text-xs font-bold uppercase tracking-widest hover:bg-white/10 transition-colors"
            >
              {{ cancelText || 'Cancel' }}
            </button>
            <button 
              @click="emit('confirm')"
              :class="[
                'flex-1 px-4 py-3 text-xs font-bold uppercase tracking-widest transition-all',
                type === 'danger' ? 'bg-accent text-white hover:bg-accent-hover' : 'bg-white text-black hover:bg-gray-200'
              ]"
            >
              {{ confirmText || 'Confirm' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.slide-up {
  animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideUp {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
</style>