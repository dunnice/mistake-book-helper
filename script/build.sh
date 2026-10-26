#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

TARGET="${1:-}"

echo "=================================================="
echo "==> [1/2] 编译前端 mistake-book-view"
echo "=================================================="
"${SCRIPT_DIR}/build-frontend.sh"

echo ""
echo "=================================================="
echo "==> [2/2] 编译后端 mistake-book-go-api"
echo "=================================================="
"${SCRIPT_DIR}/build-backend.sh" "${TARGET}"

echo ""
echo "=================================================="
echo "==> 恭喜！全栈编译完成。"
echo "==> 构建输出位于: ${ROOT_DIR}/build/"
echo "==> 前端产物位于: ${ROOT_DIR}/mistake-book-view/dist/"
echo "=================================================="
