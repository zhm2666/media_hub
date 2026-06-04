<template>
  <div class="auth-container">
    <div class="auth-card">
      <h1 class="title">登录</h1>

      <!-- 账号登录模式 -->
      <div v-if="loginMode === 'account'" class="login-mode">
        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label for="username">用户名</label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              placeholder="请输入用户名"
              required
            />
          </div>
          <div class="form-group">
            <label for="password">密码</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              required
            />
          </div>
          <div v-if="error" class="error-message">{{ error }}</div>
          <button type="submit" class="btn-primary" :disabled="loading">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </form>

        <div class="divider">
          <span>或</span>
        </div>

        <button @click="switchToScanLogin" class="btn-scan">
          <span class="scan-icon">&#x1F4F1;</span>
          微信扫码登录
        </button>
      </div>

      <!-- 扫码登录模式 -->
      <div v-else class="login-mode scan-mode">
        <div class="qrcode-container">
          <div v-if="qrcodeUrl" class="qrcode-wrapper">
            <img :src="qrcodeUrl" alt="微信扫码登录" class="qrcode-image" />
            <div v-if="scanStatus === 'scanning'" class="scan-overlay">
              <div class="spinner"></div>
              <p>等待扫码...</p>
            </div>
            <div v-if="scanStatus === 'scanned'" class="scan-success">
              <span class="success-icon">&#x2713;</span>
              <p>扫码成功</p>
            </div>
          </div>
          <div v-else class="qrcode-loading">
            <div class="spinner"></div>
            <p>加载二维码...</p>
          </div>
        </div>
        <p class="scan-tip">使用微信扫码登录</p>
        <button @click="switchToAccountLogin" class="btn-back">
          返回账号登录
        </button>
      </div>

      <p class="switch-link">
        还没有账号？<router-link to="/register">立即注册</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { storage } from '@/utils/storage'
import axios from 'axios'
import type { LoginRequest } from '@/types'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const API_BASE = 'http://localhost:8080/api'

const form = ref<LoginRequest>({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')
const loginMode = ref<'account' | 'scan'>('account')
const qrcodeUrl = ref('')
const qrcodeTicket = ref('')
const scanStatus = ref<'scanning' | 'scanned' | 'waiting'>('waiting')
let pollingTimer: number | null = null

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  try {
    await authStore.login(form.value)
    const redirect = (route.query.redirect as string) || '/home'
    router.push(redirect)
  } catch (err: any) {
    error.value = err.response?.data?.error || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}

const switchToScanLogin = async () => {
  loginMode.value = 'scan'
  await loadQrcode()
}

const switchToAccountLogin = () => {
  loginMode.value = 'account'
  stopPolling()
}

const loadQrcode = async () => {
  try {
    scanStatus.value = 'waiting'
    const response = await axios.get(`${API_BASE}/uc/login/qrcode`)
    qrcodeUrl.value = response.data.qrcode_url
    qrcodeTicket.value = response.data.ticket
    scanStatus.value = 'scanning'
    startPolling()
  } catch (err) {
    console.error('加载二维码失败:', err)
    error.value = '加载二维码失败，请重试'
  }
}

const startPolling = () => {
  stopPolling()
  pollingTimer = window.setInterval(async () => {
    try {
      const response = await axios.get(`${API_BASE}/uc/login/check-scan`, {
        params: { ticket: qrcodeTicket.value }
      })

      if (response.data.status === 'scanned') {
        scanStatus.value = 'scanned'
        stopPolling()

        // 登录成功，跳转
        setTimeout(() => {
          router.push('/home')
        }, 1000)
      }
    } catch (err) {
      console.error('检查扫码状态失败:', err)
    }
  }, 2000)
}

const stopPolling = () => {
  if (pollingTimer !== null) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

onMounted(() => {
  const savedToken = storage.getToken()
  if (savedToken) {
    router.push('/home')
  }
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
}

.auth-card {
  background: white;
  border-radius: 16px;
  padding: 40px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.title {
  text-align: center;
  color: #333;
  margin-bottom: 30px;
  font-size: 28px;
  font-weight: 600;
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

.divider {
  text-align: center;
  margin: 20px 0;
  position: relative;
}

.divider::before,
.divider::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 45%;
  height: 1px;
  background: #e0e0e0;
}

.divider::before {
  left: 0;
}

.divider::after {
  right: 0;
}

.divider span {
  color: #999;
  font-size: 14px;
}

.btn-scan {
  width: 100%;
  padding: 14px;
  background: #07c160;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-scan:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(7, 193, 96, 0.4);
}

.scan-icon {
  font-size: 20px;
}

/* 扫码模式样式 */
.scan-mode {
  text-align: center;
}

.qrcode-container {
  margin: 20px 0;
}

.qrcode-wrapper {
  position: relative;
  display: inline-block;
  border: 4px solid #07c160;
  border-radius: 12px;
  overflow: hidden;
}

.qrcode-image {
  display: block;
  width: 200px;
  height: 200px;
}

.qrcode-loading {
  padding: 60px;
  color: #999;
}

.scan-overlay,
.scan-success {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.scan-success {
  background: rgba(7, 193, 96, 0.9);
  color: white;
}

.success-icon {
  font-size: 48px;
  background: white;
  color: #07c160;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e0e0e0;
  border-top-color: #07c160;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.scan-tip {
  color: #666;
  margin: 15px 0;
  font-size: 14px;
}

.btn-back {
  background: transparent;
  color: #667eea;
  border: 2px solid #667eea;
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-back:hover {
  background: #667eea;
  color: white;
}

.switch-link {
  text-align: center;
  margin-top: 20px;
  color: #666;
}

.switch-link a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.switch-link a:hover {
  text-decoration: underline;
}
</style>
