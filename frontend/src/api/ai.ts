import http from '@/utils/http'
import type { ApiResponse } from '@/types'

export interface AIRecommendation {
  difficulty_score: number
  timing_advice: string
  success_probability: number
  recommendations: string[]
  difficulty_reason: string
}

export interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
}

export interface ChatResponse {
  session_id: string
  response: string
  related_data?: any
}

export interface StreamChunk {
  type: 'thinking' | 'content' | 'done' | 'error'
  content: string
  done?: boolean
}

export const getRecommendation = (flashSaleId: number) => {
  return http.get<any, ApiResponse<{ analysis: AIRecommendation }>>(`/ai/recommendations/${flashSaleId}`)
}

export const sendChatMessage = (sessionId: string, message: string) => {
  return http.post<any, ApiResponse<ChatResponse>>('/ai/chat', { session_id: sessionId, message })
}

export const sendChatMessageStream = async (
  sessionId: string,
  message: string,
  token: string,
  onChunk: (chunk: StreamChunk) => void
): Promise<void> => {
  const response = await fetch('/api/v1/ai/chat/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ session_id: sessionId, message })
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const reader = response.body?.getReader()
  if (!reader) {
    throw new Error('No reader available')
  }

  const decoder = new TextDecoder()
  let buffer = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })
    const lines = buffer.split('\n')
    buffer = lines.pop() || ''

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        try {
          const data = JSON.parse(line.slice(6)) as StreamChunk
          onChunk(data)
        } catch (e) {
        }
      }
    }
  }
}

export const getChatHistory = (sessionId: string) => {
  return http.get<any, ApiResponse<{ history: any[] }>>(`/ai/chat/history`, { params: { session_id: sessionId } })
}