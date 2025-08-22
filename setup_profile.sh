#!/bin/bash

# 个人中心功能设置脚本

echo "🚀 开始设置个人中心功能..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go"
    exit 1
fi

echo "✅ Go环境检查通过"

# 检查数据库文件
if [ ! -f "aiforum.db" ]; then
    echo "📝 创建数据库文件..."
    sqlite3 aiforum.db < init_db.sql
    echo "✅ 数据库初始化完成"
else
    echo "✅ 数据库文件已存在"
fi

# 添加个人中心相关表结构
echo "📝 添加个人中心数据库表..."
sqlite3 aiforum.db < profile_tables.sql
echo "✅ 个人中心表结构添加完成"

# 检查必要的目录
echo "📁 检查必要目录..."
mkdir -p images/avatars
mkdir -p static
echo "✅ 目录检查完成"

# 检查模板文件
if [ ! -f "templates/profile.html" ]; then
    echo "❌ 错误: 未找到个人中心模板文件 templates/profile.html"
    exit 1
fi

# 检查样式文件
if [ ! -f "static/profile.js" ]; then
    echo "❌ 错误: 未找到个人中心JavaScript文件 static/profile.js"
    exit 1
fi

echo "✅ 文件检查完成"

# 检查依赖
echo "📦 检查Go依赖..."
go mod tidy
echo "✅ 依赖检查完成"

# 编译项目
echo "🔨 编译项目..."
go build -o aiforum main.go
if [ $? -eq 0 ]; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败"
    exit 1
fi

echo ""
echo "🎉 个人中心功能设置完成！"
echo ""
echo "📋 功能说明："
echo "   • 两栏布局设计，左侧导航，右侧内容"
echo "   • 支持头像上传、个人资料编辑"
echo "   • 包含11个功能模块：动态、提问、回答、分享等"
echo "   • 支持关注、收藏、消息等社交功能"
echo "   • 响应式设计，支持移动端"
echo ""
echo "🚀 启动服务器："
echo "   ./aiforum"
echo ""
echo "🌐 访问地址："
echo "   http://localhost:8080/profile"
echo ""
echo "📖 详细文档："
echo "   README_PROFILE.md"
echo ""
echo "🧪 运行测试："
echo "   python3 test_profile.py"
echo "" 