import http from '@/utils/http'
import type { LoginResponse, User, ApiResponse } from '@/types'

// 登录
export const login = (data: any) => {
  return http.post<any, ApiResponse<LoginResponse>>('/auth/login', data)
}

// 注册
export const register = (data: any) => {
  return http.post<any, ApiResponse<LoginResponse>>('/auth/register', data)
}

// 获取当前用户信息
export const getMe = () => {
  return http.get<any, ApiResponse<User>>('/auth/me')
}

// 获取验证码图片 (返回 Blob)
export const getCaptcha = (identifier: string) => {
  return http.get(`/captcha?identifier=${identifier}`, { responseType: 'blob' })
}

// 检查是否需要验证码
export const checkCaptcha = (identifier: string) => {
  return http.get<any, { needs_captcha: boolean; is_locked: boolean }>(`/captcha/check?identifier=${identifier}`)
}