#!/bin/bash
set -e

echo "🔨 构建 AI Novel Backend..."

# 设置目标平台 (阿里云 FC 使用 Linux)
export GOOS=linux
export GOARCH=amd64

# 清理旧构建
rm -f ai-novel-backend

# 构建
echo "📦 编译中..."
go build -o ai-novel-backend ./cmd/server

echo "✅ 构建完成！输出文件：ai-novel-backend"

# 验证
if [ -f "ai-novel-backend" ]; then
    echo "📊 二进制文件信息:"
    ls -lh ai-novel-backend
    file ai-novel-backend
else
    echo "❌ 构建失败"
    exit 1
fi
