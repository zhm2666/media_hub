<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>AI 聊天助手</h1>
        <p>基于 Kimi 大模型的智能对话</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="token">访问令牌</label>
          <input
            id="token"
            v-model="token"
            type="text"
            placeholder="请输入访问令牌"
            required
          />
          <p class="hint">从用户中心获取访问令牌</p>
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="success" class="success-message">{{ success }}</div>

        <button type="submit" class="btn-login" :disabled="loading">
          {{ loading ? '验证中...' : '登录' }}
        </button>
      </form>

      <div class="login-footer">
        <p>登录即表示同意我们的服务条款</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores'
import { authApi } from '@/utils/api'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const token = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

const handleLogin = async () => {
  if (!token.value.trim()) {
    error.value = '请输入访问令牌'
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const response = await authApi.checkToken(token.value.trim())

    if (response.status === 200 && response.data) {
      const userData = {
        id: response.data.id,
        name: response.data.name,
        avatar_url: response.data.avatar_url || '',
      }

      authStore.setAuth(token.value.trim(), userData)
      success.value = '登录成功！'

      setTimeout(() => {
        const redirect = (route.query.redirect as string) || '/chat'
        router.push(redirect)
      }, 500)
    }
  } catch (err: any) {
    if (err.response?.status === 401) {
      error.value = '令牌无效或已过期'
    } else if (err.response?.data?.message) {
      error.value = err.response.data.message
    } else {
      error.value = '登录失败，请检查令牌是否正确'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 20px;
  padding: 50px 40px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h1 {
  font-size: 28px;
  color: #333;
  margin-bottom: 10px;
}

.login-header p {
  color: #666;
  font-size: 14px;
}

.form-group {
  margin-bottom: 25px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #555;
  font-weight: 500;
  font-size: 14px;
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
  border-color: #667eea;
}

.form-group input::placeholder {
  color: #aaa;
}

.hint {
  margin-top: 8px;
  font-size: 12px;
  color: #999;
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

.success-message {
  color: #27ae60;
  background: #e8f8f0;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
  text-align: center;
  font-size: 14px;
}

.btn-login {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-login:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-footer {
  margin-top: 30px;
  text-align: center;
}

.login-footer p {
  font-size: 12px;
  color: #999;
}
</style>
