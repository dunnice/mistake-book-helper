<!-- 感谢提交 Pull Request！请尽可能完整地填写以下信息以帮助我们快速审查你的代码。 -->

## 变更类型 (Type of Change)
<!-- 请勾选本次变更适用的类型（使用 [x] 标记） -->
- [ ] 🚀 新功能 (Feature)
- [ ] 🐛 Bug 修复 (Bug Fix)
- [ ] 🎨 UI/样式优化 (UI Enhancement)
- [ ] 📝 文档更新 (Documentation)
- [ ] ⚡ 性能优化 (Performance)
- [ ] ♻️ 代码重构 (Refactoring)
- [ ] 🔧 构建/工具/配置变更 (Chore)

---

## 解决的问题 / 需求背景 (Related Issue / Motivation)
<!-- 本次 PR 解决了哪个 Issue？或源于什么业务背景？例如：Fixes #123 或 解决手机端扫码提交后状态未即时更新的问题 -->
- 关联 Issue / 讨论：#
- 背景简述：

---

## 修改内容说明 (Description of Changes)
<!-- 请详细说明这次具体修改了什么，做了哪些调整？ -->
1. 
2. 
3. 

---

## 本地自测情况 (How Has This Been Tested?)
<!-- 请说明你在本地如何验证这些修改？是否保证测试通过？ -->
- [ ] 后端编译通过（`cd mistake-book-go-api && go build ./cmd/server`）
- [ ] 前端类型检查与构建通过（`cd mistake-book-view && npm run build`）
- [ ] 手动端到端功能验证通过（如：错题录入、列表筛选、详情抽屉、扫码答题等）
- **测试环境与浏览器说明**（例如：macOS / Chrome 120 / Go 1.22）：

---

## 安全与规范自查清单 (Checklist)
<!-- 提交前请逐项确认并勾选 -->
- [ ] 我已阅读并遵守项目的 [CONTRIBUTING.md](CONTRIBUTING.md) 贡献指南。
- [ ] 我的代码遵循项目的代码规范，无多余未使用的引入与调试日志。
- [ ] **严禁敏感信息泄漏**：代码、文档中**绝无**硬编码的真实服务器 IP、端口、账号、密码、Token 或 API Key（已全部提取到 `.env` 并在 `.gitignore` 中屏蔽）。
- [ ] **隐私保护**：未包含任何真实的个人隐私数据（如学生真实姓名、私人知识库 ID 等）。
- [ ] 文档已同步更新（若涉及接口、环境变量或主要交互变动，文档位于 `00-doc/错题本-*.md` 且顶层无重复一级标题）。
