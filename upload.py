#!/usr/bin/env python3
import os
import sys
import shutil
import hashlib
import json
import subprocess
from datetime import datetime
from PIL import Image

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
    if os.path.exists(RECORD_FILE):
        with open(RECORD_FILE, "r") as f:
            return json.load(f)
    return []

def save_records(data):
    with open(RECORD_FILE, "w") as f:
        json.dump(data, f, indent=2)

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

def main(img_path):
    ensure_dirs()

    if not os.path.exists(img_path):
        print(f"[ERROR] 文件不存在: {img_path}")
        return

    ext = os.path.splitext(img_path)[1].lower()
    if ext not in [".jpg", ".jpeg", ".png", ".webp"]:
        print(f"[WARNING] 不支持的图片格式: {ext}")
        return

    records = load_records()
    h = calc_hash(img_path)
    if any(r["hash"] == h for r in records):
        print(f"[SKIP] 已存在相同图片（hash: {h[:8]}...），跳过上传")
        return

    base_filename = os.path.splitext(os.path.basename(img_path))[0]
    new_filename = f"{base_filename}.{h[:8]}{ext}"
    img_dst = os.path.join(IMAGES_DIR, new_filename)
    thumb_dst = os.path.join(THUMBS_DIR, new_filename)

    shutil.copyfile(img_path, img_dst)
    generate_thumbnail(img_path, thumb_dst)

    width, height, size_kb = get_image_info(img_dst)
    now = int(datetime.now().timestamp())

    record = {
        "filename": new_filename,
        "url": f"{CDN_PREFIX}/images/{new_filename}",
        "thumb_url": f"{CDN_PREFIX}/thumbs/{new_filename}",
        "path": f"images/{new_filename}",
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
    print(f"[SUCCESS] 成功上传: {new_filename}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("用法: ./upload.py path/to/image.jpg")
    else:
        main(sys.argv[1])
