<script setup lang='ts'>
import { computed, onMounted } from 'vue'
import { NLayout, NLayoutContent, NButton, NAvatar, NDropdown, NIcon } from 'naive-ui'
import { useRouter } from 'vue-router'
import { useBasicLayout } from '@/hooks/useBasicLayout'
import { useAppStore, useChatStore, useAuthStore } from '@/store'
import { fetchCheckAuth } from '@/api'

const router = useRouter()
const appStore = useAppStore()
const chatStore = useChatStore()
const authStore = useAuthStore()

router.replace({ name: 'Chat', params: { uuid: chatStore.active } })

const { isMobile } = useBasicLayout()

const collapsed = computed(() => appStore.siderCollapsed)

const getMobileClass = computed(() => {
  if (isMobile.value)
    return ['rounded-none', 'shadow-none']
  return ['border', 'rounded-md', 'shadow-md', 'dark:border-neutral-800']
})

const getContainerClass = computed(() => {
  return [
    'h-full',
  ]
})

// 检查是否已登录
const isLoggedIn = computed(() => !!authStore.userInfo)

// 初始化时检查本地 token
onMounted(async () => {
  const token = localStorage.getItem('access_token')
  if (token && !authStore.userInfo) {
    try {
      const userInfo = await fetchCheckAuth(token)
      authStore.setToken(token)
      authStore.setUserInfo(userInfo)
    }
    catch {
      localStorage.removeItem('access_token')
      authStore.removeToken()
    }
  }
})

// 处理登录 - 跳转到用户中心
async function handleLogin() {
  // 从环境变量获取用户中心地址，使用 ?sys=ai 标识本系统
  const userCenter = import.meta.env.VITE_USER_CENTER || 'http://localhost:8082?sys=ai'
  // 直接跳转到用户中心登录页面
  window.location.href = userCenter
}

// 处理登出
function handleLogout() {
  localStorage.removeItem('access_token')
  authStore.removeToken()
}

// 下拉菜单选项
const userMenuOptions = [
  {
    label: '退出登录',
    key: 'logout',
  },
]

function handleUserMenuSelect(key: string) {
  if (key === 'logout') {
    handleLogout()
  }
}
</script>

<template>
  <div class="h-full dark:bg-[#24272e] transition-all" :class="[isMobile ? 'p-0' : 'p-4']">
    <div class="h-full overflow-hidden" :class="getMobileClass">
      <NLayout class="z-40 transition" :class="getContainerClass" has-sider>
        <!-- 顶部导航栏 -->
        <div class="absolute top-4 right-4 z-50 flex items-center gap-3">
          <template v-if="isLoggedIn">
            <NDropdown
              :options="userMenuOptions"
              @select="handleUserMenuSelect"
            >
              <div class="flex items-center gap-2 cursor-pointer">
                <NAvatar
                  :size="36"
                  :src="authStore.userInfo?.avatar_url"
                  round
                />
                <span class="text-sm text-gray-600 dark:text-gray-300">
                  {{ authStore.userInfo?.name }}
                </span>
              </div>
            </NDropdown>
          </template>
          <template v-else>
            <NButton type="primary" size="small" @click="handleLogin">
              登录
            </NButton>
          </template>
        </div>

        <NLayoutContent class="h-full">
          <RouterView v-slot="{ Component, route }">
            <component :is="Component" :key="route.fullPath" />
          </RouterView>
        </NLayoutContent>
      </NLayout>
    </div>
  </div>
</template>
