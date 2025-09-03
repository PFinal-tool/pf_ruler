/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化规则目录结构",
	Long: `在当前项目目录下创建标准化的规则管理目录（.ruler），
自动处理 .gitignore 配置，并通过交互式引导收集项目需求，生成基础规则文件。

示例：
  pf_ruler init  # 在当前目录初始化规则管理结构
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 创建 .ruler 目录结构
		createRulerDirs()

		// 2. 处理 .gitignore 文件
		handleGitignore()

		// 3. 交互式收集项目需求
		techStacks := collectProjectRequirements()

		// 4. 生成基础配置文件
		generateConfigFile()

		// 5. 根据技术栈生成全局规则文件
		generateGlobalRules(techStacks)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// 这里可以添加 init 命令的标志
}

// 创建 .ruler 目录结构
func createRulerDirs() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		redBold("❌ 获取当前目录失败：", err)
		os.Exit(1)
	}

	// 定义要创建的目录结构
	dirs := []string{
		filepath.Join(currentDir, ".ruler"),
		filepath.Join(currentDir, ".ruler", "global"),
		filepath.Join(currentDir, ".ruler", "project"),
		filepath.Join(currentDir, ".ruler", "templates"),
	}

	// 创建目录
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				redBold("❌ 创建目录失败：", err)
				os.Exit(1)
			}
		}
	}

	greenBold("✅ .ruler 目录结构已创建（包含 global/project/templates 子目录）")
}

// 处理 .gitignore 文件
func handleGitignore() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		redBold("❌ 获取当前目录失败：", err)
		os.Exit(1)
	}

	gitignorePath := filepath.Join(currentDir, ".gitignore")

	// 检查 .gitignore 文件是否存在
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		// 文件不存在，给出警告
		yellowBold("⚠️  未检测到 .gitignore，规则文件可能被提交至版本库")
		return
	}

	// 读取 .gitignore 文件内容
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		redBold("❌ 读取 .gitignore 文件失败：", err)
		os.Exit(1)
	}

	// 检查是否已包含 .ruler/ 忽略规则
	if strings.Contains(string(content), ".ruler/") {
		// 已包含，无需操作
		return
	}

	// 追加 .ruler/ 忽略规则
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		redBold("❌ 打开 .gitignore 文件失败：", err)
		os.Exit(1)
	}
	defer file.Close()

	// 确保文件末尾有换行符
	if len(content) > 0 && content[len(content)-1] != '\n' {
		if _, err := file.WriteString("\n"); err != nil {
			redBold("❌ 写入 .gitignore 文件失败：", err)
			os.Exit(1)
		}
	}

	// 写入 .ruler/ 忽略规则
	if _, err := file.WriteString(".ruler/\n"); err != nil {
		redBold("❌ 写入 .gitignore 文件失败：", err)
		os.Exit(1)
	}

	greenBold("✅ .gitignore 已添加 .ruler/ 忽略规则")
}

// 交互式收集项目需求
func collectProjectRequirements() []string {
	// 获取当前工作目录名称作为项目名称默认值
	currentDir, err := os.Getwd()
	if err != nil {
		redBold("❌ 获取当前目录失败：", err)
		os.Exit(1)
	}

	projectNameDefault := filepath.Base(currentDir)

	// 定义交互式问题
	var projectName string
	projectNamePrompt := &survey.Input{
		Message: "请输入项目名称：",
		Default: projectNameDefault,
	}

	var techStacks []string
	techStackPrompt := &survey.MultiSelect{
		Message: "请选择技术栈（可多选）：",
		Options: []string{
			"Go+Gin",
			"PHP+Laravel",
			"PHP+ThinkPHP",
			"PHP+Slim",
			"React+TypeScript",
			"Vue.js",
			"Java+SpringBoot",
			"Python+Django",
			"Python+Flask",
			"Node.js+Express",
			"Node.js+Koa",
			"MySQL",
			"PostgreSQL",
			"MongoDB",
			"Redis",
			"Memcached",
			"JWT",
			"OAuth2",
			"Docker",
			"Kubernetes",
		},
	}

	var codeStandards string
	codeStandardsPrompt := &survey.Input{
		Message: "请输入代码规范要求：",
		Default: "函数命名采用 snake_case，每行代码不超过 80 字符",
	}

	var securityConstraints string
	securityConstraintsPrompt := &survey.Input{
		Message: "请输入安全约束：",
		Default: "敏感数据（如密码）需加密存储",
	}

	var aiEditors []string
	aiEditorsPrompt := &survey.MultiSelect{
		Message: "请选择目标 AI 编辑器（可多选）：",
		Options: []string{
			"Trae",
			"Cursor",
			"GitHub Copilot X",
		},
	}

	// 执行交互式问答
	survey.AskOne(projectNamePrompt, &projectName)
	survey.AskOne(techStackPrompt, &techStacks)
	survey.AskOne(codeStandardsPrompt, &codeStandards)
	survey.AskOne(securityConstraintsPrompt, &securityConstraints)
	survey.AskOne(aiEditorsPrompt, &aiEditors)

	// 生成 requirements.md 文件
	requirementsContent := fmt.Sprintf(`# %s 项目需求文档

