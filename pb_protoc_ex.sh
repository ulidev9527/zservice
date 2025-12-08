#!/bin/bash

# 获取当前脚本的绝对路径
script_path=$(realpath "$0")

# 使用 dirname 提取目录部分
script_dir=$(dirname "$script_path")

cd $script_dir

go run $script_dir/tools/protoc_ex/*.go $1