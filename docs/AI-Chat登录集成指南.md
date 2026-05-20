# AI Chat 系统登录集成指南

## 一、概述

本文档说明如何为 AI Chat 系统（`ai-chat-web` + `gin_py`）添加用户中心登录功能。

---

## 二、集成架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       AI Chat 登录集成架构                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌──────────────────┐         ┌──────────────────┐                       │
│   │   AI Chat 前端    │         │   用户中心        │                       │
│   │   ai-chat-web     │ ──────→ │   user-web       │                       │
│   │   localhost:5174  │ ?sys=ai │   localhost:8082 │                       │
│   └──────────────────┘         └──────────────────┘                       │
│           │                              │                                  │
│           │ Cookie: JWT Token            │                                  │
│           │                              │                                  │
│           ↓                              ↓                                  │
│   ┌──────────────────┐         ┌──────────────────┐                       │
│   │   AI Chat 后端     │         │   用户中心后端    │                       │
│   │   gin_py          │ ──────→ │   user 服务       │                       │
│   │   localhost:7080   │ 验证Token │  localhost:8082 │                       │
│   └──────────────────┘         └──────────────────┘                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 三、修改的文件清单

### 3.1 前端文件（ai-chat-web）

| 文件路径 | 修改内容 |
|----------|----------|
| `ai-chat-web/src/utils/cookie/index.ts` | 添加 `getCookieByKey` 函数别名 |
| `ai-chat-web/src/utils/request/axios.ts` | 修改请求拦截器，正确处理 JWT Token |
| `ai-chat-web/src/views/chat/layout/Layout.vue` | 简化登录跳转逻辑 |
| `ai-chat-web/src/views/chat/index.vue` | 简化登录跳转逻辑 |
| `ai-chat-web/.env.develop` | 添加配置注释 |
| `ai-chat-web/.env.production` | 配置生产环境用户中心地址 |

### 3.2 后端文件（gin_py）

| 文件路径 | 修改内容 |
|----------|----------|
| `gin_py/internal/middleware/jwt.go` | 新增 JWT 验证中间件 |
| `gin_py/internal/service/user.go` | 新增用户服务客户端 |
| `gin_py/internal/handler/chat.go` | 添加登录验证逻辑 |
| `gin_py/config.example.yaml` | 新增配置文件模板 |

### 3.3 用户中心配置

| 文件路径 | 修改内容 |
|----------|----------|
| `user-master/user/dev.config.yaml` | 添加 `ai` 系统跳转配置 |
| `user-master/user/test.config.yaml` | 添加 `ai` 系统跳转配置 |

---

## 四、详细代码说明

### 4.1 前端 Cookie 工具函数

**文件**：`ai-chat-web/src/utils/cookie/index.ts`

添加 `getCookieByKey` 函数，用于 axios 请求拦截器中引用：

```typescript:1:21:ai-chat-web/src/utils/cookie/index.ts
export function getCookieValue(key: string) {
  const cookies = document.cookie.split(';')
  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i].trim()
    if (cookie.startsWith(`${key}=`))
      return cookie.substring(key.length + 1)
  }
  return null
}

// getCookieByKey 是 getCookieValue 的别名，用于兼容 axios.ts 中的引用
export function getCookieByKey(key: string) {
  return getCookieValue(key)
}
```

### 4.2 前端请求拦截器

**文件**：`ai-chat-web/src/utils/request/axios.ts`

修改请求拦截器，从 Cookie 获取 JWT Token 并添加到请求头：

