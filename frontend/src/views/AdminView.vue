<script setup lang="ts">
import { ref, watch } from 'vue'
import { createProduct, createFlashSale } from '@/api/admin'
import { getProducts } from '@/api/product'
import type { Product } from '@/types'
import { Loader2, Plus, Package, Zap, RefreshCw } from 'lucide-vue-next'

const activeTab = ref<'product' | 'flash_sale'>('product')
const products = ref<Product[]>([])
const loading = ref(false)
const message = ref('')
const messageType = ref<'success' | 'error'>('success')

// Product Form
const productForm = ref({
  name: '',
  description: '',
  original_price: 0,
  image_url: ''
})

// Flash Sale Form
const flashSaleForm = ref({
  product_id: 0,
  flash_price: 0,
  total_stock: 100,
  per_user_limit: 1,
  start_time: '',
  end_time: ''
})

const showMessage = (msg: string, type: 'success' | 'error') => {
  message.value = msg
  messageType.value = type
  setTimeout(() => message.value = '', 3000)
}

const handleCreateProduct = async () => {
  loading.value = true
  try {
    const res = await createProduct(productForm.value)
    if (res.code === 0) {
      showMessage(`Product created! ID: ${res.data.id}`, 'success')
      // Auto fill flash sale product id for convenience
      flashSaleForm.value.product_id = res.data.id
      // Reset form but keep image url as it might be reused
      const img = productForm.value.image_url
      productForm.value = { name: '', description: '', original_price: 0, image_url: img }
    } else {
      showMessage(res.message, 'error')
    }
  } catch (e: any) {
    showMessage(e.message || 'Failed to create product', 'error')
  } finally {
    loading.value = false
  }
}

const fetchProducts = async () => {
  try {
    const res = await getProducts(1, 100)
    if (res.code === 0) {
      products.value = res.data.products
    }
  } catch (e) {
    console.error(e)
  }
}

watch(activeTab, (val) => {
  if (val === 'flash_sale') {
    fetchProducts()
  }
})

