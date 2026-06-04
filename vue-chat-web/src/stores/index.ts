import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, ChatMessage } from '@/types'
import { storage } from '@/utils/storage'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(storage.getToken())
  const user = ref<User | null>(storage.getUser())

  const isLoggedIn = computed(() => !!token.value && !!user.value)

  const setAuth = (newToken: string, newUser: User) => {
    token.value = newToken
    user.value = newUser
    storage.setToken(newToken)
    storage.setUser(newUser)
  }

  const logout = () => {
    token.value = null
    user.value = null
    storage.clear()
  }

  return {
    token,
    user,
    isLoggedIn,
    setAuth,
    logout,
  }
})

export const useChatStore = defineStore('chat', () => {
  const messages = ref<ChatMessage[]>([])
  const isGenerating = ref(false)
  const currentGeneratingId = ref<string | null>(null)

  const addMessage = (message: ChatMessage) => {
    messages.value.push(message)
  }

  const updateMessage = (id: string, content: string, isGenerating = false) => {
    const message = messages.value.find(m => m.id === id)
    if (message) {
      message.content = content
      message.isGenerating = isGenerating
    }
  }

  const appendToMessage = (id: string, delta: string) => {
    const message = messages.value.find(m => m.id === id)
    if (message) {
      message.content += delta
    }
  }

  const clearMessages = () => {
    messages.value = []
  }

  const setGenerating = (generating: boolean, id: string | null = null) => {
    isGenerating.value = generating
    currentGeneratingId.value = id
  }

  return {
    messages,
    isGenerating,
    currentGeneratingId,
    addMessage,
    updateMessage,
    appendToMessage,
    clearMessages,
    setGenerating,
  }
})
