#!/bin/bash

echo "========================================"
echo "   启动 Demo 用户登录系统"
echo "========================================"

echo ""
echo "[1/3] 启动后端服务 (端口 8080)..."
cd "$(dirname "$0")/demo-server"
go run main.go &
SERVER_PID=$!

sleep 2

echo "[2/3] 启动前端服务 (端口 3000)..."
cd "$(dirname "$0")/demo-web"
npm run dev &
FRONTEND_PID=$!

echo "[3/3] 完成!"
echo ""
echo "服务地址:"
echo "  - 前端: http://localhost:3000"
echo "  - 后端: http://localhost:8080"
echo ""
echo "按 Ctrl+C 停止所有服务"

trap "kill $SERVER_PID $FRONTEND_PID 2>/dev/null" EXIT

wait
