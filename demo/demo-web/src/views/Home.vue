<template>
  <div class="home-container">
    <header class="header">
      <h1>欢迎回来，{{ authStore.user?.username }}！</h1>
      <nav class="nav">
        <router-link to="/home" class="nav-link active">首页</router-link>
        <router-link to="/profile" class="nav-link">个人资料</router-link>
        <button @click="handleLogout" class="btn-logout">退出登录</button>
      </nav>
    </header>

    <main class="main-content">
      <div class="welcome-card">
        <h2>用户登录系统</h2>
        <p>这是一个简单的前后端分离登录示例项目。</p>
        <div class="features">
          <div class="feature">
            <span class="icon">🔐</span>
            <span>JWT 认证</span>
          </div>
          <div class="feature">
            <span class="icon">🛡️</span>
            <span>密码加密</span>
          </div>
          <div class="feature">
            <span class="icon">📱</span>
            <span>响应式设计</span>
          </div>
        </div>
      </div>

      <div class="info-card">
        <h3>当前登录用户信息</h3>
        <div class="user-info">
          <p><strong>用户ID：</strong>{{ authStore.user?.id }}</p>
          <p><strong>用户名：</strong>{{ authStore.user?.username }}</p>
          <p><strong>邮箱：</strong>{{ authStore.user?.email || '未设置' }}</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.home-container {
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
  max-width: 800px;
  margin: 40px auto;
  padding: 0 20px;
}

.welcome-card,
.info-card {
  background: white;
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.welcome-card h2 {
  color: #333;
  margin-bottom: 15px;
}

.welcome-card > p {
  color: #666;
  margin-bottom: 25px;
}

.features {
  display: flex;
  gap: 30px;
  flex-wrap: wrap;
}

.feature {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  background: #f8f9ff;
  border-radius: 8px;
}

.feature .icon {
  font-size: 20px;
}

.info-card h3 {
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 2px solid #f0f0f0;
}

.user-info p {
  padding: 10px 0;
  color: #555;
  border-bottom: 1px solid #f0f0f0;
}

.user-info p:last-child {
  border-bottom: none;
}
</style>
