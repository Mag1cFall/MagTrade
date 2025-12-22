import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { getMe, login as apiLogin, register as apiRegister } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('access_token'))
  const user = ref<User | null>(null)
  
  // 从 localStorage 恢复用户信息（如果有）
  const storedUser = localStorage.getItem('user_info')
  if (storedUser) {
    try {
      user.value = JSON.parse(storedUser)
    } catch (e) {
      localStorage.removeItem('user_info')
    }
  }

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  const setToken = (accessToken: string, refreshToken: string) => {
    token.value = accessToken
    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('refresh_token', refreshToken)
  }

  const setUser = (userData: User) => {
    user.value = userData
    localStorage.setItem('user_info', JSON.stringify(userData))
  }

  const login = async (loginForm: any) => {
    const res = await apiLogin(loginForm)
    if (res.code === 0) {
      setToken(res.data.access_token, res.data.refresh_token)
      setUser(res.data.user)
    }
    return res
  }

  const register = async (registerForm: any) => {
    const res = await apiRegister(registerForm)
    if (res.code === 0) {
      setToken(res.data.access_token, res.data.refresh_token)
      setUser(res.data.user)
    }
    return res
  }

  const fetchUser = async () => {
    if (!token.value) return
    try {
      const res = await getMe()
      if (res.code === 0) {
        setUser(res.data)
      }
    } catch (e) {
      // 这里的错误通常由 http interceptor 处理，但如果是初始加载失败，可以清除状态
      if (!user.value) {
        token.value = null
        localStorage.removeItem('access_token')
      }
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_info')
    window.location.href = '/login'
  }

  return {
    token,
    user,
    isAuthenticated,
    isAdmin,
    login,
    register,
    fetchUser,
    logout
  }
})