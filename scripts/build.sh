#!/bin/bash

# 设置项目根目录
PROJECT_ROOT=$(pwd)

# 生成依赖注入代码

services=("auth" "user" "organization" "tag" "item" "gateway")

# 编译所有
echo "编译所有服务..."
for service in "${services[@]}"; do
    if [ -d "$PROJECT_ROOT/app/$service/cmd/linkbox_$service" ]; then
        echo "编译 $service 服务..."
        cd "$PROJECT_ROOT/app/$service/cmd/linkbox_$service" && go build -o linkbox_$service
    fi
done

# 启动依赖服务
echo "启动依赖服务..."
cd "$PROJECT_ROOT/deploy" && docker-compose up -d

# 启动所有服务
echo "启动所有服务..."
for service in "${services[@]}"; do
    if [ -f "$PROJECT_ROOT/app/$service/cmd/linkbox_$service/linkbox_$service" ]; then
        echo "启动 $service 服务..."
        cd "$PROJECT_ROOT/app/$service" && ./cmd/linkbox_$service/linkbox_$service &
    fi
done

echo "所有服务已启动"