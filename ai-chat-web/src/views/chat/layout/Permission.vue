<script setup lang='ts'>
import { computed, onMounted, ref } from 'vue'
import { NButton, NModal, NSpin, useMessage } from 'naive-ui'
import { useAuthStore } from '@/store'
import Icon403 from '@/icons/403.vue'
import { fetchCheckAuth, redirectToGitLabLogin } from '@/api'

interface Props {
  visible: boolean
}

defineProps<Props>()

const authStore = useAuthStore()
const ms = useMessage()

const loading = ref(false)

// 初始化时检查本地 token 和处理 OAuth 回调
onMounted(async () => {
  console.log('[Permission] 组件挂载，开始检查登录状态...')
  console.log('[Permission] 当前 URL:', window.location.href)

  // 解析 URL 参数
  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('access_token')
  const error = urlParams.get('error')
  const errorDescription = urlParams.get('error_description')

  console.log('[Permission] URL 参数:', { token: token ? '有值' : '无', error, errorDescription })

  // 处理 OAuth 错误
  if (error) {
    const msg = errorDescription || error
    ms.error(`登录失败: ${msg}`)
    console.error('[Permission] OAuth 错误:', error, errorDescription)
    // 清除 URL 中的错误参数
    window.history.replaceState({}, '', window.location.pathname)
    return
  }

  // 处理 OAuth 回调的 token
  if (token) {
    console.log('[Permission] 检测到 access_token，开始验证...')
    try {
      loading.value = true

      // 保存 token
      localStorage.setItem('access_token', token)
      authStore.setToken(token)

      // 验证 token 并获取用户信息
      const userInfo = await fetchCheckAuth(token)
      authStore.setUserInfo(userInfo)

      console.log('[Permission] 登录成功，用户:', userInfo)
      ms.success(`登录成功，欢迎 ${userInfo.name}`)

      // 清除 URL 中的 token 参数并刷新
      window.history.replaceState({}, '', window.location.pathname)
      window.location.reload()
    }
    catch (e) {
      console.error('[Permission] Token 验证失败:', e)
      localStorage.removeItem('access_token')
      authStore.removeToken()
      ms.error('登录验证失败，请重新登录')
    }
    finally {
      loading.value = false
    }
    return
  }

  // 检查本地存储的 token
  const localToken = localStorage.getItem('access_token')
  if (localToken) {
    console.log('[Permission] 发现本地 token，开始验证...')
    try {
      loading.value = true
      const userInfo = await fetchCheckAuth(localToken)
      authStore.setToken(localToken)
      authStore.setUserInfo(userInfo)
      ms.success(`欢迎回来，${userInfo.name}`)
    }
    catch (e) {
      console.error('[Permission] 本地 token 已过期:', e)
      localStorage.removeItem('access_token')
      authStore.removeToken()
    }
    finally {
      loading.value = false
    }
  }
})

function handleGitLabLogin() {
  console.log('[Permission] 点击了 GitLab 登录按钮')
  redirectToGitLabLogin('ai')
}
</script>

<template>
  <NModal :show="visible" style="width: 90%; max-width: 480px">
    <div class="p-10 bg-white rounded dark:bg-slate-800">
      <NSpin :show="loading">
        <div class="space-y-6">
          <header class="space-y-2">
            <h2 class="text-2xl font-bold text-center text-slate-800 dark:text-neutral-200">
              请先登录
            </h2>
            <p class="text-base text-center text-slate-500 dark:text-slate-500">
              使用 GitLab 账号登录 AI 聊天
            </p>
            <Icon403 class="w-[200px] m-auto" />
          </header>

          <div class="space-y-3">
            <NButton
              type="primary"
              block
              size="large"
              @click="handleGitLabLogin"
            >
              <template #icon>
                <svg class="w-5 h-5" viewBox="0 0 16 16" fill="currentColor">
                  <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
                </svg>
              </template>
              使用 GitLab 账号登录
            </NButton>

            <p class="text-sm text-center text-slate-400">
              登录后将跳转回 AI 聊天页面
            </p>
          </div>
        </div>
      </NSpin>
    </div>
  </NModal>
</template>
