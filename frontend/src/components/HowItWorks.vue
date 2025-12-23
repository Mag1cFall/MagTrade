<template>
  <div class="py-20 border-b border-white/5">
    <div class="max-w-6xl mx-auto px-4">
      <div class="text-center mb-16">
        <h2 class="text-3xl md:text-5xl font-bold text-white mb-4">How It Works</h2>
        <p class="text-secondary max-w-xl mx-auto">
          Get started in three simple steps and dominate every flash sale.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div 
          v-for="(step, index) in steps" 
          :key="index"
          class="relative group"
        >
          <div class="absolute -top-4 -left-4 w-16 h-16 rounded-full bg-accent/10 flex items-center justify-center text-accent text-2xl font-bold border border-accent/30 group-hover:scale-110 transition-transform">
            {{ index + 1 }}
          </div>
          
          <div class="bg-surface/50 border border-white/10 rounded-xl p-8 pt-12 hover:border-accent/30 transition-colors h-full">
            <component :is="step.icon" class="w-10 h-10 text-accent mb-4" />
            <h3 class="text-xl font-bold text-white mb-3">{{ step.title }}</h3>
            <p class="text-secondary">{{ step.description }}</p>
          </div>

          <div v-if="index < 2" class="hidden md:block absolute top-1/2 -right-4 text-white/20">
            <ChevronRight class="w-8 h-8" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { UserPlus, Search, Zap, ChevronRight } from 'lucide-vue-next'

const steps = [
  {
    icon: UserPlus,
    title: 'Create Account',
    description: 'Sign up in seconds. Verify your email and set up your payment method for lightning-fast checkout.'
  },
  {
    icon: Search,
    title: 'Browse & Research',
    description: 'Explore upcoming drops. Use our AI advisor to analyze difficulty and get timing recommendations.'
  },
  {
    icon: Zap,
    title: 'Snipe & Win',
    description: 'When the countdown hits zero, our ultra-low-latency engine gives you the edge to secure your items.'
  }
]
</script>
