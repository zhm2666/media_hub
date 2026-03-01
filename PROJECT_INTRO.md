# MediaHub 项目介绍文档

## 一、项目概述

MediaHub 是一个基于 Go 语言和 Vue.js 构建的媒体资源管理平台，主要提供图片上传托管和短链接生成服务。该项目采用现代化的前后端分离架构，支持 Docker 容器化部署，适用于个人或企业的媒体资源管理需求。

---

## 二、技术架构

### 2.1 后端技术栈

后端服务采用 Go 1.20 开发，使用了以下核心技术：

| 技术组件 | 版本 | 说明 |
|---------|------|------|
| Gin | v1.10.0 | 高性能 HTTP Web 框架 |
| MySQL | 8.x | 主数据存储 |
| Redis | v9.7.0 | 缓存数据库 |
| 腾讯云 COS | v0.7.55 | 对象存储服务 |
| gRPC | v1.62.1 | RPC 通信框架 |
| Viper | v1.19.0 | 配置管理 |
| Logrus | v1.9.3 | 日志系统 |
| Lumberjack | v2.2.1 | 日志轮转 |

### 2.2 前端技术栈

前端采用现代化的 Vue 3 技术栈：

| 技术组件 | 版本 | 说明 |
|---------|------|------|
| Vue | v3.3.4 | 渐进式 JavaScript 框架 |
| Element Plus | v2.4.1 | 企业级 Vue 3 组件库 |
| Vite | v4.4.5 | 下一代前端构建工具 |
| TypeScript | v5.0.2 | 类型安全的 JavaScript |
| Axios | v1.5.1 | HTTP 客户端 |

### 2.3 基础设施

- **容器化**: Docker - 应用打包和部署
- **编排**: Docker Swarm - 服务编排（支持多副本部署）
- **镜像仓库**: Harbor - 镜像仓库管理

---

## 三、核心功能

### 3.1 媒体文件上传

系统支持用户上传图片文件，目前支持以下格式：
- JPG/JPEG
- PNG
- GIF

**上传流程**：
1. 用户通过前端界面上传图片文件
2. 后端对文件格式进行校验
3. 计算文件 MD5 值作为唯一标识
4. 上传到腾讯云 COS 对象存储
5. 调用短链接服务生成短 URL
6. 返回文件访问短链接

### 3.2 短链接生成

每当用户上传文件后，系统会自动将文件 URL 转换为短链接，支持：
- **公开链接**：无需登录即可访问
- **私有链接**：需登录认证，仅用户本人可见

### 3.3 首页展示

首页采用随机展示机制：
- 顶部 Banner 轮播图：随机展示 3 张图片
- 下方图片瀑布流：随机展示 10 张图片（分两列展示）
- 每次刷新页面展示内容都会变化

### 3.4 用户认证

系统集成了 JWT 认证机制：
- CORS 中间件处理跨域请求
- Auth 中间件验证用户身份
- 支持公开访问和登录访问两种模式

---

## 四、项目结构

```
mediahub-master/
├── mediahub/                    # 后端服务 (Go)
│   ├── controller/              # 控制器层
│   │   ├── file.go              # 文件上传控制器
│   │   └── home.go              # 首页控制器
│   ├── middleware/              # 中间件
│   │   ├── auth.go              # 认证中间件
│   │   └── cors.go              # 跨域中间件
│   ├── pkg/                     # 公共包
│   │   ├── config/              # 配置管理
│   │   ├── constants/           # 常量定义
│   │   ├── db/                   # 数据库封装
│   │   │   ├── mysql/            # MySQL 驱动
│   │   │   └── redis/            # Redis 驱动
│   │   ├── grpc-client-pool/     # gRPC 连接池
│   │   ├── log/                  # 日志系统
│   │   ├── storage/              # 存储抽象
│   │   │   └── cos/              # 腾讯云 COS
│   │   ├── utils/                # 工具函数
│   │   └── zerror/               # 错误处理
│   ├── routers/                  # 路由定义
│   ├── services/                 # 业务逻辑层
│   │   └── shorturl/             # 短链接服务客户端
│   ├── main.go                   # 入口文件
│   ├── dev.config.yaml           # 开发配置
│   └── test.config.yaml          # 测试配置
│
├── mediahub-web/                 # 前端应用 (Vue 3)
│   ├── src/
│   │   ├── api/                  # API 接口定义
│   │   ├── components/           # 公共组件
│   │   ├── request/               # HTTP 请求封装
│   │   ├── utils/                 # 工具函数
│   │   ├── views/                 # 页面视图
│   │   ├── App.vue                # 根组件
│   │   └── main.ts                # 入口文件
│   ├── package.json              # 依赖配置
│   ├── vite.config.ts             # Vite 配置
│   └── Dockerfile                 # 前端构建镜像
│
├── sql/                          # 数据库脚本
│   └── create_db.sql             # 数据库初始化脚本
│
├── docs/                         # 文档资料
│   ├── mediahub系统架构图.png
│   ├── mediahub技术架构图.png
│   ├── mediahub整体架构图.png
│   ├── mediahub业务流程图1.png
│   ├── mediahub业务流程图2.png
│   ├── crontab语法.png
│   └── [其他技术文档 PDF]
│
├── Dockerfile                     # 主 Dockerfile（多阶段构建）
├── README.md                      # 项目说明
└── curl-amd64                     # 健康检查工具
```

