#!/bin/bash

echo "🚀 启动AI论坛..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go，请先安装Go 1.21+"
    exit 1
fi

# 检查MySQL是否运行
if ! command -v mysql &> /dev/null; then
    echo "⚠️  警告: 未找到MySQL客户端，请确保MySQL服务正在运行"
fi

# 安装依赖
echo "📦 安装Go依赖..."
go mod tidy

# 检查环境配置文件
if [ ! -f "config.env" ]; then
    echo "❌ 错误: 未找到config.env配置文件"
    echo "请创建config.env文件并配置数据库连接信息"
    exit 1
fi

# 加载环境变量
echo "⚙️  加载环境配置..."
export $(cat config.env | xargs)

# 检查数据库连接
echo "🔍 检查数据库连接..."
if command -v mysql &> /dev/null; then
    mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD -e "USE $DB_NAME;" 2>/dev/null
    if [ $? -ne 0 ]; then
        echo "❌ 错误: 无法连接到数据库"
        echo "请检查config.env中的数据库配置"
        exit 1
    fi
    echo "✅ 数据库连接成功"
fi

# 启动服务器
echo "🌐 启动Web服务器..."
echo "访问地址: http://localhost:$SERVER_PORT"
echo "按 Ctrl+C 停止服务器"
echo ""

go run main.go 