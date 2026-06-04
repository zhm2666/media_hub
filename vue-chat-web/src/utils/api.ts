import axios, { type AxiosInstance, type AxiosError } from 'axios'
import type { ApiResponse, ChatResponse, User, LoginResponse } from '@/types'
import router from '@/router'

const api: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 300000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('chat_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => {
    return response
  },
  (error: AxiosError<ApiResponse>) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('chat_token')
      localStorage.removeItem('chat_user')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export const chatApi = {
  sendMessage: (data: { prompt: string; systemMessage?: string }) => {
    return api.post<any>('/chat-process', data)
  },

  getConfig: () => api.get<ApiResponse<{ model: string }>>('/config'),
}

export const authApi = {
  checkToken: (accessToken: string) => {
    return axios.get(`http://localhost:8082/api/v1/login/check/auth?access_token=${accessToken}`)
  },
}

export default api