```typescript:1:26:ai-chat-web/src/utils/request/axios.ts
import axios, { type AxiosResponse } from 'axios'
import { deleteCookieByKey, getCookieByKey } from '../cookie/index'

const service = axios.create({
  baseURL: import.meta.env.VITE_GLOB_API_URL,
})

service.interceptors.request.use(
  (config) => {
    // 优先从 Cookie 获取 JWT Token
    const access_token = getCookieByKey('sso_0voice_access_token')
    if (access_token) {
      config.headers.Authorization = `Bearer ${access_token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error.response)
  },
)
```

### 4.3 前端登录跳转逻辑

**文件**：`ai-chat-web/src/views/chat/layout/Layout.vue`

简化登录跳转逻辑，直接跳转到用户中心：

```typescript:51:57:ai-chat-web/src/views/chat/layout/Layout.vue
// 处理登录 - 跳转到用户中心
async function handleLogin() {
  // 从环境变量获取用户中心地址，使用 ?sys=ai 标识本系统
  const userCenter = import.meta.env.VITE_USER_CENTER || 'http://localhost:8082?sys=ai'
  // 直接跳转到用户中心登录页面
  window.location.href = userCenter
}
```

### 4.4 前端发送消息时检查登录

**文件**：`ai-chat-web/src/views/chat/index.vue`

用户发送消息前检查登录状态，未登录则提示并跳转：

```typescript:62:84:ai-chat-web/src/views/chat/index.vue
// 检查是否已登录，未登录则提示登录
function checkLogin(): boolean {
  if (!authStore.userInfo) {
    ms.warning('请先登录后再发送消息')
    return false
  }
  return true
}

// 跳转到登录页面
async function goToLogin() {
  // 从环境变量获取用户中心地址，使用 ?sys=ai 标识本系统
  const userCenter = import.meta.env.VITE_USER_CENTER || 'http://localhost:8082?sys=ai'
  window.location.href = userCenter
}

