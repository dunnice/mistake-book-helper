> **出品人**：ice7  
> 感谢你对 **错题本助手 (MistakeBook Helper)** 的关注与支持！我们欢迎一切形式的贡献，无论是报告 Bug、提出新功能建议、改进文档还是提交代码补丁。

为了让大家的协作更顺畅，请在参与前仔细阅读以下指南。

---

## 参与流程（Pull Request 工作流）

我们遵循标准的 GitHub Flow 分支工作模式：

1. **Fork 本仓库**：点击 GitHub 页面右上角的 `Fork` 按钮，将项目复制到你的个人账户下。
2. **克隆到本地**：
   ```bash
   git clone https://github.com/<your-username>/mistake-book-helper.git
   cd mistake-book-helper
   ```
3. **创建独立特性分支**：切勿直接在 `main` 或 `master` 分支上修改。
   ```bash
   git checkout -b feat/your-feature-name
   # 或修复 bug:
   git checkout -b fix/issue-summary
   ```
4. **进行开发与本地验证**：
   - 参考下方的[本地开发与调试](#本地开发与调试)完成代码编写；
   - **提交 PR 前请确保本地测试通过**（运行 `go build`、前端 `npm run build` 无任何错误）。
5. **提交 Commit 并推送到你的分支**：
   ```bash
   git add .
   git commit -m "feat(module): add concise description of changes"
   git push origin feat/your-feature-name
   ```
6. **创建 Pull Request (PR)**：
   - 回到 GitHub 原仓库，点击 `Compare & pull request`；
   - 按照系统自动填入的 [PR 模板](.github/PULL_REQUEST_TEMPLATE.md) 认真填写修改说明、关联 Issue 以及自测核对清单。

---

## 本地开发与调试

### 1. 环境依赖
- **后端**：Go 1.22+
- **前端**：Node.js 18+ (推荐 20+), npm 9+
- **系统工具**：Make, Git, Bash

### 2. 环境变量配置（安全第一）
在启动服务前，先基于模板建立本地私有配置文件：
```bash
cp .env.example .env
```
> [!IMPORTANT]
> **切勿提交 `.env` 文件！** 本项目的 `.gitignore` 已配置忽略 `.env` 及所有私有配置文件。本地修改 API Key、数据库连接密码、服务器 IP 等，必须全部保存在本地 `.env` 中，严禁将真实凭据硬编码提交到代码库。

### 3. 一键启动本地联调
根目录提供了便捷的 Makefile：
```bash
# 启动本地开发联调（后端监听 8080，前端 Vite 开发服务器监听 5173 具备热重载与 API 代理）
make dev
```
若需单独启动：
- **启动后端**：
  ```bash
  cd mistake-book-go-api
  go run ./cmd/server
  ```
- **启动前端**：
  ```bash
  cd mistake-book-view
  npm install
  npm run dev
  ```

---

## 代码与提交规范

### 1. Commit Message 规范
推荐采用 Conventional Commits 格式：
- `feat`: 新增功能（如 `feat(view): add quick capsule filter`）
- `fix`: 修复缺陷（如 `fix(api): fix review count calculation`）
- `docs`: 文档变更（如 `docs: update deployment manual`）
- `style`: 代码格式、分号缺失等不影响运行的修改
- `refactor`: 重构代码（既非修复 bug 也非新增功能）
- `perf`: 性能优化
- `test`: 增加或修改测试用例
- `chore`: 构建过程、依赖变更或辅助工具变动

### 2. 代码质量要求
- **后端 Go**：保持代码格式化（运行 `gofmt -s -w .`），避免无错误处理的静默 panic，严格遵循 Gin 与标准库惯例。
- **前端 Vue/TS**：优先使用 `<script setup lang="ts">` 组合式 API，严格保持类型检查（`npm run build` 会运行 `vue-tsc` 进行静态类型校验）。
- **无 Docker 依赖**：本项目以原生高性能运行为主，支持本地与服务器无容器化部署，请勿引入破坏原生运行的额外重型封装。

---

## 文档贡献规范

- **命名规范**：所有新增或修改的架构设计、使用手册等文档统一放置在 `00-doc/` 目录下，且文件名必须以 `错题本-` 开头（例如：`00-doc/错题本-某新功能设计.md`）。
- **Obsidian 兼容规范**：在编写 Markdown 文档时，**避免在正文最上方编写重复的一级标题**（如 `# 标题`）。由于在 Obsidian 等笔记工具中打开时界面已自动展示文档标题，重复一级标题会导致视图双标题冗余。
- **AI Agent 智能体输入支持**：由于本项目核心依赖各类 AI Agent 智能助手进行错题批改与结构化数据入库，文档涉及接口契约与指令说明时，需保持精准详尽，并同步更新技能卡片。

---

## 安全与隐私核查清单

在向 GitHub 提交代码和 PR 前，请务必进行安全自查：
- [ ] 检查代码与文档中**绝无**硬编码的公网服务器真实 IP 地址；
- [ ] 检查**绝无**任何真实账号、密码、Token、API Key 残留；
- [ ] 检查**绝无**个人隐私（真实学生姓名、学校班级、私人云盘或知识库 ID 等）；
- [ ] 确保所有的本地私有配置均已提取至 `.env`，且 `.env` 未被 git 追踪。

再次感谢你为 **错题本助手** 添砖加瓦！
