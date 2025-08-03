#!/usr/bin/env python3
"""
测试配置加载
验证.env文件中的配置是否正确加载
"""
from config import config

def test_config():
    """测试配置加载"""
    print("=== 配置测试 ===")
    
    # 打印当前配置
    config.print_config()
    
    print("\n=== 配置验证 ===")
    
    # 验证关键配置
    assert config.GITHUB_NAME != 'yourname', "请设置正确的GITHUB_NAME"
    assert config.GITHUB_REPO != 'photo-bed', "请设置正确的GITHUB_REPO"
    
    print("[SUCCESS] 配置验证通过")
    print(f"[SUCCESS] GitHub用户名: {config.GITHUB_NAME}")
    print(f"[SUCCESS] GitHub仓库: {config.GITHUB_REPO}")
    print(f"[SUCCESS] CDN URL: {config.get_cdn_url()}")

if __name__ == "__main__":
    test_config() 