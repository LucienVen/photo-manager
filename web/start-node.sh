#!/bin/bash

# 图床管理工具前端页面启动脚本 (Node.js版本)

echo "🚀 启动图床管理工具前端展示页面..."

# 检查是否安装了Node.js
if command -v node &> /dev/null; then
    echo "📦 使用Node.js启动本地服务器..."
    echo "🌐 页面将在 http://localhost:8000 打开"
    echo "📁 请确保records目录下有JSON数据文件"
    echo ""
    echo "按 Ctrl+C 停止服务器"
    echo ""
    
    # 检查是否有package.json，如果有则使用npm start
    if [ -f "package.json" ]; then
        echo "📋 检测到package.json，使用npm start启动..."
        # 从项目根目录启动，这样web和records目录都可以访问
        cd ..
        npx http-server -p 8000
    else
        echo "📋 使用npx http-server启动..."
        # 从项目根目录启动，这样web和records目录都可以访问
        cd ..
        npx http-server -p 8000
    fi
else
    echo "❌ 未找到Node.js"
    echo "请安装Node.js："
    echo "  - macOS: brew install node"
    echo "  - Windows: 下载 https://nodejs.org/"
    echo "  - Linux: sudo apt install nodejs npm"
    echo ""
    echo "或者使用Python启动："
    echo "  python3 -m http.server 8000"
    exit 1
fi
