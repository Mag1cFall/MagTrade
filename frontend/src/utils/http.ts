import axios from 'axios'
import type { AxiosInstance, AxiosResponse } from 'axios'

const http: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
http.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Queue for requests waiting for token refresh
let isRefreshing = false
let requestsQueue: ((token: string) => void)[] = []

// Response interceptor
http.interceptors.response.use(
  (response: AxiosResponse) => {
    if (response.config.responseType === 'blob' || response.config.responseType === 'arraybuffer') {
      return response
    }

    const { code, message } = response.data

    // Business logic error handling
    if (code !== undefined && code !== 0) {
      return Promise.reject(new Error(message || 'Unknown Error'))
    }

    return response.data
  },
  async (error) => {
    const originalRequest = error.config

    // 排除登录接口的 401 错误，交由登录组件处理密码错误
    if (originalRequest.url.includes('/auth/login')) {
      return Promise.reject(error)
    }

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve) => {
          requestsQueue.push((token: string) => {
            originalRequest.headers.Authorization = `Bearer ${token}`
            resolve(http(originalRequest))
          })
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const refreshToken = localStorage.getItem('refresh_token')
        if (!refreshToken) {
          throw new Error('No refresh token')
        }

        // Call refresh token endpoint directly using axios to avoid interceptor loop
        const response = await axios.post('/api/v1/auth/refresh', {
          refresh_token: refreshToken,
        })

        const { access_token, refresh_token: new_refresh_token } = response.data.data

        localStorage.setItem('access_token', access_token)
        if (new_refresh_token) {
          localStorage.setItem('refresh_token', new_refresh_token)
        }

        // Process queue
        requestsQueue.forEach((cb) => cb(access_token))
        requestsQueue = []

        originalRequest.headers.Authorization = `Bearer ${access_token}`
        return http(originalRequest)
      } catch (refreshError) {
        // Refresh failed, logout
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user_info')

        if (!window.location.pathname.includes('/login')) {
          window.location.href = '/login'
        }
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }
    return Promise.reject(error)
  }
)

export default http
