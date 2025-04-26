echo "打包为docker镜像"
cd ..
PROJECT_ROOT=$(pwd)
# 为每个服务构建Docker镜像
for service in $(ls -d $PROJECT_ROOT/app/*/); do
    service_name=$(basename $service)

    # 检查是否存在runner可执行文件
    if [ -f "$service/runner" ]; then
        echo "构建 $service_name 的Docker镜像..."

        # 构建并推送镜像
        docker build -t xyq777/linkbox-$service_name -f $service/Dockerfile $service
       docker push xyq777/linkbox-$service_name
        echo "  - 镜像构建完成: xyq777/linkbox-$service_name"
    fi
done

echo "所有服务镜像构建完成"
