<template>
  <Teleport to="body">
    <div 
      class="fixed bottom-6 right-6 z-50 flex flex-col items-end"
      :class="{ 'pointer-events-none': !isOpen }"
    >
      <Transition name="chat">
        <div 
          v-if="isOpen"
          class="w-96 h-[500px] bg-surface border border-white/10 rounded-2xl shadow-2xl flex flex-col overflow-hidden mb-4 pointer-events-auto"
        >
          <div class="p-4 border-b border-white/10 flex items-center justify-between bg-black/50">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center">
                <Bot class="w-4 h-4 text-accent" />
              </div>
              <div>
                <div class="text-white font-semibold text-sm">AI 助手</div>
                <div class="text-xs text-green-400 flex items-center gap-1">
                  <span class="w-1.5 h-1.5 bg-green-400 rounded-full animate-pulse"></span>
                  Online
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <button @click="clearChat" class="text-secondary hover:text-accent p-1" title="清空对话">
                <Trash2 class="w-4 h-4" />
              </button>
              <button @click="isOpen = false" class="text-secondary hover:text-white p-1">
                <X class="w-5 h-5" />
              </button>
            </div>
          </div>

          <div ref="messagesRef" class="flex-grow overflow-y-auto p-4 space-y-4">
            <div 
              v-for="(msg, index) in messages" 
              :key="index"
              class="flex"
              :class="msg.role === 'user' ? 'justify-end' : 'justify-start'"
            >
              <div 
                class="max-w-[85%] px-4 py-2 rounded-xl text-sm"
                :class="msg.role === 'user' 
                  ? 'bg-accent text-white rounded-br-sm' 
                  : 'bg-white/10 text-white rounded-bl-sm'"
              >
                <div v-if="msg.role === 'assistant' && msg.thinking" class="text-xs text-secondary mb-1 italic">
                  Thinking...
                </div>
                <div 
                  v-if="msg.role === 'assistant'" 
                  class="markdown-content prose prose-sm prose-invert max-w-none"
                  v-html="renderMarkdown(msg.content)"
                ></div>
                <span v-else>{{ msg.content }}</span>
              </div>
            </div>
            
            <div v-if="isLoading" class="flex justify-start">
              <div class="bg-white/10 px-4 py-2 rounded-xl text-sm text-white rounded-bl-sm">
                <Loader2 class="w-4 h-4 animate-spin" />
              </div>
            </div>
          </div>

          <div class="p-4 border-t border-white/10 bg-black/30">
            <form @submit.prevent="sendMessage" class="flex gap-2">
              <input 
                v-model="input"
                type="text"
                placeholder="Ask about products, strategies..."
                class="flex-grow px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white text-sm placeholder-gray-500 focus:border-accent focus:outline-none"
                :disabled="isLoading"
              />
              <button 
                type="submit"
                :disabled="!input.trim() || isLoading"
                class="px-4 py-2 bg-accent text-white rounded-lg hover:bg-accent/90 transition-colors disabled:opacity-50"
              >
                <Send class="w-4 h-4" />
              </button>
            </form>
          </div>
        </div>
      </Transition>

      <button 
        @click="isOpen = !isOpen"
        class="w-14 h-14 bg-accent rounded-full flex items-center justify-center shadow-lg hover:scale-110 transition-transform pointer-events-auto"
        :class="{ 'ring-2 ring-accent ring-offset-2 ring-offset-background': hasNewMessage }"
      >
        <Bot v-if="!isOpen" class="w-6 h-6 text-white" />
        <X v-else class="w-6 h-6 text-white" />
      </button>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { Bot, X, Send, Loader2, Trash2 } from 'lucide-vue-next'
import { sendChatMessageStream } from '@/api/ai'
import { useAuthStore } from '@/stores/auth'
import { marked } from 'marked'

interface Message {
  role: 'user' | 'assistant'
  content: string
  thinking?: boolean
}

marked.setOptions({
  breaks: true,
  gfm: true
})

const renderMarkdown = (content: string) => {
  if (!content) return ''
  return marked.parse(content)
}

const isOpen = ref(false)
const input = ref('')
const isLoading = ref(false)
const hasNewMessage = ref(false)
const messagesRef = ref<HTMLElement | null>(null)

const sessionId = ref(`session_${Date.now()}`)

const messages = ref<Message[]>([
  { role: 'assistant', content: 'Hi! I\'m your AI shopping assistant. Ask me about flash sales, product recommendations, or trading strategies!' }
])

const authStore = useAuthStore()

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}

watch(messages, scrollToBottom, { deep: true })

const sendMessage = async () => {
  if (!input.value.trim() || isLoading.value) return
  
  const userMessage = input.value.trim()
  messages.value.push({ role: 'user', content: userMessage })
  input.value = ''
  isLoading.value = true

  try {
    let assistantContent = ''
    messages.value.push({ role: 'assistant', content: '', thinking: true })
    
    await sendChatMessageStream(
      sessionId.value,
      userMessage,
      authStore.token || '',
      (chunk) => {
        if (chunk.type === 'content') {
          assistantContent += chunk.content
          messages.value[messages.value.length - 1].content = assistantContent
          messages.value[messages.value.length - 1].thinking = false
        } else if (chunk.type === 'thinking') {
          messages.value[messages.value.length - 1].thinking = true
        }
      }
    )
    
    if (!assistantContent) {
      messages.value[messages.value.length - 1].content = 'I apologize, but I encountered an issue. Please try again.'
    }
  } catch (error) {
    messages.value.push({ 
      role: 'assistant', 
      content: 'Sorry, I\'m having trouble connecting. Please check your login status and try again.' 
    })
  } finally {
    isLoading.value = false
    messages.value[messages.value.length - 1].thinking = false
  }
}

const clearChat = () => {
  sessionId.value = `session_${Date.now()}`
  messages.value = [
    { role: 'assistant', content: 'Hi! I\'m your AI shopping assistant. Ask me about flash sales, product recommendations, or trading strategies!' }
  ]
}
</script>

<style scoped>
.chat-enter-active,
.chat-leave-active {
  transition: all 0.3s ease;
}

.chat-enter-from,
.chat-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

.markdown-content :deep(p) {
  margin: 0.25em 0;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin: 0.5em 0;
  padding-left: 1.5em;
}

.markdown-content :deep(code) {
  background: rgba(0,0,0,0.3);
  padding: 0.1em 0.3em;
  border-radius: 4px;
  font-size: 0.9em;
}

.markdown-content :deep(pre) {
  background: rgba(0,0,0,0.4);
  padding: 0.5em;
  border-radius: 6px;
  overflow-x: auto;
  margin: 0.5em 0;
}

.markdown-content :deep(pre code) {
  background: none;
  padding: 0;
}

.markdown-content :deep(strong) {
  color: #fff;
}

.markdown-content :deep(a) {
  color: #e33535;
  text-decoration: underline;
}
</style>
