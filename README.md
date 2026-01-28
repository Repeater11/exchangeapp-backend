# Forum Backend

一个基于 Go 的论坛后端服务，实现了注册/登录、发帖/回帖、点赞等完整链路，并包含分页、缓存与异步回写等工程化设计。

## 功能概览
- 用户：注册 / 登录（JWT）
- 帖子：创建 / 列表 / 详情 / 更新 / 删除
- 回复：创建 / 列表 / 更新 / 删除
- 点赞：赞 / 取消赞 / 点赞状态
- 分页：offset 与 cursor 两种方式（推荐 cursor）
- 健康检查：`/healthz`

## 技术栈
- Go + Gin
- GORM + MySQL
- Redis
- JWT

## 目录结构
- `cmd/server`：服务入口
- `internal/app`：服务装配与生命周期（含后台 worker）
- `internal/handler`：HTTP 层
- `internal/service`：业务逻辑
- `internal/repository`：数据访问与缓存
- `internal/models` / `internal/dto`：模型与 DTO
- `configs/config.yml`：默认配置
- `docs/openapi.yaml`：接口文档

## 本地运行
### 1. 启动依赖
准备 MySQL 与 Redis（本地或容器均可）。

### 2. 启动服务
```bash
go run ./cmd/server
```
默认读取 `configs/config.yml`，也可通过环境变量覆盖（见下文）。

## 配置说明
配置文件：`configs/config.yml`

示例：
```yml
app:
  name: "forum-backend"
  port: 3000

database:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: ""
  name: "exchangeapp"

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

jwt:
  secret: "your_secret"
  expire_minutes: 60

like_worker:
  batch: 200
  interval_seconds: 1
```

环境变量前缀：`EXCHANGEAPP_`，支持覆盖配置文件字段：
```bash
export EXCHANGEAPP_DATABASE_HOST=127.0.0.1
export EXCHANGEAPP_DATABASE_PORT=3306
export EXCHANGEAPP_DATABASE_USER=root
export EXCHANGEAPP_DATABASE_PASSWORD=your_password
export EXCHANGEAPP_DATABASE_NAME=exchangeapp

export EXCHANGEAPP_REDIS_ADDR=127.0.0.1:6379
export EXCHANGEAPP_REDIS_PASSWORD=
export EXCHANGEAPP_REDIS_DB=0

export EXCHANGEAPP_JWT_SECRET=your_jwt_secret
export EXCHANGEAPP_JWT_EXPIRE_MINUTES=60

export EXCHANGEAPP_LIKE_WORKER_BATCH=200
export EXCHANGEAPP_LIKE_WORKER_INTERVAL_SECONDS=1
```

## 分页说明
### 1) cursor 分页（推荐）
- cursor 格式：`<created_at_unix_nano>_<id>`
- 使用方式：
```
GET /threads?size=20
=> 返回 next_cursor

GET /threads?size=20&cursor=1700000000000000000_123
```

### 2) offset 分页
- 仍可用，但数据量大时性能会明显下降
- 推荐在前端统一使用 cursor

## 点赞计数策略
- 点赞写入：只更新 Redis 计数 + 标记 dirty
- 后台 worker 定期回写 MySQL（最终一致）
- 可通过 `like_worker.batch` / `like_worker.interval_seconds` 调整回写频率与批量大小

## 性能优化要点
- 列表改为游标分页，避免 offset 深分页性能退化
- 建立与排序一致的复合索引（`created_at desc, id desc`）
- 详情页使用 Redis 缓存 + singleflight 防击穿
- 点赞计数采用 Redis 写入 + 异步回写 MySQL

## 测试
```bash
GOCACHE=./.gocache go test ./...
```

## 接口概览
- `POST /register` 用户注册
- `POST /login` 用户登录
- `GET /threads` 帖子列表（支持 cursor / page）
- `GET /threads/:id` 帖子详情
- `GET /threads/:id/replies` 回复列表
- `POST /api/threads` 发帖（需登录）
- `POST /api/threads/:id/replies` 回复（需登录）
- `POST /api/threads/:id/like` 点赞（需登录）
- `DELETE /api/threads/:id/like` 取消点赞（需登录）

完整接口见：`docs/openapi.yaml`

## 备注
- 数据库迁移在启动时自动执行（AutoMigrate）
- 点赞数最终一致（worker 异步回写）
