<template>
  <div class="min-h-screen bg-background">
    <div class="max-w-4xl mx-auto px-6 py-20">
      <h1 class="text-4xl font-bold text-white mb-8">Contact Support</h1>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-12">
        <div class="space-y-6">
          <p class="text-secondary text-lg">
            Have questions or need help? Reach out to our support team.
          </p>

          <div class="space-y-4">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-full bg-accent/20 flex items-center justify-center">
                <Mail class="w-5 h-5 text-accent" />
              </div>
              <div>
                <div class="text-sm text-secondary">Email</div>
                <div class="text-white font-mono">support@magtrade.io</div>
              </div>
            </div>

            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-full bg-accent/20 flex items-center justify-center">
                <MessageCircle class="w-5 h-5 text-accent" />
              </div>
              <div>
                <div class="text-sm text-secondary">Live Chat</div>
                <div class="text-white">Available 24/7</div>
              </div>
            </div>

            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-full bg-accent/20 flex items-center justify-center">
                <Clock class="w-5 h-5 text-accent" />
              </div>
              <div>
                <div class="text-sm text-secondary">Response Time</div>
                <div class="text-white">Usually within 2 hours</div>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-surface/50 border border-white/10 rounded-xl p-6">
          <h2 class="text-xl font-semibold text-white mb-6">Send a Message</h2>

          <form class="space-y-4" @submit.prevent="submitForm">
            <div>
              <label class="block text-sm text-secondary mb-2">Name</label>
              <input
                v-model="form.name"
                type="text"
                class="w-full px-4 py-3 bg-black/50 border border-white/10 rounded text-white focus:border-accent focus:outline-none"
                placeholder="Your name"
              />
            </div>

            <div>
              <label class="block text-sm text-secondary mb-2">Email</label>
              <input
                v-model="form.email"
                type="email"
                class="w-full px-4 py-3 bg-black/50 border border-white/10 rounded text-white focus:border-accent focus:outline-none"
                placeholder="your@email.com"
              />
            </div>

            <div>
              <label class="block text-sm text-secondary mb-2">Message</label>
              <textarea
                v-model="form.message"
                rows="4"
                class="w-full px-4 py-3 bg-black/50 border border-white/10 rounded text-white focus:border-accent focus:outline-none resize-none"
                placeholder="How can we help?"
              ></textarea>
            </div>

            <button
              type="submit"
              class="w-full py-3 bg-accent text-white font-semibold rounded hover:bg-accent/90 transition-colors"
            >
              Send Message
            </button>
          </form>
        </div>
      </div>

      <div class="mt-12">
        <router-link to="/" class="text-accent hover:underline">← Back to Home</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { Mail, MessageCircle, Clock } from 'lucide-vue-next'

const form = reactive({
  name: '',
  email: '',
  message: '',
})

const submitForm = () => {
  alert('Message sent! We will get back to you soon.')
  form.name = ''
  form.email = ''
  form.message = ''
}
</script>