---

## 五、数据库设计

### 5.1 url_map 表（公共 URL 映射表）

存储所有公开的短链接映射关系：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT(20) | 主键 ID，自增 |
| short_key | VARCHAR(45) | 短链接 key |
| original_url | VARCHAR(512) | 原始 URL |
| times | INT | 访问次数 |
| create_at | BIGINT(64) | 创建时间戳 |
| update_at | BIGINT(64) | 更新时间戳 |

### 5.2 url_map_user 表（用户 URL 映射表）

存储用户个人的短链接：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT(20) | 主键 ID，自增 |
| user_id | BIGINT(20) | 用户 ID |
| short_key | VARCHAR(45) | 短链接 key |
| original_url | VARCHAR(512) | 原始 URL |
| times | INT | 访问次数 |
| create_at | BIGINT(64) | 创建时间戳 |
| update_at | BIGINT(64) | 更新时间戳 |

---

## 六、API 接口

### 6.1 文件上传

```
POST /api/v1/file/upload
Content-Type: multipart/form-data

参数:
- file: 文件数据（图片）

响应:
{
    "url": "生成的短链接 URL"
}
```

### 6.2 首页数据

```
GET /api/v1/home

响应:
{
    "banners": ["banner1.jpg", "banner2.jpg", "banner3.jpg"],
    "images1": ["img1.jpg", "img2.jpg", "img3.jpg", "img4.jpg", "img5.jpg"],
    "images2": ["img6.jpg", "img7.jpg", "img8.jpg", "img9.jpg", "img10.jpg"]
}
```

### 6.3 健康检查

```
GET /health

响应: HTTP 200 OK
```

---

## 七、配置文件说明

### 开发环境配置 (dev.config.yaml)

```yaml
http:
  ip: 0.0.0.0
  port: 8080
  mode: debug

redis:
  host: "192.168.239.161"
  port: 6379
  pwd: "123456"

mysql:
  dsn: "root:123456@tcp(192.168.239.161:3306)/mediahub?collation=utf8mb4_unicode_ci&charset=utf8mb4"
  maxLifeTime: 3600
  maxOpenConn: 10
  maxIdleConn: 10

log:
  level: "info"
  logPath: "runtime/logs/app.log"

cos:
  secretId: AK
  secretKey: HA
  cdnDomain: "https://mediahubdev.0voice.com"
  bucketUrl: "https://mediahubdev-1256487221.cos.ap-guangzhou.myqcloud.com"

dependOn:
  shortUrl:
    address: "localhost:50051"
    accessToken: "4cdQRSk678lTe015hpqWXYZAab2VrwxyGHIJKjmn"
  user:
    address: "http://localhost:8082"
```

---

## 八、Docker 部署

### 8.1 多阶段构建

项目使用 Docker 多阶段构建技术，同时构建前端和后端：

1. **Stage 0**: 编译 Go 后端服务
2. **Stage 1**: 构建 Vue 3 前端应用
3. **Stage 2**: 整合最终镜像

### 8.2 部署步骤

**构建镜像**：
```bash
docker build -t 2410/mediahub:1.0.0 -t 192.168.239.161:5000/2410/mediahub:1.0.0 .
```

**推送镜像**：
```bash
docker push 192.168.239.161:5000/2410/mediahub:1.0.0
```

**创建配置**：
```bash
docker config create 2410-mediahub-conf dev.config.yaml
```

**启动服务**：
```bash
docker service create --name 2410-mediahub \
-p 8080:8080 \
--config src=2410-mediahub-conf,target=/app/config.yaml \
--replicas 2 \
--health-cmd "curl -f http://localhost:8080/health" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
--network 2410-shorturl-net \
192.168.239.161:5000/2410/mediahub:1.0.0
```

---

## 九、系统特性

### 9.1 高可用性
- 支持 Docker Swarm 多副本部署
- 配置健康检查机制
- 自动故障恢复

### 9.2 安全性
- JWT 用户认证
- CORS 跨域保护
- 文件格式校验

### 9.3 性能优化
- Redis 缓存加速
- CDN 内容分发加速（腾讯云）
- gRPC 高效通信

### 9.4 可扩展性
- 模块化设计
- 配置驱动
- 支持服务拆分

### 9.5 易用性
- 简洁的前端界面
- 拖拽式上传
- 响应式设计

---

## 十、依赖外部服务

| 服务 | 说明 |
|------|------|
| 腾讯云 COS | 对象存储 |
| 腾讯云 CDN | 内容分发网络 |
| gRPC 短链接服务 | 独立部署的短链接服务 |
| MySQL | 数据持久化 |
| Redis | 缓存服务 |

---

## 十一、总结

MediaHub 是一个功能完整的媒体资源管理平台，展示了现代化的 Go 后端开发实践和 Vue 3 前端开发技术的有机结合。项目具备以下亮点：

1. **前后端分离**：清晰的架构边界
2. **微服务思维**：通过 gRPC 与短链接服务通信
3. **云原生支持**：完整的 Docker 部署方案
4. **最佳实践**：日志、配置、错误处理等完善

项目代码结构清晰，模块划分合理，是学习微服务架构和前后端分离开发的好参考。

---

*文档生成时间：2026-02-27*
