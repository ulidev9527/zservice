# ------------------- 打包
FROM golang:1.22.3 AS build-step

# 工作目录
WORKDIR /app

# 依赖准备
# 复制文件 / 文件夹
COPY go.mod go.sum .

# 执行命令
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod download

# 打包
COPY zservice ./zservice
COPY service ./service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o __SERVICE_NAME__ ./service/__SERVICE_NAME__/__SERVICE_NAME__.service.go

# ------------------- 运行准备
FROM alpine:latest AS runtime-step

# 工作目录
WORKDIR /app

# 复制文件
COPY --from=build-step /app/__SERVICE_NAME__ .

# fix: runtime: not found
# https://www.cnblogs.com/yangzp/p/14609641.html
# 创建lib64的软连接
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENTRYPOINT ["/app/__SERVICE_NAME__"]