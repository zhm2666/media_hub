# Demo - 用户登录系统

这是一个简单的前后端分离项目示例，演示了完整的用户认证流程。

## 项目结构

```
demo/
├── demo-server/          # 后端 Go 项目
│   ├── main.go          # 入口文件
│   ├── go.mod           # Go 模块文件
│   ├── handlers/        # 处理器
│   │   ├── auth.go      # 认证相关处理器
│   │   └── health.go    # 健康检查
│   ├── middleware/       # 中间件
│   │   └── auth.go      # 认证中间件
│   ├── models/          # 数据模型
│   │   └── user.go      # 用户模型
│   └── utils/           # 工具函数
│       └── jwt.go       # JWT 工具
│
└── demo-web/            # 前端 Vue3 项目
    ├── src/
    │   ├── main.ts      # 入口文件
    │   ├── App.vue      # 根组件
    │   ├── style.css    # 全局样式
    │   ├── router/      # 路由配置
    │   ├── stores/      # Pinia 状态管理
    │   ├── utils/       # 工具函数
    │   ├── views/       # 页面组件
    │   └── types/       # TypeScript 类型
    └── package.json     # 项目依赖
```

## 技术栈

### 后端
- **Go 1.21+**
- **Gin** - Web 框架
- **JWT** - JSON Web Token 认证
- **bcrypt** - 密码加密

### 前端
- **Vue 3.4** - 渐进式 JavaScript 框架
- **Vue Router 4** - Vue.js 官方路由
- **Pinia** - Vue.js 状态管理
- **Axios** - HTTP 客户端
- **TypeScript** - JavaScript 超集
- **Vite** - 下一代前端构建工具

## 功能特性

- ✅ 用户注册
- ✅ 用户登录
- ✅ JWT Token 认证
- ✅ 密码加密存储 (bcrypt)
- ✅ 路由守卫
- ✅ 持久化登录状态
- ✅ 个人资料管理
- ✅ 退出登录
- ✅ 响应式设计

## API 接口

### 认证接口

| 方法 | 路径 | 描述 | 需要认证 |
|------|------|------|----------|
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| POST | /api/auth/logout | 退出登录 | 是 |
| GET | /api/auth/me | 获取当前用户 | 是 |

### 用户接口

| 方法 | 路径 | 描述 | 需要认证 |
|------|------|------|----------|
| GET | /api/user/profile | 获取用户资料 | 是 |
| PUT | /api/user/profile | 更新用户资料 | 是 |

### 健康检查

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /health | 服务健康检查 |

## 快速开始

### 1. 启动后端服务

```bash
cd demo/demo-server

# 安装依赖
go mod tidy

# 启动服务 (默认端口 8080)
go run main.go
```

### 2. 启动前端服务

```bash
cd demo/demo-web

# 安装依赖
npm install

# 启动开发服务器 (默认端口 3000)
npm run dev
```

### 3. 访问应用

打开浏览器访问 http://localhost:3000

## 使用说明

1. **注册账号**：访问注册页面，填写用户名、邮箱、密码进行注册
2. **登录**：使用注册的账号登录系统
3. **访问首页**：登录成功后自动跳转到首页
4. **个人资料**：在个人资料页面可以更新邮箱
5. **退出登录**：点击退出登录按钮清除登录状态

## 请求/响应示例

### 注册请求

```bash
POST /api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

### 注册响应

```json
{
  "message": "注册成功",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

### 登录请求

```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

### 登录响应

```json
{
  "message": "登录成功",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

## 开发说明

### 前端代理配置

前端开发服务器配置了代理，将 `/api` 请求代理到后端服务：

```typescript
// vite.config.ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
},
```

### Token 存储

- Token 存储在浏览器的 LocalStorage 中
- 请求时通过 `Authorization: Bearer <token>` 头部发送
- Token 有效期为 24 小时

### 路由守卫

- 未登录用户访问需要认证的页面会跳转到登录页
- 已登录用户访问登录/注册页会跳转到首页

## 生产部署

### 后端部署

```bash
cd demo/demo-server

# 构建
go build -o server main.go

# 运行
./server
```

### 前端部署

```bash
cd demo/demo-web

# 构建生产版本
npm run build

# 构建产物在 dist 目录
```

## 注意事项

1. 本示例使用内存存储用户数据，重启服务后会丢失
2. 生产环境请更换 JWT 密钥
3. 生产环境请使用数据库存储用户数据
4. 建议添加 HTTPS 支持
5. 建议添加请求频率限制

## License

MIT
