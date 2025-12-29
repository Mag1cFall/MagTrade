import http from '@/utils/http'
import type { LoginResponse, User, ApiResponse } from '@/types'

export const login = (data: any) => {
  return http.post<any, ApiResponse<LoginResponse>>('/auth/login', data)
}

export const register = (data: any) => {
  return http.post<any, ApiResponse<LoginResponse>>('/auth/register', data)
}

export const getMe = () => {
  return http.get<any, ApiResponse<User>>('/auth/me')
}

export const sendEmailCode = (email: string) => {
  return http.post<any, ApiResponse<{ message: string }>>('/auth/send-code', { email })
}

export const getCaptcha = (identifier: string) => {
  return http.get(`/captcha?identifier=${identifier}`, { responseType: 'blob' })
}

export const checkCaptcha = (identifier: string) => {
  return http.get<any, { needs_captcha: boolean; is_locked: boolean }>(
    `/captcha/check?identifier=${identifier}`
  )
}
