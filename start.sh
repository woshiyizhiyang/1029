#!/bin/bash

# 动画形象AI对话系统 - 快速启动脚本

set -e

echo "=================================="
echo "动画形象AI对话系统 - 快速启动"
echo "=================================="

# 检查依赖
check_dependencies() {
    echo "检查依赖..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo "错误: 未找到Go,请先安装Go 1.21+"
        exit 1
    fi
    echo "✓ Go已安装: $(go version)"
    
    # 检查MySQL
    if ! command -v mysql &> /dev/null; then
        echo "警告: 未找到MySQL客户端"
    else
        echo "✓ MySQL客户端已安装"
    fi
    
    echo ""
}

# 配置检查
check_config() {
    echo "检查配置文件..."
    
    if [ ! -f "config/app.json" ]; then
        echo "错误: 配置文件不存在"
        echo "请复制 config/app.example.json 为 config/app.json 并填入真实配置"
        exit 1
    fi
    echo "✓ 配置文件存在"
    echo ""
}

# 数据库检查
check_database() {
    echo "检查数据库..."
    echo "请确保已执行 database/init.sql 创建数据库"
    echo "按 Enter 继续,或 Ctrl+C 退出..."
    read
}

# 启动后端
start_backend() {
    echo "=================================="
    echo "启动后端服务..."
    echo "=================================="
    
    cd backend
    
    # 安装依赖
    echo "安装Go依赖..."
    go mod tidy
    
    # 创建日志目录
    mkdir -p logs
    
    # 启动服务
    echo "启动后端服务 (端口: 8080)..."
    go run cmd/main.go -config ../config/app.json -log logs/app.log &
    BACKEND_PID=$!
    
    echo "后端服务已启动 (PID: $BACKEND_PID)"
    echo "WebSocket地址: ws://localhost:8080/ws"
    echo "健康检查: http://localhost:8080/health"
    echo ""
    
    cd ..
    
    # 等待服务启动
    sleep 3
}

# 主函数
main() {
    check_dependencies
    check_config
    check_database
    start_backend
    
    echo "=================================="
    echo "✓ 后端服务启动成功!"
    echo "=================================="
    echo ""
    echo "前端启动步骤:"
    echo "1. cd frontend"
    echo "2. npm install"
    echo "3. 将5张口型图片放到 public/images/ 目录"
    echo "   - mouth-0.jpg (闭口)"
    echo "   - mouth-1.jpg (微张)"
    echo "   - mouth-2.jpg (半张)"
    echo "   - mouth-3.jpg (大张)"
    echo "   - mouth-4.jpg (完全张口)"
    echo "4. npm run dev"
    echo "5. 访问 http://localhost:3000"
    echo ""
    echo "按 Ctrl+C 停止后端服务"
    
    # 等待中断
    wait $BACKEND_PID
}

main
