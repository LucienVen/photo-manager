# photo-manager

> github 搭建个人图床，图床简单管理工具

[toc]

## 希望实现效果：

> tags 标签信息比 desc 描述信息重要

```shell
# pic-go 图片本地路径 tags:标签[tag1,tag2...] desc:描述note
$ pic-go demo1.jpg tags:风景,广州塔,金融城 desc:拍摄于广州
$ pic-go demo1.jpg 风景,广州塔,金融城 拍摄于广州
$ pic-go demo1.jpg 风景,广州塔,金融城 (只有一个参数，识别为 tags)
$ pic-go demo1.jpg 风景(如果单独一个句子，则识别为 tags)
$ pic-go demo1.jpg tags:风景(增加 tags标签：则识别为标签)
$ pic-go demo1.jpg desc:风景(增加 desc标签)

$ [返回 json 信息]

```

## 关于存储：

- 目前实现，使用 json，写入文件

  - 多个按月/年分割的 JSON 文件
  - records/2025-08.json

- 后续：首选考虑 Elasticsearch 可以支持全文搜索的平台

  - 考虑基础 sql 数据库,（pg mysql）

- 使用 json 记录

  ```json
  [
    {
      "filename": "sunset.<hash>.jpg", // 组装 hash
      "url": "https://cdn.jsdelivr.net/gh/your-username/photo-bed/images/2025/08/sunset.jpg",
      "thumb_url": "https://cdn.jsdelivr.net/gh/your-username/photo-bed/thumbs/2025/08/sunset.jpg",
      "path": "images/2025/08/sunset.jpg",
      "uploaded_at": "2025-08-03T23:30:00+08:00",
      "tags": ["sunset", "nature", "evening"],
      "desc": "拍摄于云南香格里拉，使用富士XT4",
      "size_kb": 320,
      "width": 1920,
      "height": 1080,
      "hash": "e4d909c290d0fb1ca068ffaddf22cbd0"
    }
  ]
  ```

  | 字段名        | 类型     | 说明                                                         |
  | ------------- | -------- | ------------------------------------------------------------ |
  | `filename`    | string   | 原始文件名或别名，便于识别 , 组装 hash , 如 sunset.`<hash>`.jpg |
  | `url`         | string   | CDN 图片地址（如 jsDelivr）                                  |
  | thumb_url     | string   | 略缩图地址                                                   |
  | path          | string   | 仓库中的路径（相对路径）                                     |
  | `uploaded_at` | int      | 上传时间，时间戳                                             |
  | `tags`        | string[] | 标签，用于搜索、分类                                         |
  | `desc`        | string   | 可选备注，描述来源或用途                                     |
  | `size_kb`     | number   | 文件大小（KB）                                               |
  | `width`       | number   | 图片宽度（像素）                                             |
  | `height`      | number   | 图片高度（像素）                                             |
  | `hash`        | string   | 文件内容的 hash（如 MD5 或 SHA256），避免重复上传            |

# 项目代码结构(TODO)：

```shell
- web
- records
	- 2025-08.json
	- 2025-09.json
- clean-script.sh # 清理历史图片
- docker-compose.yml # docker 打包发布
```

# 环境准备

- ~~python~~
- npm
- golang
- docker 用于打包生成可执行文件

```shell
$ npm install -g picgo
$ picgo -v
1.5.9
```

### picgo 配置

```shell
$ picgo set uploader

? Choose a(n) uploader github
设定仓库名 格式：username/repo yourname/photo-bed # 你的图床仓库
设定分支名 例如：main main
? token: [hidden]
设定存储路径 例如：test/ images # 仓库存储图片的文件夹路径名称
设定自定义域名 例如：https://test.com https://cdn.jsdelivr.net/gh/yourname/photo-bed@main
[PicGo SUCCESS]: Configure config successfully!
```

### 测试

```shell
$ picgo upload Snipaste_2024-03-29_22-51-21.png
[PicGo INFO]: Before transform
[PicGo INFO]: Transforming... Current transformer is [path]
[PicGo INFO]: Before upload
[PicGo INFO]: Uploading... Current uploader is [github]
[PicGo SUCCESS]:
https://cdn.jsdelivr.net/gh/yourname/photo-bed@main/imagesSnipaste_2024-03-29_22-51-21.png
```

# 工作流程简述

1. 脚本包装 picgo 上传图片，构建图床信息

   1. photo-manager 存储图床信息

2. 上传图床（githun photo-bed ）触发 github action

# golang 脚本工作流程

1. 读取 env 配置
2. 接收图片路径
3. 生成文件 hash（SHA256 推荐）
4. 重命名并生成缩略图（使用 xxxx）
5. 将两者复制到 photo-bed 中对应目录
6. 更新 photo-manager 中的记录文件（JSON、YAML、TOML 均可选）
7. 支持跳过重复上传（通过 hash 检查）

## ~~Python 脚本打包为可执行文件：~~

```shell
# 可以使用 pyinstaller 打包为 .app 或 .exe：
pip install pyinstaller
pyinstaller --onefile upload.py

```
