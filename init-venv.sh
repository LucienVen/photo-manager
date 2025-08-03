#!/bin/bash

# 1. 设置项目目录
PROJECT_DIR=$(pwd)
VENV_DIR="$PROJECT_DIR/.venv"

# 2. 检查是否已存在 venv
if [ -d "$VENV_DIR" ]; then
  echo "✅ 虚拟环境已存在：$VENV_DIR"
else
  echo "📦 创建虚拟环境..."
  python3 -m venv "$VENV_DIR"
  echo "✅ 虚拟环境创建完成：$VENV_DIR"
fi

# 3. 提示如何激活虚拟环境
echo ""
echo "👉 要激活虚拟环境，请运行："
echo "source $VENV_DIR/bin/activate"

# 4. 安装依赖（如果有 requirements.txt）
if [ -f "$PROJECT_DIR/requirements.txt" ]; then
  echo ""
  echo "📦 安装 requirements.txt 中的依赖..."
  source "$VENV_DIR/bin/activate"
  pip install --upgrade pip
  pip install -r requirements.txt
  deactivate
  echo "✅ 依赖安装完成"
fi

# 5. 添加 .venv 到 .gitignore（如果适用）
if [ -d .git ]; then
  if ! grep -q '^.venv$' .gitignore 2>/dev/null; then
    echo ".venv" >> .gitignore
    echo "📄 已添加 .venv 到 .gitignore"
  fi
fi

