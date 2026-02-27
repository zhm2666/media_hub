# MediaHub 新手入门指南

## 前言

欢迎使用 MediaHub！本指南专为第一次接触该项目的开发者编写，帮助你快速了解项目并搭建开发环境。

---

## 一、准备工作

### 1.1 需要安装的软件

在开始之前，你需要安装以下软件：

| 软件 | 版本要求 | 说明 |
|------|---------|------|
| **Go** | 1.20+ | 后端开发语言 |
| **Node.js** | 18.x | 前端开发环境 |
| **MySQL** | 8.0+ | 数据库 |
| **Redis** | 6.0+ | 缓存数据库 |
| **Git** | 任意版本 | 代码版本管理 |
| **Docker** | 最新版 | 容器化部署（可选） |
| **VS Code** | 任意版本 | 推荐 IDE |

### 1.2 安装指南

**Go 安装（Windows）**：
1. 访问 https://go.dev/dl/ 下载 Windows 安装包
2. 运行安装程序，勾选 "Add to PATH"
3. 打开终端，输入 `go version` 验证安装

**Node.js 安装**：
1. 访问 https://nodejs.org/ 下载 LTS 版本
2. 运行安装程序
3. 打开终端，输入 `node -v` 和 `npm -v` 验证安装

**MySQL 安装**：
1. 访问 https://dev.mysql.com/downloads/mysql/ 下载
2. 或使用 MySQL Installer
3. 记住设置的 root 密码

**Redis 安装**：
1. 访问 https://redis.io/download 下载 Windows 版本
2. 或使用 Memurai/Redis Windows 替代品
3. 记住设置的密码

---

## 二、获取项目

### 2.1 克隆项目

打开终端，执行以下命令：

```bash
# 切换到项目目录（根据你的实际路径）
cd F:\GOMaster\AI-chat

# 克隆项目
git clone https://github.com/your-repo/mediahub-master.git

# 进入项目目录
cd mediahub-master
```

### 2.2 目录结构初览

```
mediahub-master/
├── mediahub/        # 后端代码（Go）
├── mediahub-web/    # 前端代码（Vue）
├── sql/             # 数据库脚本
├── docs/            # 文档资料
├── Dockerfile       # Docker 构建文件
└── README.md       # 项目说明
```

---

## 三、数据库配置

### 3.1 创建数据库

1. 打开 MySQL 客户端（如 MySQL Workbench、Navicat 或命令行）

2. 执行 `sql/create_db.sql` 文件中的 SQL 语句：

```sql
CREATE SCHEMA `mediahub` DEFAULT CHARACTER SET utf8mb4;
USE mediahub;

-- 这会自动创建 url_map 和 url_map_user 两张表
```

### 3.2 验证数据库

```sql
-- 查看数据库
SHOW DATABASES;

-- 使用数据库
USE mediahub;

-- 查看表
SHOW TABLES;
```

你应该能看到 `url_map` 和 `url_map_user` 两张表。

---

## 四、后端配置

### 4.1 安装 Go 依赖

```bash
cd mediahub

# 初始化 Go 模块（如果需要）
go mod tidy

# 下载依赖
go mod download
```

### 4.2 修改配置文件

编辑 `mediahub/dev.config.yaml`，修改为你本地的配置：

```yaml
http:
  ip: 0.0.0.0
  port: 8080
  mode: debug

redis:
  host: "localhost"      # 修改为你的 Redis 地址
  port: 6379
  pwd: "123456"         # 修改为你的 Redis 密码

mysql:
  dsn: "root:123456@tcp(localhost:3306)/mediahub?collation=utf8mb4_unicode_ci&charset=utf8mb4"
  # 修改为你的 MySQL 用户名:密码@地址/数据库名

log:
  level: "debug"        # 开发阶段建议用 debug
  logPath: "runtime/logs/app.log"

cos:
  secretId: "your-secret-id"      # 替换为你的腾讯云密钥
  secretKey: "your-secret-key"   # 替换为你的腾讯云密钥
  cdnDomain: "https://your-cdn.com"
  bucketUrl: "https://your-bucket.cos.region.myqcloud.com"

dependOn:
  shortUrl:
    address: "localhost:50051"
    accessToken: "your-token"
  user:
    address: "http://localhost:8082"
```

> ⚠️ **注意**：如果你没有腾讯云 COS 账号，可以暂时注释掉 COS 相关配置，或者使用本地存储替代。

### 4.3 启动后端服务

```bash
cd mediahub

# 运行项目
go run main.go --config dev.config.yaml
```

如果看到类似以下输出，说明启动成功：

```
[GIN-debug] POST   /api/v1/file/upload      ...
[GIN-debug] GET    /api/v1/home             ...
[GIN-debug] GET    /health                  ...
[GIN-debug] Environment: debug
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8080
```

---

## 五、前端配置

### 5.1 安装依赖

```bash
cd mediahub-web

# 安装项目依赖
npm install
```

