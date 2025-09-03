# pf_ruler 使用指南

> AI编辑器规则统一管理工具 - 让AI助手更懂你的项目

## 📖 目录

- [项目简介](#项目简介)
- [快速开始](#快速开始)
- [安装方法](#安装方法)
- [核心功能](#核心功能)
- [详细使用](#详细使用)
- [配置说明](#配置说明)
- [常见问题](#常见问题)
- [进阶用法](#进阶用法)

## 🚀 项目简介

`pf_ruler` 是一款基于 Go 语言的命令行工具，用于统一管理各 AI 编辑器的规则配置。它解决了在不同 AI 编辑器间切换时需要重复配置规则的问题，让你可以：

- 📝 **一次配置，多平台使用** - 在 `.ruler` 目录中维护统一的规则
- 🔄 **自动转换格式** - 自动转换为 Trae、Cursor 等平台的原生格式
- 🎯 **项目化规则管理** - 每个项目独立的规则配置
- 🛠️ **智能规则生成** - 根据技术栈自动生成基础规则

## ⚡ 快速开始

### 1. 安装工具

```bash
# 下载对应平台的二进制文件
# macOS ARM64
curl -L -o pf_ruler https://github.com/pfinal/pf_ruler/releases/latest/download/pf_ruler-darwin_arm64
chmod +x pf_ruler

# macOS AMD64
curl -L -o pf_ruler https://github.com/pfinal/pf_ruler/releases/latest/download/pf_ruler-darwin_amd64
chmod +x pf_ruler

# Linux
curl -L -o pf_ruler https://github.com/pfinal/pf_ruler/releases/latest/download/pf_ruler-linux_amd64
chmod +x pf_ruler

# Windows
# 下载 pf_ruler-windows_amd64.exe
```

### 2. 初始化项目

```bash
# 进入你的项目目录
cd your-project

# 初始化规则管理结构
./pf_ruler init
```

### 3. 生成规则

```bash
# 生成默认平台规则（Trae）
./pf_ruler generate

# 生成指定平台规则
./pf_ruler generate --platform=cursor
```

## 📦 安装方法

### 方法一：直接下载（推荐）

1. 访问 [GitHub Releases](https://github.com/pfinal/pf_ruler/releases)
2. 下载对应平台的二进制文件
3. 添加执行权限：`chmod +x pf_ruler`
4. 移动到 PATH 目录：`sudo mv pf_ruler /usr/local/bin/`

### 方法二：从源码编译

```bash
git clone https://github.com/pfinal/pf_ruler.git
cd pf_ruler
go build -o pf_ruler main.go
```

### 方法三：使用 Go 安装

```bash
go install github/pfinal/pf_ruler@latest
```

## 🎯 核心功能

### 1. 项目初始化 (`init`)

创建标准化的规则管理目录结构：

```
.ruler/
├── config.yaml          # 配置文件
├── global/              # 全局规则
│   ├── go_rules.md
│   ├── php_rules.md
│   └── frontend_rules.md
├── project/             # 项目特定规则
│   ├── requirements.md
│   └── tech_stack.yaml
└── templates/           # 规则模板
```

**特性：**
- 🎨 交互式项目需求收集
- 🔧 自动生成技术栈规则
- 📝 智能 .gitignore 配置
- 🏗️ 标准化目录结构

### 2. 规则生成 (`generate`)

将统一规则转换为各平台原生格式：

**支持的平台：**
- **Trae**: 生成 `.trae/rules/project_rules.md`
- **Cursor**: 生成 `.cursor/rules.json`

**特性：**
- 🔄 自动格式转换
- 🎯 平台特定优化
- ⚡ 增量更新支持
- 🛡️ 冲突检测与处理

## 📚 详细使用

### 初始化命令详解

```bash
pf_ruler init
```

**执行流程：**
1. 创建 `.ruler` 目录结构
2. 更新 `.gitignore` 文件
3. 交互式收集项目信息：
   - 项目名称
   - 技术栈选择
   - 代码规范要求
   - 安全约束
   - 目标 AI 编辑器
4. 生成配置文件
5. 创建技术栈规则文件

**示例输出：**
```
✅ .ruler 目录结构已创建（包含 global/project/templates 子目录）
✅ .gitignore 已添加 .ruler/ 忽略规则
✅ 项目需求已写入 .ruler/project/requirements.md
✅ 技术栈信息已写入 .ruler/project/tech_stack.yaml
✅ 基础配置文件 .ruler/config.yaml 已创建
✅ 技术栈全局规则文件已生成
```

### 生成命令详解

```bash
# 基本用法
pf_ruler generate

# 指定平台
pf_ruler generate --platform=cursor

# 强制覆盖
pf_ruler generate --platform=cursor --force
```

**参数说明：**
- `--platform, -p`: 目标平台 (trae, cursor)
- `--force, -f`: 强制覆盖现有文件

**执行流程：**
1. 验证平台参数
2. 加载统一规则文件
3. 执行跨平台转换
4. 输出到目标目录

**示例输出：**
```
✅ 已加载规则（项目规则 3 条 + 全局规则 15 条）
✅ 已完成 cursor 规则格式转换
✅ cursor 规则已生成: .cursor/rules.json
```

## ⚙️ 配置说明

### 配置文件结构

```yaml
# .ruler/config.yaml
default_platform: "trae"           # 默认目标平台
rule_priority: ["project", "global", "templates"]  # 规则优先级
last_init_time: "2025-01-20 10:30:00"            # 最后初始化时间
```

### 规则优先级

1. **project**: 项目特定规则（最高优先级）
2. **global**: 全局技术栈规则
3. **templates**: 基础模板规则（最低优先级）

### 自定义规则

你可以在 `.ruler` 目录中手动编辑规则文件：

```markdown
# .ruler/global/custom_rules.md

## 项目特定规范
- 所有 API 响应必须包含 status 字段
- 错误信息使用中文
- 日志格式：`[时间] [级别] 消息`

## 安全要求
- 所有用户输入必须验证
- 敏感操作需要二次确认
- 定期更新依赖包
```

## 🔧 常见问题

### Q1: 下载的二进制文件无法运行

**解决方案：**
```bash
# 添加执行权限
chmod +x pf_ruler-darwin_amd64

# 检查文件类型
file pf_ruler-darwin_amd64

# 检查文件权限
ls -la pf_ruler-darwin_amd64
```

### Q2: 初始化时提示目录已存在

**解决方案：**
```bash
# 删除现有目录重新初始化
rm -rf .ruler
pf_ruler init
```

### Q3: 生成规则时提示文件已存在

**解决方案：**
```bash
# 使用强制覆盖
pf_ruler generate --platform=cursor --force

# 或手动删除后重试
rm .cursor/rules.json
pf_ruler generate --platform=cursor
```

### Q4: 如何添加新的 AI 编辑器支持

**解决方案：**
1. 在 `pkg/platform/` 目录下创建新的适配器
2. 实现 `PlatformAdapter` 接口
3. 在 `cmd/generate.go` 中注册新平台

## 🚀 进阶用法

### 1. 批量项目处理

```bash
# 为多个项目批量初始化
for project in project1 project2 project3; do
    cd $project
    pf_ruler init
    pf_ruler generate --platform=cursor
    cd ..
done
```

### 2. 自定义规则模板

在 `.ruler/templates/` 目录下创建自定义模板：

```markdown
# .ruler/templates/company_standards.md

## 公司编码规范
- 提交信息格式：`type(scope): description`
- 分支命名：`feature/功能名` 或 `fix/问题描述`
- 代码审查必须通过
```

### 3. 集成到 CI/CD

```yaml
# .github/workflows/rules.yml
name: Update AI Rules
on:
  push:
    paths: ['.ruler/**']
    branches: [main]

jobs:
  update-rules:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Download pf_ruler
        run: |
          curl -L -o pf_ruler https://github.com/pfinal/pf_ruler/releases/latest/download/pf_ruler-linux_amd64
          chmod +x pf_ruler
      - name: Generate Rules
        run: |
          ./pf_ruler generate --platform=cursor
          ./pf_ruler generate --platform=trae
      - name: Commit Changes
        run: |
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add .
          git commit -m "chore: update AI editor rules" || exit 0
          git push
```

### 4. 规则版本管理

```bash
# 创建规则快照
cp -r .ruler .ruler.backup.$(date +%Y%m%d)

# 恢复规则版本
cp -r .ruler.backup.20250120 .ruler

# 比较规则差异
diff -r .ruler .ruler.backup.20250120
```

## 📞 获取帮助

```bash
# 查看帮助信息
pf_ruler --help

# 查看子命令帮助
pf_ruler init --help
pf_ruler generate --help

# 查看版本信息
pf_ruler version
```

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

**让 AI 助手更懂你的项目，从 pf_ruler 开始！** 🚀
