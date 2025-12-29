<template>
  <div class="bg-surface/50 border border-white/10 rounded-xl p-6">
    <div class="flex items-center gap-3 mb-6">
      <div class="w-10 h-10 rounded-lg bg-accent/20 flex items-center justify-center">
        <Wand2 class="w-5 h-5 text-accent" />
      </div>
      <div>
        <h3 class="text-lg font-bold text-white">AI 商品生成器</h3>
        <p class="text-xs text-secondary">Let AI help you create product listings</p>
      </div>
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm text-secondary mb-2">Product Name / Keywords</label>
        <input
          v-model="productName"
          type="text"
          placeholder="e.g., iPhone 15 Pro Max, Limited Edition Sneakers"
          class="w-full px-4 py-3 bg-black/50 border border-white/10 rounded-lg text-white placeholder-gray-500 focus:border-accent focus:outline-none"
        />
      </div>

      <div>
        <label class="block text-sm text-secondary mb-2">Category</label>
        <select
          v-model="category"
          class="w-full px-4 py-3 bg-black/50 border border-white/10 rounded-lg text-white focus:border-accent focus:outline-none"
        >
          <option value="">Select category</option>
          <option value="electronics">Electronics</option>
          <option value="fashion">Fashion</option>
          <option value="home">Home & Living</option>
          <option value="sports">Sports</option>
          <option value="collectibles">Collectibles</option>
        </select>
      </div>

      <button
        :disabled="!productName || isGenerating"
        class="w-full py-3 bg-accent text-white font-semibold rounded-lg hover:bg-accent/90 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
        @click="generateContent"
      >
        <Loader2 v-if="isGenerating" class="w-4 h-4 animate-spin" />
        <Sparkles v-else class="w-4 h-4" />
        {{ isGenerating ? 'Generating...' : 'Generate with AI' }}
      </button>

      <Transition name="fade">
        <div v-if="generatedContent" class="mt-6 space-y-4">
          <div class="p-4 bg-black/30 rounded-lg border border-accent/30">
            <h4 class="text-sm font-semibold text-accent mb-2">Generated Title</h4>
            <p class="text-white">{{ generatedContent.title }}</p>
          </div>

          <div class="p-4 bg-black/30 rounded-lg border border-white/10">
            <h4 class="text-sm font-semibold text-secondary mb-2">Description</h4>
            <p class="text-white text-sm">{{ generatedContent.description }}</p>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="p-4 bg-black/30 rounded-lg border border-white/10">
              <h4 class="text-sm font-semibold text-secondary mb-2">Suggested Price</h4>
              <p class="text-white font-bold text-xl">¥{{ generatedContent.price }}</p>
            </div>
            <div class="p-4 bg-black/30 rounded-lg border border-white/10">
              <h4 class="text-sm font-semibold text-secondary mb-2">Flash Price</h4>
              <p class="text-accent font-bold text-xl">¥{{ generatedContent.flashPrice }}</p>
            </div>
          </div>

          <div class="p-4 bg-black/30 rounded-lg border border-white/10">
            <h4 class="text-sm font-semibold text-secondary mb-2">Suggested Tags</h4>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="tag in generatedContent.tags"
                :key="tag"
                class="px-2 py-1 bg-white/10 text-white text-xs rounded"
              >
                {{ tag }}
              </span>
            </div>
          </div>

          <button
            class="w-full py-3 bg-white text-black font-semibold rounded-lg hover:bg-white/90 transition-colors"
          >
            Use This Content
          </button>
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Wand2, Sparkles, Loader2 } from 'lucide-vue-next'

const productName = ref('')
const category = ref('')
const isGenerating = ref(false)

interface GeneratedContent {
  title: string
  description: string
  price: number
  flashPrice: number
  tags: string[]
}

const generatedContent = ref<GeneratedContent | null>(null)

const generateContent = async () => {
  isGenerating.value = true
  generatedContent.value = null

  await new Promise((resolve) => setTimeout(resolve, 2000))

  const templates: Record<string, GeneratedContent> = {
    electronics: {
      title: `Premium ${productName.value} - Limited Edition`,
      description: `Experience cutting-edge technology with this ${productName.value}. Features the latest specifications, premium build quality, and exceptional performance. Perfect for tech enthusiasts who demand the best. Comes with full warranty and exclusive accessories.`,
      price: 8999,
      flashPrice: 6999,
      tags: ['Electronics', 'Premium', 'Limited', 'Hot Sale', 'New Arrival'],
    },
    fashion: {
      title: `Designer ${productName.value} - Exclusive Drop`,
      description: `Make a statement with this exclusive ${productName.value}. Crafted with premium materials and attention to detail. Limited quantity available. Authentic guarantee with certificate of authenticity included.`,
      price: 2999,
      flashPrice: 1999,
      tags: ['Fashion', 'Designer', 'Exclusive', 'Limited Edition', 'Trending'],
    },
    default: {
      title: `${productName.value} - Flash Sale Special`,
      description: `Don't miss this incredible deal on ${productName.value}. High quality, great value, and limited availability. Order now before it's gone! Fast shipping and hassle-free returns.`,
      price: 1999,
      flashPrice: 999,
      tags: ['Flash Sale', 'Hot Deal', 'Limited Stock', 'Best Seller'],
    },
  }

  generatedContent.value = templates[category.value] || templates.default
  isGenerating.value = false
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
