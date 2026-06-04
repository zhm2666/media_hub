<template>
  <div class="profile-container">
    <header class="header">
      <h1>个人资料</h1>
      <nav class="nav">
        <router-link to="/home" class="nav-link">首页</router-link>
        <router-link to="/profile" class="nav-link active">个人资料</router-link>
        <button @click="handleLogout" class="btn-logout">退出登录</button>
      </nav>
    </header>

    <main class="main-content">
      <div class="profile-card">
        <h2>编辑个人资料</h2>

        <form @submit.prevent="handleUpdateProfile">
          <div class="form-group">
            <label for="username">用户名</label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              disabled
            />
            <small>用户名不可修改</small>
          </div>

          <div class="form-group">
            <label for="email">邮箱</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="请输入邮箱"
            />
          </div>

          <div v-if="message" class="success-message">{{ message }}</div>
          <div v-if="error" class="error-message">{{ error }}</div>

          <button type="submit" class="btn-primary" :disabled="loading">
            {{ loading ? '保存中...' : '保存修改' }}
          </button>
        </form>
      </div>

      <div class="danger-zone">
        <h3>危险区域</h3>
        <p>退出登录后，您需要重新输入用户名和密码才能访问系统。</p>
        <button @click="handleLogout" class="btn-danger">退出登录</button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { userApi } from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  email: '',
})

const loading = ref(false)
const message = ref('')
const error = ref('')

onMounted(() => {
  if (authStore.user) {
    form.value.username = authStore.user.username
    form.value.email = authStore.user.email || ''
  }
})

const handleUpdateProfile = async () => {
  loading.value = true
  message.value = ''
  error.value = ''

  try {
    await userApi.updateProfile({ email: form.value.email })
    message.value = '个人资料更新成功！'
    if (authStore.user) {
      authStore.user.email = form.value.email
    }
  } catch (err: any) {
    error.value = err.response?.data?.error || '更新失败，请重试'
  } finally {
    loading.value = false
  }
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.header {
  background: white;
  padding: 20px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.header h1 {
  color: #333;
  font-size: 24px;
}

.nav {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-link {
  color: #666;
  text-decoration: none;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.3s;
}

.nav-link:hover,
.nav-link.active {
  color: #667eea;
  background: #f0f2ff;
}

.btn-logout {
  padding: 8px 20px;
  background: #f56565;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.3s;
}

.btn-logout:hover {
  background: #e53e3e;
}

.main-content {
  max-width: 600px;
  margin: 40px auto;
  padding: 0 20px;
}

.profile-card,
.danger-zone {
  background: white;
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.profile-card h2,
.danger-zone h3 {
  color: #333;
  margin-bottom: 25px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #555;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid #e0e0e0;
  border-radius: 10px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-group input:disabled {
  background: #f5f5f5;
  color: #999;
  cursor: not-allowed;
}

.form-group small {
  display: block;
  margin-top: 6px;
  color: #999;
  font-size: 12px;
}

.success-message {
  color: #27ae60;
  background: #e8f8f0;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
  text-align: center;
  font-size: 14px;
}

.error-message {
  color: #e74c3c;
  background: #fde8e8;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
  text-align: center;
  font-size: 14px;
}

.btn-primary {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.danger-zone {
  border: 2px solid #feb2b2;
}

.danger-zone p {
  color: #666;
  margin-bottom: 15px;
}

.btn-danger {
  padding: 12px 24px;
  background: #f56565;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.3s;
}

.btn-danger:hover {
  background: #e53e3e;
}
</style>
