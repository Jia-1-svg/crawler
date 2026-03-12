#!/bin/bash

echo "测试负载均衡..."
echo "===================="
echo ""

# 发送10个请求，观察是否分发到不同的服务实例
for i in {1..10}; do
  echo "请求 $i:"
  
  # 创建订单
  response=$(curl -s -X POST http://localhost:8080/orders \
    -H "Content-Type: application/json" \
    -d "{\"user_id\": 100, \"total_price\": 99.99}" \
    -w "\nHTTP Status: %{http_code}\n" \
    2>&1)
  
  echo "$response"
  echo "--------------------"
  
  # 间隔1秒
  sleep 1
done

echo ""
echo "测试完成！"
