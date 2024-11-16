# 使用官方 Go 语言镜像作为基础镜像
FROM golang:1.23.3 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目源代码
COPY . .

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o ogimg ./cmd/server/main.go

# 使用轻量级的基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /root/
COPY config/ ./config/

# 复制构建好的可执行文件
COPY --from=builder /app/ogimg .
# 暴露应用程序的端口
EXPOSE 8888

# 运行可执行文件
CMD ["./ogimg"]
