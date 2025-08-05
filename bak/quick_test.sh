#!/bin/bash
# 快速测试脚本
# 一键执行所有测试

echo "🚀 开始执行快速测试..."

# 检查Python环境
echo "1. 检查Python环境..."
python --version
if [ $? -ne 0 ]; then
    echo "[ERROR] Python未安装或不在PATH中"
    exit 1
fi

# 检查依赖
echo "2. 检查依赖..."
python -c "import dotenv, PIL" 2>/dev/null
if [ $? -ne 0 ]; then
    echo "[WARNING] 缺少依赖，正在安装..."
    pip install -r requirements.txt
fi

# 检查配置文件
echo "3. 检查配置文件..."
if [ ! -f ".env" ]; then
    echo "[WARNING] .env文件不存在，正在创建..."
    cp env.example .env
    echo "[INFO] 请编辑.env文件设置正确的配置"
fi

# 执行配置测试
echo "4. 执行配置测试..."
python test_config.py
if [ $? -ne 0 ]; then
    echo "[ERROR] 配置测试失败"
    exit 1
fi

# 执行功能测试
echo "5. 执行功能测试..."
python test_upload.py
if [ $? -ne 0 ]; then
    echo "[ERROR] 功能测试失败"
    exit 1
fi

echo ""
echo "🎉 所有测试通过！"
echo ""
echo "📋 测试总结："
echo "  ✅ Python环境正常"
echo "  ✅ 依赖安装完成"
echo "  ✅ 配置文件正确"
echo "  ✅ 配置加载正常"
echo "  ✅ 上传功能正常"
echo ""
echo "📖 更多测试信息请查看 TEST_GUIDE.md" 