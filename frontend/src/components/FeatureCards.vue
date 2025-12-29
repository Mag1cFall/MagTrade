<template>
  <div class="py-20 border-b border-white/5">
    <div class="max-w-6xl mx-auto px-4">
      <div class="text-center mb-16">
        <h2 class="text-3xl md:text-5xl font-bold text-white mb-4">Powerful Features</h2>
        <p class="text-secondary max-w-xl mx-auto">
          Every tool you need to outperform the competition
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="feature in features"
          :key="feature.title"
          class="group bg-surface/30 border border-white/5 rounded-xl p-6 hover:border-accent/50 hover:bg-surface/50 transition-all duration-300"
        >
          <div
            class="w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center mb-4 group-hover:bg-accent/20 transition-colors"
          >
            <component :is="feature.icon" class="w-6 h-6 text-accent" />
          </div>
          <h3 class="text-lg font-bold text-white mb-2">{{ feature.title }}</h3>
          <p class="text-secondary text-sm leading-relaxed">{{ feature.description }}</p>

          <div v-if="feature.badge" class="mt-4">
            <span class="px-2 py-1 text-xs bg-accent/20 text-accent rounded">{{
              feature.badge
            }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Zap, Shield, Brain, Globe, Clock, TrendingUp } from 'lucide-vue-next'

const features = [
  {
    icon: Zap,
    title: 'Lightning Execution',
    description:
      'Sub-10ms order processing with our Obsidian matching engine. Be first in line, every time.',
    badge: '< 10ms',
  },
  {
    icon: Brain,
    title: 'AI-Powered Analysis',
    description:
      'Get real-time difficulty scores, timing advice, and success probability predictions for every drop.',
    badge: 'GPT-4 Powered',
  },
  {
    icon: Shield,
    title: 'Bank-Grade Security',
    description:
      'End-to-end encryption, 2FA, and cold storage for assets. Audited by leading security firms.',
    badge: null,
  },
  {
    icon: Globe,
    title: 'Global Edge Network',
    description: '35+ edge nodes worldwide ensure minimal latency regardless of your location.',
    badge: null,
  },
  {
    icon: Clock,
    title: 'Precision Timing',
    description:
      'Atomic clock synchronization ensures your countdown is accurate to the millisecond.',
    badge: null,
  },
  {
    icon: TrendingUp,
    title: 'Real-Time Analytics',
    description:
      'Live market data, competitor analysis, and historical success rates all in one dashboard.',
    badge: 'New',
  },
]
</script>