function handleSubmit() {
  // 检查登录状态
  if (!checkLogin()) {
    goToLogin()
    return
  }
  onConversation()
}
```

### 4.5 后端 JWT 验证中间件

**文件**：`gin_py/internal/middleware/jwt.go`

```go:1:52:gin_py/internal/middleware/jwt.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 验证 JWT Token 的中间件
// 从 Authorization header 中获取 Token 并验证
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "缺少 Authorization header",
			})
			c.Abort()
			return
		}

		// 检查 Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization 格式错误，应为: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Token 不能为空",
			})
			c.Abort()
			return
		}

		// 将 token 传递给后续处理函数
		c.Set("token", tokenString)

		c.Next()
	}
}
```

### 4.6 后端用户服务客户端

**文件**：`gin_py/internal/service/user.go`

```go:1:80:gin_py/internal/service/user.go
package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// UserServiceClient 用户中心服务客户端
type UserServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// UserInfo 用户信息结构
type UserInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// NewUserServiceClient 创建用户服务客户端
func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyToken 验证 JWT Token 并获取用户信息
// 调用 user 服务的 /api/v1/login/check/auth 接口
func (c *UserServiceClient) VerifyToken(token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/api/v1/login/check/auth?access_token=%s", c.baseURL, token)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求用户中心失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("验证失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Code    int       `json:"code"`
		Message string    `json:"message"`
		Data    *UserInfo `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("验证失败: %s", result.Message)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("用户信息为空")
	}

	return result.Data, nil
}
```

### 4.7 后端聊天处理添加登录验证

**文件**：`gin_py/internal/handler/chat.go`

```go:28:61:gin_py/internal/handler/chat.go
// ChatProcess 处理流式请求（带登录验证）
func (h *ChatHandler) ChatProcess(c *gin.Context) {
	startTime := time.Now()
	fmt.Printf("[Handler] 收到请求: %s\n", startTime.Format("15:04:05.000"))

	// ========== 登录验证 START ==========
	// 从 Authorization header 获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "请先登录后再使用聊天功能",
		})
		return
	}

	// 解析 Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authorization 格式错误",
		})
		return
	}

	token := parts[1]

	// 调用用户中心验证 token
	userInfo, err := h.userService.VerifyToken(token)
	if err != nil {
		fmt.Printf("[Handler] Token 验证失败: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "登录已过期，请重新登录",
		})
		return
	}
	fmt.Printf("[Handler] 用户 %s (ID: %d) 请求聊天\n", userInfo.Name, userInfo.ID)
	// ========== 登录验证 END ==========
```

---

## 五、配置说明

### 5.1 前端环境变量

**文件**：`ai-chat-web/.env.develop`

```bash
# 用户中心地址 - 跳转到用户中心登录，登录后返回 ai-chat 页面
# ?sys=ai 标识本系统，后端会根据这个参数返回对应的 redirect_url
VITE_USER_CENTER="http://localhost:8082?sys=ai"
```

**文件**：`ai-chat-web/.env.production`

```bash
VITE_USER_CENTER="https://user.0voice.com?sys=ai"
```

### 5.2 用户中心配置

**文件**：`user-master/user/dev.config.yaml`

```yaml
internalSystemEntry:
  mediahub: "http://localhost:8080"
  transform: "http://localhost:8084"
  tunnel: "http://localhost:5173"
  ai: "http://localhost:5174"  # 新增：AI Chat 系统的跳转地址
```

---

## 六、登录流程

```
┌──────────┐         ┌──────────┐         ┌──────────┐
│   用户   │         │AI Chat   │         │ 用户中心  │
│          │         │  前端    │         │          │
└────┬─────┘         └────┬─────┘         └────┬─────┘
     │                    │                    │
     │  1. 访问 AI Chat   │                    │
     │──────────────────→│                    │
     │                    │                    │
     │                    │  2. 无 Token，显示登录按钮
     │←──────────────────│                    │
     │                    │                    │
     │  3. 点击"登录"      │                    │
     │──────────────────→│                    │
     │                    │                    │
     │                    │  4. 跳转 ?sys=ai   │
     │                    │──────────────────→│
     │                    │                    │
     │  5. 显示登录页面    │                    │
     │←──────────────────────────────────────│
     │                    │                    │
     │  6. 微信扫码登录    │                    │
     │───────────────────────────────────────│
     │                    │                    │
     │                    │                    │  7. 验证成功
     │                    │                    │     生成 JWT
     │                    │                    │     设置 Cookie
     │                    │                    │
     │                    │  8. 跳转 redirect_url
     │                    │    (http://localhost:5174)
     │←──────────────────────────────────────│
     │                    │                    │
     │  9. 页面加载，读取 Cookie               │
     │──────────────────→│                    │
     │                    │                    │
     │                    │  10. 验证 Token    │
     │                    │    获取用户信息     │
     │                    │                    │
     │  11. 显示用户头像  │                    │
     │←──────────────────│                    │
     │                    │                    │
     │  12. 发送消息      │                    │
     │──────────────────→│                    │
     │                    │                    │
     │                    │  13. 携带 Token    │
     │                    │     调用后端 API   │
     │                    │──────────────────→│
     │                    │                    │
     │                    │  14. 验证 Token    │
     │                    │←──────────────────│ │
     │                    │                    │
     │                    │  15. AI 响应       │
     │←──────────────────│                    │
     │                    │                    │
```

---

## 七、测试步骤

1. **启动用户中心服务**
   ```bash
   cd user-master/user
   go run main.go
   ```

2. **启动 AI Chat 后端**
   ```bash
   cd gin_py
   go run cmd/main.go
   ```

3. **启动 AI Chat 前端**
   ```bash
   cd ai-chat-web
   pnpm dev
   ```

4. **测试登录流程**
   - 访问 http://localhost:5174
   - 点击右上角"登录"按钮
   - 页面跳转到用户中心
   - 使用微信扫码登录
   - 登录成功后自动跳转回 AI Chat 页面
   - 右上角显示用户头像和名称

5. **测试聊天功能**
   - 未登录时发送消息会提示"请先登录"
   - 登录后可以正常发送消息
   - 后端控制台会显示用户信息

---

## 八、注意事项

1. **Cookie 域名**：确保用户中心设置的 Cookie 域名与 AI Chat 域名一致
2. **CORS 配置**：如果前后端域名不同，需要在后端配置 CORS
3. **Token 过期**：JWT Token 过期后需要重新登录
4. **端口配置**：确保各服务的端口配置正确
