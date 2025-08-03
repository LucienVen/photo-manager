# 测试执行指南

## 快速开始

### 一键测试
```bash
./quick_test.sh
```

这个脚本会自动：
- 检查Python环境
- 安装依赖
- 创建配置文件
- 执行所有测试

### 手动测试

#### 1. 配置测试
```bash
python test_config.py
```

#### 2. 功能测试
```bash
python test_upload.py
```

## 测试内容

### 配置测试
- ✅ 环境变量加载
- ✅ GitHub配置验证
- ✅ CDN URL生成

### 功能测试
- ✅ 文件不存在处理
- ✅ 不支持格式处理
- ✅ 正常上传流程
- ✅ 重复上传检测

## 测试结果示例

```
=== 配置测试 ===
当前配置:
  GitHub用户名: example-user
  GitHub仓库: my-photo-bed
  CDN前缀: https://cdn.jsdelivr.net/gh
  缩略图宽度: 320
  调试模式: False
  完整CDN URL: https://cdn.jsdelivr.net/gh/example-user/my-photo-bed

=== 配置验证 ===
[SUCCESS] 配置验证通过

=== 功能测试 ===
[ERROR] 文件不存在: /path/to/nonexistent/image.jpg
[WARNING] 不支持的图片格式: .txt
[SUCCESS] 成功上传: test-image.abc12345.jpg
[SKIP] 已存在相同图片，跳过上传
```

## 故障排除

### 常见问题

1. **依赖缺失**
   ```bash
   pip install -r requirements.txt
   ```

2. **配置文件问题**
   ```bash
   cp env.example .env
   # 编辑.env文件设置正确的值
   ```

3. **权限问题**
   ```bash
   chmod +x quick_test.sh
   ```

## 详细文档

更多测试信息请查看 `TEST_GUIDE.md` 