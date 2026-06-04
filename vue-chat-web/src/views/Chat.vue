<template>
  <div class="chat-container">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <h2>AI 助手</h2>
        <button @click="toggleSidebar" class="btn-toggle">
          {{ sidebarCollapsed ? '展开' : '收起' }}
        </button>
      </div>

      <button @click="clearChat" class="btn-new-chat">
        <span class="icon">+</span>
        新建对话
      </button>

      <div class="chat-list">
        <div class="chat-item" :class="{ active: true }">
          <span class="chat-title">当前对话</span>
        </div>
      </div>

      <div class="sidebar-footer">
        <div class="user-info">
          <div class="avatar">{{ user?.name?.charAt(0) || 'U' }}</div>
          <div class="user-details">
            <span class="username">{{ user?.name || 'User' }}</span>
          </div>
        </div>
        <button @click="handleLogout" class="btn-logout">退出</button>
      </div>
    </aside>

    <!-- 主聊天区域 -->
    <main class="chat-main">
      <!-- 消息列表 -->
      <div class="messages-container" ref="messagesContainer">
        <div v-if="messages.length === 0" class="welcome-message">
          <div class="welcome-icon">&#x1F4AC;</div>
          <h2>你好，我是 AI 助手</h2>
          <p>有什么我可以帮助你的吗？</p>
        </div>

        <div
          v-for="message in messages"
          :key="message.id"
          class="message"
          :class="message.role"
        >
          <div class="message-avatar">
            {{ message.role === 'user' ? 'U' : 'AI' }}
          </div>
          <div class="message-content">
            <MarkdownRenderer
              v-if="message.role === 'assistant'"
              :content="message.content"
            />
            <div v-else class="user-message">{{ message.content }}</div>
            <span
              v-if="message.isGenerating"
              class="typing-indicator"
            >
              <span></span><span></span><span></span>
            </span>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="input-container">
        <div class="input-wrapper">
          <textarea
            v-model="inputMessage"
            @keydown.enter.exact.prevent="handleSend"
            placeholder="输入消息... (Enter 发送，Shift+Enter 换行)"
            :disabled="isGenerating"
            rows="1"
            ref="textareaRef"
          ></textarea>
          <button
            @click="handleSend"
            :disabled="!inputMessage.trim() || isGenerating"
            class="btn-send"
          >
            <span v-if="isGenerating" class="spinner"></span>
            <span v-else>&#x27A4;</span>
          </button>
        </div>
        <p class="input-hint">
          AI 助手可能会产生不准确的信息，请仔细核对重要内容。
        </p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore, useChatStore } from '@/stores'
import { chatApi } from '@/utils/api'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import type { ChatMessage } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const chatStore = useChatStore()

const user = authStore.user
const messages = ref<ChatMessage[]>([])
const inputMessage = ref('')
const isGenerating = ref(false)
const sidebarCollapsed = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)
const textareaRef = ref<HTMLTextAreaElement | null>(null)

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const handleSend = async () => {
  const content = inputMessage.value.trim()
  if (!content || isGenerating.value) return

  const userMessage: ChatMessage = {
    id: `user-${Date.now()}`,
    role: 'user',
    content,
    createdAt: new Date(),
  }

  messages.value.push(userMessage)
  inputMessage.value = ''
  scrollToBottom()

  const assistantMessage: ChatMessage = {
    id: `assistant-${Date.now()}`,
    role: 'assistant',
    content: '',
    createdAt: new Date(),
    isGenerating: true,
  }
  messages.value.push(assistantMessage)
  isGenerating.value = true
  scrollToBottom()

  try {
    const token = authStore.token
    if (!token) {
      throw new Error('未登录')
    }

    const response = await chatApi.sendMessage({
      prompt: content,
      systemMessage: '你是一个有帮助的AI助手，请用Markdown格式回答问题。',
    })

    const reader = response.data.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        if (line.startsWith('event: ')) {
          const event = line.slice(7).trim()
          continue
        }
        if (line.startsWith('data: ')) {
          const data = line.slice(6).trim()
          if (data === '[DONE]' || data === '') continue

          try {
            const parsed = JSON.parse(data)
            if (parsed.delta) {
              assistantMessage.content += parsed.delta
              scrollToBottom()
            }
          } catch (e) {
            // ignore parse errors
          }
        }
      }
    }
  } catch (error: any) {
    console.error('发送消息失败:', error)
    assistantMessage.content = error.response?.data?.message || error.message || '发送消息失败，请重试'
  } finally {
    assistantMessage.isGenerating = false
    isGenerating.value = false
    scrollToBottom()
  }
}

const clearChat = () => {
  messages.value = []
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

onMounted(() => {
  if (!authStore.isLoggedIn) {
    router.push('/login')
  }
})
</script>

<style scoped>
.chat-container {
  display: flex;
  height: 100vh;
  background: #f5f7fa;
}

/* 侧边栏 */
.sidebar {
  width: 280px;
  background: white;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
}

.sidebar.collapsed {
  width: 0;
  overflow: hidden;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header h2 {
  font-size: 18px;
  color: #333;
}

.btn-toggle {
  padding: 6px 12px;
  background: #f0f0f0;
  border-radius: 6px;
  font-size: 12px;
  color: #666;
}

.btn-new-chat {
  margin: 16px 16px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 10px;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-new-chat:hover {
  opacity: 0.9;
}

.btn-new-chat .icon {
  font-size: 18px;
}

.chat-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 8px;
}

.chat-item {
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 4px;
  color: #666;
  font-size: 14px;
}

.chat-item:hover {
  background: #f5f5f5;
}

.chat-item.active {
  background: #f0f0ff;
  color: #667eea;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
}

.user-details {
  display: flex;
  flex-direction: column;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.btn-logout {
  padding: 8px 16px;
  background: #fde8e8;
  color: #e74c3c;
  border-radius: 6px;
  font-size: 12px;
}

.btn-logout:hover {
  background: #fbd;
}

/* 主聊天区域 */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.welcome-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  color: #666;
}

.welcome-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.welcome-message h2 {
  font-size: 24px;
  color: #333;
  margin-bottom: 10px;
}

.message {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  flex-shrink: 0;
}

.message.user .message-avatar {
  background: #667eea;
  color: white;
}

.message.assistant .message-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.message-content {
  max-width: 75%;
  min-width: 100px;
}

.user-message {
  background: #667eea;
  color: white;
  padding: 12px 16px;
  border-radius: 16px;
  border-top-right-radius: 4px;
  line-height: 1.5;
}

.message.assistant .message-content {
  background: white;
  padding: 16px;
  border-radius: 16px;
  border-top-left-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.typing-indicator {
  display: inline-flex;
  gap: 4px;
  padding: 8px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: #999;
  border-radius: 50%;
  animation: bounce 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(1) { animation-delay: 0s; }
.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.5; }
  40% { transform: scale(1); opacity: 1; }
}

/* 输入区域 */
.input-container {
  padding: 16px 20px 20px;
  background: #f5f7fa;
}

.input-wrapper {
  display: flex;
  gap: 12px;
  background: white;
  border-radius: 16px;
  padding: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.input-wrapper textarea {
  flex: 1;
  border: none;
  resize: none;
  padding: 12px;
  font-size: 15px;
  line-height: 1.5;
  max-height: 150px;
  font-family: inherit;
}

.input-wrapper textarea::placeholder {
  color: #aaa;
}

.btn-send {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}

.btn-send:hover:not(:disabled) {
  transform: scale(1.05);
}

.btn-send:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.input-hint {
  text-align: center;
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}
</style>
