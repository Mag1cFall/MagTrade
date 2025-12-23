<template>
  <div class="min-h-screen bg-background text-primary selection:bg-accent selection:text-white overflow-x-hidden">
    <div class="pt-32 pb-16 px-6 md:px-12 border-b border-white/5 bg-black/80 backdrop-blur-3xl sticky top-0 z-30">
      <div class="container mx-auto flex flex-col lg:flex-row justify-between items-end gap-12">
        <div class="space-y-4">
          <div class="flex items-center gap-3">
            <div class="w-1.5 h-6 bg-accent"></div>
            <span class="text-xs font-mono tracking-[0.4em] text-accent uppercase">Archive Registry v2.0.4</span>
          </div>
          <h1 class="text-7xl md:text-9xl font-black tracking-tighter text-white leading-none hero-title">
            MAG<span class="text-accent">STORE</span>
          </h1>
          <p class="text-secondary tracking-[0.3em] uppercase text-xs font-medium max-w-md">HIGH-VELOCITY ASSET ACQUISITION INTERFACE // ENCRYPTED ACCESS ONLY</p>
        </div>
        
        <div class="flex flex-col items-end gap-8 w-full lg:w-auto">
          <div class="flex items-center gap-6 w-full lg:w-auto">
            <div class="relative group w-full lg:w-80">
              <div class="absolute inset-0 bg-accent/5 scale-x-0 group-focus-within:scale-x-100 transition-transform origin-left duration-500"></div>
              <input 
                v-model="searchQuery"
                type="text"
                placeholder="ID / NAME / PROTOCOL..."
                class="w-full bg-transparent border-b border-white/10 py-3 pl-0 pr-10 text-white placeholder-gray-700 focus:border-accent focus:outline-none transition-all uppercase font-mono text-sm tracking-widest"
              />
              <Search class="absolute right-0 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-700 group-focus-within:text-accent transition-colors" />
            </div>
            
            <button 
              @click="cartOpen = true"
              class="relative p-3 group bg-white/5 hover:bg-accent transition-all duration-300 rounded-sm"
            >
              <ShoppingCart class="w-6 h-6 text-white group-hover:scale-110 transition-transform" />
              <div v-if="cartCount > 0" class="absolute -top-2 -right-2 w-5 h-5 bg-accent text-white rounded-full flex items-center justify-center text-[10px] font-bold border-2 border-background animate-pulse">
                {{ cartCount }}
              </div>
            </button>
          </div>
          
          <div class="flex flex-wrap justify-end gap-3 font-mono">
             <button 
              v-for="cat in categories" 
              :key="cat.id"
              @click="selectedCategory = cat.id"
              :class="[
                'px-6 py-2 border text-[10px] font-bold tracking-[0.2em] transition-all duration-500 uppercase rounded-full',
                selectedCategory === cat.id 
                  ? 'border-accent bg-accent text-white shadow-[0_0_15px_rgba(227,53,53,0.4)]' 
                  : 'border-white/10 text-secondary hover:border-white/40 hover:text-white bg-black/20'
              ]"
            >
              {{ cat.name }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="container mx-auto px-6 md:px-12 py-20">
      <div v-if="loading" class="flex flex-col items-center justify-center py-60 gap-6">
        <div class="relative">
          <Loader2 class="w-16 h-16 animate-spin text-accent" />
          <div class="absolute inset-0 flex items-center justify-center">
            <Zap class="w-6 h-6 text-white animate-pulse" />
          </div>
        </div>
        <span class="text-xs font-mono tracking-[0.5em] text-secondary animate-pulse">SYNCING DATASTREAMS...</span>
      </div>

      <div v-else-if="flashSales.length === 0" class="py-60 text-center relative">
        <div class="absolute inset-0 flex items-center justify-center opacity-[0.02] pointer-events-none select-none">
          <Ban class="w-96 h-96" />
        </div>
        <h3 class="text-4xl font-black text-white mb-4 tracking-tighter">NULL_RESPONSE</h3>
        <p class="text-secondary font-mono text-sm tracking-widest">COORDINATES LEAD TO VOID // RE-ENTER PARAMETERS</p>
        <button @click="searchQuery = ''; selectedCategory = 'all'" class="mt-8 px-8 py-3 bg-white text-black font-bold uppercase tracking-widest hover:bg-accent hover:text-white transition-all">Clear Filters</button>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-x-12 gap-y-24">
        <div 
          v-for="(item, index) in flashSales" 
          :key="item.id"
          class="group relative shop-card-entry"
          :style="{ '--delay': index * 0.1 + 's' }"
        >
          <FlowBorderCard 
            :active="item.status === 1" 
            :color="item.status === 1 ? '#FF3B30' : '#333'"
            border-radius="0px"
            border-width="1px"
            duration="3"
            class="transition-all duration-700"
          >
            <div class="relative aspect-[3/4] bg-surface-light overflow-hidden flex flex-col">
              <div class="absolute inset-0 bg-gradient-to-br from-accent/5 to-transparent pointer-events-none z-10"></div>
              
              <div class="w-full h-full p-8 flex items-center justify-center relative overflow-hidden">
                <div v-if="item.product?.image_url" class="w-full h-full relative z-10 transition-transform duration-1000 ease-out group-hover:scale-110">
                  <SmartImage
                    :src="item.product.image_url"
                    :alt="item.product.name"
                    class-name="w-full h-full object-contain"
                  />
                </div>
                <div v-else class="text-9xl opacity-20 grayscale transition-all duration-700 group-hover:grayscale-0 group-hover:opacity-60 relative z-10">
                  {{ getProductIcon(item.product?.name) }}
                </div>
                
                <div class="absolute inset-0 pointer-events-none opacity-0 group-hover:opacity-10 transition-opacity duration-700">
                  <div class="h-full w-full bg-[radial-gradient(circle_at_center,white_0%,transparent_70%)]"></div>
                </div>
              </div>

              <div class="absolute inset-x-0 bottom-0 p-8 translate-y-full group-hover:translate-y-0 transition-transform duration-500 ease-out z-20 bg-black/60 backdrop-blur-xl border-t border-white/10">
                 <button 
                  @click="addToCart(item)"
                  class="w-full h-16 bg-white text-black font-black uppercase tracking-[0.2em] hover:bg-accent hover:text-white transition-all flex items-center justify-center gap-3 active:scale-[0.98]"
                >
                  <Zap class="w-5 h-5 fill-current" />
                  Acquire Asset — ¥{{ item.flash_price }}
                </button>
              </div>

              <div class="absolute top-6 left-6 z-20 flex flex-col gap-3">
                <div v-if="item.status === 1" class="flex items-center gap-2 px-3 py-1 bg-accent text-white text-[10px] font-black uppercase tracking-widest animate-pulse">
                  <span class="w-1.5 h-1.5 bg-white rounded-full"></span>
                  Live Connection
                </div>
                <div class="px-3 py-1 bg-white text-black text-[10px] font-black uppercase tracking-widest">
                  -{{ calculateDiscount(item) }}% Reduced
                </div>
              </div>
              
              <div class="absolute top-6 right-6 z-20 font-mono text-[10px] text-secondary tracking-widest uppercase text-right opacity-40 group-hover:opacity-100 transition-opacity">
                Ref_ID: {{ item.id.toString().padStart(6, '0') }}<br/>
                Stock_LVL: {{ item.available_stock }}
              </div>
            </div>
          </FlowBorderCard>

          <div class="mt-8 flex justify-between items-end pl-2 pr-2">
            <div class="space-y-3 flex-1">
              <div class="flex items-center gap-2">
                <div class="w-1 h-4 bg-accent/40"></div>
                <h3 class="text-2xl md:text-3xl font-black text-white tracking-tighter uppercase group-hover:text-accent transition-colors">
                  {{ item.product?.name || 'GENERIC PROTOCOL' }}
                </h3>
              </div>
              <p class="text-secondary text-sm line-clamp-2 max-w-md font-medium leading-relaxed opacity-60 group-hover:opacity-100 transition-opacity">
                {{ item.product?.description || 'NO ADDITIONAL DATA LOGS FOUND FOR THIS REGISTRY ENTRY.' }}
              </p>
              
              <div class="flex gap-6 items-center pt-2">
                <div class="flex flex-col">
                  <span class="text-[9px] text-tertiary uppercase tracking-widest mb-1">Stock Ratio</span>
                  <div class="w-32 h-1 bg-white/5 relative overflow-hidden">
                    <div 
                      class="absolute top-0 left-0 h-full bg-accent transition-all duration-1000"
                      :style="{ width: (item.available_stock / item.total_stock * 100) + '%' }"
                    ></div>
                  </div>
                </div>
                <div class="flex flex-col border-l border-white/10 pl-6">
                  <span class="text-[9px] text-tertiary uppercase tracking-widest mb-1">Status</span>
                  <span class="text-[10px] font-mono font-bold" :class="item.status === 1 ? 'text-green-500' : 'text-orange-500'">
                    {{ item.status === 1 ? 'AVAILABLE' : 'STANDBY' }}
                  </span>
                </div>
              </div>
            </div>
            
            <div class="text-right flex flex-col items-end gap-1">
              <div class="text-xs text-tertiary font-mono line-through tracking-tighter">¥{{ item.product?.original_price }}</div>
              <div class="text-4xl font-black text-white font-mono tracking-tighter">¥{{ item.flash_price }}</div>
              <div class="w-12 h-1 bg-accent mt-2"></div>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-40 border-t border-white/5 pt-16 flex flex-col md:flex-row justify-between items-center gap-12 font-mono text-[10px] tracking-[0.4em] uppercase">
        <button 
          @click="currentPage = Math.max(1, currentPage - 1)"
          :disabled="currentPage === 1"
          class="flex items-center gap-4 text-white group disabled:opacity-20 transition-all cursor-pointer hover:text-accent"
        >
          <ChevronLeft class="w-4 h-4 group-hover:-translate-x-2 transition-transform" />
          [ Load Prev Registry ]
        </button>
        
        <div class="flex items-center gap-8">
           <span class="text-tertiary">Page {{ currentPage }} // {{ totalPages }}</span>
           <div class="flex gap-1">
              <div v-for="p in totalPages" :key="p" :class="['w-1.5 h-1.5', p === currentPage ? 'bg-accent' : 'bg-white/10']"></div>
           </div>
        </div>
        
        <button 
          @click="currentPage = Math.min(totalPages, currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="flex items-center gap-4 text-white group disabled:opacity-20 transition-all cursor-pointer hover:text-accent"
        >
          [ Load Next Registry ]
          <ChevronRight class="w-4 h-4 group-hover:translate-x-2 transition-transform" />
        </button>
      </div>
    </div>

    <CartSidebar :open="cartOpen" @close="cartOpen = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getFlashSales } from '@/api/flash-sale'
import type { FlashSale } from '@/types'
import CartSidebar from '@/components/CartSidebar.vue'
import FlowBorderCard from '@/components/FlowBorderCard.vue'
import SmartImage from '@/components/SmartImage.vue'
import { Loader2, ShoppingCart, Search, ChevronLeft, ChevronRight, Zap, Ban } from 'lucide-vue-next'
import { useCartStore } from '@/stores/cart'

const router = useRouter()
const loading = ref(true)
const flashSales = ref<FlashSale[]>([])
const currentPage = ref(1)
const totalPages = ref(8)
const cartOpen = ref(false)
const searchQuery = ref('')
const selectedCategory = ref('all')

const categories = [
  { id: 'all', name: 'ALL_PROTOCOLS' },
  { id: 'electronics', name: 'HARDWARE_LINK' },
  { id: 'fashion', name: 'NEURAL_WEAR' },
  { id: 'collectibles', name: 'EXOTIC_DATA' },
]

const cartStore = useCartStore()
const cartCount = computed(() => cartStore.totalItems)

const getProductIcon = (name: string = '') => {
  const n = name.toLowerCase()
  if (n.includes('phone')) return '📱'
  if (n.includes('mac') || n.includes('laptop')) return '💻'
  if (n.includes('pod') || n.includes('headphone') || n.includes('audio')) return '🎧'
  if (n.includes('chip') || n.includes('neural')) return '🧠'
  return '📦'
}

const calculateDiscount = (item: FlashSale) => {
  const original = item.product?.original_price || item.flash_price
  if (original === 0) return 0
  return Math.round((1 - item.flash_price / original) * 100)
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getFlashSales(currentPage.value, 10, 1)
    if (res.code === 0) {
      flashSales.value = res.data.flash_sales
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

watch([currentPage, selectedCategory], () => fetchData())

onMounted(() => fetchData())

const addToCart = (item: FlashSale) => {
  cartStore.addItem(item)
  cartOpen.value = true
}
</script>

<style scoped>
.hero-title {
  font-size: clamp(3rem, 12vw, 12rem);
}

.shop-card-entry {
  opacity: 0;
  transform: translateY(30px);
  animation: slideUpIn 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  animation-delay: var(--delay);
}

@keyframes slideUpIn {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
