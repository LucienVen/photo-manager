#!/bin/bash

# 图床管理工具启动脚本

echo "🚀 启动图床管理工具..."

# 切换到项目根目录
cd "$(dirname "$0")/.."
echo "📁 当前工作目录: $(pwd)"

# 检查是否安装了 http-server
if ! command -v http-server &> /dev/null; then
    echo "📦 正在安装 http-server..."
    npm install -g http-server
fi

# 启动服务器
echo "🌐 启动本地服务器在 http://localhost:8080"
echo "📁 服务目录: $(pwd)"
echo "📂 可访问的目录: web/, records/"
echo "⏹️  按 Ctrl+C 停止服务器"
echo ""

http-server -c-1 -p 8080 --cors
