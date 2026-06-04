import type { Router } from 'vue-router'
import { useAuthStoreWithout } from '@/store/modules/auth'
import { fetchCheckAuth } from '@/api'

export function setupPageGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStoreWithout()

    // 检查 URL 中是否有 OAuth 回调的 access_token
    const urlParams = new URLSearchParams(window.location.search)
    const token = urlParams.get('access_token')
    const error = urlParams.get('error')

    if (error) {
      console.error('OAuth error:', error)
      window.history.replaceState({}, '', window.location.pathname)
    }

    if (token && !authStore.userInfo) {
      try {
        // 保存 token
        localStorage.setItem('access_token', token)
        authStore.setToken(token)

        // 获取用户信息
        const userInfo = await fetchCheckAuth(token)
        authStore.setUserInfo(userInfo)

        // 清除 URL 中的 token 参数
        window.history.replaceState({}, '', window.location.pathname)
      }
      catch {
        localStorage.removeItem('access_token')
        authStore.removeToken()
      }
    }

    next()
  })
}
