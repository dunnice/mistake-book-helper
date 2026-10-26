#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BACKEND_DIR="${ROOT_DIR}/mistake-book-go-api"
OUT_DIR="${ROOT_DIR}/build"

TARGET_OS="${1:-}"

mkdir -p "${OUT_DIR}"

echo "=================================================="
echo "==> 正在编译后端工程: mistake-book-go-api"
echo "=================================================="

if [ ! -d "${BACKEND_DIR}" ]; then
  echo "错误: 后端目录不存在: ${BACKEND_DIR}" >&2
  exit 1
fi

cd "${BACKEND_DIR}"

if [ "${TARGET_OS}" = "--linux" ] || [ "${TARGET_OS}" = "linux" ]; then
  OUTPUT_BIN="${OUT_DIR}/mistakebook-linux-amd64"
  echo "==> 目标系统: Linux (amd64)"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "${OUTPUT_BIN}" ./cmd/server
  echo "==> Linux 产物编译完成: ${OUTPUT_BIN}"
else
  OUTPUT_BIN="${OUT_DIR}/mistakebook-api"
  echo "==> 目标系统: 本地宿主系统 ($(uname -s)/$(uname -m))"
  CGO_ENABLED=0 go build -ldflags="-s -w" -o "${OUTPUT_BIN}" ./cmd/server
  echo "==> 本地产物编译完成: ${OUTPUT_BIN}"
fi
