#!/bin/bash

# 启动多个 order_srv 实例（负载均衡测试）

echo "启动 Consul..."
docker-compose up -d consul
sleep 3

echo ""
echo "启动 order_srv 实例1..."
SERVICE_ID="order-service-1" \
SERVICE_PORT="50051" \
go run order_srv/cmd/main.go &
ORDER_PID1=$!

echo "启动 order_srv 实例2..."
SERVICE_ID="order-service-2" \
SERVICE_PORT="50052" \
go run order_srv/cmd/main.go &
ORDER_PID2=$!

echo "启动 order_srv 实例3..."
SERVICE_ID="order-service-3" \
SERVICE_PORT="50053" \
go run order_srv/cmd/main.go &
ORDER_PID3=$!

echo ""
echo "等待3秒让服务注册到 Consul..."
sleep 3

echo ""
echo "检查 Consul 中的服务实例:"
curl -s http://localhost:8500/v1/health/service/order-service | jq -r '.[] | "ServiceID: \(.Service.ID), Address: \(.Service.Address):\(.Service.Port), Status: \(.Checks[0].Status)"'

echo ""
echo "启动 h5_bff 服务..."
go run h5_bff/cmd/main.go &
BFF_PID=$!

echo ""
echo "=========================================="
echo "服务启动完成!"
echo "order-service-1: localhost:50051 (PID: $ORDER_PID1)"
echo "order-service-2: localhost:50052 (PID: $ORDER_PID2)"
echo "order-service-3: localhost:50053 (PID: $ORDER_PID3)"
echo "h5_bff: localhost:8080 (PID: $BFF_PID)"
echo "Consul: http://localhost:8500"
echo "=========================================="
echo ""
echo "按 Ctrl+C 停止所有服务..."

# 捕获 Ctrl+C 信号，停止所有服务
trap "echo '停止所有服务...'; kill $ORDER_PID1 $ORDER_PID2 $ORDER_PID3 $BFF_PID; exit" INT

# 保持脚本运行
wait
