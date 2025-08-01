# Dockerfile - 应用程序容器化配置文件
# 使用多阶段构建优化镜像大小

# 第一阶段：构建阶段
# 使用官方 Go 镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
# 先复制依赖文件，利用 Docker 缓存机制
COPY go.mod go.sum ./

# 下载依赖
# 安装项目所需的 Go 模块
RUN go mod download

# 复制源代码
# 复制所有源代码到容器中
COPY . .

# 构建应用
# 编译 Go 代码生成可执行文件
# CGO_ENABLED=0 禁用 CGO，生成静态链接的二进制文件
# GOOS=linux 指定目标操作系统
# -a -installsuffix cgo 强制重新编译所有包
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lemon-tree-core .

# 第二阶段：运行阶段
# 使用轻量级的 alpine 镜像作为运行环境
FROM alpine:latest

# 安装 ca-certificates 用于 HTTPS 请求
# 确保应用程序可以访问 HTTPS 服务
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
# 只复制编译好的可执行文件，不包含源代码
COPY --from=builder /app/lemon-tree-core .

# 复制配置文件
# 复制应用程序的配置文件
COPY --from=builder /app/config ./config

# 暴露端口
# 声明应用程序监听的端口
EXPOSE 8080

# 运行应用
# 容器启动时执行的命令
CMD ["./lemon-tree-core"]