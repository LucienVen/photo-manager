#!/usr/bin/env python3
import os
import sys
import shutil
import hashlib
import json
import subprocess
from datetime import datetime
from PIL import Image
import glob
from collections import defaultdict

# 导入配置模块
from config import config

# 配置路径
HOME = os.path.expanduser("~")
BASE_DIR = os.path.abspath(os.path.dirname(__file__))
PHOTO_BED_DIR = os.path.join(BASE_DIR, "photo-bed")
IMAGES_DIR = os.path.join(PHOTO_BED_DIR, "images")
THUMBS_DIR = os.path.join(PHOTO_BED_DIR, "thumbs")
RECORD_FILE = os.path.join(BASE_DIR, "photo-manager", "images.json")

# 从配置文件获取CDN前缀和缩略图宽度
CDN_PREFIX = config.get_cdn_url()
THUMB_WIDTH = config.THUMB_WIDTH
RECORD_DIR = config.RECORD_DIR

def calc_hash(filepath, algo="sha256"):
    h = hashlib.new(algo)
    with open(filepath, "rb") as f:
        while chunk := f.read(8192):
            h.update(chunk)
    return h.hexdigest()

def ensure_dirs():
    os.makedirs(IMAGES_DIR, exist_ok=True)
    os.makedirs(THUMBS_DIR, exist_ok=True)
    os.makedirs(os.path.dirname(RECORD_FILE), exist_ok=True)

def load_records():
    records = []
    record_dir = RECORD_DIR
    if not os.path.exists(record_dir):
        return records
    for file in glob.glob(os.path.join(record_dir, "*.json")):
        with open(file, "r") as f:
            try:
                data = json.load(f)
                if isinstance(data, list):
                    records.extend(data)
            except Exception as e:
                print(f"[WARNING] 读取 {file} 失败: {e}")
    return records

def save_records(data):
    record_dir = RECORD_DIR
    os.makedirs(record_dir, exist_ok=True)
    groups = defaultdict(list)
    for record in data:
        ts = record.get("uploaded_at")
        if not ts:
            print(f"[WARNING] 记录缺少 uploaded_at 字段: {record}")
            continue
        dt = datetime.fromtimestamp(ts)
        key = f"{dt.year:04d}-{dt.month:02d}"
        groups[key].append(record)
    for ym, records in groups.items():
        file_path = os.path.join(record_dir, f"{ym}.json")
        with open(file_path, "w") as f:
            json.dump(records, f, indent=2, ensure_ascii=False)

def get_image_info(filepath):
    with Image.open(filepath) as img:
        width, height = img.size
    size_kb = os.path.getsize(filepath) / 1024
    return width, height, round(size_kb, 2)

def generate_thumbnail(src_path, dst_path):
    """使用PIL生成缩略图"""
    try:
        with Image.open(src_path) as img:
            # 计算等比例缩放后的尺寸
            width, height = img.size
            ratio = width / height
            new_width = THUMB_WIDTH
            new_height = int(new_width / ratio)
            
            # 生成缩略图
            thumb = img.resize((new_width, new_height), Image.Resampling.LANCZOS)
            thumb.save(dst_path, quality=85)
    except Exception as e:
        print(f"[WARNING] 缩略图生成失败: {e}")
        # 如果缩略图生成失败，复制原图作为缩略图
        shutil.copyfile(src_path, dst_path)

def upload_with_picgo(img_path):
    """使用 PicGo 上传图片并返回 URL"""
    try:
        # 调用 PicGo CLI 上传图片
        result = subprocess.run(
            ["picgo", "upload", img_path],
            capture_output=True, text=True, check=True
        )
        output = result.stdout.strip()
        
        # PicGo 返回的通常是图片 URL，一行一个
        if output:
            # 取第一个 URL（如果有多行）
            url = output.split('\n')[0].strip()
            if url.startswith('http'):
                return url
        
        print(f"[ERROR] PicGo 返回格式异常: {output}")
        return None
    except subprocess.CalledProcessError as e:
        print(f"[ERROR] PicGo 上传失败: {e.stderr}")
        return None
    except FileNotFoundError:
        print("[ERROR] 未找到 PicGo，请确保已安装并配置 PicGo CLI")
        return None
    except Exception as e:
        print(f"[ERROR] PicGo 调用异常: {e}")
        return None

def main(img_path):
    if not os.path.exists(img_path):
        print(f"[ERROR] 文件不存在: {img_path}")
        return

    ext = os.path.splitext(img_path)[1].lower()
    if ext not in [".jpg", ".jpeg", ".png", ".webp"]:
        print(f"[WARNING] 不支持的图片格式: {ext}")
        return

    # 检查重复图片
    records = load_records()
    h = calc_hash(img_path)
    if any(r["hash"] == h for r in records):
        print(f"[SKIP] 已存在相同图片（hash: {h[:8]}...），跳过上传")
        return

    # 使用 PicGo 上传图片
    print(f"[INFO] 正在通过 PicGo 上传: {os.path.basename(img_path)}")
    url = upload_with_picgo(img_path)
    if not url:
        print("[ERROR] PicGo 上传失败，请检查 PicGo 配置")
        return

    # 获取图片信息
    width, height, size_kb = get_image_info(img_path)
    now = int(datetime.now().timestamp())
    base_filename = os.path.splitext(os.path.basename(img_path))[0]

    # 创建记录
    record = {
        "filename": f"{base_filename}.{h[:8]}{ext}",
        "url": url,
        "thumb_url": url,  # 暂时使用原图 URL，后续可通过 PicGo 插件生成缩略图
        "path": url,       # 使用完整 URL
        "uploaded_at": now,
        "tags": [],
        "note": "",
        "size_kb": size_kb,
        "width": width,
        "height": height,
        "hash": h
    }

    records.append(record)
    save_records(records)
    print(f"[SUCCESS] PicGo 上传成功: {url}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("用法: ./upload.py path/to/image.jpg")
    else:
        main(sys.argv[1])