### 5.2 启动开发服务器

```bash
# 启动开发模式
npm run dev
```

启动成功后，打开浏览器访问 `http://localhost:5173`（端口可能不同，看终端输出）。

### 5.3 构建生产版本

```bash
# 构建前端
npm run build
```

构建后的文件会生成在 `mediahub-web/dist` 目录。

---

## 六、快速体验功能

### 6.1 访问首页

打开浏览器，访问：`http://localhost:8080`

你应该能看到：
- 顶部 Banner 轮播图
- 随机展示的图片列表

### 6.2 上传图片

1. 找到上传按钮
2. 选择一张 JPG/PNG/GIF 图片
3. 上传成功后，会返回短链接

### 6.3 测试 API

使用 Postman 或 curl 测试接口：

```bash
# 健康检查
curl http://localhost:8080/health

# 获取首页数据
curl http://localhost:8080/api/v1/home
```

---

## 七、代码结构解析

### 7.1 后端目录结构

```
mediahub/
├── main.go              # 程序入口
├── controller/          # 控制器（处理请求）
│   ├── file.go          # 文件上传
│   └── home.go          # 首页数据
├── middleware/          # 中间件
│   ├── auth.go          # 认证
│   └── cors.go          # 跨域
├── pkg/                 # 公共包
│   ├── config/          # 配置管理
│   ├── db/              # 数据库
│   ├── log/             # 日志
│   ├── storage/         # 存储
│   └── utils/           # 工具函数
├── routers/             # 路由定义
└── services/           # 业务逻辑
```

### 7.2 前端目录结构

```
mediahub-web/src/
├── main.ts              # 入口文件
├── App.vue              # 根组件
├── api/                 # API 接口
├── views/               # 页面视图
├── components/          # 组件
└── request/            # 请求封装
```

### 7.3 请求流程

```
用户操作
   ↓
前端 (Vue) → HTTP 请求
   ↓
后端 (Gin) → 路由 → 控制器 → 服务 → 数据库/存储
   ↓
响应 ← JSON 数据
   ↓
前端渲染
```

---

## 八、常见问题

### 8.1 数据库连接失败

**问题**：报错 `Connection refused`

**解决**：
1. 确认 MySQL 服务已启动
2. 检查 `dev.config.yaml` 中的 MySQL 配置
3. 确认数据库名称 `mediahub` 已创建

### 8.2 Redis 连接失败

**问题**：报错 `Redis connection refused`

**解决**：
1. 确认 Redis 服务已启动
2. 检查 `dev.config.yaml` 中的 Redis 配置
3. 验证密码是否正确

### 8.3 前端无法访问后端

**问题**：CORS 跨域错误

**解决**：
1. 后端已配置 CORS 中间件
2. 检查后端是否正常运行
3. 确认端口 8080 未被占用

### 8.4 Go 依赖下载慢

**问题**：下载依赖超时

**解决**：
```bash
# 设置国内镜像
go env -w GOPROXY=https://goproxy.cn,direct
```

### 8.5 npm 安装依赖慢

**问题**：npm 安装超时

**解决**：
```bash
# 设置淘宝镜像
npm config set registry https://mirrors.huaweicloud.com/repository/npm/
```

---

## 九、下一步学习建议

### 9.1 理解核心代码

建议按以下顺序阅读代码：

1. **main.go** - 了解程序入口和初始化流程
2. **routers/routers.go** - 了解所有 API 路由
3. **controller/file.go** - 了解文件上传逻辑
4. **pkg/storage/cos/** - 了解云存储集成

### 9.2 尝试修改功能

- 修改首页展示的图片数量
- 添加新的 API 接口
- 更换图片存储位置

### 9.3 学习相关技术

- **Gin 框架**：https://gin-gonic.com/
- **Vue 3 文档**：https://vuejs.org/
- **Docker 入门**：https://www.docker.com/

---

## 十、Docker 部署（进阶）

如果你想使用 Docker 运行项目：

### 10.1 构建镜像

```bash
docker build -t mediahub:latest .
```

### 10.2 运行容器

```bash
docker run -d -p 8080:8080 \
  -v ./dev.config.yaml:/app/config.yaml \
  mediahub:latest
```

---

## 十一、联系与支持

如果在学习过程中遇到问题：

1. 查看 `docs/` 目录下的架构图和技术文档
2. 阅读 `README.md` 了解部署说明
3. 检查代码中的注释和日志

---

## 总结

恭喜你完成了 MediaHub 的入门学习！现在你应该已经：

- ✅ 安装了所有必需的开发工具
- ✅ 克隆并获取了项目代码
- ✅ 配置并启动了数据库和缓存
- ✅ 运行了后端服务
- ✅ 运行了前端应用
- ✅ 体验了核心功能
- ✅ 了解了代码结构

祝你开发愉快！🎉

---

*文档生成时间：2026-02-27*
*适合版本：MediaHub v1.0.0*
