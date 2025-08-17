#!/bin/bash

echo "🔍 图床管理工具状态检查"
echo "=========================="

# 检查服务器是否运行
echo "1. 检查服务器状态..."
if curl --noproxy "*" -s -o /dev/null -w "%{http_code}" http://localhost:8080/ | grep -q "200"; then
    echo "   ✅ 服务器正在运行 (http://localhost:8080)"
else
    echo "   ❌ 服务器未运行"
    exit 1
fi

# 检查 web 目录
echo "2. 检查 web 目录..."
if [ -d "/Users/liangliangtoo/code/photo-manager/web" ]; then
    echo "   ✅ web 目录存在"
    echo "   📁 文件列表:"
    ls -la /Users/liangliangtoo/code/photo-manager/web/ | grep -E "\.(html|js|css)$"
else
    echo "   ❌ web 目录不存在"
fi

# 检查 records 目录
echo "3. 检查 records 目录..."
if [ -d "/Users/liangliangtoo/code/photo-manager/records" ]; then
    echo "   ✅ records 目录存在"
    echo "   📄 数据文件:"
    ls -la /Users/liangliangtoo/code/photo-manager/records/
else
    echo "   ❌ records 目录不存在"
fi

# 测试数据访问
echo "4. 测试数据访问..."
if curl --noproxy "*" -s http://localhost:8080/records/index.json | grep -q "2025-08.json"; then
    echo "   ✅ index.json 可访问"
else
    echo "   ❌ index.json 访问失败"
fi

# 测试图片数据
echo "5. 测试图片数据..."
if curl --noproxy "*" -s http://localhost:8080/records/2025-08.json | grep -q "filename"; then
    echo "   ✅ 图片数据可访问"
else
    echo "   ❌ 图片数据访问失败"
fi

echo ""
echo "🎯 访问地址:"
echo "   主页: http://localhost:8080/web/"
echo "   测试页: http://localhost:8080/web/test.html"
echo ""
echo "📋 使用说明:"
echo "   1. 在浏览器中打开 http://localhost:8080/web/"
echo "   2. 应用会自动加载图片数据"
echo "   3. 可以使用搜索功能筛选图片"
echo "   4. 点击图片查看详细信息"
