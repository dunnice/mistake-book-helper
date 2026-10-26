> **MistakeBook Helper** · 初中全科目智能错题收集与复习巩固系统  
> **出品人**：ice7  
> **核心机制**：支持各类 AI Agent 智能批改试卷与提取错题，本系统负责错题的**存储、归因剖析、纸质打印（带二维码）、手机扫码重做打卡与易错知识点看板**。

## 项目架构与工程结构

本项目经过系统化工程重构，采用清晰的前后端模块化分层，**零容器负担**，本地开箱即调，线上原生服务托管：

```
mistake-book-helper/
├── 00-doc/                               # 【文档目录】统一集中收纳项目所有文档
│   ├── 错题本-工程优化与系统架构.md       # 架构与重构设计说明
│   ├── 错题本-使用手册-WorkBuddy工作流.md  # WorkBuddy 拍照录题工作流
│   ├── 错题本-最小闭环实施方案.md         # 方案与实施细节
│   ├── 错题本-设计文档.md                # 原系统设计文档
│   ├── 错题本-部署与运维手册.md           # 生产部署与运维指南
│   └── 错题本-WorkBuddy指令卡.md          # AI 提示词与指令卡
├── script/                               # 【脚本目录】统一编译与开发自动化
│   ├── build.sh                          # 全栈本地编译
│   ├── build-frontend.sh                 # 编译前端 Vue 项目
│   ├── build-backend.sh                  # 编译 Go 后端（支持 --linux 参数）
│   ├── dev.sh                            # 本地一键启动前后端协同开发
│   └── deploy-aliyun.sh                  # 一键部署至阿里云（编译+上传+systemd重启+探活）
├── mistake-book-view/                    # 【前端工程】Vue 3 + TypeScript + Vite
│   ├── src/                              # 遵循 teamwork-front-view 设计规范 (Dock/毛玻璃/卡片)
│   └── dist/                             # 前端打包输出
├── mistake-book-go-api/                  # 【后端工程】Go + Gin + GORM
│   ├── cmd/server/main.go                # API 服务入口与静态托管
│   └── internal/                         # 模型与业务处理器
├── build/                                # 编译二进制与产物输出
└── Makefile                              # 便捷命令驱动入口
```

---

## 本地开发与调试 (直接 IDE 运行)

### 1. 后端 API 调试
直接在 IDE（Goland / VSCode / Cursor）中打开 `mistake-book-go-api/cmd/server/main.go` 点击运行，或通过命令行：
```bash
make dev-api
# 或 cd mistake-book-go-api && go run ./cmd/server
```
- 默认监听 `8080` 端口；
- 默认自动使用轻量级 SQLite (`data/mistakebook.db`)，无需配置任何外部数据库依赖。

### 2. 前端界面调试
在 `mistake-book-view/` 目录运行：
```bash
make dev-web
# 或 cd mistake-book-view && npm run dev
```
- 前端运行在 `http://127.0.0.1:5173/`，Vite 极速热更新，已配置自动代理 `/api` 请求到本地 `:8080` 后端。

### 3. 一键双端启动
```bash
make dev
```
一条命令同时拉起后端 Go API 与前端 Vite 开发服务器。

---

## 远端阿里云部署 (直接原生部署)

无需 Docker，通过自动化脚本直接将 Linux amd64 二进制与前端产物部署至阿里云：

```bash
make deploy
# 或 ./script/deploy-aliyun.sh
```

**脚本自动执行流程**：
1. 本机交叉编译 Linux amd64 二进制 `mistakebook-linux-amd64`；
2. 构建前端工程为适配反代子路径 `/cuoti/` 的静态产物；
3. 通过 SSH/SCP 将二进制与 `dist/` 流式同步至阿里云 `/opt/mistakebook/`；
4. 远端自动重启 systemd 守护进程 `systemctl restart mistakebook`；
5. 自动验证内网 8090 及公网 5173 接口探活响应。

---

## 环境变量配置

请参考 [.env.example](.env.example) 复制生成本地 `.env` 文件：
```bash
cp .env.example .env
```
根据本地开发或生产环境实际情况调整端口、`API_KEY`、数据库连接串或部署目标机地址。**`.env` 已被 `.gitignore` 严格忽略，切勿提交到版本库！**

---

## 参与贡献

欢迎提交 Issue 和 Pull Request！详细规范请查阅 [CONTRIBUTING.md](CONTRIBUTING.md)。
提交 PR 时请配合遵循 [PR 模板](.github/PULL_REQUEST_TEMPLATE.md) 填写变更内容与自测项。

---

## 文档索引

更详细的实施与配置说明请参阅 [00-doc/](00-doc/)：
- [错题本 - 工程优化与系统架构](00-doc/错题本-工程优化与系统架构.md)
- [错题本 - 部署与运维手册](00-doc/错题本-部署与运维手册.md)
- [错题本 - 使用手册 (WorkBuddy 工作流)](00-doc/错题本-使用手册-WorkBuddy工作流.md)
- [错题本 - 最小闭环实施方案](00-doc/错题本-最小闭环实施方案.md)
- [错题本 - 设计文档](00-doc/错题本-设计文档.md)
- [错题本 - WorkBuddy 指令卡](00-doc/错题本-WorkBuddy指令卡.md)
