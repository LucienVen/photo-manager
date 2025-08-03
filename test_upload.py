#!/usr/bin/env python3
"""
上传功能测试脚本
测试图片上传的各种场景
"""
import os
import tempfile
import shutil
from PIL import Image, ImageDraw
from upload import main

def create_test_image(filename, size=(800, 600), color='blue'):
    """创建测试图片"""
    img = Image.new('RGB', size, color=color)
    draw = ImageDraw.Draw(img)
    draw.text((size[0]//2, size[1]//2), 'Test Image', fill='white')
    img.save(filename)
    return filename

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
    """测试正常上传"""
    print("=== 测试正常上传 ===")
    
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
    """测试重复上传"""
    print("=== 测试重复上传 ===")
    
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
    
    print()

def run_all_tests():
    """运行所有测试"""
    print("开始执行上传功能测试...\n")
    
    # 运行各种测试
    test_file_not_exists()
    test_unsupported_format()
    test_normal_upload()
    test_duplicate_upload()
    
    print("=== 测试完成 ===")
    print("[SUCCESS] 所有测试执行完成")
    
    # 询问是否清理测试文件
    response = input("\n是否清理测试生成的文件？(y/N): ")
    if response.lower() in ['y', 'yes']:
        cleanup_test_files()

if __name__ == "__main__":
    run_all_tests() 