# user-service

基于 Go 语言开发的用户管理与地理围栏好友发现微服务，采用 Gin + GORM 架构，具备完整的可观测性（日志、指标、链路追踪）和安全防护体系。

---

## 目录

- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
  - [环境要求](#环境要求)
  - [本地编译运行](#本地编译运行)
  - [Docker Compose 一键启动](#docker-compose-一键启动)
- [配置说明](#配置说明)
  - [基础配置](#基础配置)
  - [环境变量覆盖](#环境变量覆盖)
- [API 接口](#api-接口)
- [可观测性](#可观测性)
  - [日志](#日志)
  - [指标监控](#指标监控)
  - [链路追踪](#链路追踪)
- [安全设计](#安全设计)
- [测试](#测试)
- [部署](#部署)
- [运维手册](#运维手册)

---

## 功能特性

- **用户管理**：注册、登录、查询、修改、删除，JWT Bearer 认证
- **好友系统**：双向添加好友、好友列表分页查询
- **附近的人**：基于 Geohash 算法的地理位置好友发现
- **安全防护**：bcrypt 密码哈希、安全响应头、CORS、请求限流、请求体大小限制
- **可观测性**：Zap 结构化滚动日志、Prometheus 指标、OpenTelemetry 链路追踪
- **生产就绪**：健康检查探针、优雅关闭、数据库连接池、TLS 支持
- **容器化**：多阶段 Docker 构建、Kubernetes 部署清单、Helm Chart

---

## 技术栈

| 类别 | 技术 |
|------|------|
| 语言 | Go 1.25 |
| Web 框架 | Gin v1.12 |
| ORM | GORM v2 + MySQL 8.0 |
| 认证 | JWT (golang-jwt/v5) + bcrypt |
| 参数校验 | go-playground/validator |
| 日志 | Zap + lumberjack 滚动切片 |
| 指标 | Prometheus client_golang |
| 追踪 | OpenTelemetry (OTLP gRPC / Stdout) |
| 配置 | Viper + embed 嵌入配置 + 环境变量覆盖 |
| API 文档 | Swagger (swaggo/swag) |

---

## 项目结构

```
├── cmd/
│   ├── main.go                    # 应用入口，启动流程编排
│   └── makefile                   # 编译/测试/清理 Makefile
├── constant/
│   └── ErrorCode.go               # 统一错误码定义
├── dao/                           # 数据访问层（GORM 实现）
│   ├── UserDao.go                 # 用户 DAO 接口与实现
│   └── FriendsDao.go              # 好友 DAO 接口与实现
├── deploy/                        # 可观测性配置
│   ├── prometheus/rules.yml       # Prometheus 告警规则
│   └── grafana/user-service-dashboard.json  # Grafana 仪表盘
├── docs/                          # Swagger 自动生成文档
├── helm/user-service/             # Helm Chart
├── k8s/                           # Kubernetes 部署清单
│   ├── namespace.yaml             # 命名空间
│   ├── serviceaccount.yaml        # 服务账号
│   ├── configmap.yaml             # 配置映射
│   ├── secret.yaml                # 密钥
│   ├── deployment.yaml            # 部署定义
│   ├── service.yaml               # 服务暴露
│   ├── hpa.yaml                   # 水平自动扩缩
│   ├── pdb.yaml                   # Pod 中断预算
│   └── kustomization.yaml         # Kustomize 入口
├── model/                         # 数据模型（GORM 实体 + DTO）
│   ├── User.go
│   └── Friends.go
├── pkg/                           # 通用工具包
│   ├── config/                    # 配置管理（Viper + embed）
│   ├── dbconn/                    # MySQL 连接（含重试）
│   ├── logger/                    # Logger 接口 + Zap 实现
│   └── util/                      # 响应结构、错误类型、JSON 时间
├── service/                       # 业务逻辑 + 中间件 + 路由
│   ├── Service.go                 # DI 容器（Service 结构体）
│   ├── Routes.go                  # 路由注册与中间件链
│   ├── middleware.go              # RequestID、CORS、限流、安全头、审计日志
│   ├── auth.go                    # JWT 签发与解析
│   ├── auth_helper.go             # bcrypt 密码哈希
│   ├── UserService.go             # 用户 CRUD 处理
│   ├── FriendsServer.go           # 好友/附近好友处理
│   ├── request.go                 # 请求 DTO + 校验
│   ├── ResponseHandler.go         # 响应封装、日志中间件、健康探针
│   ├── metrics.go                 # Prometheus 指标定义
│   ├── dbmetrics.go               # 数据库连接池指标采集器
│   └── tracing.go                 # OpenTelemetry 初始化
├── sql/
│   └── init.sql                   # 数据库初始化 DDL
├── test/                          # 单元测试（mock 实现）
├── docker-compose.yml             # 本地开发环境
├── Dockerfile                     # 多阶段构建
├── .dockerignore
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## 快速开始

### 环境要求

- Go 1.25+
- MySQL 8.0+
- （可选）Docker & Docker Compose

### 本地编译运行

```bash
# 1. 克隆项目
git clone https://github.com/binbin3828/user.git
cd user

# 2. 修改数据库连接配置（或通过环境变量设置）
vi pkg/config/config.yaml

# 3. 初始化数据库表
mysql -u root -p < sql/init.sql

# 4. 编译
cd cmd && make build

# 5. 启动
./UserServer
```

### Docker Compose 一键启动

```bash
# 启动 MySQL + user-service（自动建表）
docker-compose up -d

# 查看日志
docker-compose logs -f user-service

# 验证
curl http://localhost:8080/healthz

# 停止
docker-compose down
```

---

## 配置说明

### 基础配置

配置文件嵌入在二进制中（`pkg/config/config.yaml`）：

```yaml
mysql:
  driveName: mysql
  dataSourceName: ${MYSQL_USER:root}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST:127.0.0.1}:${MYSQL_PORT:3306})/${MYSQL_DATABASE:bobby_test}?charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true&timeout=5s&readTimeout=2s&writeTimeout=2s
  maxIdle: 25
  maxOpen: 50
  maxLifetime: 1

cors:
  allowedOrigins: "*"

jwt:
  secret: dev-change-this-in-production

tls:
  enabled: false
  certFile: "cert.pem"
  keyFile: "key.pem"
```

### 环境变量覆盖

| 环境变量 | 说明 | 示例 |
|----------|------|------|
| `MYSQL_DSN` | 完整 DSN（优先级最高，覆盖其他 MySQL 变量） | `user:pass@tcp(host:3306)/db?...` |
| `MYSQL_HOST` | 数据库地址（默认 127.0.0.1） | `mysql.example.com` |
| `MYSQL_PORT` | 数据库端口（默认 3306） | `3306` |
| `MYSQL_USER` | 数据库用户名（默认 root） | `xgame` |
| `MYSQL_PASSWORD` | 数据库密码（**生产必设**） | `mysecret` |
| `MYSQL_DATABASE` | 数据库名（默认 bobby_test） | `bobby_test` |
| `JWT_SECRET` | JWT 签名密钥（**生产必设**） | `openssl rand -base64 32` |
| `LOG_LEVEL` | 日志级别：debug / info / warn / error | `info` |
| `REDIS_ADDR` | Redis 地址（设了则启用 Redis 分布式限流） | `redis:6379` |
| `REDIS_PASSWORD` | Redis 密码（可选） | `""` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP 导出地址（为空则输出到 stdout） | `localhost:4317` |
| `OTEL_TRACE_SAMPLE_RATE` | 采样率 0.0-1.0（默认 0.1） | `0.5` |

> Redis 未配置时使用内存级限流（单实例可用），配置后自动切换为 Redis 分布式限流（多实例共享计数）。

---

## API 接口

### 系统端点（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/healthz` | 存活探针 |
| GET | `/readyz` | 就绪探针（含数据库检查） |
| GET | `/metrics` | Prometheus 指标暴露 |
| GET | `/swagger/*any` | Swagger UI 文档 |

### 业务端点

| 方法 | 路径 | 认证 | 限流 | 说明 |
|------|------|------|------|------|
| POST | `/v1/auth/login` | 否 | 是（10次/分钟/IP） | 登录获取 JWT |
| POST | `/v1/auth/forgot-password` | 否 | 否 | 忘记密码，发送重置令牌 |
| POST | `/v1/auth/reset-password` | 否 | 否 | 凭令牌重置密码 |
| GET | `/v1/user/:uid` | Bearer | 否 | 获取用户信息 |
| POST | `/v1/user` | Bearer | 否 | 创建用户 |
| PUT | `/v1/user` | Bearer | 否 | 修改用户信息 |
| DELETE | `/v1/user/:uid` | Bearer | 否 | 删除用户（仅本人） |
| POST | `/v1/friends` | Bearer | 否 | 直接添加好友（双向） |
| GET | `/v1/friends/:uid` | Bearer | 否 | 好友列表（分页） |
| GET | `/v1/nearbyfriends/:uid` | Bearer | 否 | 附近好友（分页，Geohash） |
| GET | `/v1/nearby-users/:uid` | Bearer | 否 | 附近陌生人推荐（分页，Geohash） |
| POST | `/v1/friend-requests` | Bearer | 否 | 发起好友请求 |
| GET | `/v1/friend-requests/incoming` | Bearer | 否 | 收到的好友请求（分页） |
| GET | `/v1/friend-requests/outgoing` | Bearer | 否 | 发出的好友请求（分页） |
| PUT | `/v1/friend-requests/:id/accept` | Bearer | 否 | 同意好友请求 |
| PUT | `/v1/friend-requests/:id/reject` | Bearer | 否 | 拒绝好友请求 |
| POST | `/v1/blacklist` | Bearer | 否 | 拉黑用户 |
| DELETE | `/v1/blacklist/:uid` | Bearer | 否 | 取消拉黑 |
| GET | `/v1/blacklist` | Bearer | 否 | 黑名单列表（分页） |
| GET | `/v1/users/online` | Bearer | 否 | 批量查询用户在线状态 |

### 统一响应格式

**成功：**

```json
{
  "code": 0,
  "data": { ... }
}
```

**分页成功：**

```json
{
  "code": 0,
  "data": [ ... ],
  "pagination": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

**错误：**

```json
{
  "code": -1,
  "msg": "param uid not set"
}
```

**错误码：**

| 错误码 | 含义 | HTTP 状态码 |
|--------|------|------------|
| 0 | 成功 | 200 |
| -1 | 参数错误 | 400 |
| -2 | 服务内部错误 | 500 |
| -3 | 认证失败 | 401 |
| -4 | 权限不足 | 403 |
| -5 | 已经是好友 | 409 |
| -6 | 好友请求已存在 | 409 |
| -7 | 好友请求不存在 | 404 |
| -8 | 好友请求已非待处理状态 | 409 |
| -9 | 已被对方拉黑 | 409 |

---

## 可观测性

### 日志

- 使用 **Zap** JSON 格式输出，注入 trace_id / span_id 实现日志-追踪关联
- **lumberjack** 滚动切片（单文件 1MB，保留 300 个备份，30 天过期，自动压缩）
- GORM SQL 日志自动记录执行时间和影响行数

### 指标监控

**业务指标：**

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `login_attempts_total` | CounterVec | 登录尝试（分 success/fail） |
| `user_creations_total` | Counter | 用户注册数 |
| `friend_additions_total` | Counter | 好友添加数 |
| `friend_requests_sent_total` | Counter | 好友请求发送数 |
| `friend_requests_accepted_total` | Counter | 好友请求通过数 |

**HTTP 指标：**

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `http_requests_total` | CounterVec | 请求总数（method/path/status） |
| `http_request_duration_seconds` | HistogramVec | 请求耗时分布 |
| `http_requests_in_flight` | Gauge | 当前并发请求数 |

**数据库连接池指标：**

| 指标名 | 类型 |
|--------|------|
| `db_connections_max_open` | Gauge |
| `db_connections_open` | Gauge |
| `db_connections_in_use` | Gauge |
| `db_connections_idle` | Gauge |
| `db_connections_wait_count_total` | Counter |
| `db_connections_wait_duration_seconds_total` | Counter |
| `db_connections_max_idle_closed_total` | Counter |
| `db_connections_max_lifetime_closed_total` | Counter |

**Go 运行时指标**通过 `NewGoCollector()` 自动注册（goroutine、内存、GC 等）。

### 链路追踪

- 基于 **OpenTelemetry**，支持 OTLP gRPC 导出（Jaeger/Tempo 等兼容）或 Stdout 导出
- 采样率通过 `OTEL_TRACE_SAMPLE_RATE` 控制（默认 10%）
- DAO 层每个方法生成独立 Span，携带 table、user.id 等属性
- Gin 请求自动生成 Span（通过 otelgin 中间件）

---

## 安全设计

| 层面 | 措施 |
|------|------|
| 传输安全 | 支持 TLS（可选），HSTS 头（TLS 开启时） |
| 认证 | JWT HS256，24 小时过期 |
| 密码存储 | bcrypt 哈希 |
| 参数校验 | go-playground/validator，最小密码长度 8 位 |
| 请求体限制 | 最大 1MB |
| Content-Type | POST/PUT 强制 `application/json` |
| CORS | 可配置允许来源 |
| 安全响应头 | X-Content-Type-Options、X-Frame-Options、X-XSS-Protection、Referrer-Policy、CSP |
| 限流 | 登录接口 10 次/分钟/IP（内存实现） |
| 权限控制 | 用户仅可操作自己的资源 |
| 请求超时 | 全局 30 秒 |
| 错误脱敏 | 内部错误不泄露到客户端响应 |
| 审计日志 | 所有请求记录 user_id、method、path、status |

---

## 测试

```bash
# 运行所有单元测试
cd cmd && make test

# 或直接使用 go test
go test -v ./test/...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out
```

测试覆盖范围：

| 测试文件 | 覆盖内容 |
|----------|----------|
| `Healthz_test.go` | 健康/就绪探针 |
| `GetUser_test.go` | 用户查询 |
| `CreateUser_test.go` | 用户创建 |
| `Login_test.go` | 登录认证 |
| `AuthMiddleware_test.go` | JWT 中间件 |
| `ModifyUser_test.go` | 用户修改 |
| `DeleteUser_test.go` | 用户删除 |
| `AddFriends_test.go` | 添加好友 |
| `GetFriendsList_test.go` | 好友列表 |
| `GetNearbyFriendList_test.go` | 附近好友 |
| `CORS_test.go` | CORS 中间件 |
| `SecurityHeaders_test.go` | 安全响应头 |
| `RateLimiter_test.go` | 限流中间件 |
| `Logger_test.go` | 日志模块 |
| `Config_test.go` | 配置模块 |

---

## 部署

### Docker

```bash
docker build -t user-service:latest .
docker run -d -p 8080:8080 \
  -e MYSQL_HOST="host" \
  -e MYSQL_USER="user" \
  -e MYSQL_PASSWORD="pass" \
  -e MYSQL_DATABASE="db" \
  -e JWT_SECRET="$(openssl rand -base64 32)" \
  user-service:latest
```

### Kubernetes

```bash
# 方式一：Kustomize
kubectl apply -k k8s/

# 方式二：Helm
helm upgrade --install user-service ./helm/user-service \
  --namespace user-service --create-namespace \
  --set secret.jwtSecret="$(openssl rand -base64 32)" \
  --set secret.mysqlHost="mysql" \
  --set secret.mysqlUser="xgame" \
  --set secret.mysqlPassword="password" \
  --set secret.mysqlDatabase="bobby_test"
```

### 交叉编译

```bash
# Linux amd64
cd cmd && make build-linux

# 手动指定平台
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o UserServer ./cmd/
```

---

## 运维手册

详细的部署、监控、告警、故障排查指南参见 [部署运维手册](docs/ops-manual.md)。

---

## License

MIT License

Copyright (c) 2022 guobin
