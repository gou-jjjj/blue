# 使用 golang:1.21-alpine 作为基础镜像进行构建
FROM golang:1.21-alpine AS builder

# 环境变量设置
ENV CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /build

# 安装 make 工具
RUN apk add --no-cache make

# 将当前目录下的文件和目录全部复制到工作目录中
COPY . .

# 运行 go mod tidy 清理依赖
RUN go mod tidy

# 根据 Makefile 文件构建应用
RUN  make bin

# 使用 alpine 作为运行时镜像的基础
FROM alpine:latest

# 复制构建出的二进制文件从构建镜像到运行时镜像
COPY --from=builder /build/bin/linux/blue-server /blue-server
COPY --from=builder /build/bin/linux/blue-client /blue-client
COPY --from=builder /build/bin/linux/blue-cli.conf /blue-cli.conf
COPY --from=builder /build/bin/linux/blue-server.json /blue-server.json

EXPOSE 13140
# 设置容器启动时执行的命令
CMD ["./blue-server"]
