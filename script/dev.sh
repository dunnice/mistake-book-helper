#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

BACKEND_DIR="${ROOT_DIR}/mistake-book-go-api"
FRONTEND_DIR="${ROOT_DIR}/mistake-book-view"

cleanup() {
  echo ""
  echo "==> 正在停止开发环境..."
  kill $(jobs -p) 2>/dev/null || true
  exit 0
}
trap cleanup SIGINT SIGTERM EXIT

echo "=================================================="
echo "==> 启动本地全栈协同开发环境"
echo "==> 后端 API: http://127.0.0.1:8080"
echo "==> 前端界面: http://127.0.0.1:5173"
echo "=================================================="

# 启动后端
cd "${BACKEND_DIR}"
PORT=8080 go run ./cmd/server &
BACKEND_PID=$!

# 等待后端稍微就绪
sleep 1

# 启动前端
cd "${FRONTEND_DIR}"
npm run dev &
FRONTEND_PID=$!

wait
