# Soundwave-Go 服务发现中心

## 功能特性
- 服务注册与发现
- 服务健康检查
- 负载均衡
- 服务状态监控
- 配置管理

## 快速开始
...

## API文档
...

## 配置说明
...

## 开发指南
...

## 性能指标
...

一个轻量级的服务注册与发现中心，提供以下功能：
- 服务注册
- 服务发现
- 服务状态查询

## 项目结构 


# 注册用户服务
curl -X POST http://localhost:7777/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "user-service",
    "id": "user-service-1",
    "hostname": "host-1",
    "ip": "192.168.1.100",
    "port": 8081,
    "metadata": {
      "version": "1.0.0",
      "environment": "development"
    }
  }'

# 注册订单服务
curl -X POST http://localhost:7777/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "order-service",
    "id": "order-service-1",
    "hostname": "host-2",
    "ip": "192.168.1.101",
    "port": 8082,
    "metadata": {
      "version": "1.0.0",
      "environment": "development"
    }
  }'

# 注册支付服务
curl -X POST http://localhost:7777/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "payment-service",
    "id": "payment-service-1",
    "hostname": "host-3",
    "ip": "192.168.1.102",
    "port": 8083,
    "metadata": {
      "version": "1.0.0",
      "environment": "development"
    }
  }'

# 注册同一服务的不同实例
curl -X POST http://localhost:7777/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "user-service",
    "id": "user-service-2",
    "hostname": "host-4",
    "ip": "192.168.1.103",
    "port": 8084,
    "metadata": {
      "version": "1.0.0",
      "environment": "development"
    }
  }'

# 查询用户服务
curl http://localhost:7777/services/user-service | jq '.'

# 查询订单服务
curl http://localhost:7777/services/order-service | jq '.'

# 查询支付服务
curl http://localhost:7777/services/payment-service | jq '.'

# 查询所有服务
curl http://localhost:7777/services | jq '.'

# 用户服务心跳
curl -X PUT http://localhost:7777/services/user-service/user-service-1/heartbeat | jq '.'

# 用户服务第二个实例心跳
curl -X PUT http://localhost:7777/services/user-service/user-service-2/heartbeat | jq '.'

# 订单服务心跳
curl -X PUT http://localhost:7777/services/order-service/order-service-1/heartbeat | jq '.'

# 支付服务心跳
curl -X PUT http://localhost:7777/services/payment-service/payment-service-1/heartbeat | jq '.'

# 获取用户服务统计信息
curl http://localhost:7777/services/user-service/stats | jq '.'

# 获取订单服务统计信息
curl http://localhost:7777/services/order-service/stats | jq '.'

# 获取支付服务统计信息
curl http://localhost:7777/services/payment-service/stats | jq '.'

# 注意：如果没有安装jq，可以使用以下命令安装：
# Ubuntu/Debian:
# sudo apt-get install jq
#
# CentOS/RHEL:
# sudo yum install jq
#
# macOS:
# brew install jq
