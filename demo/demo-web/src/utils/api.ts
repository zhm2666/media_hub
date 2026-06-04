import axios, { type AxiosInstance, type AxiosError } from 'axios'
import type { AuthResponse, LoginRequest, RegisterRequest, User, ApiResponse } from '@/types'
import router from '@/router'

const api: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

const externalApi = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
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
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  register: (data: RegisterRequest) => api.post<AuthResponse>('/auth/register', data),

  login: (data: LoginRequest) => api.post<AuthResponse>('/auth/login', data),

  logout: () => api.post<ApiResponse>('/auth/logout'),

  getCurrentUser: () => api.get<User>('/auth/me'),
}

export const userApi = {
  getProfile: () => api.get<User>('/user/profile'),

  updateProfile: (data: { email?: string }) => api.put<User>('/user/profile', data),
}

export const userCenterApi = {
  getLoginMethods: () => externalApi.get('/uc/login/methods'),

  getWxQrcode: () => externalApi.get('/uc/login/qrcode'),

  checkScanStatus: (ticket: string) => externalApi.get('/uc/login/check-scan', {
    params: { ticket }
  }),

  checkToken: (accessToken: string) => externalApi.get('/uc/login/check-token', {
    params: { access_token: accessToken }
  }),
}

export default api
