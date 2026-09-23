# Panel

书签导航面板，基于 Go + Vue 3 构建的全栈 Web 应用。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.25、Gin、GORM、SQLite、JWT |
| 前端 | Vue 3、Vite 5、Element Plus、Pinia、Vue Router |
| 部署 | Docker 多阶段构建、Alpine Linux |

## 功能特性

- 书签分组管理与拖拽排序
- 用户认证与权限控制（JWT）
- 书签图标自动获取
- 响应式前端界面
- Docker 一键部署，支持自定义用户权限（PUID/PGID）

## 快速开始

### Docker 部署（推荐）

```bash
docker run -d \
  --name panel \
  -p 5678:5678 \
  -v panel-data:/app/data \
  -e JWT_SECRET=your-secret-key \
  jiangming/panel:latest
```

访问 `http://localhost:5678` 即可使用。

### Docker Compose

```yaml
version: "3"
services:
  panel:
    image: jiangming/panel:latest
    container_name: panel
    ports:
      - "5678:5678"
    volumes:
      - panel-data:/app/data
    environment:
      - JWT_SECRET=your-secret-key
      - GIN_MODE=release
    restart: unless-stopped

volumes:
  panel-data:
```

### 本地开发

**前提条件**：Go 1.25+、Node.js 22+

```bash
# 克隆仓库
git clone https://github.com/jiangbaoming/panel.git
cd panel

# 启动后端
go mod download
go run ./cmd/server

# 启动前端（新终端）
cd web
npm install
npm run dev
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `5678` | 服务监听端口 |
| `GIN_MODE` | `release` | Gin 运行模式（debug/release） |
| `JWT_SECRET` | `dev-jwt-secret-change-in-production` | JWT 签名密钥，生产环境务必修改 |
| `DATA_DIR` | `./data` | 数据存储目录 |
| `UPLOAD_DIR` | `./data/uploads` | 上传文件目录 |
| `DIST_DIR` | `./dist` | 前端静态资源目录 |
| `PUID` | `1000` | 容器内运行用户 UID |
| `PGID` | `1000` | 容器内运行用户 GID |
| `LOG_LEVEL` | `info` | 日志级别（debug/info/warn/error） |
| `LOG_JSON` | `false` | 是否启用 JSON 格式日志 |

## 构建镜像

```bash
# 默认 latest 标签
./build.sh

# 指定版本标签
./build.sh v1.0.0
```

## 项目结构

```
panel/
├── cmd/server/        # 后端入口
├── config/            # 配置加载
├── db/                # 数据库初始化
├── logger/            # 日志封装
├── middleware/        # 中间件（认证等）
├── model/             # 数据模型
├── router/            # 路由注册
├── service/           # 业务逻辑
├── web/               # 前端源码
│   ├── src/           # Vue 组件与页面
│   └── dist/          # 构建产物（.gitignore 排除）
├── Dockerfile         # 多阶段构建
├── entrypoint.sh      # 容器入口脚本
└── build.sh           # 镜像构建脚本
```

## 许可证

MIT