# 🚀 快速启动指南

## 最简单的方式

```bash
# 进入web目录
cd web

# 启动服务器（推荐）
./start-node.sh

# 或者使用npm
npm start
```

## 访问地址

启动成功后，在浏览器中访问：
**http://localhost:8000/web/**

## 如果遇到问题

### 1. 404 错误

确保服务器从项目根目录启动，这样 web 和 records 目录都可以访问。

### 2. 数据加载失败

检查 `records/` 目录下是否有 JSON 文件，如 `2025-08.json`。

### 3. Node.js 未安装

```bash
# macOS
brew install node

# Windows
# 下载 https://nodejs.org/

# Linux
sudo apt install nodejs npm
```

## 文件结构

```
photo-manager/
├── web/           # 前端页面
│   ├── index.html
│   ├── styles.css
│   ├── app.js
│   └── start-node.sh
└── records/       # 数据文件
    └── 2025-08.json
```
