export interface User {
  id: number
  username: string
  email: string
  created_at?: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  email: string
}

export interface AuthResponse {
  message: string
  token: string
  user: User
}

export interface ApiResponse<T = any> {
  data?: T
  error?: string
  message?: string
}
