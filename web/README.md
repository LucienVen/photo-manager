# 图床管理工具 - Vue3 前端

一个基于 Vue3 的图床图片展示和管理工具。

## 功能特性

- 📸 图片展示和管理
- 🔍 多条件搜索（文件名、标签、描述）
- 📊 统计信息显示
- 📱 响应式设计
- 🎨 现代化界面
- 📋 一键复制图片链接

## 快速开始

### 方法一：使用启动脚本（推荐）

```bash
cd web
./start.sh
```

### 方法二：手动启动

```bash
# 在项目根目录下
cd web
npm install
npm run dev
```

### 方法三：直接使用 http-server

```bash
# 在项目根目录下
npx http-server -c-1 -p 8080 --cors
```

## 访问地址

启动后访问：http://localhost:8080/web/

## 目录结构

```
web/
├── index.html      # 主页面
├── app.js          # Vue3 应用代码
├── styles.css      # 样式文件
├── start.sh        # 启动脚本
├── package.json    # 依赖配置
└── README.md       # 说明文档
```

## 数据格式

应用会读取 `../records/` 目录下的 JSON 文件：

- `index.json` - 索引文件，包含所有记录文件名
- `*.json` - 图片记录文件，包含图片信息

## 技术栈

- Vue 3 (Composition API)
- 原生 JavaScript
- CSS3 (Grid/Flexbox)
- HTTP Server

## 浏览器支持

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+
