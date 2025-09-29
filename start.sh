#!/bin/bash

# Script để khởi động server và tự động xử lý port conflict
PORT=8080

echo "🚀 Starting Dish API Server..."

# Kiểm tra và dừng tiến trình đang sử dụng port 8080
if lsof -ti:$PORT > /dev/null 2>&1; then
    echo "⚠️  Port $PORT đang được sử dụng. Đang dừng tiến trình cũ..."
    lsof -ti:$PORT | xargs kill -9
    sleep 1
    echo "✅ Đã dừng tiến trình cũ"
fi

# Khởi động server
echo "🔥 Khởi động server trên port $PORT..."
go run cmd/main.go
