import { post } from '@/utils/request'
import { useSettingStore } from '@/store'

export function fetchChatAPI<T = any>(
  prompt: string,
  options?: { conversationId?: string; parentMessageId?: string },
  signal?: AbortSignal,
) {
  return post<T>({
    url: '/chat',
    data: { prompt, options },
    signal,
  })
}

export function fetchChatConfig<T = any>() {
  return post<T>({
    url: '/config',
  })
}

// =============================================
// 以下为 SSE 流式请求实现 (替换原来的 fetchChatAPIProcess)
// =============================================

export interface SSEProgressEvent {
  text: string        // 累计文本
  delta: string       // 本次增量文本
  id: string          // 消息ID
  detail: any         // 完整响应详情
}

export interface SSEOptions {
  prompt: string
  options?: { conversationId?: string; parentMessageId?: string }
  signal?: AbortSignal
  onMessage: (event: SSEProgressEvent) => void       // 每收到一条消息的回调
  onDone?: () => void                                  // 流结束回调
  onError?: (error: Error) => void                    // 错误回调
}

/**
 * 使用 SSE (Server-Sent Events) 协议发送流式请求
 *
 * 【修改原因】
 * 1. 原实现使用 Axios 的 onDownloadProgress，但 Axios 不适合处理流式响应
 * 2. Axios 会尝试缓冲整个响应，无法实现真正的实时流式更新
 * 3. SSE 是标准的服务器推送事件协议，更适合流式数据传输
 * 4. 使用 fetch + ReadableStream 可以正确解析 SSE 事件格式
 *
 * 【SSE 响应格式】
 * event: message
 * data: {"text":"你","delta":"你","id":"xxx","detail":{...}}
 *
 * event: done
 * data: {"status":"complete"}
 */
export function fetchChatAPIProcessSSE(params: SSEOptions) {
  const settingStore = useSettingStore()

  // 1. 构建请求体
  const requestBody = {
    prompt: params.prompt,
    options: params.options,
    systemMessage: settingStore.systemMessage,
  }

  // 2. 获取基础 URL
  const baseURL = import.meta.env.VITE_API_BASE_URL || ''
  const url = baseURL + '/api/chat-process'

  // 3. 发送 SSE 请求
  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(requestBody),
    signal: params.signal,
  })
    .then(async (response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const reader = response.body?.getReader()
      if (!reader) {
        throw new Error('无法获取响应流')
      }

      const decoder = new TextDecoder()
      let buffer = ''
      let currentEvent = ''
      let currentData = ''

      try {
        while (true) {
          const { done, value } = await reader.read()

          if (done) {
            // 流正常结束
            console.log('[SSE] 流结束')
            break
          }

          buffer += decoder.decode(value, { stream: true })

          // 按行分割
          let lines = buffer.split('\n')
          buffer = lines.pop() || ''

          for (let i = 0; i < lines.length; i++) {
            const line = lines[i]

            // 空行表示一个事件结束
            if (line === '') {
              if (currentEvent && currentData) {
                try {
                  const data = JSON.parse(currentData)

                  if (currentEvent === 'message') {
                    params.onMessage({
                      text: data.text || '',
                      delta: data.delta || '',
                      id: data.id || '',
                      detail: data.detail,
                    })
                  }
                  else if (currentEvent === 'done') {
                    console.log('[SSE] 收到 done 事件')
                    params.onDone?.()
                    return // 正常结束
                  }
                  else if (currentEvent === 'error') {
                    console.log('[SSE] 收到 error 事件:', data.error)
                    params.onError?.(new Error(data.error || '未知错误'))
                    return
                  }
                }
                catch (e) {
                  console.warn('[SSE] JSON 解析失败，跳过:', currentData, e)
                }
              }
              currentEvent = ''
              currentData = ''
              continue
            }

            // 解析事件类型: event: <type>
            if (line.startsWith('event:')) {
              currentEvent = line.slice(6).trim()
              continue
            }

            // 解析数据: data: <json>
            if (line.startsWith('data:')) {
              const dataContent = line.slice(5).trim()
              // SSE 数据可能跨多行(用空行分隔)，这里简单处理单行
              if (dataContent) {
                currentData = dataContent
              }
              continue
            }

            // 如果是 data: 之后的多行数据
            if (currentData && !line.startsWith('event:') && !line.startsWith('data:')) {
              currentData += '\n' + line
            }
          }
        }

        // 处理缓冲区中可能剩余的数据
        if (buffer.trim()) {
          const lastLine = buffer.trim()
          if (lastLine.startsWith('data:')) {
            currentData = lastLine.slice(5).trim()
            if (currentEvent === 'message' && currentData) {
              try {
                const data = JSON.parse(currentData)
                params.onMessage({
                  text: data.text || '',
                  delta: data.delta || '',
                  id: data.id || '',
                  detail: data.detail,
                })
              }
              catch (e) {
                console.warn('[SSE] 缓冲区 JSON 解析失败:', e)
              }
            }
          }
        }
      }
      catch (e) {
        // 区分是 AbortError 还是其他错误
        if ((e as Error).name === 'AbortError' || params.signal?.aborted) {
          console.log('[SSE] 请求被取消')
          params.onDone?.()
        }
        else {
          console.error('[SSE] 流读取错误:', e)
          params.onError?.(e as Error)
        }
      }
    })
    .catch((error) => {
      console.error('[SSE] fetch 错误:', error)
      if (error.name === 'AbortError') {
        params.onDone?.()
      }
      else {
        params.onError?.(error)
      }
    })
}

