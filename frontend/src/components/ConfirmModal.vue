<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4" @click.self="onCancel">
        <div class="absolute inset-0 bg-black/80 backdrop-blur-sm"></div>
        
        <div class="relative bg-surface border border-white/10 p-8 max-w-md w-full">
          <div class="absolute inset-0 border border-accent/20 pointer-events-none"></div>
          
          <div class="space-y-6">
            <div class="text-center space-y-2">
              <div class="flex justify-center">
                <component :is="icon" class="w-12 h-12" :class="iconClass" />
              </div>
              <h3 class="text-xl font-black text-white uppercase tracking-widest">{{ title }}</h3>
              <p class="text-secondary text-sm">{{ message }}</p>
            </div>
            
            <div class="flex gap-4">
              <button 
                v-if="showCancel"
                @click="onCancel"
                class="flex-1 h-12 bg-white/5 border border-white/10 text-white font-bold uppercase tracking-widest hover:bg-white/10 transition-all"
              >
                取消
              </button>
              <button 
                @click="onConfirm"
                class="flex-1 h-12 bg-accent text-white font-bold uppercase tracking-widest hover:bg-white hover:text-black transition-all"
              >
                {{ confirmText }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { CheckCircle, AlertTriangle, Info } from 'lucide-vue-next'
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  show: boolean
  title: string
  message: string
  type?: 'success' | 'warning' | 'info'
  confirmText?: string
  showCancel?: boolean
}>(), {
  type: 'info',
  confirmText: '确定',
  showCancel: false
})

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const icon = computed(() => {
  switch (props.type) {
    case 'success': return CheckCircle
    case 'warning': return AlertTriangle
    default: return Info
  }
})

const iconClass = computed(() => {
  switch (props.type) {
    case 'success': return 'text-green-500'
    case 'warning': return 'text-yellow-500'
    default: return 'text-accent'
  }
})

const onConfirm = () => emit('confirm')
const onCancel = () => emit('cancel')
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active > div:last-child,
.modal-leave-active > div:last-child {
  transition: transform 0.3s ease;
}

.modal-enter-from > div:last-child,
.modal-leave-to > div:last-child {
  transform: scale(0.9);
}
</style>
