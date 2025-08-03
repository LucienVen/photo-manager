#!/usr/bin/env python3
"""
配置管理模块
用于加载和管理环境变量配置
"""
import os
from dotenv import load_dotenv

# 加载.env文件
load_dotenv()

class Config:
    """配置类，用于管理所有环境变量"""
    
    # GitHub 配置
    GITHUB_NAME = os.getenv('GITHUB_NAME', 'yourname')
    GITHUB_REPO = os.getenv('GITHUB_REPO', 'photo-bed')
    
    # CDN 配置
    CDN_PREFIX = os.getenv('CDN_PREFIX', 'https://cdn.jsdelivr.net/gh')
    
    # 缩略图配置
    THUMB_WIDTH = int(os.getenv('THUMB_WIDTH', '320'))
    
    # 调试模式
    DEBUG = os.getenv('DEBUG', 'false').lower() == 'true'
    
    @classmethod
    def get_cdn_url(cls):
        """获取完整的CDN URL前缀"""
        return f"{cls.CDN_PREFIX}/{cls.GITHUB_NAME}/{cls.GITHUB_REPO}"
    
    @classmethod
    def print_config(cls):
        """打印当前配置信息"""
        print("当前配置:")
        print(f"  GitHub用户名: {cls.GITHUB_NAME}")
        print(f"  GitHub仓库: {cls.GITHUB_REPO}")
        print(f"  CDN前缀: {cls.CDN_PREFIX}")
        print(f"  缩略图宽度: {cls.THUMB_WIDTH}")
        print(f"  调试模式: {cls.DEBUG}")
        print(f"  完整CDN URL: {cls.get_cdn_url()}")

# 创建全局配置实例
config = Config() 