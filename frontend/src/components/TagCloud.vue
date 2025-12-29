<template>
  <div
    class="bg-black/40 backdrop-blur border border-white/10 rounded-xl p-6 relative overflow-hidden"
  >
    <div
      class="absolute inset-0 bg-gradient-to-b from-accent/5 to-transparent pointer-events-none"
    ></div>

    <h3 class="text-xs font-bold text-gray-500 uppercase tracking-widest mb-4 relative">
      TRENDING SEARCHES
    </h3>

    <div class="flex flex-wrap gap-2 relative">
      <span
        v-for="tag in tags"
        :key="tag"
        class="tag-pill"
        :class="{ active: selectedTag === tag }"
        @click="selectTag(tag)"
      >
        {{ tag }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  select: [tag: string | null]
}>()

const tags = ref([
  'Drone',
  'Camera',
  'Sneakers',
  'Console',
  'Vintage',
  'Crypto',
  'NFT',
  'AirPods',
  'Gaming',
  'Watch',
  'Designer',
  'Limited Edition',
  'GPU',
  'AI',
  'Tech',
])

const selectedTag = ref<string | null>(null)

const selectTag = (tag: string) => {
  if (selectedTag.value === tag) {
    selectedTag.value = null // Toggle off
    emit('select', null)
  } else {
    selectedTag.value = tag
    emit('select', tag)
  }
}
</script>

<style scoped>
.tag-cloud {
  background: #111;
  border: 1px solid #222;
  border-radius: 12px;
  padding: 20px;
  height: 100%;
}

.cloud-header {
  font-size: 12px;
  color: #555;
  font-weight: bold;
  letter-spacing: 1px;
  margin-bottom: 16px;
}

.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-pill {
  font-size: 12px;
  color: #888;
  background: #1a1a1a;
  padding: 6px 12px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid transparent;
}

.tag-pill:hover,
.tag-pill.active {
  background: #e33535;
  color: #fff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(227, 53, 53, 0.3);
}

.tag-pill.active {
  border-color: #ff9999;
}

.hash {
  color: #444;
  margin-right: 2px;
  transition: color 0.3s;
}

.tag-pill:hover .hash,
.tag-pill.active .hash {
  color: rgba(255, 255, 255, 0.5);
}
</style>
