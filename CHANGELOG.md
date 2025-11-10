# 更新日志

## [1.1.0] - 2025-11-10

### 新增功能
- ✨ 支持从 `requirements.md` 中自动提取技术栈信息
- ✨ 支持加载 `project/` 目录下的所有 `.md` 规则文件（除 `requirements.md` 外）
- ✨ 新增 `loadProjectFiles()` 函数，自动扫描并加载项目规则文件

### 修复问题
- 🐛 修复全局规则文件内容丢失的问题（`parseMarkdownRules` 现在会收集所有非标题行内容）
- 🐛 修复技术栈信息在项目信息中显示为空的问题
- 🐛 修复规则解析时只收集列表项，导致其他格式内容丢失的问题

### 改进
- 🔧 优化规则加载逻辑，支持更灵活的规则文件组织方式
- 🔧 改进技术栈提取逻辑，支持从多个来源获取技术栈信息
- 🔧 优化 Markdown 解析，支持更多内容格式

### 技术细节
- 修改 `LoadProjectRules()` 函数，增加对 `project/` 目录下所有 `.md` 文件的扫描
- 修改 `LoadMetadata()` 函数，支持从 `requirements.md` 提取技术栈
- 修改 `getProjectTechStacks()` 函数，支持从 `requirements.md` 提取技术栈
- 优化 `parseMarkdownRules()` 函数，收集所有非标题行内容

---

## [1.0.0] - 2025-11-05

### 初始版本
- 🎉 首次发布
- ✨ 支持项目初始化（`init` 命令）
- ✨ 支持规则生成（`generate` 命令）
- ✨ 支持 Trae 和 Cursor 平台
- ✨ 交互式项目需求收集
- ✨ 技术栈规则自动生成

