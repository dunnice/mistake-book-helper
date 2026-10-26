#!/usr/bin/env bash
# 错题本助手 · 阿里云原生一键部署脚本
# 目标机: 由 .env 中的 DEPLOY_SERVER_HOST 指定 (SSH 免密登录)
# 服务目录: /opt/mistakebook
# systemd: mistakebook (监听 127.0.0.1:8090, BASE_PATH=/cuoti)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

# 读取本地 .env 环境变量（如果存在）
if [ -f "${ROOT_DIR}/.env" ]; then
  set -a
  source "${ROOT_DIR}/.env"
  set +a
fi

TARGET_HOST="${DEPLOY_SERVER_HOST:-${1:-}}"
TARGET_USER="${DEPLOY_SERVER_USER:-root}"
REMOTE_DIR="${DEPLOY_SERVER_PATH:-/opt/mistakebook}"

if [ -z "${TARGET_HOST}" ]; then
  echo "❌ 错误: 未检测到目标服务器地址！"
  echo "请在根目录 .env 文件中配置 DEPLOY_SERVER_HOST=your-server-ip，或作为参数传入：./script/deploy-aliyun.sh <server_ip>"
  exit 1
fi

HOST="${TARGET_USER}@${TARGET_HOST}"
SERVICE_NAME="mistakebook"
FRONT_DIR="${ROOT_DIR}/mistake-book-view"
BACK_DIR="${ROOT_DIR}/mistake-book-go-api"
BUILD_DIR="${ROOT_DIR}/build"

echo "=================================================="
echo "==> 1. 本机交叉编译 Go 后端 (Linux amd64)"
echo "=================================================="
mkdir -p "${BUILD_DIR}"
cd "${BACK_DIR}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "${BUILD_DIR}/mistakebook-linux-amd64" ./cmd/server
echo "    后端二进制大小: $(ls -lh "${BUILD_DIR}/mistakebook-linux-amd64" | awk '{print $5}')"

echo ""
echo "=================================================="
echo "==> 2. 构建前端工程 (子路径: /cuoti/)"
echo "=================================================="
cd "${FRONT_DIR}"
if [ ! -d "node_modules" ]; then
  npm install --no-fund --no-audit
fi
VITE_BASE_PATH="/cuoti/" npm run build

echo ""
echo "=================================================="
echo "==> 3. 上传产物至阿里云 (${HOST}:${REMOTE_DIR})"
echo "=================================================="
ssh "${HOST}" "mkdir -p ${REMOTE_DIR}/dist ${REMOTE_DIR}/data/images"

# 上传后端二进制
echo "--> 上传 Linux 二进制..."
scp -q "${BUILD_DIR}/mistakebook-linux-amd64" "${HOST}:${REMOTE_DIR}/mistakebook-linux-amd64.new"
ssh "${HOST}" "chmod +x ${REMOTE_DIR}/mistakebook-linux-amd64.new && mv -f ${REMOTE_DIR}/mistakebook-linux-amd64.new ${REMOTE_DIR}/mistakebook-linux-amd64"

# 同步前端静态资源 dist
echo "--> 同步前端 dist 静态文件..."
tar -C "${FRONT_DIR}/dist" \
  --exclude '._*' \
  --exclude '.DS_Store' \
  -czf - . \
  | ssh "${HOST}" "tar -xzf - -C ${REMOTE_DIR}/dist && find ${REMOTE_DIR}/dist -name '._*' -delete"

echo ""
echo "=================================================="
echo "==> 4. 重启远端 systemd 服务并探活"
echo "=================================================="
ssh "${HOST}" bash -s <<REMOTE
set -euo pipefail
echo "--> 重启 ${SERVICE_NAME} 服务..."
systemctl restart "${SERVICE_NAME}"
sleep 2

if systemctl is-active --quiet "${SERVICE_NAME}"; then
  echo "    服务启动正常: active (running)"
else
  echo "错误: 服务启动失败，查看日志:" >&2
  journalctl -u "${SERVICE_NAME}" -n 30 --no-pager
  exit 1
fi

echo "--> 本机 8090 内网探活:"
curl -sf http://127.0.0.1:8090/cuoti/healthz && echo " -> ok"
REMOTE

echo ""
echo "--> 公网 Nginx (端口 5173 /cuoti/) 冒烟验证:"
PUBLIC_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://${DEPLOY_SERVER_HOST}:5173/cuoti/healthz" || true)
echo "    http://${DEPLOY_SERVER_HOST}:5173/cuoti/healthz -> HTTP ${PUBLIC_CODE}"

if [ "${PUBLIC_CODE}" = "200" ]; then
  echo ""
  echo "=================================================="
  echo "🎉 部署大功告成！"
  echo "错题本工作台访问地址: http://${DEPLOY_SERVER_HOST}:5173/cuoti/"
  echo "=================================================="
else
  echo "警告: 公网探活返回状态码 ${PUBLIC_CODE}，请检查远端 Nginx 配置或防火墙"
fi
