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

// 获取策略推荐
export const getRecommendation = (flashSaleId: number) => {
  return http.get<any, ApiResponse<{ analysis: AIRecommendation }>>(`/ai/recommendations/${flashSaleId}`)
}

// 发送聊天消息
export const sendChatMessage = (sessionId: string, message: string) => {
  return http.post<any, ApiResponse<ChatResponse>>('/ai/chat', { session_id: sessionId, message })
}

// 获取聊天历史
export const getChatHistory = (sessionId: string) => {
  return http.get<any, ApiResponse<{ history: any[] }>>(`/ai/chat/history`, { params: { session_id: sessionId } })
}