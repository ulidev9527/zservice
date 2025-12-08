#!/bin/bash

echo "OS: $OSTYPE"

root_dir=$(dirname "$0")
echo "Root: $root_dir"

_compiler_dir=$(realpath "./tools/fb_protoc")
_compiler_runtime=$(realpath "./tools/fb_protoc")

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    _compiler_runtime="${_compiler_dir}/flatc"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    _compiler_runtime="${_compiler_dir}/flatc"
elif [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    _compiler_runtime="${_compiler_dir}/flatc.exe"
else
    echo "Unsupported OS type: $OSTYPE"
    exit 1
fi

echo "Runtime: $_compiler_runtime"

if [ ! -f "$_compiler_runtime" ]; then
    echo "Error: Runtime '$_compiler_runtime' does not exist."
    exit 1
fi

# 判断输入路径是否是文件夹
if [ ! -d "$1" ]; then
    echo "Error: '$1' is not a directory."
    exit 1
fi

# 创建临时路径，临时路径就在$0所在文件夹下创建一个文件夹,文件夹名称使用 __时间戳 格式生成，如果失败则报错退出
script_dir="$(cd $root_dir && pwd)"
timestamp=$(date +%s)
tmp_dir="${script_dir}/__fb_protoc_tmpdir_${timestamp}"

if ! mkdir -p "$tmp_dir"; then
    echo "Error: Failed to create temporary directory '$tmp_dir'."
    exit 1
fi

echo "tempDir: $tmp_dir"

inputPath="${1}/*.fbs"

# 执行 fb_protoc 编译器
runCmd="$_compiler_runtime -o $tmp_dir --gen-object-api --json --go --csharp $inputPath"
echo "CMD: ${runCmd}"
$runCmd

if [ $? -ne 0 ]; then
    echo "Error: Compilation failed."
    exit 1
else
    echo "Compilation succeeded."
fi

# 收集临时文件夹和文件夹子路径中所有文件路径存储到 fileList
fileList=()
while IFS= read -r -d $'\0' file; do
    fileList+=("$file")
done < <(find "$tmp_dir" -type f -print0)

# 执行扩展
echo "Run Go Extend"
go run $_compiler_dir/*.go $tmp_dir $1

# ---------------- 输出
echo "All File Move To: $1"

# 删除临时文件夹
rm -rf "$tmp_dir"
echo "Temporary directory '$tmp_dir' deleted."