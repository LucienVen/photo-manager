# 测试指南

## 概述
本文档说明如何执行项目中的各种测试。

## 前置条件

### 1. 安装依赖
```bash
pip install -r requirements.txt
```

### 2. 配置环境变量
```bash
# 复制示例配置文件
cp env.example .env

# 编辑配置文件，设置你的实际值
vim .env
```

## 测试类型

### 1. 配置测试
测试环境变量配置是否正确加载。

```bash
python test_config.py
```

**预期输出：**
```
=== 配置测试 ===
当前配置:
  GitHub用户名: your-github-username
  GitHub仓库: your-repo-name
  CDN前缀: https://cdn.jsdelivr.net/gh
  缩略图宽度: 320
  调试模式: False
  完整CDN URL: https://cdn.jsdelivr.net/gh/your-github-username/your-repo-name

=== 配置验证 ===
[SUCCESS] 配置验证通过
[SUCCESS] GitHub用户名: your-github-username
[SUCCESS] GitHub仓库: your-repo-name
[SUCCESS] CDN URL: https://cdn.jsdelivr.net/gh/your-github-username/your-repo-name
```

### 2. 功能测试
测试图片上传功能。

#### 2.1 测试文件不存在的情况
```bash
python upload.py /path/to/nonexistent/image.jpg
```

**预期输出：**
```
[ERROR] 文件不存在: /path/to/nonexistent/image.jpg
```

#### 2.2 测试不支持的格式
```bash
python upload.py /path/to/document.txt
```

**预期输出：**
```
[WARNING] 不支持的图片格式: .txt
```

#### 2.3 测试正常上传
```bash
# 需要准备一个测试图片文件
python upload.py /path/to/test-image.jpg
```

**预期输出：**
```
[SUCCESS] 成功上传: test-image.abc12345.jpg
```

### 3. 集成测试
测试整个工作流程。

```bash
# 1. 确保配置正确
python test_config.py

# 2. 测试上传功能
python upload.py /path/to/test-image.jpg

# 3. 检查生成的文件
ls -la photo-bed/images/
ls -la photo-bed/thumbs/
cat photo-manager/images.json
```

## 测试数据准备

### 创建测试图片
如果没有测试图片，可以使用以下方法创建：

```bash
# 使用 ImageMagick 创建测试图片
convert -size 800x600 xc:red test-image.jpg

# 或者使用 Python 创建
python -c "
from PIL import Image, ImageDraw
img = Image.new('RGB', (800, 600), color='blue')
draw = ImageDraw.Draw(img)
draw.text((400, 300), 'Test Image', fill='white')
img.save('test-image.jpg')
"
```

## 故障排除

### 常见问题

1. **ModuleNotFoundError: No module named 'dotenv'**
   ```bash
   pip install python-dotenv
   ```

2. **ModuleNotFoundError: No module named 'PIL'**
   ```bash
   pip install Pillow
   ```

3. **配置验证失败**
   - 检查 `.env` 文件是否存在
   - 确保设置了正确的 `GITHUB_NAME` 和 `GITHUB_REPO`
   - 确保没有使用默认值 `yourname` 和 `photo-bed`

4. **图片处理失败**
   - 确保安装了 ImageMagick：`brew install imagemagick`
   - 或者使用 PIL 替代 ImageMagick

### 调试模式
在 `.env` 文件中设置：
```env
DEBUG=true
```

## 自动化测试

### 创建测试脚本
```bash
#!/bin/bash
# test_all.sh

echo "开始执行所有测试..."

# 1. 配置测试
echo "1. 执行配置测试..."
python test_config.py
if [ $? -ne 0 ]; then
    echo "[ERROR] 配置测试失败"
    exit 1
fi

# 2. 功能测试
echo "2. 执行功能测试..."
# 这里可以添加更多测试用例

echo "[SUCCESS] 所有测试通过"
```

### 运行自动化测试
```bash
chmod +x test_all.sh
./test_all.sh
```

## 持续集成

### GitHub Actions 示例
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Python
      uses: actions/setup-python@v2
      with:
        python-version: 3.9
    - name: Install dependencies
      run: |
        pip install -r requirements.txt
    - name: Run tests
      run: |
        python test_config.py
```

## 测试覆盖率

### 安装覆盖率工具
```bash
pip install coverage
```

### 运行覆盖率测试
```bash
coverage run test_config.py
coverage report
coverage html  # 生成HTML报告
```

## 性能测试

### 测试大文件处理
```bash
# 创建大尺寸测试图片
convert -size 4000x3000 xc:green large-test-image.jpg

# 测试上传性能
time python upload.py large-test-image.jpg
```

### 测试批量处理
```bash
# 创建多个测试文件
for i in {1..10}; do
    convert -size 800x600 xc:blue test-$i.jpg
done

# 批量上传测试
for file in test-*.jpg; do
    python upload.py "$file"
done
``` 