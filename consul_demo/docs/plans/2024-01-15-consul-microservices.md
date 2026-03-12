# Consul 微服务系统实施计划

> **实施方式:** 使用 dispatching-parallel-agents 并行执行独立任务

**目标:** 在 `/Users/hyx/work/gowork/src/consul_demo` 目录下完成基于 Consul + gRPC 的微服务系统

**架构:** 采用微服务架构，Consul 作为服务注册发现中心，gRPC 作为服务间通信协议，MySQL 存储订单数据

**技术栈:** Go 1.21+, gRPC, Consul, MySQL 8, Protocol Buffers

---

## 任务列表

### 任务 1: Docker Compose 部署 Consul

**目标:** 创建 docker-compose.yml 完成 Consul 单节点部署

**文件:**
- 创建: `docker-compose.yml`

**依赖:** 无

**验收标准:**
- Consul UI 可通过 http://localhost:8500 访问
- 服务健康检查正常

---

### 任务 2: Consul SDK 开发

**目标:** 在 `pkg/consul` 目录下创建 Consul 客户端 SDK，包含服务注册、发现、健康检查功能

**文件:**
- 创建: `pkg/consul/client.go` - Consul 客户端封装
- 创建: `pkg/consul/registry.go` - 服务注册功能
- 创建: `pkg/consul/discovery.go` - 服务发现功能
- 创建: `pkg/consul/go.mod` - 模块定义

**依赖:** 无

**验收标准:**
- 可以成功连接到 Consul 服务器
- 支持服务注册（含健康检查）
- 支持服务发现（按服务名查询）
- 代码通过 go mod tidy

---

### 任务 3: Protocol Buffers 定义

**目标:** 创建 order.proto 定义订单服务接口，并生成 Go 代码

**文件:**
- 创建: `proto/order.proto` - 订单服务 protobuf 定义
- 创建: `proto/generate.sh` - 代码生成脚本
- 生成: `proto/order/order.pb.go` - Go protobuf 代码
- 生成: `proto/order/order_grpc.pb.go` - Go gRPC 代码
- 创建: `proto/go.mod` - 模块定义

**依赖:** 无（可并行于任务 2）

**验收标准:**
- proto 文件包含 Order 消息定义和 CRUD 接口
- 成功生成 .pb.go 文件
- 代码可编译通过

---

### 任务 4: Order 服务开发

**目标:** 在 `order_srv` 目录下实现订单微服务，包含 CRUD 功能和 Consul 注册

**文件:**
- 创建: `order_srv/main.go` - 服务入口
- 创建: `order_srv/service/order.go` - 订单业务逻辑
- 创建: `order_srv/model/order.go` - 订单数据模型
- 创建: `order_srv/config/config.go` - 配置管理
- 创建: `order_srv/go.mod` - 模块定义

**依赖:**
- 任务 2: Consul SDK
- 任务 3: protobuf 生成的代码

**数据库信息:**
- Host: localhost (通过 Docker)
- Port: 3306
- Database: crmeb
- User: root
- Password: 123456
- 表名: eb_store_order

**验收标准:**
- 服务启动后自动注册到 Consul
- 实现订单的增删改查 gRPC 接口
- 支持健康检查
- 可从数据库读取/写入订单数据

---

### 任务 5: H5 BFF 网关开发

**目标:** 在 `h5_bff` 目录下实现 BFF 层，通过 Consul 发现并调用订单服务

**文件:**
- 创建: `h5_bff/main.go` - 服务入口
- 创建: `h5_bff/handler/order.go` - 订单 HTTP 接口
- 创建: `h5_bff/client/order.go` - 订单服务 gRPC 客户端
- 创建: `h5_bff/config/config.go` - 配置管理
- 创建: `h5_bff/go.mod` - 模块定义

**依赖:**
- 任务 2: Consul SDK
- 任务 3: protobuf 生成的代码
- 任务 4: Order 服务（运行时依赖）

**验收标准:**
- BFF 服务可启动
- 通过 Consul 服务发现获取 order_srv 地址
- 提供 HTTP RESTful API 调用 gRPC 订单服务
- HTTP API 包含完整的订单 CRUD 端点

---

## 实施策略

**并行阶段:**
- 阶段 A（可并行）:
  - 任务 1: Consul Docker 部署
  - 任务 2: Consul SDK 开发
  - 任务 3: Protocol Buffers 定义

- 阶段 B（依赖阶段 A）:
  - 任务 4: Order 服务开发

- 阶段 C（依赖阶段 A 和 B）:
  - 任务 5: H5 BFF 网关开发

**验证流程:**
1. 启动 Consul: `docker-compose up -d`
2. 启动 MySQL: （用户已部署）
3. 启动 Order 服务: `cd order_srv && go run main.go`
4. 启动 H5 BFF: `cd h5_bff && go run main.go`
5. 测试 HTTP API: `curl http://localhost:8080/orders`

---

## 技术规范

### Consul SDK API 设计

```go
// Client 封装 Consul 客户端
type Client struct {
    client *api.Client
}

// RegisterService 注册服务
func (c *Client) RegisterService(serviceID, serviceName, address string, port int) error

// DeregisterService 注销服务
func (c *Client) DeregisterService(serviceID string) error

// DiscoverService 发现服务
func (c *cClient) DiscoverService(serviceName string) ([]*api.ServiceEntry, error)
```

### Order Proto 定义

```protobuf
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrder(GetOrderRequest) returns (Order);
  rpc UpdateOrder(UpdateOrderRequest) returns (UpdateOrderResponse);
  rpc DeleteOrder(DeleteOrderRequest) returns (DeleteOrderResponse);
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
}

message Order {
  int64 id = 1;
  string order_no = 2;
  int64 user_id = 3;
  double total_price = 4;
  int32 status = 5;
  string create_time = 6;
  string update_time = 7;
}
```

### Order 数据库表结构

表名: `eb_store_order`
字段:
- id: bigint, 主键
- order_no: varchar, 订单编号
- user_id: bigint, 用户ID
- total_price: decimal, 订单总价
- status: tinyint, 订单状态
- create_time: datetime, 创建时间
- update_time: datetime, 更新时间

### H5 BFF HTTP API

- GET    /orders           - 获取订单列表
- GET    /orders/:id       - 获取订单详情
- POST   /orders           - 创建订单
- PUT    /orders/:id       - 更新订单
- DELETE /orders/:id       - 删除订单
