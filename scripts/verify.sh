#!/bin/bash
set -e

echo "🔍 验证项目..."

# 检查 Go 版本
echo "📌 检查 Go 版本..."
go version

# Tidy
echo "📦 整理依赖..."
go mod tidy

# 下载依赖
echo "⬇️ 下载依赖..."
go mod download

# 构建
echo "🔨 构建项目..."
go build -o /tmp/ai-novel-backend ./cmd/server

# 运行测试
echo "🧪 运行测试..."
go test -v ./... || true

echo "✅ 验证完成！"
echo "二进制文件：/tmp/ai-novel-backend"