// =============================================
// 保留原 Axios 实现，作为备用 (已注释)
//
// 原实现问题:
// 1. Axios 的 onDownloadProgress 回调中，xhr.responseText 是增量更新的
// 2. 通过 lastIndexOf('\n') 分割的方式不够可靠
// 3. 无法正确识别 SSE 事件边界
// 4. 错误处理不够完善
// =============================================

export function fetchSession<T>() {
  return post<T>({
    url: '/session',
  })
}

export function fetchVerify<T>(token: string) {
  return post<T>({
    url: '/verify',
    data: { token },
  })
}

export function fetchCode<T>(phone: string) {
  return post<T>({
    url: '/v1/sms/send/code',
    data: { phone },
  })
}

export function login<T>(phone: string, code: string) {
  return post<T>({
    url: '/v1/user/login',
    data: { user_name: phone, pwd: code, type: 1 },
  })
}

// =============================================
// user-master 服务集成 - 用户认证 API
// =============================================

export interface LoginMethodsResponse {
  gitlab: string
  wx_qrcode: string
}

export interface UserInfo {
  id: number
  name: string
  avatar_url: string
}

export interface CheckAuthResponse extends UserInfo {}

/**
 * 获取登录方式（GitLab OAuth 和微信二维码）
 * @param sys 系统标识，默认 'ai'
 */
export function fetchLoginMethods<T = LoginMethodsResponse>(sys = 'ai') {
  // 开发环境走代理，生产环境直接调用
  const baseURL = import.meta.env.DEV 
    ? '' 
    : import.meta.env.VITE_USER_CENTER?.split('?')[0] || 'http://localhost:8082'
  return fetch(`${baseURL}/api/v1/login/methods?sys=${sys}`)
    .then(res => res.json())
    .then(data => data as T)
}

/**
 * 验证 JWT Token 并获取用户信息
 * @param token JWT token
 */
export function fetchCheckAuth<T = CheckAuthResponse>(token: string) {
  // 开发环境走代理，生产环境直接调用
  const baseURL = import.meta.env.DEV 
    ? '' 
    : import.meta.env.VITE_USER_CENTER?.split('?')[0] || 'http://localhost:8082'
  return fetch(`${baseURL}/api/v1/login/check/auth?access_token=${token}`)
    .then(res => res.json())
    .then(data => data as T)
}

/**
 * 跳转到 GitLab OAuth 登录
 */
export function redirectToGitLabLogin(sys = 'ai') {
  // GitLab OAuth 配置（来自 user-master/dev.config.yaml）
  const gitlabConfig = {
    domain: 'https://gitlab.0voice.com',
    clientId: '35f0948c772e5fe99dc147dc23026b12cf569f87b3164fd4fd01ccf307603b0b',
    redirectUri: encodeURIComponent('http://localhost:8082/api/v1/login/gitlab/redirect'),
  }

  // 构建 GitLab OAuth 授权 URL
  const gitlabAuthUrl = `${gitlabConfig.domain}/oauth/authorize?client_id=${gitlabConfig.clientId}&redirect_uri=${gitlabConfig.redirectUri}&response_type=code&state=${sys}`

  console.log('GitLab 授权 URL:', gitlabAuthUrl)
  console.log('即将跳转到:', gitlabAuthUrl)

  window.location.href = gitlabAuthUrl
}