## 项目基本信息
- 项目名称：%s
- 初始化时间：%s

## 技术栈
%s

## 代码规范
%s

## 安全约束
%s

## 目标 AI 编辑器
%s
`,
		projectName,
		projectName,
		time.Now().Format("2006-01-02 15:04:05"),
		formatList(techStacks, "- "),
		codeStandards,
		securityConstraints,
		formatList(aiEditors, "- "),
	)

	// 生成 tech_stack.yaml 文件
	techStackData := map[string]interface{}{
		"project_name": projectName,
		"tech_stacks":  techStacks,
		"ai_editors":   aiEditors,
		"created_at":   time.Now().Format("2006-01-02 15:04:05"),
	}

	techStackYaml, err := yaml.Marshal(techStackData)
	if err != nil {
		redBold("❌ 生成技术栈 YAML 文件失败：", err)
		os.Exit(1)
	}

	// 写入文件
	projectDir := filepath.Join(currentDir, ".ruler", "project")

	if err := os.WriteFile(filepath.Join(projectDir, "requirements.md"), []byte(requirementsContent), 0644); err != nil {
		redBold("❌ 写入项目需求文件失败：", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filepath.Join(projectDir, "tech_stack.yaml"), techStackYaml, 0644); err != nil {
		redBold("❌ 写入技术栈文件失败：", err)
		os.Exit(1)
	}

	greenBold("✅ 项目需求已写入 .ruler/project/requirements.md")
	greenBold("✅ 技术栈信息已写入 .ruler/project/tech_stack.yaml")

	return techStacks
}

// formatList 格式化列表为Markdown格式
func formatList(items []string, prefix string) string {
	if len(items) == 0 {
		return "暂无配置"
	}

	var result strings.Builder
	for _, item := range items {
		result.WriteString(fmt.Sprintf("%s%s\n", prefix, item))
	}
	return strings.TrimSpace(result.String())
}

// 生成基础配置文件
func generateConfigFile() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		redBold("❌ 获取当前目录失败：", err)
		os.Exit(1)
	}

	// 定义配置内容
	configData := map[string]interface{}{
		"default_platform": "trae",
		"rule_priority":    [3]string{"project", "global", "templates"},
		"last_init_time":   time.Now().Format("2006-01-02 15:04:05"),
	}

	// 转换为 YAML
	configYaml, err := yaml.Marshal(configData)
	if err != nil {
		redBold("❌ 生成配置文件失败：", err)
		os.Exit(1)
	}

	// 写入文件
	configPath := filepath.Join(currentDir, ".ruler", "config.yaml")

	if err := os.WriteFile(configPath, configYaml, 0644); err != nil {
		redBold("❌ 写入配置文件失败：", err)
		os.Exit(1)
	}

	greenBold("✅ 基础配置文件 .ruler/config.yaml 已创建")
}

// 根据技术栈生成全局规则文件
func generateGlobalRules(techStacks []string) {
	if len(techStacks) == 0 {
		return
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		redBold("❌ 获取当前目录失败：", err)
		return
	}

	globalDir := filepath.Join(currentDir, ".ruler", "global")

	// 为每个技术栈生成对应的规则文件
	for _, tech := range techStacks {
		generateTechStackRules(globalDir, tech)
	}

	greenBold("✅ 技术栈全局规则文件已生成")
}

// 为特定技术栈生成规则文件
func generateTechStackRules(globalDir, tech string) {
	var rules []string
	var filename string

	switch {
	case strings.Contains(tech, "PHP"):
		filename = "php_rules.md"
		rules = generatePHPRules(tech)
	case strings.Contains(tech, "Go"):
		filename = "go_rules.md"
		rules = generateGoRules(tech)
	case strings.Contains(tech, "Java"):
		filename = "java_rules.md"
		rules = generateJavaRules(tech)
	case strings.Contains(tech, "Python"):
		filename = "python_rules.md"
		rules = generatePythonRules(tech)
	case strings.Contains(tech, "Node.js"):
		filename = "nodejs_rules.md"
		rules = generateNodeRules(tech)
	case strings.Contains(tech, "React") || strings.Contains(tech, "Vue"):
		filename = "frontend_rules.md"
		rules = generateFrontendRules(tech)
	case strings.Contains(tech, "MySQL") || strings.Contains(tech, "PostgreSQL"):
		filename = "database_rules.md"
		rules = generateDatabaseRules(tech)
	case strings.Contains(tech, "Redis") || strings.Contains(tech, "Memcached"):
		filename = "cache_rules.md"
		rules = generateCacheRules(tech)
	case strings.Contains(tech, "Docker") || strings.Contains(tech, "Kubernetes"):
		filename = "devops_rules.md"
		rules = generateDevOpsRules(tech)
	default:
		return
	}

	if len(rules) > 0 {
		filePath := filepath.Join(globalDir, filename)
		content := strings.Join(rules, "\n")
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			redBold(fmt.Sprintf("❌ 写入 %s 失败：%v", filename, err))
		}
	}
}

// 生成PHP规则
func generatePHPRules(tech string) []string {
	rules := []string{
		"# PHP 开发规范与最佳实践",
		"",
		"## 基础规则",
		"- 遵循 PSR-12 代码规范",
		"- 所有 PHP 文件必须以 <?php 开头，不使用关闭标签 ?>",
		"- 使用 4 个空格缩进，不使用 Tab",
		"- 行宽不超过 120 字符",
		"",
		"## 文件组织",
		"- 类名与文件名保持一致，使用 PascalCase",
		"- 函数名与方法名使用 camelCase",
		"- 常量名使用全大写加下划线",
		"- 目录结构遵循 PSR-4 自动加载标准",
		"",
		"## 注释与文档",
		"- 公共方法必须写 PHPDoc，包括 @param 和 @return",
		"- 代码逻辑复杂处需要行内注释",
		"",
		"## 安全与实践",
		"- 避免使用 mysql_*，统一使用 PDO 或框架自带的数据库层",
		"- 避免硬编码敏感信息，使用配置文件或环境变量",
		"- 异常处理要用 try/catch，不允许裸 die/exit",
		"",
		"## 测试",
		"- 所有新类必须配套 PHPUnit 单元测试",
		"- 测试文件放在 tests/ 目录下",
	}

	// 根据具体框架添加特定规则
	if strings.Contains(tech, "Laravel") {
		rules = append(rules, []string{
			"",
			"## Laravel 特定规范",
			"- 使用 Eloquent ORM 进行数据库操作",
			"- 遵循 MVC 架构模式",
			"- 使用 Artisan 命令生成代码",
			"- 启用 CSRF 保护",
			"- 使用 Laravel 的验证器进行数据验证",
			"- 使用 Laravel 的缓存系统",
		}...)
	}

	return rules
}

// 生成Go规则
func generateGoRules(tech string) []string {
	return []string{
		"# Go 开发规范与最佳实践",
		"",
		"## 代码规范",
		"- 使用 gofmt 格式化代码",
		"- 遵循 Go 官方命名约定",
		"- 使用 go mod 管理依赖",
		"- 包名使用小写字母",
		"- 接口名以 er 结尾",
		"",
		"## 错误处理",
		"- 始终检查错误返回值",
		"- 使用 errors.Wrap 包装错误",
		"- 避免忽略错误",
		"- 自定义错误类型实现 Error() 方法",
		"",
		"## 性能优化",
		"- 使用 sync.Pool 复用对象",
		"- 避免在循环中分配内存",
		"- 使用 strings.Builder 进行字符串拼接",
		"- 合理使用 goroutine 和 channel",
	}
}

// 生成Java规则
func generateJavaRules(tech string) []string {
	return []string{
		"# Java 开发规范与最佳实践",
		"",
		"## 代码规范",
		"- 遵循 Java 命名约定",
		"- 使用 Lombok 减少样板代码",
		"- 启用代码检查工具",
		"- 类名使用 PascalCase",
		"- 方法名使用 camelCase",
		"",
		"## Spring Boot 规范",
		"- 使用 Spring Boot 自动配置",
		"- 遵循 RESTful API 设计原则",
		"- 使用 Spring Security 进行安全控制",
		"- 使用 Spring Data JPA 进行数据访问",
	}
}

// 生成Python规则
func generatePythonRules(tech string) []string {
	return []string{
		"# Python 开发规范与最佳实践",
		"",
		"## 代码规范",
		"- 遵循 PEP 8 规范",
		"- 使用类型提示",
		"- 使用虚拟环境管理依赖",
		"- 函数名和变量名使用 snake_case",
		"- 类名使用 PascalCase",
		"",
		"## 最佳实践",
		"- 使用 list comprehension 和 generator",
		"- 使用 with 语句管理资源",
		"- 使用 dataclass 简化类定义",
		"- 编写 docstring 文档",
	}
}

// 生成Node.js规则
func generateNodeRules(tech string) []string {
	return []string{
		"# Node.js 开发规范与最佳实践",
		"",
		"## 安全规范",
		"- 使用 helmet 中间件",
		"- 验证所有输入",
		"- 使用 bcrypt 加密密码",
		"- 定期更新依赖",
		"- 使用 HTTPS",
		"",
		"## 代码规范",
		"- 使用 ESLint 进行代码检查",
		"- 使用 Prettier 格式化代码",
		"- 遵循异步编程最佳实践",
		"- 使用 async/await 而不是回调",
	}
}

// 生成前端规则
func generateFrontendRules(tech string) []string {
	return []string{
		"# 前端开发规范与最佳实践",
		"",
		"## 安全规范",
		"- 使用 HTTPS",
		"- 验证用户输入",
		"- 防止 XSS 攻击",
		"- 使用 CSP 策略",
		"- 避免在客户端存储敏感信息",
		"",
		"## 代码规范",
		"- 使用 ESLint 和 Prettier",
		"- 遵循组件化开发原则",
		"- 使用 TypeScript 进行类型检查",
		"- 编写单元测试",
	}
}

// 生成数据库规则
func generateDatabaseRules(tech string) []string {
	return []string{
		"# 数据库开发规范与最佳实践",
		"",
		"## 安全规范",
		"- 使用参数化查询防止 SQL 注入",
		"- 限制数据库用户权限",
		"- 定期备份数据",
		"- 加密敏感数据",
		"",
		"## 性能优化",
		"- 合理设计索引",
		"- 避免 SELECT *",
		"- 使用连接池",
		"- 定期分析慢查询",
	}
}

// 生成缓存规则
func generateCacheRules(tech string) []string {
	return []string{
		"# 缓存使用规范与最佳实践",
		"",
		"## 使用规范",
		"- 设置合理的过期时间",
		"- 避免缓存穿透",
		"- 使用缓存预热",
		"- 监控缓存命中率",
		"",
		"## 注意事项",
		"- 缓存数据一致性",
		"- 缓存雪崩防护",
		"- 合理设置内存限制",
		"- 定期清理过期数据",
	}
}

// 生成DevOps规则
func generateDevOpsRules(tech string) []string {
	return []string{
		"# DevOps 规范与最佳实践",
		"",
		"## 容器安全",
		"- 使用非 root 用户运行容器",
		"- 定期更新基础镜像",
		"- 扫描镜像漏洞",
		"- 限制容器权限",
		"",
		"## 部署规范",
		"- 使用 CI/CD 流水线",
		"- 自动化测试",
		"- 蓝绿部署或金丝雀发布",
		"- 监控和日志收集",
	}
}
