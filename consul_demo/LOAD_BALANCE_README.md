# Consul 负载均衡配置说明

## 功能概述

本项目已实现基于 Consul 的服务注册与发现，以及客户端负载均衡功能。

## 核心特性

1. **多实例部署**：支持启动多个 order_srv 实例，使用相同的服务名称，不同的服务ID和端口
2. **自动注册**：每个实例启动时自动注册到 Consul，并维持 TTL 健康检查
3. **客户端负载均衡**：BFF 层自动发现所有健康实例，使用轮询算法分发请求
4. **动态扩缩容**：支持运行时添加或移除服务实例，无需重启 BFF 层

## 快速开始

### 1. 启动 Consul
```bash
docker-compose up -d consul
```

### 2. 启动多个 order_srv 实例

#### 方式一：使用启动脚本（推荐）
```bash
./start_services.sh
```

该脚本会自动启动：
- order-service-1: localhost:50051
- order-service-2: localhost:50052
- order-service-3: localhost:50053
- h5_bff: localhost:8080

#### 方式二：手动启动

**实例1：**
```bash
SERVICE_ID="order-service-1" \
SERVICE_PORT="50051" \
go run order_srv/cmd/main.go
```

**实例2：**
```bash
SERVICE_ID="order-service-2" \
SERVICE_PORT="50052" \
go run order_srv/cmd/main.go
```

**实例3：**
```bash
SERVICE_ID="order-service-3" \
SERVICE_PORT="50053" \
go run order_srv/cmd/main.go
```

### 3. 启动 BFF 服务
```bash
go run h5_bff/cmd/main.go
```

### 4. 验证负载均衡

查看 Consul 中的服务实例：
```bash
curl http://localhost:8500/v1/health/service/order-service
```

发送多个请求测试负载分发：
```bash
./test_load_balance.sh
```

## 配置说明

### order_srv 环境变量

- `SERVICE_ID`：服务唯一标识（默认：order-service-1）
- `SERVICE_PORT`：服务监听端口（默认：50051）
- `SERVICE_NAME`：服务名称（在代码中配置为 order-service）

### 负载均衡算法

当前实现：**轮询算法（Round Robin）**

BFF 层会按顺序循环使用所有健康的 order_srv 实例。

可扩展的负载均衡算法：
- 随机算法（Random）
- 加权轮询（Weighted Round Robin）
- 一致性哈希（Consistent Hash）
- 最少连接数（Least Connections）

## 架构图

```
┌─────────────────┐
│   客户端请求    │
│  (localhost:8080)│
└────────┬────────┘
         │
         ▼
┌─────────────────────────────┐
│         BFF 层              │
│  (负载均衡 + HTTP/gRPC转换)  │
└────────┬────────────────────┘
         │
         │ 轮询分发
         ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ order-service-1 │  │ order-service-2 │  │ order-service-3 │
│ (localhost:50051)│  │ (localhost:50052)│  │ (localhost:50053)│
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │                    │                    │
         └───────────┬────────┴──────────┬─────────┘
                     ▼                   ▼
            ┌─────────────────┐
            │   MySQL 数据库   │
            │   (数据共享)     │
            └─────────────────┘
```

## 注意事项

1. **服务名称必须一致**：所有 order_srv 实例必须使用相同的服务名称（order-service）
2. **服务ID必须唯一**：每个实例的服务ID必须不同（order-service-1, order-service-2, order-service-3）
3. **端口不能冲突**：每个实例使用不同的端口（50051, 50052, 50053）
4. **共享数据库**：所有实例连接同一个 MySQL 数据库
5. **健康检查**：每个实例每 25 秒向 Consul 报告一次健康状态

## 故障处理

### 实例宕机
如果某个 order_srv 实例异常退出：
1. Consul 会在 60 秒后自动注销该实例
2. BFF 层会自动剔除不健康的实例
3. 请求会自动分发到剩余的健康实例

### 手动下线实例
```bash
# 查看服务实例
curl http://localhost:8500/v1/catalog/service/order-service

# 手动注销实例
curl --request PUT http://localhost:8500/v1/agent/service/deregister/order-service-2
```

## 性能优化建议

1. **连接池**：每个 gRPC 连接可配置连接池参数
2. **超时控制**：为每个 RPC 调用设置合理的超时时间
3. **重试机制**：对失败的请求实现指数退避重试
4. **熔断降级**：当错误率过高时自动熔断，返回降级数据
5. **监控告警**：监控服务实例的健康状态和请求分发情况

## 日志输出示例

BFF 层启动时会显示连接的所有实例：
```
Connected to order service: localhost:50051
Connected to order service: localhost:50052
Connected to order service: localhost:50053
Successfully connected to 3 order service instances
```

每个请求会自动轮询分发到不同实例。
