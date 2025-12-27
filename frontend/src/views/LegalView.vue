<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Shield, FileText, HelpCircle } from 'lucide-vue-next'

const router = useRouter()
const activeTab = ref<'terms' | 'privacy' | 'support'>('terms')

const tabs = [
  { id: 'terms', label: 'Terms of Service', icon: FileText },
  { id: 'privacy', label: 'Privacy Policy', icon: Shield },
  { id: 'support', label: 'Support Center', icon: HelpCircle }
]
</script>

<template>
  <div class="min-h-screen bg-background pt-24 px-6 md:px-12 pb-20">
    <div class="max-w-4xl mx-auto">
      <button @click="router.push('/')" class="text-sm text-secondary hover:text-white flex items-center gap-2 mb-8 uppercase tracking-widest transition-colors group">
        <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" /> Back to Home
      </button>

      <div class="flex flex-wrap gap-4 mb-12 border-b border-white/10 pb-4">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id as any"
          :class="[
            'flex items-center gap-2 px-6 py-3 text-sm font-bold uppercase tracking-widest transition-all duration-300 rounded-md',
            activeTab === tab.id 
              ? 'bg-accent text-white shadow-[0_0_20px_rgba(227,53,53,0.3)]' 
              : 'text-secondary hover:text-white hover:bg-white/5'
          ]"
        >
          <component :is="tab.icon" class="w-4 h-4" />
          {{ tab.label }}
        </button>
      </div>

      <div class="bg-surface border border-white/5 p-8 md:p-12 rounded-lg relative overflow-hidden">
        <div class="absolute top-0 left-0 w-1 h-full bg-accent"></div>
        
        <div v-if="activeTab === 'terms'" class="space-y-8 animate-fade-in">
          <h1 class="text-3xl font-bold text-white mb-6">Terms of Service</h1>
          <p class="text-secondary leading-relaxed">Last updated: December 2025</p>
          
          <section class="space-y-4">
            <h2 class="text-xl font-bold text-white">1. Acceptance of Terms</h2>
            <p class="text-secondary leading-relaxed">
              By accessing and using MagTrade ("the Platform"), you agree to comply with and be bound by these Terms of Service. The Platform leverages the Obsidian Velocity Protocol to ensure high-frequency trading capabilities. Users must acknowledge that network latency may vary based on geographical location.
            </p>
          </section>

          <section class="space-y-4">
            <h2 class="text-xl font-bold text-white">2. High-Frequency Trading Risks</h2>
            <p class="text-secondary leading-relaxed">
              MagTrade employs advanced edge computing to minimize latency. However, you acknowledge that:
              <br/>a) Flash sale outcomes are determined by millisecond-level timestamps.
              <br/>b) Network jitter on the client side is beyond our control.
              <br/>c) Use of automated bots not authorized by our API protocols may result in immediate account suspension.
            </p>
          </section>

          <section class="space-y-4">
            <h2 class="text-xl font-bold text-white">3. Digital Assets & Settlements</h2>
            <p class="text-secondary leading-relaxed">
              All transactions are final. The Platform ensures atomic swaps for inventory management using Lua scripts on Redis clusters. In the event of a system rollback, refunds will be processed to the original payment method within 5-10 business days.
            </p>
          </section>
        </div>

        <div v-if="activeTab === 'privacy'" class="space-y-8 animate-fade-in">
          <h1 class="text-3xl font-bold text-white mb-6">Privacy Policy</h1>
          
          <section class="space-y-4">
            <h2 class="text-xl font-bold text-white">1. Data Collection</h2>
            <p class="text-secondary leading-relaxed">
              We collect minimal data necessary for transaction processing and anti-fraud analysis. This includes IP addresses, device fingerprints, and transaction logs. Our AI Anomaly Detector analyzes this data in real-time to prevent sybil attacks.
            </p>
          </section>

          <section class="space-y-4">
            <h2 class="text-xl font-bold text-white">2. Data Encryption</h2>
            <p class="text-secondary leading-relaxed">
              All sensitive data is encrypted using AES-256 standards at rest and TLS 1.3 in transit. We do not store raw credit card numbers; all payments are processed through secure third-party gateways compliant with PCI-DSS.
            </p>
          </section>
        </div>

        <div v-if="activeTab === 'support'" class="space-y-8 animate-fade-in">
          <h1 class="text-3xl font-bold text-white mb-6">Support Center</h1>
          <p class="text-secondary text-lg">
            Need help? Our dedicated support team is available 24/7.
          </p>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
            <div class="p-6 bg-surface-light border border-white/5 rounded hover:border-accent/50 transition-colors cursor-pointer group">
              <h3 class="font-bold text-white mb-2 group-hover:text-accent">Live Chat</h3>
              <p class="text-sm text-secondary">Connect with our AI agents instantly for order inquiries.</p>
            </div>
            <div class="p-6 bg-surface-light border border-white/5 rounded hover:border-accent/50 transition-colors cursor-pointer group">
              <h3 class="font-bold text-white mb-2 group-hover:text-accent">Email Support</h3>
              <p class="text-sm text-secondary">chaemsxadmanph@gmail.com</p>
              <p class="text-xs text-tertiary mt-2">Response time: < 2 hours</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>