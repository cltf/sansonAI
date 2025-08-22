#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
个人中心功能测试脚本
"""

import requests
import json
import time

# 测试配置
BASE_URL = "http://localhost:8080"
TEST_USER = {
    "username": "testuser",
    "email": "test@example.com",
    "password": "testpass123"
}

class ProfileTester:
    def __init__(self):
        self.session = requests.Session()
        self.user_id = None
        self.token = None
    
    def login(self):
        """登录测试用户"""
        print("正在登录...")
        login_data = {
            "username": TEST_USER["username"],
            "password": TEST_USER["password"]
        }
        
        response = self.session.post(f"{BASE_URL}/auth/login", data=login_data)
        if response.status_code == 200:
            print("✓ 登录成功")
            return True
        else:
            print(f"✗ 登录失败: {response.status_code}")
            return False
    
    def test_profile_page(self):
        """测试个人中心页面访问"""
        print("\n测试个人中心页面...")
        response = self.session.get(f"{BASE_URL}/profile")
        if response.status_code == 200:
            print("✓ 个人中心页面访问成功")
            return True
        else:
            print(f"✗ 个人中心页面访问失败: {response.status_code}")
            return False
    
    def test_user_activity(self):
        """测试用户动态API"""
        print("\n测试用户动态API...")
        response = self.session.get(f"{BASE_URL}/api/user/activity")
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                print("✓ 用户动态API调用成功")
                print(f"  返回 {len(data.get('activities', []))} 条动态")
                return True
            else:
                print(f"✗ 用户动态API返回错误: {data.get('error')}")
                return False
        else:
            print(f"✗ 用户动态API调用失败: {response.status_code}")
            return False
    
    def test_user_questions(self):
        """测试用户提问API"""
        print("\n测试用户提问API...")
        response = self.session.get(f"{BASE_URL}/api/user/questions")
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                print("✓ 用户提问API调用成功")
                print(f"  返回 {len(data.get('questions', []))} 个提问")
                return True
            else:
                print(f"✗ 用户提问API返回错误: {data.get('error')}")
                return False
        else:
            print(f"✗ 用户提问API调用失败: {response.status_code}")
            return False
    
    def test_user_answers(self):
        """测试用户回答API"""
        print("\n测试用户回答API...")
        response = self.session.get(f"{BASE_URL}/api/user/answers")
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                print("✓ 用户回答API调用成功")
                print(f"  返回 {len(data.get('answers', []))} 个回答")
                return True
            else:
                print(f"✗ 用户回答API返回错误: {data.get('error')}")
                return False
        else:
            print(f"✗ 用户回答API调用失败: {response.status_code}")
            return False
    
    def test_user_favorites(self):
        """测试用户收藏API"""
        print("\n测试用户收藏API...")
        response = self.session.get(f"{BASE_URL}/api/user/favorites")
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                print("✓ 用户收藏API调用成功")
                print(f"  返回 {len(data.get('favorites', []))} 个收藏")
                return True
            else:
                print(f"✗ 用户收藏API返回错误: {data.get('error')}")
                return False
        else:
            print(f"✗ 用户收藏API调用失败: {response.status_code}")
            return False
    
    def test_user_messages(self):
        """测试用户消息API"""
        print("\n测试用户消息API...")
        response = self.session.get(f"{BASE_URL}/api/user/messages")
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                print("✓ 用户消息API调用成功")
                print(f"  返回 {len(data.get('messages', []))} 条消息")
                return True
            else:
                print(f"✗ 用户消息API返回错误: {data.get('error')}")
                return False
        else:
            print(f"✗ 用户消息API调用失败: {response.status_code}")
            return False
    
    def test_update_profile(self):
        """测试更新个人资料API"""
        print("\n测试更新个人资料API...")
        profile_data = {
            "username": TEST_USER["username"],
            "email": TEST_USER["email"],
            "bio": "这是一个测试用户简介",
            "phone": "13800138000",
            "website": "https://example.com",
            "profile_public": True,
            "show_email": False,
            "show_phone": False
        }
        
        response = self.session.put(
            f"{BASE_URL}/api/user/profile",
            json=profile_data,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                print("✓ 更新个人资料API调用成功")
                return True
            else:
                print(f"✗ 更新个人资料API返回错误: {data.get('error')}")
                return False
        else:
            print(f"✗ 更新个人资料API调用失败: {response.status_code}")
            return False
    
    def test_change_password(self):
        """测试修改密码API"""
        print("\n测试修改密码API...")
        password_data = {
            "current_password": TEST_USER["password"],
            "new_password": "newpass123"
        }
        
        response = self.session.put(
            f"{BASE_URL}/api/user/password",
            json=password_data,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                print("✓ 修改密码API调用成功")
                # 改回原密码
                password_data = {
                    "current_password": "newpass123",
                    "new_password": TEST_USER["password"]
                }
                self.session.put(
                    f"{BASE_URL}/api/user/password",
                    json=password_data,
                    headers={"Content-Type": "application/json"}
                )
                return True
            else:
                print(f"✗ 修改密码API返回错误: {data.get('error')}")
                return False
        else:
            print(f"✗ 修改密码API调用失败: {response.status_code}")
            return False
    
    def test_notification_settings(self):
        """测试通知设置API"""
        print("\n测试通知设置API...")
        settings_data = {
            "email_notifications": True,
            "browser_notifications": True,
            "question_notifications": True,
            "follow_notifications": True
        }
        
        response = self.session.put(
            f"{BASE_URL}/api/user/notifications",
            json=settings_data,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                print("✓ 通知设置API调用成功")
                return True
            else:
                print(f"✗ 通知设置API返回错误: {data.get('error')}")
                return False
        else:
            print(f"✗ 通知设置API调用失败: {response.status_code}")
            return False
    
    def run_all_tests(self):
        """运行所有测试"""
        print("开始个人中心功能测试...")
        print("=" * 50)
        
        tests = [
            self.test_profile_page,
            self.test_user_activity,
            self.test_user_questions,
            self.test_user_answers,
            self.test_user_favorites,
            self.test_user_messages,
            self.test_update_profile,
            self.test_change_password,
            self.test_notification_settings,
        ]
        
        passed = 0
        total = len(tests)
        
        for test in tests:
            try:
                if test():
                    passed += 1
            except Exception as e:
                print(f"✗ 测试异常: {e}")
            time.sleep(0.5)  # 避免请求过快
        
        print("\n" + "=" * 50)
        print(f"测试完成: {passed}/{total} 通过")
        
        if passed == total:
            print("🎉 所有测试通过！个人中心功能正常。")
        else:
            print("⚠️  部分测试失败，请检查相关功能。")

def main():
    """主函数"""
    tester = ProfileTester()
    
    # 先尝试登录
    if not tester.login():
        print("请确保测试用户已注册且服务器正在运行")
        return
    
    # 运行所有测试
    tester.run_all_tests()

if __name__ == "__main__":
    main() 