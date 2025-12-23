<template>
  <div class="py-20 bg-gradient-to-b from-black to-surface relative overflow-hidden">
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_bottom,_var(--tw-gradient-stops))] from-accent/10 via-transparent to-transparent"></div>
    
    <div class="max-w-4xl mx-auto px-4 text-center relative z-10">
      <h2 class="text-3xl md:text-5xl font-bold text-white mb-6">
        Ready to Dominate?
      </h2>
      <p class="text-secondary text-lg mb-10 max-w-xl mx-auto">
        Join over 2.5 million traders who trust MagTrade for their flash sale needs. 
        Get started free, no credit card required.
      </p>

      <div class="flex flex-col sm:flex-row gap-4 justify-center items-center mb-12">
        <input 
          v-model="email"
          type="email" 
          placeholder="Enter your email"
          class="w-full sm:w-80 px-6 py-4 bg-black/50 border border-white/20 rounded-lg text-white placeholder-gray-500 focus:border-accent focus:outline-none"
        />
        <button 
          @click="subscribe"
          class="w-full sm:w-auto px-8 py-4 bg-accent text-white font-bold rounded-lg hover:bg-accent/90 transition-colors flex items-center justify-center gap-2"
        >
          <Rocket class="w-5 h-5" />
          Get Early Access
        </button>
      </div>

      <div class="flex flex-wrap justify-center gap-8 text-sm text-secondary">
        <div class="flex items-center gap-2">
          <CheckCircle class="w-4 h-4 text-accent" />
          Free tier available
        </div>
        <div class="flex items-center gap-2">
          <CheckCircle class="w-4 h-4 text-accent" />
          No credit card required
        </div>
        <div class="flex items-center gap-2">
          <CheckCircle class="w-4 h-4 text-accent" />
          Cancel anytime
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Rocket, CheckCircle } from 'lucide-vue-next'

const email = ref('')

const subscribe = () => {
  if (email.value) {
    alert(`Thanks for subscribing! We'll reach out to ${email.value} soon.`)
    email.value = ''
  }
}
</script>
