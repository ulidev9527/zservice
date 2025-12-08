# ------------------- 打包
FROM golang:1.24.5-alpine3.22 AS build-step

# 工作目录
WORKDIR /app

# mod准备
COPY go.mod go.mod
COPY go.sum go.sum
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod download

# 全量复制
COPY . .

# 打包
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o __SERVICE_NAME__ ./service/__SERVICE_NAME__/__SERVICE_NAME__.service.go

# ------------------- 运行准备
FROM busybox:1.37.0 AS runtime-step

# 工作目录
WORKDIR /app

# 复制文件
COPY --from=build-step /app/__SERVICE_NAME__ .

ENTRYPOINT ["/app/__SERVICE_NAME__"]