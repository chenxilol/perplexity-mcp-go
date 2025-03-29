FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum* ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY *.go ./

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o perplexity-search-mcp

# 使用轻量级镜像
FROM alpine:latest

# 安装CA证书（用于HTTPS请求）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从builder阶段复制编译好的应用
COPY --from=builder /app/perplexity-search-mcp .

# 运行应用
ENTRYPOINT ["./perplexity-search-mcp"] 