const handleCreateFlashSale = async () => {
  loading.value = true
  try {
    // Ensure time format is RFC3339
    const payload = {
      ...flashSaleForm.value,
      start_time: new Date(flashSaleForm.value.start_time).toISOString(),
      end_time: new Date(flashSaleForm.value.end_time).toISOString()
    }
    
    const res = await createFlashSale(payload)
    if (res.code === 0) {
      showMessage(`Flash Sale created! ID: ${res.data.id}`, 'success')
    } else {
      showMessage(res.message, 'error')
    }
  } catch (e: any) {
    showMessage(e.message || 'Failed to create flash sale', 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen pb-20">
    <div class="py-12 border-b border-white/5 mb-8">
      <h1 class="text-4xl font-bold text-white tracking-tighter mb-2">MERCHANT <span class="text-accent">CONSOLE</span>.</h1>
      <p class="text-secondary text-sm tracking-widest uppercase">Manage inventory & events</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <!-- Sidebar -->
      <div class="space-y-2">
        <button 
          @click="activeTab = 'product'"
          :class="[
            'w-full text-left px-4 py-3 text-sm font-bold uppercase tracking-wider transition-colors border-l-2',
            activeTab === 'product' ? 'border-accent text-white bg-white/5' : 'border-transparent text-secondary hover:text-white hover:bg-white/5'
          ]"
        >
          <div class="flex items-center gap-3">
            <Package class="w-4 h-4" /> Create Product
          </div>
        </button>
        <button 
          @click="activeTab = 'flash_sale'"
          :class="[
            'w-full text-left px-4 py-3 text-sm font-bold uppercase tracking-wider transition-colors border-l-2',
            activeTab === 'flash_sale' ? 'border-accent text-white bg-white/5' : 'border-transparent text-secondary hover:text-white hover:bg-white/5'
          ]"
        >
          <div class="flex items-center gap-3">
            <Zap class="w-4 h-4" /> Create Flash Sale
          </div>
        </button>
      </div>

      <!-- Content -->
      <div class="lg:col-span-3 bg-surface border border-white/5 p-8">
        <!-- Message Alert -->
        <div v-if="message" :class="[
          'mb-6 p-4 text-sm font-bold tracking-wide border',
          messageType === 'success' ? 'bg-green-900/20 text-green-500 border-green-900' : 'bg-red-900/20 text-red-500 border-red-900'
        ]">
          > {{ message }}
        </div>

        <!-- Product Form -->
        <form v-if="activeTab === 'product'" @submit.prevent="handleCreateProduct" class="space-y-6 max-w-2xl">
          <div class="space-y-2">
            <label class="text-xs text-secondary uppercase tracking-widest">Product Name</label>
            <input v-model="productForm.name" type="text" required class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none" />
          </div>
          
          <div class="space-y-2">
            <label class="text-xs text-secondary uppercase tracking-widest">Description</label>
            <textarea v-model="productForm.description" rows="3" class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none"></textarea>
          </div>

          <div class="grid grid-cols-2 gap-6">
            <div class="space-y-2">
              <label class="text-xs text-secondary uppercase tracking-widest">Original Price (¥)</label>
              <input v-model.number="productForm.original_price" type="number" step="0.01" required class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none" />
            </div>
            <div class="space-y-2">
              <label class="text-xs text-secondary uppercase tracking-widest">Digital Asset</label>
              
              <!-- Visual Selector -->
              <div class="grid grid-cols-4 gap-3 mb-3">
                <button
                  type="button"
                  v-for="type in ['phone', 'laptop', 'headphone', 'box']"
                  :key="type"
                  @click="productForm.image_url = `local:${type}`"
                  :class="[
                    'py-3 border text-xs font-bold uppercase tracking-wider transition-all duration-300',
                    productForm.image_url === `local:${type}`
                      ? 'bg-accent text-white border-accent shadow-[0_0_15px_rgba(255,59,48,0.3)]'
                      : 'bg-black/30 border-white/10 text-secondary hover:bg-white/5 hover:border-white/30'
                  ]"
                >
                  {{ type }}
                </button>
              </div>

              <div class="relative group">
                <input
                  v-model="productForm.image_url"
                  type="text"
                  placeholder="Select above or paste URL..."
                  class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none text-xs font-mono tracking-wide"
                />
              </div>
            </div>
          </div>

          <button type="submit" :disabled="loading" class="px-8 py-3 bg-white text-black font-bold uppercase tracking-widest hover:bg-gray-200 transition-colors flex items-center gap-2">
            <Loader2 v-if="loading" class="w-4 h-4 animate-spin" />
            <Plus v-else class="w-4 h-4" />
            Create Product
          </button>
        </form>

        <!-- Flash Sale Form -->
        <form v-if="activeTab === 'flash_sale'" @submit.prevent="handleCreateFlashSale" class="space-y-6 max-w-2xl">
          <div class="space-y-2">
            <div class="flex justify-between items-center">
              <label class="text-xs text-secondary uppercase tracking-widest">Select Product</label>
              <button type="button" @click="fetchProducts" class="text-xs text-accent hover:text-white flex items-center gap-1 transition-colors">
                <RefreshCw class="w-3 h-3" /> Refresh List
              </button>
            </div>
            <div class="relative">
              <select
                v-model.number="flashSaleForm.product_id"
                required
                class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none appearance-none cursor-pointer"
              >
                <option :value="0" disabled>Select a product...</option>
                <option v-for="p in products" :key="p.id" :value="p.id">
                  {{ p.name }} (ID: {{ p.id }}) - ¥{{ p.original_price }}
                </option>
              </select>
              <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-secondary text-xs">
                ▼
              </div>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-6">
            <div class="space-y-2">
              <label class="text-xs text-secondary uppercase tracking-widest">Flash Price (¥)</label>
              <input v-model.number="flashSaleForm.flash_price" type="number" step="0.01" required class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none" />
            </div>
            <div class="space-y-2">
              <label class="text-xs text-secondary uppercase tracking-widest">Total Stock</label>
              <input v-model.number="flashSaleForm.total_stock" type="number" required class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none" />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-6">
            <div class="space-y-2">
              <label class="text-xs text-secondary uppercase tracking-widest">Start Time</label>
              <input v-model="flashSaleForm.start_time" type="datetime-local" required class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none" />
            </div>
            <div class="space-y-2">
              <label class="text-xs text-secondary uppercase tracking-widest">End Time</label>
              <input v-model="flashSaleForm.end_time" type="datetime-local" required class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none" />
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-xs text-secondary uppercase tracking-widest">Per User Limit</label>
            <input v-model.number="flashSaleForm.per_user_limit" type="number" min="1" required class="w-full bg-black/50 border border-white/10 px-4 py-2 text-white focus:border-accent outline-none" />
          </div>

          <button type="submit" :disabled="loading" class="px-8 py-3 bg-accent text-white font-bold uppercase tracking-widest hover:bg-accent-hover transition-colors flex items-center gap-2">
            <Loader2 v-if="loading" class="w-4 h-4 animate-spin" />
            <Zap v-else class="w-4 h-4" />
            Launch Event
          </button>
        </form>
      </div>
    </div>
  </div>
</template>