#!/bin/bash

# 设置项目根目录
PROJECT_ROOT=$(dirname "$(pwd)")
IDL_DIR=$PROJECT_ROOT/idl


# 生成各个服务的 RPC 代码
echo "生成用户服务..."
kitex -module github.com/linkbox-group/linkbox-server/rpc-gen \
      -I $IDL_DIR \
-gen-path . \
      $IDL_DIR/user/service.proto

 echo "生成内容服务..."
 kitex -module github.com/linkbox-group/linkbox-server/rpc-gen \
       -I $IDL_DIR \
-gen-path . \
       $IDL_DIR/content/service.proto

 echo "生成组织服务..."
 kitex -module github.com/linkbox-group/linkbox-server/rpc-gen \
       -I $IDL_DIR \
-gen-path . \
       $IDL_DIR/organization/service.proto

 echo "生成标签服务..."
 kitex -module github.com/linkbox-group/linkbox-server/rpc-gen \
       -I $IDL_DIR \
-gen-path . \
       $IDL_DIR/tag/service.proto

 echo "生成搜索服务..."
 kitex -module github.com/linkbox-group/linkbox-server/rpc-gen \
       -I $IDL_DIR \
-gen-path . \
       $IDL_DIR/search/service.proto

 echo "生成分享服务..."
 kitex -module github.com/linkbox-group/linkbox-server/rpc-gen \
       -I $IDL_DIR \
-gen-path . \
       $IDL_DIR/sharing/service.proto

 echo "生成分析服务..."
 kitex -module github.com/linkbox-group/linkbox-server/rpc-gen \
       -I $IDL_DIR \
-gen-path . \
       $IDL_DIR/analytics/service.proto

 echo "生成AI服务..."
 kitex -module github.com/linkbox-group/linkbox-server/rpc-gen \
       -I $IDL_DIR \
-gen-path . \
       $IDL_DIR/ai/service.proto

 echo "所有代码生成完成！"
