#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
FRONTEND_DIR="${ROOT_DIR}/mistake-book-view"

echo "=================================================="
echo "==> 正在编译前端工程: mistake-book-view"
echo "=================================================="

if [ ! -d "${FRONTEND_DIR}" ]; then
  echo "错误: 前端目录不存在: ${FRONTEND_DIR}" >&2
  exit 1
fi

cd "${FRONTEND_DIR}"

if [ ! -d "node_modules" ]; then
  echo "==> 检测到依赖尚未安装，正在执行 npm install..."
  npm install
fi

echo "==> 开始执行 Vite 生产构建..."
npm run build

echo "==> 前端构建完成！产物路径: ${FRONTEND_DIR}/dist"
