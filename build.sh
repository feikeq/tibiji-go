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
go build -o ./dist/ .


