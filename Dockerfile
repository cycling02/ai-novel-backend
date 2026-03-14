# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 安装 build-essential 和 git
RUN apk add --no-cache build-base git

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ai-novel-backend ./cmd/server

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates 用于 HTTPS
RUN apk add --no-cache ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/ai-novel-backend .

# 复制配置文件
COPY config.yaml .
COPY .env.example .env

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./ai-novel-backend"]
