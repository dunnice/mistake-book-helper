.PHONY: help dev dev-api dev-web build build-web build-api build-linux deploy clean

help:
	@echo "错题本助手 (MistakeBook) 工程化指令:"
	@echo "  make dev          - 本地一键启动前后端协同开发"
	@echo "  make dev-api      - 本地启动 Go 后端 API (:8080，使用 SQLite)"
	@echo "  make dev-web      - 本地启动 Vue 3 前端开发服务器 (:5173，带热更)"
	@echo "  make build        - 全栈编译 (前端 Vue 产物 + 本地 Go 二进制)"
	@echo "  make build-linux  - 交叉编译 Linux amd64 二进制"
	@echo "  make deploy       - 一键部署到阿里云服务器 (编译 + 上传 + systemd 重启 + 探活)"
	@echo "  make clean        - 清理本地 build 产物"

dev:
	@./script/dev.sh

dev-api:
	@cd mistake-book-go-api && go run ./cmd/server

dev-web:
	@cd mistake-book-view && npm run dev

build:
	@./script/build.sh

build-web:
	@./script/build-frontend.sh

build-api:
	@./script/build-backend.sh

build-linux:
	@./script/build-backend.sh --linux

deploy:
	@./script/deploy-aliyun.sh

clean:
	@rm -rf build mistake-book-view/dist
	@echo "清理完成"
