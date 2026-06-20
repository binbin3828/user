# 部署运维手册

## 目录

- [部署架构](#部署架构)
- [环境准备](#环境准备)
- [部署方式](#部署方式)
  - [方式一：二进制部署](#方式一二进制部署)
  - [方式二：Docker 部署](#方式二docker-部署)
  - [方式三：Docker Compose 部署](#方式三docker-compose-部署)
  - [方式四：Kubernetes 部署](#方式四kubernetes-部署)
- [配置管理](#配置管理)
- [监控告警](#监控告警)
  - [Prometheus 指标采集](#prometheus-指标采集)
  - [Grafana 仪表盘](#grafana-仪表盘)
  - [告警规则](#告警规则)
- [日志管理](#日志管理)
- [数据库运维](#数据库运维)
- [日常运维操作](#日常运维操作)
- [故障排查](#故障排查)
- [备份恢复](#备份恢复)
- [安全加固](#安全加固)
- [升级策略](#升级策略)
- [容量规划](#容量规划)

---

## 部署架构

```
                         ┌──────────────┐
                         │  Load Balancer│
                         │  (Nginx/ALB) │
                         └──────┬───────┘
                                │
          ┌─────────────────────┼─────────────────────┐
          │                     │                     │
  ┌───────▼───────┐   ┌───────▼───────┐   ┌───────▼───────┐
  │  user-service │   │  user-service │   │  user-service │
  │    Pod-1      │   │    Pod-2      │   │    Pod-3      │
  └───────┬───────┘   └───────┬───────┘   └───────┬───────┘
          │                     │                     │
          └─────────────────────┼─────────────────────┘
                                │
                    ┌───────────▼───────────┐
                    │      MySQL 8.0        │
                    │   (主库 / 主从集群)    │
                    └───────────────────────┘
                                │
                    ┌───────────▼───────────┐
                    │  可观测性平台          │
                    │  Prometheus + Grafana  │
                    │  Jaeger / Tempo        │
                    └───────────────────────┘
```

---

## 环境准备

### 硬件要求

| 环境 | CPU | 内存 | 磁盘 |
|------|-----|------|------|
| 开发/测试 | 1 核 | 512MB | 1GB |
| 生产（单实例） | 2 核 | 512MB | 10GB（含日志） |
| MySQL | 2 核 | 4GB | 100GB+ SSD |

### 软件依赖

| 组件 | 版本要求 | 说明 |
|------|----------|------|
| Go | 1.25+ | 编译需要 |
| MySQL | 8.0+ | 数据存储 |
| Docker | 20.10+ | 容器化部署（可选） |
| Kubernetes | 1.25+ | K8s 部署（可选） |
| Helm | 3.12+ | Helm 部署（可选） |

### 网络要求

| 端口 | 用途 | 方向 |
|------|------|------|
| 8080 | HTTP API | 入站 |
| 443 | HTTPS API（TLS 启用时） | 入站 |
| 3306 | MySQL | 出站（到数据库服务器） |
| 4317 | OTLP gRPC（链路追踪） | 出站 |

---

## 部署方式

### 方式一：二进制部署

适用于传统虚拟机/物理机部署。

```bash
# 1. 创建运行用户
sudo useradd -r -s /bin/false -m user-service

# 2. 准备目录
sudo mkdir -p /opt/user-service/log
sudo chown -R user-service:user-service /opt/user-service

# 3. 编译
cd /path/to/user
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o UserServer ./cmd/

# 4. 部署二进制
sudo cp UserServer /opt/user-service/
sudo cp -r pkg/config/config.yaml /opt/user-service/pkg/config/

# 5. 修改配置（必须修改 JWT Secret 和数据库连接）
sudo vi /opt/user-service/pkg/config/config.yaml

# 6. 配置 systemd 服务
sudo tee /etc/systemd/system/user-service.service > /dev/null <<'EOF'
[Unit]
Description=User Service
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=user-service
Group=user-service
WorkingDirectory=/opt/user-service
Environment="DB_PASSWORD=your_db_password"
Environment="JWT_SECRET=your_long_random_secret"
Environment="LOG_LEVEL=info"
ExecStart=/opt/user-service/UserServer
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

# 7. 启动
sudo systemctl daemon-reload
sudo systemctl enable user-service
sudo systemctl start user-service

# 8. 验证
curl http://localhost:8080/healthz
```

### 方式二：Docker 部署

```bash
# 1. 构建镜像
docker build -t user-service:1.0.0 .

# 2. 运行容器
docker run -d \
  --name user-service \
  --restart unless-stopped \
  -p 8080:8080 \
  -e DB_PASSWORD="your_db_password" \
  -e JWT_SECRET="$(openssl rand -base64 32)" \
  -e LOG_LEVEL="info" \
  -v /opt/user-service/log:/app/log \
  user-service:1.0.0

# 3. 验证
curl http://localhost:8080/healthz

# 4. 查看日志
docker logs -f user-service
```

### 方式三：Docker Compose 部署

适用于开发、测试环境（含 MySQL）。

```bash
# 准备环境变量文件
cat > .env <<EOF
MYSQL_ROOT_PASSWORD=root123
MYSQL_PASSWORD=xgame123
JWT_SECRET=$(openssl rand -base64 32)
LOG_LEVEL=debug
EOF

# 启动
docker-compose up -d

# 验证
curl http://localhost:8080/healthz
```

**注意**：生产环境不建议在 docker-compose 中运行 MySQL，应使用云数据库或独立部署的数据库。

### 方式四：Kubernetes 部署

**前置条件：**

- 已有可用的 MySQL 实例
- 集群中已安装 Metrics Server（HPA 需要）

#### Kustomize 部署

```bash
# 1. 修改密钥配置
vi k8s/secret.yaml
# 将 MYSQL_DSN 和 JWT_SECRET 替换为实际值

# 2. 一键部署
kubectl apply -k k8s/

# 3. 查看状态
kubectl -n user-service get pods,svc,hpa

# 4. 端口转发测试
kubectl -n user-service port-forward svc/user-service 8080:80
curl http://localhost:8080/healthz
```

#### Helm 部署

```bash
# 1. 安装（生产环境必须指定 JWT Secret）
helm upgrade --install user-service ./helm/user-service \
  --namespace user-service --create-namespace \
  --set secret.jwtSecret="$(openssl rand -base64 32)" \
  --set secret.mysqlDsn="xgame:password@tcp(mysql-host:3306)/bobby_test?charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true&timeout=5s" \
  --set config.logLevel="info" \
  --set autoscaling.enabled=true \
  --set autoscaling.minReplicas=3 \
  --set autoscaling.maxReplicas=10

# 2. 查看部署状态
helm -n user-service status user-service

# 3. 升级配置
helm upgrade user-service ./helm/user-service \
  --namespace user-service \
  --reuse-values \
  --set config.logLevel="debug"

# 4. 回滚
helm -n user-service rollback user-service

# 5. 卸载
helm -n user-service uninstall user-service
```

---

## 配置管理

### 配置优先级

环境变量 > 配置文件（`pkg/config/config.yaml` 嵌入到二进制）

### 必须修改的配置项

| 配置项 | 说明 | 生成命令 |
|--------|------|----------|
| `jwt.secret` | JWT 签名密钥，默认值启动会 panic | `openssl rand -base64 32` |
| `mysql.dataSourceName` | 数据库连接串 | 联系 DBA 获取 |
| `tls.certFile / tls.keyFile` | 如需 HTTPS，配置证书路径 | - |

### 数据库连接池调优

```yaml
mysql:
  maxIdle: 25      # 空闲连接数，建议设为 maxOpen 的 50%
  maxOpen: 50      # 最大连接数，根据 MySQL max_connections 设定
  maxLifetime: 1   # 连接最大存活时间（小时），应小于 MySQL wait_timeout
```

**调优建议：**

| 场景 | maxIdle | maxOpen |
|------|---------|---------|
| 开发环境 | 5 | 10 |
| 生产低负载 | 10 | 25 |
| 生产高负载 | 25 | 50-100 |

---

## 监控告警

### Prometheus 指标采集

**指标暴露地址：** `http://<host>:8080/metrics`

**Prometheus 采集配置：**

```yaml
scrape_configs:
  - job_name: user-service
    metrics_path: /metrics
    static_configs:
      - targets:
          - user-service-1:8080
          - user-service-2:8080
          - user-service-3:8080
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance
```

**Kubernetes 自动发现：**

```yaml
scrape_configs:
  - job_name: user-service
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - user-service
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
```

### Grafana 仪表盘

预置仪表盘文件：`deploy/grafana/user-service-dashboard.json`

**导入方式：**

```bash
# 方式一：Grafana UI 导入
# Dashboard → Import → Upload JSON file

# 方式二：API 导入
curl -X POST http://grafana:3000/api/dashboards/db \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <api_key>" \
  -d @deploy/grafana/user-service-dashboard.json
```

**仪表盘包含面板：**

1. HTTP 请求速率（QPS）
2. P99 响应延迟
3. 错误率（5xx）
4. 业务指标（登录次数、注册数）
5. 数据库连接池状态
6. Goroutine 数量
7. 内存使用
8. GC 次数与暂停时间
9. 当前并发请求数
10. 数据库连接等待时间

### 告警规则

告警规则文件：`deploy/prometheus/rules.yml`

| 告警名称 | 级别 | 触发条件 | 持续时间 |
|----------|------|----------|----------|
| **HighErrorRate** | 🔴 critical | 5xx 错误率 > 5% | 3 分钟 |
| **InstanceDown** | 🔴 critical | 实例不可达 | 1 分钟 |
| **DBNotReady** | 🔴 critical | 数据库连接数 = 0 | 1 分钟 |
| **HighLatency** | 🟡 warning | P99 延迟 > 2s | 5 分钟 |
| **DBConnectionPoolExhaustion** | 🟡 warning | 连接池使用率 > 80% | 2 分钟 |
| **HighGoroutineCount** | 🟡 warning | Goroutine > 10000 | 5 分钟 |
| **HighLoginFailureRate** | 🟡 warning | 登录失败率 > 50% | 5 分钟 |

**告警处理流程：**

1. **InstanceDown** → 立即检查 Pod/进程状态，查看是否 OOM 或 CrashLoop
2. **DBNotReady** → 检查数据库可达性、网络策略、MySQL 进程状态
3. **HighErrorRate** → 查看日志定位错误原因
4. **HighLatency** → 检查数据库慢查询、下游依赖、GC 暂停
5. **HighLoginFailureRate** → 可能存在暴力破解攻击，检查审计日志

---

## 日志管理

### 日志位置

| 部署方式 | 日志路径 |
|----------|----------|
| 二进制部署 | `/opt/user-service/log/user.log` |
| Docker | 容器内 `/app/log/user.log`（建议挂载宿主机目录） |
| Kubernetes | 标准输出（`kubectl logs`） + 挂载卷 |

### 日志滚动策略

- **单文件最大**：1MB
- **最大备份数**：300
- **保留天数**：30 天
- **自动压缩**：是（gzip）

### 日志级别

通过 `LOG_LEVEL` 环境变量控制：

```bash
LOG_LEVEL=debug   # 开发环境，输出 SQL 语句和请求体
LOG_LEVEL=info    # 生产环境（默认）
LOG_LEVEL=warn    # 仅警告和错误
LOG_LEVEL=error   # 仅错误
```

**生产环境建议使用 `info` 级别**，`debug` 会输出请求体和 SQL 语句，存在安全风险和性能影响。

### 日志格式

```json
{
  "level": "INFO",
  "ts": "2024-01-15T10:30:00.000+0800",
  "caller": "service/ResponseHandler.go:39",
  "msg": "[trace_id=abc123] [span_id=def456] [audit] request_id=uuid method=GET path=/v1/user/1 status=200"
}
```

### 日志分析

```bash
# 查看最近错误
grep '"level":"ERROR"' log/user.log | tail -50

# 统计 5xx 错误次数
grep 'status=5' log/user.log | wc -l

# 按 API 路径统计请求量
grep 'audit' log/user.log | awk -F'path=' '{print $2}' | awk '{print $1}' | sort | uniq -c | sort -rn

# 提取慢请求（超时/取消）
grep 'timeout\|canceled' log/user.log
```

---

## 数据库运维

### 初始化建表

```bash
mysql -h <host> -u <user> -p <database> < sql/init.sql
```

Kubernetes 可以通过 InitContainer 或 Job 执行：

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: user-db-init
spec:
  template:
    spec:
      containers:
        - name: init-db
          image: mysql:8.0
          command:
            - mysql
            - -h
            - mysql-service
            - -u
            - xgame
            - -p$(MYSQL_PASSWORD)
            - bobby_test
            - -e
            - "$(cat sql/init.sql)"
      restartPolicy: Never
```

### 数据库备份

```bash
# 全量备份
mysqldump -h <host> -u <user> -p \
  --single-transaction \
  --routines \
  --triggers \
  bobby_test > backup_$(date +%Y%m%d_%H%M%S).sql

# 仅备份用户和好友表
mysqldump -h <host> -u <user> -p \
  --single-transaction \
  bobby_test user friends > backup_tables_$(date +%Y%m%d).sql
```

### 数据库恢复

```bash
# 恢复全量备份
mysql -h <host> -u <user> -p bobby_test < backup_20240115.sql

# Kubernetes 环境
kubectl -n user-service exec -i mysql-0 -- mysql -u xgame -p bobby_test < backup.sql
```

### 慢查询排查

```sql
-- 查看当前慢查询
SHOW FULL PROCESSLIST;

-- 查看是否有未使用索引的查询
SELECT * FROM information_schema.statistics
WHERE table_schema = 'bobby_test' AND seq_in_index = 1;

-- 检查表大小
SELECT
  table_name,
  ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.tables
WHERE table_schema = 'bobby_test';
```

### 日常巡检 SQL

```sql
-- 用户总数
SELECT COUNT(*) FROM user;

-- 好友关系总数
SELECT COUNT(*) FROM friends;

-- 今日新增用户
SELECT COUNT(*) FROM user WHERE DATE(create_at) = CURDATE();

-- 今日新增好友关系
SELECT COUNT(*) FROM friends WHERE DATE(create_time) = CURDATE();

-- 没有位置信息的用户数
SELECT COUNT(*) FROM user WHERE loc_geohash = '';
```

---

## 日常运维操作

### 优雅重启

```bash
# systemd
sudo systemctl restart user-service

# Kubernetes（滚动更新）
kubectl -n user-service rollout restart deployment/user-service

# Docker
docker restart user-service
```

### 扩缩容

```bash
# Kubernetes 手动扩缩
kubectl -n user-service scale deployment/user-service --replicas=5

# Helm
helm upgrade user-service ./helm/user-service \
  --reuse-values \
  --set autoscaling.enabled=false \
  --set replicaCount=5
```

HPA 开启时，手动设置 replicas 会被 HPA 覆盖，需要关闭 HPA 或调整 minReplicas/maxReplicas。

### 查看运行状态

```bash
# 健康检查
curl http://<host>:8080/healthz

# 就绪检查
curl http://<host>:8080/readyz

# 版本信息（未来版本将支持）
# curl http://<host>:8080/version

# 查看实时指标
curl http://<host>:8080/metrics | grep -E "^(http_requests_total|login_attempts_total)"
```

### 临时开启 Debug 日志

```bash
# Kubernetes（通过环境变量重启）
kubectl -n user-service set env deployment/user-service LOG_LEVEL=debug

# Docker（重建容器）
docker stop user-service
docker run -d ... -e LOG_LEVEL=debug ... user-service:1.0.0
```

**注意**：debug 模式会输出请求体和 SQL，问题排查完毕后必须切回 info 级别。

---

## 故障排查

### 问题 1：服务启动失败

**症状：** 进程启动后立即退出

**排查步骤：**

```bash
# 1. 查看启动日志
journalctl -u user-service -f         # systemd
docker logs user-service              # Docker
kubectl -n user-service logs deployment/user-service  # K8s

# 2. 常见原因
# - JWT Secret 未修改（默认值导致 panic）
# - 数据库连接失败（网络不通/密码错误/数据库不存在）
# - TLS 证书不存在（但启用了 TLS）
# - 端口被占用
```

**错误信息与解决方案：**

| 错误日志 | 原因 | 解决 |
|----------|------|------|
| `FATAL: jwt.secret must be changed` | JWT Secret 为默认值 | 设置 `JWT_SECRET` 环境变量 |
| `mysql connection failed after 3 attempts` | 数据库不可达 | 检查网络、DSN、防火墙 |
| `bind: address already in use` | 端口被占用 | 修改端口或终止占用进程 |
| `TLS ... no such file` | 证书文件不存在 | 禁用 TLS 或提供证书 |
| `auto migrate failed` | 数据库权限不足 | 授予 CREATE TABLE 权限 |

### 问题 2：API 返回 5xx

**排查步骤：**

```bash
# 1. 查看最近错误日志
grep '"level":"ERROR"' log/user.log | tail -20

# 2. 检查数据库连接
curl http://<host>:8080/readyz

# 3. 查看数据库连接池指标
curl -s http://<host>:8080/metrics | grep db_connections

# 4. 检查数据库是否运行
mysql -h <host> -u <user> -p -e "SELECT 1"
```

**常见原因：**

- 数据库连接池耗尽 → 调大 `maxOpen` 或排查慢查询
- 数据库连接超时 → 检查网络延迟，调大 `readTimeout`/`writeTimeout`
- MySQL 服务不可用 → 检查 MySQL 进程、磁盘空间

### 问题 3：P99 延迟偏高

**排查步骤：**

```bash
# 1. 查看数据库慢查询
mysql> SHOW FULL PROCESSLIST;

# 2. 检查索引使用情况
mysql> EXPLAIN SELECT * FROM user WHERE name = 'xxx';

# 3. 查看 Goroutine 是否泄漏
curl -s http://<host>:8080/metrics | grep go_goroutines

# 4. 检查 GC 暂停时间
curl -s http://<host>:8080/metrics | grep go_gc_duration_seconds
```

**优化方向：**

- 添加缺失的数据库索引
- 检查 `friends` 表大数据量下的 JOIN 查询效率
- 调整 Geohash 查询精度（降低 precision 参数）
- 考虑引入 Redis 缓存热点数据

### 问题 4：登录接口被限流

**症状：** 返回 `{"code": -1, "msg": "too many requests"}`

**排查：**

- 当前限流策略：同一 IP 每分钟最多 10 次登录请求
- 限流实现方式：内存级（进程重启后重置）

**临时解除：**

```bash
# 重启服务清除限流计数
kubectl -n user-service rollout restart deployment/user-service
```

### 问题 5：内存持续增长

```bash
# 查看内存指标
curl -s http://<host>:8080/metrics | grep -E "(go_memstats_alloc_bytes|go_memstats_heap)"

# 触发 GC 观察回收效果
curl -s http://<host>:8080/metrics | grep go_gc

# 使用 pprof 分析（需要添加 pprof 端点）
# go tool pprof http://<host>:8080/debug/pprof/heap
```

---

## 备份恢复

### 备份策略

| 备份类型 | 频率 | 保留 |
|----------|------|------|
| 数据库全量备份 | 每日 | 30 天 |
| 数据库增量备份 | 每小时 | 7 天 |
| 配置文件备份 | 变更时 | 永久 |

### 备份脚本示例

```bash
#!/bin/bash
# backup.sh - 数据库备份脚本

BACKUP_DIR="/backup/user-service"
DB_HOST="mysql-host"
DB_USER="xgame"
DB_PASS="${DB_PASSWORD}"
DB_NAME="bobby_test"
RETENTION_DAYS=30

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/${DB_NAME}_${DATE}.sql.gz"

mkdir -p "${BACKUP_DIR}"

# 全量备份
mysqldump -h "${DB_HOST}" -u "${DB_USER}" -p"${DB_PASS}" \
  --single-transaction \
  --routines \
  --triggers \
  "${DB_NAME}" | gzip > "${BACKUP_FILE}"

# 清理过期备份
find "${BACKUP_DIR}" -name "*.sql.gz" -mtime +${RETENTION_DAYS} -delete

echo "Backup completed: ${BACKUP_FILE}"
```

### 灾难恢复检查清单

- [ ] 确认备份文件完整性（`gunzip -t backup.sql.gz`）
- [ ] 准备一个空数据库实例
- [ ] 恢复备份
- [ ] 验证用户数据完整性
- [ ] 验证好友关系数据完整性
- [ ] 修改 DNS / LoadBalancer 指向恢复后的实例

---

## 安全加固

### 生产环境必须执行

```bash
# 1. JWT Secret 必须修改（否则启动 panic）
export JWT_SECRET="$(openssl rand -base64 32)"

# 2. 数据库密码使用强密码
export DB_PASSWORD="$(openssl rand -base64 24)"

# 3. 启用 TLS
# 修改 config.yaml 中 tls.enabled = true
# 提供有效的 certFile 和 keyFile

# 4. 数据库连接使用 TLS
# DSN 中添加 &tls=true

# 5. 限制 CORS 来源
# 将 allowedOrigins 从 "*" 改为具体域名
```

### 安全加固检查清单

| 检查项 | 说明 |
|--------|------|
| ✅ JWT Secret | 使用 `openssl rand -base64 32` 生成的强随机字符串 |
| ✅ 数据库密码 | 使用强密码，定期轮换 |
| ✅ TLS 加密 | 生产环境必须启用 |
| ✅ CORS 白名单 | 限制为具体域名，不要使用 `*` |
| ✅ 非 root 运行 | Docker/K8s 使用非 root 用户（uid 1001） |
| ✅ 只读文件系统 | K8s 部署中 `readOnlyRootFilesystem: true` |
| ✅ 日志脱敏 | 生产环境使用 `LOG_LEVEL=info`（避免打印请求体） |
| ✅ 网络隔离 | K8s NetworkPolicy 限制出站流量 |

### 安全事件响应

**可疑登录行为：**

```bash
# 检查登录失败率
curl -s http://<host>:8080/metrics | grep login_attempts_total

# 如果失败率异常高（>50% 持续 5 分钟）→ 可能存在暴力破解
# 应对措施：
# 1. 检查审计日志，识别攻击来源 IP
# 2. 在 LoadBalancer/防火墙封禁攻击 IP
# 3. 联系安全团队
```

---

## 升级策略

### 灰度发布（Kubernetes）

```bash
# 1. 构建新版本镜像
docker build -t user-service:1.1.0 .
docker push user-service:1.1.0

# 2. 更新 Deployment 镜像
kubectl -n user-service set image deployment/user-service \
  user-service=user-service:1.1.0

# 3. 监控滚动更新
kubectl -n user-service rollout status deployment/user-service

# 4. 如果异常，立即回滚
kubectl -n user-service rollout undo deployment/user-service
```

### 数据库迁移

数据库结构变更时：

```sql
-- 迁移脚本命名规范：YYYYMMDD_HHMMSS_description.sql
-- 示例：20240115_100000_add_user_status.sql

ALTER TABLE user ADD COLUMN status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=正常 0=禁用';
```

```bash
# 执行迁移前先备份
mysqldump -h <host> -u <user> -p bobby_test > pre_migration_backup.sql

# 执行迁移
mysql -h <host> -u <user> -p bobby_test < migrations/20240115_100000_add_user_status.sql

# 验证
mysql -h <host> -u <user> -p -e "DESC user" bobby_test
```

### 回滚方案

| 回滚内容 | 操作 |
|----------|------|
| 应用版本 | `kubectl rollout undo` / `helm rollback` |
| 数据库结构 | 执行逆向迁移脚本 |
| 配置变更 | `helm rollback` / 重新 apply 旧 ConfigMap |
| 数据误删 | 从备份恢复 |

---

## 容量规划

### 数据增长估算

| 指标 | 估算 |
|------|------|
| 每用户数据量 | ~1KB（含好友关系） |
| 每好友关系 | ~20 字节 |
| 日志（info 级别） | ~10MB/天（千级 QPS） |
| 日志（debug 级别） | ~100MB/天 |

### 扩容阈值

| 指标 | 触发条件 | 措施 |
|------|----------|------|
| CPU | 持续 > 70% | HPA 自动扩容 / 手动增加副本 |
| 内存 | 持续 > 80% | 增加 memory limit |
| 数据库连接 | 使用率 > 80% | 调大 maxOpen 或增加只读副本 |
| 磁盘（日志） | 使用率 > 80% | 清理旧日志或调整保留策略 |
| API P99 延迟 | > 1s 持续 5 分钟 | 增加副本 / 优化查询 |

### 性能基准

| 场景 | 预期 QPS（单实例） |
|------|---------------------|
| 健康检查 | 10000+ |
| 用户查询（主键） | 3000+ |
| 用户创建 | 500+ |
| 好友列表（分页 20 条） | 2000+ |
| 附近好友（Geohash 前缀） | 1000+ |
| 登录 | 200+（受限流影响） |

---

## 附录

### A. 常用命令速查

```bash
# 服务状态
systemctl status user-service
kubectl -n user-service get pods

# 查看日志
journalctl -u user-service -f --since "10 min ago"
kubectl -n user-service logs -f deployment/user-service

# 健康检查
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

### B. 环境变量完整列表

| 变量 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `DB_PASSWORD` | 是* | - | 替换 DSN 中的 `${DB_PASSWORD}` |
| `MYSQL_DSN` | 是* | - | 完整覆盖 DSN（优先级更高） |
| `JWT_SECRET` | 是 | - | JWT 签名密钥 |
| `LOG_LEVEL` | 否 | info | debug/info/warn/error |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | 否 | - | 链路追踪导出地址（空=stdout） |
| `OTEL_TRACE_SAMPLE_RATE` | 否 | 0.1 | 采样率 |

> *DB_PASSWORD 和 MYSQL_DSN 至少设置一个

### C. 紧急联系

如遇线上紧急问题，按照以下顺序联系：

1. 查看 [告警规则](#告警规则) 确认告警级别
2. 按照 [故障排查](#故障排查) 章节定位问题
3. 如需升级，联系开发团队负责人
