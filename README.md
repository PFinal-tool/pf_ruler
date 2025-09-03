# pf_ruler - AI编辑器规则统一管理工具

pf_ruler 是一款基于 Go 语言的命令行工具，用于统一管理各 AI 编辑器的规则配置。

## 🚀 核心功能

- **一键初始化规则目录结构** - 自动创建标准化的规则管理目录
- **交互式收集项目需求** - 通过问答引导收集项目技术栈、规范要求等
- **跨平台规则生成** - 支持将统一规则自动转换为各 AI 编辑器的适配格式
- **降低使用门槛** - 无需手动编写复杂配置，通过命令行交互即可完成全流程操作

## 📋 支持平台

- **Trae** - 生成 `.trae/rules/project_rules.md` 文件
- **Cursor** - 生成 `.cursor/rules.json` 文件

## 🛠️ 安装

```bash
# 克隆项目
git clone https://github.com/pfinal/pf_ruler.git
cd pf_ruler

# 构建项目
go build -o pf_ruler .

# 或者直接运行
go run main.go
```

## 📖 使用方法

### 1. 初始化规则目录（`init` 命令）

```bash
# 进入项目根目录
cd ~/projects/your-project

# 执行初始化
./pf_ruler init
```

初始化过程会：
- 创建 `.ruler` 目录结构（包含 `global`、`project`、`templates` 子目录）
- 自动处理 `.gitignore` 配置
- 通过交互式问答收集项目需求
- 生成基础配置文件

### 2. 生成跨平台规则（`generate` 命令）

```bash
# 生成默认平台（Trae）规则
./pf_ruler generate

# 生成指定平台规则
./pf_ruler generate --platform=cursor

# 强制覆盖现有文件
./pf_ruler generate --platform=cursor --force
```

## 🏗️ 项目结构

```
your-project/
├── .ruler/                    # 规则管理目录
│   ├── config.yaml           # 工具配置文件
│   ├── global/               # 全局通用规则
│   ├── project/              # 项目特定规则
│   │   ├── requirements.md   # 项目需求文档
│   │   └── tech_stack.yaml  # 技术栈信息
│   └── templates/            # 自定义规则模板
├── .trae/                    # Trae 平台规则输出
│   └── rules/
│       └── project_rules.md
├── .cursor/                  # Cursor 平台规则输出
│   └── rules.json
└── pf_ruler                  # 工具可执行文件
```

## 🔧 配置说明

### 配置文件 (.ruler/config.yaml)

```yaml
default_platform: trae          # 默认生成平台
rule_priority:                  # 规则优先级
  - project                     # 项目规则（最高优先级）
  - global                      # 全局规则（次优先级）
  - templates                   # 模板规则（可选）
last_init_time: "2025-09-03 09:02:19"  # 最后初始化时间
```

### 技术栈配置 (.ruler/project/tech_stack.yaml)

```yaml
project_name: "your-project"    # 项目名称
tech_stacks:                    # 技术栈列表
  - "Go+Gin"
  - "MySQL"
  - "Redis"
ai_editors:                     # 目标AI编辑器
  - "Trae"
  - "Cursor"
created_at: "2025-09-03 09:02:19"  # 创建时间
```

## 🎯 使用流程示例

### 完整工作流程

```bash
# 1. 进入项目目录
cd ~/projects/ecommerce-backend

# 2. 初始化规则管理
./pf_ruler init

# 3. 交互式配置项目需求
? 请输入项目名称：电商后台系统
? 请选择技术栈（可多选）： 
  ◉ Go+Gin
  ◉ MySQL
  ◉ Redis
  ◉ JWT
? 请输入代码规范要求：函数命名采用 snake_case，每行代码不超过 80 字符
? 请输入安全约束：敏感数据（如密码）需 bcrypt 加密存储

# 4. 生成 Trae 平台规则
./pf_ruler generate --platform=trae

# 5. 生成 Cursor 平台规则
./pf_ruler generate --platform=cursor
```

## 🔌 扩展新平台

pf_ruler 支持通过"适配器接口"新增 AI 编辑器平台，无需修改核心逻辑。

### 实现步骤

1. 在 `pkg/platform/` 目录下创建新平台文件（如 `copilot.go`）
2. 实现 `PlatformAdapter` 接口：
   ```go
   type PlatformAdapter interface {
       Name() string                // 返回平台名称
       DefaultOutputPath() string   // 返回平台规则默认输出路径
       Convert(ruleSet *RuleSet) ([]byte, error)  // 将统一规则转换为平台格式
   }
   ```
3. 在工具初始化时注册适配器，即可支持 `--platform=copilot` 命令

## 🐛 故障排除

### 常见问题

1. **".ruler 目录不存在"**
   - 解决：先运行 `./pf_ruler init` 命令初始化

2. **"文件已存在"**
   - 解决：使用 `--force` 标志强制覆盖，或手动删除现有文件

3. **"不支持的平台"**
   - 解决：检查 `--platform` 参数，当前支持：`trae`、`cursor`

## 📝 开发说明

### 技术栈

- **开发语言**: Go 1.21+
- **命令行框架**: Cobra
- **交互式问答**: Survey
- **配置文件**: YAML
- **跨平台**: 支持 Windows、macOS、Linux

### 项目结构

```
pf_ruler/
├── cmd/                       # 命令行命令
│   ├── root.go               # 根命令
│   ├── init.go               # 初始化命令
│   └── generate.go           # 生成命令
├── pkg/                      # 核心包
│   ├── platform/             # 平台适配器
│   │   ├── base.go           # 基础接口
│   │   ├── trae.go           # Trae 适配器
│   │   └── cursor.go         # Cursor 适配器
│   └── rules/                # 规则管理
│       ├── types.go          # 规则类型定义
│       └── loader.go         # 规则加载器
├── main.go                   # 主程序入口
└── go.mod                    # Go 模块文件
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证。
