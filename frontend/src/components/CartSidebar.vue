<template>
  <Teleport to="body">
    <Transition name="slide">
      <div v-if="open" class="fixed inset-0 z-50 flex justify-end">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="$emit('close')"></div>
        
        <div class="relative w-full max-w-md bg-surface border-l border-white/10 h-full flex flex-col">
          <div class="p-6 border-b border-white/10 flex items-center justify-between">
            <h2 class="text-xl font-bold text-white flex items-center gap-2">
              <ShoppingCart class="w-5 h-5" />
              购物车
            </h2>
            <button @click="$emit('close')" class="text-secondary hover:text-white">
              <X class="w-5 h-5" />
            </button>
          </div>
          
          <div class="flex-grow overflow-y-auto p-6">
            <div v-if="cartStore.items.length === 0" class="text-center py-12 text-secondary">
              <ShoppingCart class="w-12 h-12 mx-auto mb-4 opacity-30" />
              <p>Your cart is empty</p>
            </div>
            
            <div v-else class="space-y-4">
              <div 
                v-for="item in cartStore.items" 
                :key="item.product.id"
                class="flex gap-4 p-4 bg-black/30 rounded-lg border border-white/5"
              >
                <div class="w-16 h-16 bg-white/5 rounded flex items-center justify-center text-2xl">
                  📦
                </div>
                <div class="flex-grow">
                  <h4 class="text-white font-medium truncate">{{ item.product.product?.name || 'Product' }}</h4>
                  <p class="text-accent font-bold">¥{{ item.product.flash_price }}</p>
                  <div class="flex items-center gap-2 mt-2">
                    <button 
                      @click="cartStore.updateQuantity(item.product.id, item.quantity - 1)"
                      class="w-6 h-6 rounded bg-white/10 text-white hover:bg-white/20"
                    >-</button>
                    <span class="text-white w-8 text-center">{{ item.quantity }}</span>
                    <button 
                      @click="cartStore.updateQuantity(item.product.id, item.quantity + 1)"
                      class="w-6 h-6 rounded bg-white/10 text-white hover:bg-white/20"
                    >+</button>
                  </div>
                </div>
                <button 
                  @click="cartStore.removeItem(item.product.id)"
                  class="text-secondary hover:text-accent"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
          
          <div class="p-6 border-t border-white/10 space-y-4">
            <div class="flex justify-between text-lg">
              <span class="text-secondary">Total</span>
              <span class="text-white font-bold">¥{{ cartStore.totalPrice.toFixed(2) }}</span>
            </div>
            <button 
              :disabled="cartStore.items.length === 0"
              class="w-full py-3 bg-accent text-white font-bold rounded hover:bg-accent/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Checkout ({{ cartStore.totalItems }} items)
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ShoppingCart, X, Trash2 } from 'lucide-vue-next'
import { useCartStore } from '@/stores/cart'

defineProps<{ open: boolean }>()
defineEmits<{ close: [] }>()

const cartStore = useCartStore()
</script>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
}

.slide-enter-from > div:last-child,
.slide-leave-to > div:last-child {
  transform: translateX(100%);
}
</style>
