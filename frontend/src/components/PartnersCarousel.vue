<template>
  <div class="py-16 bg-surface/20 border-b border-white/5 overflow-hidden">
    <div class="max-w-6xl mx-auto px-4 mb-8">
      <div class="text-center">
        <h2 class="text-2xl font-bold text-white mb-2">Trusted By Industry Leaders</h2>
        <p class="text-secondary text-sm">Partnered with the biggest names in payments and logistics</p>
      </div>
    </div>

    <div class="relative">
      <div class="flex animate-scroll-left">
        <div 
          v-for="(partner, index) in [...partners, ...partners]" 
          :key="index"
          class="flex-shrink-0 px-8 py-4 mx-4 bg-white/5 rounded-lg border border-white/10 hover:border-accent/30 transition-colors cursor-pointer group"
        >
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded bg-white/10 flex items-center justify-center text-xl">
              {{ partner.emoji }}
            </div>
            <div>
              <div class="text-white font-semibold group-hover:text-accent transition-colors">{{ partner.name }}</div>
              <div class="text-xs text-secondary">{{ partner.type }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const partners = [
  { name: 'Visa', emoji: '💳', type: 'Payments' },
  { name: 'Mastercard', emoji: '💳', type: 'Payments' },
  { name: 'Alipay', emoji: '📱', type: 'Payments' },
  { name: 'WeChat Pay', emoji: '💬', type: 'Payments' },
  { name: 'SF Express', emoji: '📦', type: 'Logistics' },
  { name: 'JD Logistics', emoji: '🚚', type: 'Logistics' },
  { name: 'Cloudflare', emoji: '☁️', type: 'Infrastructure' },
  { name: 'AWS', emoji: '🔶', type: 'Infrastructure' },
  { name: 'Google Cloud', emoji: '🌐', type: 'Infrastructure' }
]
</script>

<style scoped>
@keyframes scroll-left {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(-50%);
  }
}

.animate-scroll-left {
  animation: scroll-left 30s linear infinite;
  width: fit-content;
}

.animate-scroll-left:hover {
  animation-play-state: paused;
}
</style>
