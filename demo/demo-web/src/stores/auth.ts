import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userCenterApi } from '@/utils/api'
import { storage } from '@/utils/storage'
import type { User, LoginRequest, RegisterRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(storage.getToken())
  const user = ref<User | null>(storage.getUser())

  const isLoggedIn = computed(() => !!token.value && !!user.value)

  const login = async (credentials: LoginRequest) => {
    const response = await authApi.login(credentials)
    const { token: newToken, user: newUser } = response.data

    token.value = newToken
    user.value = newUser
    storage.setToken(newToken)
    storage.setUser(newUser)

    return response.data
  }

  const register = async (data: RegisterRequest) => {
    const response = await authApi.register(data)
    const { token: newToken, user: newUser } = response.data

    token.value = newToken
    user.value = newUser
    storage.setToken(newToken)
    storage.setUser(newUser)

    return response.data
  }

  const logout = async () => {
    try {
      await authApi.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      token.value = null
      user.value = null
      storage.clear()
    }
  }

  const fetchCurrentUser = async () => {
    try {
      const response = await authApi.getCurrentUser()
      user.value = response.data
      storage.setUser(response.data)
    } catch (error) {
      logout()
      throw error
    }
  }

  const loginWithUserCenter = async (accessToken: string) => {
    try {
      const response = await userCenterApi.checkToken(accessToken)
      const { token: newToken, user: newUser } = response.data

      token.value = newToken
      user.value = newUser
      storage.setToken(newToken)
      storage.setUser(newUser)

      return response.data
    } catch (error) {
      console.error('User center login error:', error)
      throw error
    }
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    logout,
    fetchCurrentUser,
    loginWithUserCenter,
  }
})
