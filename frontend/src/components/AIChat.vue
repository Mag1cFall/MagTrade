<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { sendChatMessageStream, type StreamChunk } from '@/api/ai'
import { useAuthStore } from '@/stores/auth'
import { Bot, X, Send, Terminal, Loader2, Brain } from 'lucide-vue-next'
import { v4 as uuidv4 } from 'uuid'

const authStore = useAuthStore()
const isOpen = ref(false)
const input = ref('')
const loading = ref(false)
const thinkingContent = ref('')
const messages = ref<{ role: 'user' | 'assistant' | 'thinking'; content: string }[]>([
  { role: 'assistant', content: '我是 MagTrade 智能战术助手。关于秒杀时机、库存状态或订单追踪，随时下达指令。' }
])
const messagesContainer = ref<HTMLElement | null>(null)

const sessionId = ref(localStorage.getItem('ai_session_id') || uuidv4())
localStorage.setItem('ai_session_id', sessionId.value)

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const toggleChat = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) scrollToBottom()
}

const handleSend = async () => {
  if (!input.value.trim() || loading.value) return
  
  if (!authStore.isAuthenticated) {
    messages.value.push({ role: 'assistant', content: '权限不足。请先登录以访问 AI 核心。' })
    scrollToBottom()
    return
  }

  const userMsg = input.value
  messages.value.push({ role: 'user', content: userMsg })
  input.value = ''
  loading.value = true
  thinkingContent.value = ''
  scrollToBottom()

  let assistantContent = ''
  let thinkingIdx = -1
  let assistantIdx = -1

  try {
    await sendChatMessageStream(
      sessionId.value, 
      userMsg, 
      authStore.token || '',
      (chunk: StreamChunk) => {
        if (chunk.type === 'thinking') {
          thinkingContent.value += chunk.content
          if (thinkingIdx === -1) {
            thinkingIdx = messages.value.length
            messages.value.push({ role: 'thinking', content: thinkingContent.value })
          } else {
            messages.value[thinkingIdx]!.content = thinkingContent.value
          }
        } else if (chunk.type === 'content') {
          assistantContent += chunk.content
          if (assistantIdx === -1) {
            assistantIdx = messages.value.length
            messages.value.push({ role: 'assistant', content: assistantContent })
          } else {
            messages.value[assistantIdx]!.content = assistantContent
          }
        }
        scrollToBottom()
      }
    )
  } catch (e) {
    messages.value.push({ role: 'assistant', content: '网络连接异常，请检查网络。' })
  } finally {
    loading.value = false
    thinkingContent.value = ''
    scrollToBottom()
  }
}
</script>

<template>
  <div class="fixed bottom-6 right-6 z-50 flex flex-col items-end">
    <transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 translate-y-4 scale-95"
      enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 translate-y-4 scale-95"
    >
      <div v-if="isOpen" class="mb-4 w-80 md:w-96 h-[500px] bg-black/90 backdrop-blur-xl border border-white/10 flex flex-col shadow-2xl shadow-black overflow-hidden rounded-lg">
        <div class="p-4 border-b border-white/10 flex justify-between items-center bg-white/5">
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span class="text-sm font-bold tracking-widest text-white">AI CORE ONLINE</span>
          </div>
          <button @click="toggleChat" class="text-secondary hover:text-white transition-colors">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4 scrollbar-hide">
          <div v-for="(msg, idx) in messages" :key="idx" :class="['flex', msg.role === 'user' ? 'justify-end' : 'justify-start']">
            <div :class="[
              'max-w-[85%] p-3 text-sm leading-relaxed border',
              msg.role === 'user' 
                ? 'bg-accent text-white border-accent' 
                : msg.role === 'thinking'
                  ? 'bg-purple-900/30 text-purple-200 border-purple-500/30 italic'
                  : 'bg-surface-light text-secondary border-white/5'
            ]">
              <div v-if="msg.role === 'thinking'" class="flex items-center gap-2 mb-1 text-xs text-purple-400 uppercase tracking-wider">
                <Brain class="w-3 h-3" /> Thinking...
              </div>
              <div v-else-if="msg.role === 'assistant'" class="flex items-center gap-2 mb-1 text-xs text-tertiary uppercase tracking-wider">
                <Terminal class="w-3 h-3" /> System
              </div>
              {{ msg.content }}
            </div>
          </div>
          <div v-if="loading && !thinkingContent && messages[messages.length-1]?.role === 'user'" class="flex justify-start">
            <div class="bg-surface-light p-3 border border-white/5">
              <Loader2 class="w-4 h-4 text-accent animate-spin" />
            </div>
          </div>
        </div>

        <div class="p-4 border-t border-white/10 bg-white/5">
          <div class="relative">
            <input 
              v-model="input" 
              @keyup.enter="handleSend"
              type="text" 
              placeholder="Enter command..." 
              class="w-full bg-black/50 border border-white/10 py-2 pl-4 pr-10 text-white text-sm focus:outline-none focus:border-accent transition-colors placeholder:text-tertiary"
            />
            <button 
              @click="handleSend"
              :disabled="!input.trim() || loading"
              class="absolute right-2 top-1/2 -translate-y-1/2 text-secondary hover:text-accent disabled:opacity-50 transition-colors"
            >
              <Send class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </transition>

    <button 
      @click="toggleChat"
      class="w-14 h-14 bg-surface hover:bg-surface-light border border-white/10 rounded-full flex items-center justify-center shadow-lg transition-all duration-300 group hover:border-accent/50"
    >
      <Bot class="w-6 h-6 text-white group-hover:text-accent transition-colors" />
      <span v-if="!isOpen" class="absolute top-0 right-0 w-3 h-3 bg-red-500 rounded-full animate-ping"></span>
      <span v-if="!isOpen" class="absolute top-0 right-0 w-3 h-3 bg-red-500 rounded-full"></span>
    </button>
  </div>
</template>