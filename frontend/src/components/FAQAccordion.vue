<template>
  <div class="py-20 border-b border-white/5">
    <div class="max-w-3xl mx-auto px-4">
      <div class="text-center mb-12">
        <h2 class="text-3xl md:text-5xl font-bold text-white mb-4">Frequently Asked Questions</h2>
        <p class="text-secondary">Everything you need to know about MagTrade</p>
      </div>

      <div class="space-y-4">
        <div 
          v-for="(faq, index) in faqs" 
          :key="index"
          class="border border-white/10 rounded-xl overflow-hidden bg-surface/30"
        >
          <button
            @click="toggleFaq(index)"
            class="w-full px-6 py-4 flex items-center justify-between text-left hover:bg-white/5 transition-colors"
          >
            <span class="font-semibold text-white">{{ faq.question }}</span>
            <ChevronDown 
              class="w-5 h-5 text-secondary transition-transform duration-300"
              :class="{ 'rotate-180': openIndex === index }"
            />
          </button>
          
          <Transition name="slide">
            <div v-if="openIndex === index" class="px-6 pb-4">
              <p class="text-secondary leading-relaxed">{{ faq.answer }}</p>
            </div>
          </Transition>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ChevronDown } from 'lucide-vue-next'

const openIndex = ref<number | null>(0)

const toggleFaq = (index: number) => {
  openIndex.value = openIndex.value === index ? null : index
}

const faqs = [
  {
    question: 'How does the flash sale sniping work?',
    answer: 'Our system uses distributed edge nodes to minimize latency between your click and our matching engine. When you hit the "snipe" button, your order is processed in under 10ms through our Obsidian engine, giving you the fastest possible execution.'
  },
  {
    question: 'Is MagTrade legal to use?',
    answer: 'Yes, MagTrade is completely legal. We simply provide a faster, more reliable interface for participating in flash sales. We do not use bots or automated scripts that violate platform terms of service.'
  },
  {
    question: 'What payment methods do you accept?',
    answer: 'We accept all major credit cards (Visa, Mastercard, Amex), as well as Alipay, WeChat Pay, and cryptocurrency payments including USDT and ETH.'
  },
  {
    question: 'How does the AI advisor help me?',
    answer: 'Our AI analyzes historical data, current demand signals, and market conditions to give you a difficulty score (1-10), optimal timing advice, and success probability percentage for each drop.'
  },
  {
    question: 'What happens if a transaction fails?',
    answer: 'Failed transactions are automatically refunded within 24 hours. Our system also tracks partial failures and can auto-retry in certain scenarios (configurable in settings).'
  },
  {
    question: 'Is there a mobile app?',
    answer: 'Our mobile app is currently in beta for iOS and Android. Sign up for the waitlist to get early access. The web app is fully responsive and works great on mobile browsers.'
  },
  {
    question: 'How do I contact support?',
    answer: 'You can reach our 24/7 support team via live chat on the website, email at support@magtrade.io, or through our Discord community where team members are always active.'
  },
  {
    question: 'Are there any subscription fees?',
    answer: 'MagTrade offers a free tier with basic features. Pro ($29/mo) unlocks AI advisor, priority queue, and advanced analytics. Enterprise plans are available for high-volume traders.'
  }
]
</script>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
}

.slide-enter-to,
.slide-leave-from {
  max-height: 200px;
  opacity: 1;
}
</style>
