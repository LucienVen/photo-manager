#!/bin/bash

# 图床管理工具前端页面启动脚本

echo "🚀 启动图床管理工具前端展示页面..."

# 检查是否安装了Node.js
if command -v node &> /dev/null; then
    echo "📦 使用Node.js启动本地服务器..."
    echo "🌐 页面将在 http://localhost:8000 打开"
    echo "📁 请确保records目录下有JSON数据文件"
    echo ""
    echo "按 Ctrl+C 停止服务器"
    echo ""
    npx http-server -p 8000
elif command -v python3 &> /dev/null; then
    echo "📦 使用Python3启动本地服务器..."
    echo "🌐 页面将在 http://localhost:8000 打开"
    echo "📁 请确保records目录下有JSON数据文件"
    echo ""
    echo "按 Ctrl+C 停止服务器"
    echo ""
    python3 -m http.server 8000
elif command -v python &> /dev/null; then
    echo "📦 使用Python启动本地服务器..."
    echo "🌐 页面将在 http://localhost:8000 打开"
    echo "📁 请确保records目录下有JSON数据文件"
    echo ""
    echo "按 Ctrl+C 停止服务器"
    echo ""
    python -m http.server 8000
else
    echo "❌ 未找到可用的HTTP服务器"
    echo "请安装以下任一工具："
    echo "  - Node.js: brew install node (macOS)"
    echo "  - Python3: brew install python3 (macOS)"
    echo ""
    echo "或者直接在浏览器中打开 index.html 文件"
    exit 1
fi
