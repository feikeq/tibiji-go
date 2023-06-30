#!/bin/bash

# 创建目录
mkdir -p ./dist

# 拷贝文档文件
cp ./README.md ./dist/README.md

# 拷贝资源文件夹
cp -r ./assets ./dist/assets

# 拷贝配置文件夹
cp -r ./config ./dist/config

# 删除源代码文件
rm ./dist/config/*.go

# 编译项目 并 指定输出文件的路径和名称
# go build -o ./dist/ .
# 移除调试信息 -ldflags 调试信息会增加可执行文件的大小您可以使用以下命令编译可执行文件时移除调试信息：
# go build -ldflags "-w -s -extldflags '-static'" -o ./dist/ . # 使用Go的静态编译。静态编译会将所有依赖库打包到可执行文件中，从而减小文件大小

# 交叉编译到Windows和Linux平台分别编译可执行文件：
GOOS=windows GOARCH=amd64 go build -o ./dist/ .
GOOS=linux GOARCH=amd64 go build -o ./dist/ .