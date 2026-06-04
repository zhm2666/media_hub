export interface User {
  id: number
  name: string
  avatar_url: string
}

export interface ChatMessage {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  createdAt: Date
  isGenerating?: boolean
}

export interface ChatRequest {
  prompt: string
  systemMessage?: string
  options?: {
    parentMessageId?: string
  }
}

export interface ChatResponse {
  id: string
  text: string
  delta: string
  role: string
  detail?: any
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ApiResponse<T = any> {
  status?: string
  message?: string
  error?: string
  data?: T
}
