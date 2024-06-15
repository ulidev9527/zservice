# ------------------- 打包
FROM golang:latest AS build-step

# 环境变量
ARG SERVERNAME=SERVERNAME

# 工作目录
WORKDIR /app

# 复制文件 / 文件夹
COPY go.mod go.sum .
COPY internal ./internal
COPY service ./service

# 执行命令
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o runtime ./service/$SERVERNAME/$SERVERNAME.go

# ------------------- 运行准备
FROM alpine:latest AS runtime-step

# 环境变量
ARG SERVERNAME=SERVERNAME

# 工作目录
WORKDIR /app

# 复制文件
COPY --from=build-step /app/$SERVERNAME .

# fix: runtime: not found
# https://www.cnblogs.com/yangzp/p/14609641.html
# 创建lib64的软连接
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENTRYPOINT ["/app/$SERVERNAME"]