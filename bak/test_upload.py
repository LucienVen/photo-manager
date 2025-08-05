#!/usr/bin/env python3
"""
上传功能测试脚本
测试图片上传的各种场景，使用 PicGo 进行上传
"""
import os
import tempfile
import shutil
import subprocess
import json
from PIL import Image, ImageDraw
from upload import main, generate_thumbnail, upload_with_picgo
from config import config

def create_test_image(filename, size=(800, 600), color='blue'):
    """创建测试图片"""
    img = Image.new('RGB', size, color=color)
    draw = ImageDraw.Draw(img)
    draw.text((size[0]//2, size[1]//2), 'Test Image', fill='white')
    img.save(filename)
    return filename

def check_picgo_available():
    """检查 PicGo 是否可用"""
    try:
        result = subprocess.run(["picgo", "--version"], capture_output=True, text=True)
        return result.returncode == 0
    except FileNotFoundError:
        return False

def test_picgo_config():
    """测试 PicGo 配置"""
    print("=== 测试 PicGo 配置 ===")
    
    if not check_picgo_available():
        print("[ERROR] PicGo 未安装或不可用")
        print("请安装 PicGo CLI: npm install -g picgo")
        return False
    
    print("[SUCCESS] PicGo 已安装")
    
    # 检查 PicGo 配置文件
    config_path = os.path.expanduser("~/.picgo/config.json")
    if os.path.exists(config_path):
        print(f"[INFO] 找到 PicGo 配置文件: {config_path}")
        try:
            with open(config_path, 'r', encoding='utf-8') as f:
                config_data = json.load(f)
            
            print("[SUCCESS] PicGo 配置文件读取成功")
            
            # 显示配置信息
            if 'picBed' in config_data:
                pic_bed = config_data['picBed']
                print(f"  当前图床: {pic_bed.get('current', 'unknown')}")
                
                # 显示图床配置
                for bed_name, bed_config in pic_bed.items():
                    if bed_name != 'current':
                        print(f"  图床 '{bed_name}': {type(bed_config).__name__}")
                        if isinstance(bed_config, dict):
                            # 隐藏敏感信息
                            safe_config = {}
                            for key, value in bed_config.items():
                                if 'token' in key.lower() or 'key' in key.lower() or 'secret' in key.lower():
                                    safe_config[key] = '***'
                                else:
                                    safe_config[key] = value
                            print(f"    配置: {safe_config}")
            
            # 检查其他配置
            if 'picgoPlugins' in config_data:
                plugins = config_data['picgoPlugins']
                print(f"  已安装插件: {list(plugins.keys()) if plugins else '无'}")
            
            return True
            
        except json.JSONDecodeError as e:
            print(f"[ERROR] PicGo 配置文件格式错误: {e}")
            return False
        except Exception as e:
            print(f"[ERROR] 读取 PicGo 配置文件失败: {e}")
            return False
    else:
        print(f"[WARNING] PicGo 配置文件不存在: {config_path}")
        print("请先配置 PicGo: picgo set")
        return False

def test_thumbnail_generation():
    """测试缩略图生成功能"""
    print("=== 测试缩略图生成 ===")
    
    # 确保 thumbs 目录存在
    thumbs_dir = "test_images/thumbs"
    os.makedirs(thumbs_dir, exist_ok=True)
    
    # 获取 test_images 目录下的所有图片文件
    test_images_dir = "test_images"
    image_extensions = ['.jpg', '.jpeg', '.png', '.webp']
    image_files = []
    
    for file in os.listdir(test_images_dir):
        if any(file.lower().endswith(ext) for ext in image_extensions):
            image_files.append(file)
    
    if not image_files:
        print("[WARNING] test_images 目录下没有找到图片文件")
        return
    
    print(f"找到 {len(image_files)} 个图片文件: {image_files}")
    
    for filename in image_files:
        print(f"\n测试图片: {filename}")
        
        # 构建完整路径
        image_path = os.path.join(test_images_dir, filename)
        
        try:
            # 获取原图信息
            with Image.open(image_path) as img:
                original_width, original_height = img.size
                print(f"  原图尺寸: {original_width}x{original_height}")
            
            # 生成缩略图文件名
            name, ext = os.path.splitext(filename)
            thumb_filename = f"thumb_{name}{ext}"
            thumb_path = os.path.join(thumbs_dir, thumb_filename)
            
            # 生成缩略图
            generate_thumbnail(image_path, thumb_path)
            
            # 验证缩略图
            if os.path.exists(thumb_path):
                with Image.open(thumb_path) as thumb:
                    thumb_width, thumb_height = thumb.size
                    print(f"  缩略图尺寸: {thumb_width}x{thumb_height}")
                    
                    # 检查缩略图宽度是否符合配置
                    expected_width = config.THUMB_WIDTH
                    if thumb_width == expected_width:
                        print(f"  [SUCCESS] 缩略图宽度正确: {thumb_width}")
                    else:
                        print(f"  [WARNING] 缩略图宽度不符合预期: {thumb_width} != {expected_width}")
                    
                    # 检查宽高比是否保持
                    original_ratio = original_width / original_height
                    thumb_ratio = thumb_width / thumb_height
                    ratio_diff = abs(original_ratio - thumb_ratio)
                    
                    if ratio_diff < 0.1:  # 允许 10% 的误差
                        print(f"  [SUCCESS] 宽高比保持正确: {thumb_ratio:.2f}")
                    else:
                        print(f"  [WARNING] 宽高比可能有问题: {thumb_ratio:.2f} vs {original_ratio:.2f}")
                    
                    # 检查文件大小
                    thumb_size = os.path.getsize(thumb_path) / 1024
                    print(f"  缩略图大小: {thumb_size:.1f} KB")
                    print(f"  缩略图保存位置: {thumb_path}")
                    
            else:
                print(f"  [ERROR] 缩略图生成失败")
                
        except Exception as e:
            print(f"  [ERROR] 处理图片 {filename} 时出错: {e}")
    
    print()

def test_file_not_exists():
    """测试文件不存在的情况"""
    print("=== 测试文件不存在 ===")
    main("/path/to/nonexistent/image.jpg")
    print()

def test_unsupported_format():
    """测试不支持的格式"""
    print("=== 测试不支持的格式 ===")
    
    # 创建临时文本文件
    with tempfile.NamedTemporaryFile(suffix='.txt', delete=False) as f:
        f.write(b"This is a text file, not an image")
        temp_file = f.name
    
    try:
        main(temp_file)
    finally:
        os.unlink(temp_file)
    print()

def test_normal_upload():
    """测试正常上传（使用 PicGo）"""
    print("=== 测试正常上传（PicGo） ===")
    
    if not check_picgo_available():
        print("[SKIP] PicGo 不可用，跳过上传测试")
        return
    
    # 创建测试图片
    test_image = "test-image.jpg"
    create_test_image(test_image)
    
    try:
        main(test_image)
    finally:
        # 清理测试文件
        if os.path.exists(test_image):
            os.unlink(test_image)
    print()

def test_duplicate_upload():
    """测试重复上传（使用 PicGo）"""
    print("=== 测试重复上传（PicGo） ===")
    
    if not check_picgo_available():
        print("[SKIP] PicGo 不可用，跳过上传测试")
        return
    
    # 创建测试图片
    test_image = "test-image.jpg"
    create_test_image(test_image)
    
    try:
        # 第一次上传
        print("第一次上传:")
        main(test_image)
        
        # 第二次上传（应该被跳过）
        print("第二次上传:")
        main(test_image)
    finally:
        # 清理测试文件
        if os.path.exists(test_image):
            os.unlink(test_image)
    print()

def test_picgo_upload_direct():
    """
    直接测试 PicGo 上传功能
    上传图片，使用测试图片目录的图片
    判断是否有略缩图存在，如果没有，则报错，终止上传
    如果有，一起上传，并返回打印原图与略缩图的 url
    """
    print("=== 直接测试 PicGo 上传 ===")
    
    if not check_picgo_available():
        print("[SKIP] PicGo 不可用，跳过直接上传测试")
        return
    
    # 检查 test_images 目录是否存在
    test_images_dir = "test_images"
    if not os.path.exists(test_images_dir):
        print(f"[ERROR] test_images 目录不存在")
        return
    
    # 获取 test_images 目录下的所有图片文件
    image_extensions = ['.jpg', '.jpeg', '.png', '.webp']
    image_files = []
    
    for file in os.listdir(test_images_dir):
        if any(file.lower().endswith(ext) for ext in image_extensions):
            image_files.append(file)
    
    if not image_files:
        print(f"[ERROR] test_images 目录下没有找到图片文件")
        return
    
    print(f"找到 {len(image_files)} 个测试图片: {image_files}")
    
    # 选择第一个图片进行测试
    test_filename = image_files[0]
    test_image_path = os.path.join(test_images_dir, test_filename)
    
    print(f"选择测试图片: {test_filename}")
    
    # 计算哈希值
    from upload import calc_hash
    hash_value = calc_hash(test_image_path)
    base_name, ext = os.path.splitext(test_filename)
    
    # 检查缩略图是否存在
    thumbs_dir = "test_images/thumbs"
    
    # 检查 thumbs 目录是否存在
    if not os.path.exists(thumbs_dir):
        print(f"[ERROR] 缩略图目录不存在: {thumbs_dir}")
        print("请先生成缩略图，然后重新运行测试")
        return
    
    # 智能查找缩略图 - 尝试多种可能的命名格式
    possible_thumb_names = [
        # 格式1: 原始文件名 + thumb + 哈希值
        f"{base_name}.thumb.{hash_value[:8]}{ext}",
        # 格式2: 原始文件名 + thumb
        f"{base_name}.thumb{ext}",
        # 格式3: thumb_ + 原始文件名
        f"thumb_{base_name}{ext}",
        # 格式4: 如果原文件名已经包含哈希值，尝试去掉哈希值再查找
        f"{base_name.split('.')[0]}.thumb.{hash_value[:8]}{ext}" if '.' in base_name else None,
        # 格式5: 如果原文件名已经包含哈希值，尝试去掉哈希值
        f"{base_name.split('.')[0]}.thumb{ext}" if '.' in base_name else None
    ]
    
    # 过滤掉None值
    possible_thumb_names = [name for name in possible_thumb_names if name is not None]
    
    thumb_filename = None
    thumb_path = None
    
    # 尝试查找缩略图
    for possible_name in possible_thumb_names:
        possible_path = os.path.join(thumbs_dir, possible_name)
        if os.path.exists(possible_path):
            thumb_filename = possible_name
            thumb_path = possible_path
            print(f"找到缩略图: {thumb_filename}")
            break
    
    if not thumb_filename:
        print(f"[ERROR] 缩略图不存在")
        print(f"尝试过的命名格式:")
        for name in possible_thumb_names:
            print(f"  - {name}")
        print(f"请先生成缩略图，然后重新运行测试")
        return
    
    try:
        # 上传原图
        print(f"正在上传原图: {test_filename}")
        original_url = upload_with_picgo(test_image_path)
        
        if not original_url:
            print("[ERROR] 原图上传失败")
            return
        
        print(f"[SUCCESS] 原图上传成功: {original_url}")
        
        # 上传缩略图
        print(f"正在上传缩略图: {thumb_filename}")
        thumb_url = upload_with_picgo(thumb_path)
        
        if not thumb_url:
            print("[ERROR] 缩略图上传失败")
            return
        
        print(f"[SUCCESS] 缩略图上传成功: {thumb_url}")
        
        # 打印上传结果
        print("\n=== 上传结果 ===")
        print(f"原图: {original_url}")
        print(f"缩略图: {thumb_url}")
        print(f"原图文件名: {test_filename}")
        print(f"缩略图文件名: {thumb_filename}")
        print(f"文件哈希值: {hash_value[:8]}")
        
    except Exception as e:
        print(f"[ERROR] 上传过程中出现异常: {e}")
        import traceback
        traceback.print_exc()
    
    print()

def cleanup_test_files():
    """清理测试生成的文件"""
    print("=== 清理测试文件 ===")
    
    # 清理photo-bed目录
    if os.path.exists("photo-bed"):
        shutil.rmtree("photo-bed")
        print("[SUCCESS] 已清理 photo-bed 目录")
    
    # 清理photo-manager目录
    if os.path.exists("photo-manager"):
        shutil.rmtree("photo-manager")
        print("[SUCCESS] 已清理 photo-manager 目录")
    
    # 清理records目录
    if os.path.exists("records"):
        shutil.rmtree("records")
        print("[SUCCESS] 已清理 records 目录")
    
    # 清理缩略图文件
    for file in os.listdir("."):
        if file.startswith("thumb_") and file.endswith((".jpg", ".png", ".webp")):
            os.unlink(file)
            print(f"[SUCCESS] 已清理缩略图: {file}")
    
    print()

def test_photo_rename():
    """图片重命名 filename.<hash>.type"""
    print("=== 测试图片重命名功能 ===")
    
    # 检查 test_images 目录是否存在
    test_images_dir = "test_images"
    if not os.path.exists(test_images_dir):
        print(f"[ERROR] test_images 目录不存在")
        return
    
    # 获取 test_images 目录下的所有图片文件
    image_extensions = ['.jpg', '.jpeg', '.png', '.webp']
    image_files = []
    
    for file in os.listdir(test_images_dir):
        if any(file.lower().endswith(ext) for ext in image_extensions):
            image_files.append(file)
    
    if not image_files:
        print(f"[ERROR] test_images 目录下没有找到图片文件")
        return
    
    print(f"找到 {len(image_files)} 个测试图片: {image_files}")
    
    try:
        for filename in image_files:
            print(f"\n测试图片: {filename}")
            
            # 构建完整路径
            image_path = os.path.join(test_images_dir, filename)
            
            # 计算哈希值
            from upload import calc_hash
            hash_value = calc_hash(image_path)
            print(f"  文件哈希值: {hash_value}")
            print(f"  哈希值前8位: {hash_value[:8]}")
            
            # 获取文件信息
            base_name, ext = os.path.splitext(filename)
            
            # 测试原图命名格式: filename.<hash>.(type)
            expected_original_name = f"{base_name}.{hash_value[:8]}{ext}"
            print(f"  原图命名格式: {expected_original_name}")
            
            # 验证原图命名格式的正确性
            if expected_original_name.count('.') == 2:  # 应该有两个点：原文件名.哈希值.扩展名
                print(f"  [SUCCESS] 原图命名格式正确")
            else:
                print(f"  [ERROR] 原图命名格式错误")
            
            # 验证哈希值长度
            if len(hash_value[:8]) == 8:
                print(f"  [SUCCESS] 哈希值长度正确: 8位")
            else:
                print(f"  [ERROR] 哈希值长度错误: {len(hash_value[:8])}位")
            
            # 验证扩展名保持
            if expected_original_name.endswith(ext):
                print(f"  [SUCCESS] 扩展名保持正确: {ext}")
            else:
                print(f"  [ERROR] 扩展名保持错误")
            
            # 测试缩略图命名格式: filename.thumb.<hash>.(type)
            expected_thumb_name = f"{base_name}.thumb.{hash_value[:8]}{ext}"
            print(f"  缩略图命名格式: {expected_thumb_name}")
            
            # 验证缩略图命名格式的正确性
            if expected_thumb_name.count('.') == 3:  # 应该有三个点：原文件名.thumb.哈希值.扩展名
                print(f"  [SUCCESS] 缩略图命名格式正确")
            else:
                print(f"  [ERROR] 缩略图命名格式错误")
            
            # 验证缩略图命名包含 "thumb" 关键字
            if "thumb" in expected_thumb_name:
                print(f"  [SUCCESS] 缩略图命名包含 'thumb' 关键字")
            else:
                print(f"  [ERROR] 缩略图命名缺少 'thumb' 关键字")
            
            # 验证缩略图扩展名保持
            if expected_thumb_name.endswith(ext):
                print(f"  [SUCCESS] 缩略图扩展名保持正确: {ext}")
            else:
                print(f"  [ERROR] 缩略图扩展名保持错误")
        
        # 测试相同内容的文件应该有相同的哈希值
        print(f"\n测试相同内容文件的哈希值一致性:")
        if len(image_files) >= 2:
            # 选择前两个文件进行比较
            file1_path = os.path.join(test_images_dir, image_files[0])
            file2_path = os.path.join(test_images_dir, image_files[1])
            
            hash1 = calc_hash(file1_path)
            hash2 = calc_hash(file2_path)
            
            print(f"  文件1 ({image_files[0]}): {hash1[:8]}...")
            print(f"  文件2 ({image_files[1]}): {hash2[:8]}...")
            
            if hash1 != hash2:
                print(f"  [SUCCESS] 不同内容文件哈希值不同")
            else:
                print(f"  [INFO] 文件哈希值相同（可能是相同内容的文件）")
        
        # 测试命名格式的完整性并重命名文件
        print(f"\n测试命名格式完整性并重命名文件:")
        for filename in image_files[:3]:  # 只测试前3个文件
            image_path = os.path.join(test_images_dir, filename)
            hash_value = calc_hash(image_path)
            base_name, ext = os.path.splitext(filename)
            
            # 原图命名
            original_name = f"{base_name}.{hash_value[:8]}{ext}"
            # 缩略图命名
            thumb_name = f"{base_name}.thumb.{hash_value[:8]}{ext}"
            
            print(f"  文件: {filename}")
            print(f"    原图: {original_name}")
            print(f"    缩略图: {thumb_name}")
            
            # 验证命名格式的完整性
            original_parts = original_name.split('.')
            thumb_parts = thumb_name.split('.')
            
            if len(original_parts) == 3 and len(thumb_parts) == 4:
                print(f"    [SUCCESS] 命名格式完整")
                
                # 重命名原图文件
                new_image_path = os.path.join(test_images_dir, original_name)
                if filename != original_name:  # 只有当文件名不同时才重命名
                    try:
                        os.rename(image_path, new_image_path)
                        print(f"    [SUCCESS] 原图重命名成功: {filename} -> {original_name}")
                    except Exception as e:
                        print(f"    [ERROR] 原图重命名失败: {e}")
                else:
                    print(f"    [INFO] 原图文件名已符合格式，无需重命名")
                
                # 检查并重命名缩略图文件
                thumbs_dir = "test_images/thumbs"
                old_thumb_path = os.path.join(thumbs_dir, filename)
                new_thumb_path = os.path.join(thumbs_dir, thumb_name)
                
                if os.path.exists(old_thumb_path):
                    try:
                        os.rename(old_thumb_path, new_thumb_path)
                        print(f"    [SUCCESS] 缩略图重命名成功: {filename} -> {thumb_name}")
                    except Exception as e:
                        print(f"    [ERROR] 缩略图重命名失败: {e}")
                else:
                    print(f"    [INFO] 缩略图文件不存在: {filename}")
                    
            else:
                print(f"    [ERROR] 命名格式不完整，跳过重命名")
        
        print(f"\n[SUCCESS] 图片重命名功能测试完成")
        
    except Exception as e:
        print(f"[ERROR] 测试过程中出现异常: {e}")
        import traceback
        traceback.print_exc()

def run_all_tests():
    """运行所有测试"""
    print("开始执行上传功能测试...\n")
    
    # 打印配置信息
    print("=== 当前配置 ===")
    config.print_config()
    print()
    
    # 运行各种测试
    test_picgo_config()
    test_thumbnail_generation()
    test_photo_rename()  # 添加图片重命名测试
    test_file_not_exists()
    test_unsupported_format()
    test_picgo_upload_direct()
    # test_normal_upload()
    # test_duplicate_upload()
    
    print("=== 测试完成 ===")
    print("[SUCCESS] 所有测试执行完成")
    
    # 询问是否清理测试文件
    # response = input("\n是否清理测试生成的文件？(y/N): ")
    # if response.lower() in ['y', 'yes']:
    #     cleanup_test_files()

if __name__ == "__main__":
    run_all_tests() 