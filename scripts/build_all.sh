#!/bin/bash

# 设置项目根目录
cd ..
PROJECT_ROOT=$(pwd)


# 仅支持从Darwin arm64到Linux amd64的交叉编译
PLATFORMS=("linux/amd64")


# 遍历app目录下的所有服务
for service in $(ls -d $PROJECT_ROOT/app/*/); do
    service_name=$(basename $service)
    
    # 检查是否存在main.go文件
    if [ -f "$service/cmd/linkbox_$service_name/main.go" ]; then
        echo "编译 $service_name 服务..."
        
        # 为每个平台交叉编译
        for platform in "${PLATFORMS[@]}"; do
            os=${platform%/*}
            arch=${platform#*/}
            
            # 创建输出目录
            
            # 进入服务目录并设置环境变量编译
            cd "$service/cmd/linkbox_$service_name"
            env GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o "$service/runner" .
            cd - > /dev/null
            
            echo "  - 编译完成: $os-$arch"
        done
    fi
done

echo "所有服务交叉编译完成"

echo "打包为docker镜像"

# 为每个服务构建Docker镜像
for service in $(ls -d $PROJECT_ROOT/app/*/); do
    service_name=$(basename $service)
    
    # 检查是否存在runner可执行文件
    if [ -f "$service/runner" ]; then
        echo "构建 $service_name 的Docker镜像..."
        
        # 构建并推送镜像
        docker build -t xyq777/linkbox-$service_name -f $service/Dockerfile $service
        docker push xyq777/linkbox-$service_name
        
        echo "  - 镜像构建完成并已推送: xyq777/linkbox-$service_name"
    fi
done

echo "所有服务镜像构建完成"

