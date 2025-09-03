package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadProjectRules 加载项目规则
func (l *FileLoader) LoadProjectRules() ([]Rule, error) {
	projectDir := filepath.Join(l.basePath, "project")

	// 检查项目目录是否存在
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return []Rule{}, nil
	}

	// 读取 requirements.md 文件
	requirementsPath := filepath.Join(projectDir, "requirements.md")
	var requirementsContent string
	if _, err := os.Stat(requirementsPath); err == nil {
		// 文件存在，读取内容
		requirementsData, err := os.ReadFile(requirementsPath)
		if err != nil {
			return nil, fmt.Errorf("读取项目需求文件失败: %w", err)
		}
		requirementsContent = string(requirementsData)
	}

	// 读取 tech_stack.yaml 文件
	techStackPath := filepath.Join(projectDir, "tech_stack.yaml")
	var techStack map[string]interface{}
	if _, err := os.Stat(techStackPath); err == nil {
		// 文件存在，读取内容
		techStackData, err := os.ReadFile(techStackPath)
		if err != nil {
			return nil, fmt.Errorf("读取技术栈文件失败: %w", err)
		}

		if err := yaml.Unmarshal(techStackData, &techStack); err != nil {
			return nil, fmt.Errorf("解析技术栈文件失败: %w", err)
		}
	}

	// 从 requirements.md 中解析所有章节内容
	requirementsRules := parseRequirementsMarkdown(requirementsContent)

	// 生成项目规则
	rules := []Rule{}

	// 技术栈规范规则
	if techStack != nil && techStack["tech_stacks"] != nil {
		techStacks := getStringSlice(techStack, "tech_stacks")
		if len(techStacks) > 0 {
			rules = append(rules, Rule{
				Title:       "技术栈规范",
				Description: "项目使用的技术栈和版本要求",
				Type:        "tech_stack",
				Content:     fmt.Sprintf("技术栈: %s", strings.Join(techStacks, ", ")),
				Priority:    5,
				Enabled:     true,
				Tags:        []string{"tech", "stack"},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
		}
	}

	// 添加从 requirements.md 解析的规则
	rules = append(rules, requirementsRules...)

	// 如果没有从文件中读取到内容，使用默认值
	if len(rules) == 0 {
		rules = append(rules, Rule{
			Title:       "代码规范",
			Description: "项目代码编写规范和要求",
			Type:        "code_style",
			Content:     "函数命名采用 snake_case，每行代码不超过 80 字符",
			Priority:    4,
			Enabled:     true,
			Tags:        []string{"code", "style", "naming"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})

		rules = append(rules, Rule{
			Title:       "安全约束",
			Description: "项目安全相关的要求和约束",
			Type:        "security",
			Content:     "敏感数据（如密码）需加密存储",
			Priority:    5,
			Enabled:     true,
			Tags:        []string{"security", "encryption"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return rules, nil
}

// parseRequirementsMarkdown 解析 requirements.md 文件，提取所有章节内容
func parseRequirementsMarkdown(content string) []Rule {
	if content == "" {
		return []Rule{}
	}

	var rules []Rule
	lines := strings.Split(content, "\n")
	
	var currentSection string
	var currentContent []string
	var inSection bool
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// 检查是否是二级标题（## 章节名）
		if strings.HasPrefix(line, "## ") {
			// 保存前一个章节的内容
			if inSection && currentSection != "" && len(currentContent) > 0 {
				rule := createRuleFromSection(currentSection, currentContent)
				if rule != nil {
					rules = append(rules, *rule)
				}
			}
			
			// 开始新章节
			currentSection = strings.TrimPrefix(line, "## ")
			currentContent = []string{}
			inSection = true
			
			// 跳过项目基本信息和技术栈章节（这些由其他方式处理）
			if currentSection == "项目基本信息" || currentSection == "技术栈" || currentSection == "目标 AI 编辑器" {
				inSection = false
			}
		} else if inSection && line != "" && !strings.HasPrefix(line, "#") {
			// 收集章节内容（跳过空行和标题行）
			currentContent = append(currentContent, line)
		}
	}
	
	// 保存最后一个章节
	if inSection && currentSection != "" && len(currentContent) > 0 {
		rule := createRuleFromSection(currentSection, currentContent)
		if rule != nil {
			rules = append(rules, *rule)
		}
	}
	
	return rules
}

// createRuleFromSection 根据章节名和内容创建规则
func createRuleFromSection(sectionName string, content []string) *Rule {
	if len(content) == 0 {
		return nil
	}
	
	// 根据章节名推断规则类型和标签
	ruleType, tags := inferRuleTypeAndTags(sectionName)
	
	// 合并内容
	contentText := strings.Join(content, "\n")
	
	// 设置优先级（根据章节类型）
	priority := 4 // 默认优先级
	if strings.Contains(strings.ToLower(sectionName), "安全") || strings.Contains(strings.ToLower(sectionName), "性能") {
		priority = 5 // 安全和性能规则优先级最高
	}
	
	return &Rule{
		Title:       sectionName,
		Description: fmt.Sprintf("项目 %s 相关的要求和规范", sectionName),
		Type:        ruleType,
		Content:     contentText,
		Priority:    priority,
		Enabled:     true,
		Tags:        tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// inferRuleTypeAndTags 根据章节名推断规则类型和标签
func inferRuleTypeAndTags(sectionName string) (string, []string) {
	sectionLower := strings.ToLower(sectionName)
	var ruleType string
	var tags []string
	
	// 推断规则类型
	switch {
	case strings.Contains(sectionLower, "代码规范") || strings.Contains(sectionLower, "编码规范"):
		ruleType = "code_style"
		tags = []string{"code", "style", "naming"}
	case strings.Contains(sectionLower, "安全约束") || strings.Contains(sectionLower, "安全规范"):
		ruleType = "security"
		tags = []string{"security", "encryption"}
	case strings.Contains(sectionLower, "性能") || strings.Contains(sectionLower, "优化"):
		ruleType = "performance"
		tags = []string{"performance", "optimization"}
	case strings.Contains(sectionLower, "测试") || strings.Contains(sectionLower, "单元测试"):
		ruleType = "testing"
		tags = []string{"testing", "unit_test"}
	case strings.Contains(sectionLower, "部署") || strings.Contains(sectionLower, "运维"):
		ruleType = "deployment"
		tags = []string{"deployment", "devops"}
	case strings.Contains(sectionLower, "数据库") || strings.Contains(sectionLower, "存储"):
		ruleType = "database"
		tags = []string{"database", "storage"}
	case strings.Contains(sectionLower, "缓存") || strings.Contains(sectionLower, "redis"):
		ruleType = "cache"
		tags = []string{"cache", "redis"}
	case strings.Contains(sectionLower, "接口") || strings.Contains(sectionLower, "api"):
		ruleType = "api"
		tags = []string{"api", "interface"}
	case strings.Contains(sectionLower, "文档") || strings.Contains(sectionLower, "注释"):
		ruleType = "documentation"
		tags = []string{"documentation", "comments"}
	default:
		ruleType = "general"
		tags = []string{"general", "requirements"}
	}
	
	// 添加章节名作为标签
	tags = append(tags, strings.ToLower(sectionName))
	
	return ruleType, tags
}

// extractCodeStandards 从 requirements.md 中提取代码规范
func extractCodeStandards(content string) string {
	if content == "" {
		return ""
	}

	// 简单的文本提取逻辑
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, "代码规范") && i+1 < len(lines) {
			// 找到下一行的内容
			nextLine := strings.TrimSpace(lines[i+1])
			if nextLine != "" && !strings.HasPrefix(nextLine, "#") {
				return nextLine
			}
		}
	}
	return ""
}

// extractSecurityConstraints 从 requirements.md 中提取安全约束
func extractSecurityConstraints(content string) string {
	if content == "" {
		return ""
	}

	// 简单的文本提取逻辑
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, "安全约束") && i+1 < len(lines) {
			// 找到下一行的内容
			nextLine := strings.TrimSpace(lines[i+1])
			if nextLine != "" && !strings.HasPrefix(nextLine, "#") {
				return nextLine
			}
		}
	}
	return ""
}

// LoadGlobalRules 加载全局规则
func (l *FileLoader) LoadGlobalRules() ([]Rule, error) {
	globalDir := filepath.Join(l.basePath, "global")
	
	// 检查全局目录是否存在
	if _, err := os.Stat(globalDir); os.IsNotExist(err) {
		return []Rule{}, nil
	}
	
	// 首先读取 global 目录中的实际文件内容
	fileRules, err := l.loadGlobalFiles(globalDir)
	if err != nil {
		return nil, fmt.Errorf("读取全局规则文件失败: %w", err)
	}
	
	// 获取项目技术栈信息
	techStacks := l.getProjectTechStacks()
	
	// 生成默认全局规则
	rules := []Rule{
		{
			Title:       "通用命名规范",
			Description: "适用于所有项目的通用命名规范",
			Type:        "naming",
			Content:     "变量和函数名应具有描述性，避免使用缩写",
			Priority:    3,
			Enabled:     true,
			Tags:        []string{"naming", "general"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "代码注释规范",
			Description: "代码注释的编写规范",
			Type:        "documentation",
			Content:     "所有公共函数和复杂逻辑都应添加注释",
			Priority:    3,
			Enabled:     true,
			Tags:        []string{"documentation", "comments"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "错误处理规范",
			Description: "错误处理的标准做法",
			Type:        "error_handling",
			Content:     "所有可能失败的操作都应进行错误处理",
			Priority:    4,
			Enabled:     true,
			Tags:        []string{"error", "handling"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	
	// 根据技术栈生成相应的规避规则
	techSpecificRules := l.generateTechSpecificRules(techStacks)
	rules = append(rules, techSpecificRules...)
	
	// 将文件规则添加到结果中
	rules = append(rules, fileRules...)
	
	return rules, nil
}

// loadGlobalFiles 读取 global 目录中的规则文件
func (l *FileLoader) loadGlobalFiles(globalDir string) ([]Rule, error) {
	var allRules []Rule
	
	// 读取目录中的所有 .md 文件
	files, err := os.ReadDir(globalDir)
	if err != nil {
		return nil, fmt.Errorf("读取全局目录失败: %w", err)
	}
	
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		
		filePath := filepath.Join(globalDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue // 跳过无法读取的文件
		}
		
		// 解析 Markdown 文件内容，提取规则
		rules := l.parseMarkdownRules(string(content), file.Name())
		allRules = append(allRules, rules...)
	}
	
	return allRules, nil
}

// parseMarkdownRules 解析 Markdown 文件内容，提取规则
func (l *FileLoader) parseMarkdownRules(content, filename string) []Rule {
	var rules []Rule
	
	lines := strings.Split(content, "\n")
	var currentRule *Rule
	var currentContent []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// 检查是否是标题行
		if strings.HasPrefix(line, "## ") {
			// 保存前一个规则
			if currentRule != nil && len(currentContent) > 0 {
				currentRule.Content = strings.Join(currentContent, "\n")
				rules = append(rules, *currentRule)
			}
			
			// 创建新规则
			title := strings.TrimPrefix(line, "## ")
			currentRule = &Rule{
				Title:       title,
				Description: fmt.Sprintf("来自 %s 的规则", filename),
				Type:        l.inferRuleType(title),
				Priority:    4, // 默认优先级
				Enabled:     true,
				Tags:        l.inferTags(title, filename),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			currentContent = []string{}
		} else if strings.HasPrefix(line, "- ") && currentRule != nil {
			// 收集规则内容
			currentContent = append(currentContent, line)
		}
	}
	
	// 保存最后一个规则
	if currentRule != nil && len(currentContent) > 0 {
		currentRule.Content = strings.Join(currentContent, "\n")
		rules = append(rules, *currentRule)
	}
	
	return rules
}

// inferRuleType 根据标题推断规则类型
func (l *FileLoader) inferRuleType(title string) string {
	title = strings.ToLower(title)
	
	switch {
	case strings.Contains(title, "安全"):
		return "security"
	case strings.Contains(title, "性能"):
		return "performance"
	case strings.Contains(title, "代码"):
		return "code_style"
	case strings.Contains(title, "测试"):
		return "testing"
	case strings.Contains(title, "部署"):
		return "deployment"
	case strings.Contains(title, "数据库"):
		return "database"
	case strings.Contains(title, "缓存"):
		return "cache"
	default:
		return "general"
	}
}

// inferTags 根据标题和文件名推断标签
func (l *FileLoader) inferTags(title, filename string) []string {
	var tags []string
	
	// 从文件名推断技术栈
	filename = strings.ToLower(filename)
	switch {
	case strings.Contains(filename, "php"):
		tags = append(tags, "php")
	case strings.Contains(filename, "go"):
		tags = append(tags, "go")
	case strings.Contains(filename, "java"):
		tags = append(tags, "java")
	case strings.Contains(filename, "python"):
		tags = append(tags, "python")
	case strings.Contains(filename, "nodejs"):
		tags = append(tags, "nodejs")
	case strings.Contains(filename, "frontend"):
		tags = append(tags, "frontend")
	case strings.Contains(filename, "database"):
		tags = append(tags, "database")
	case strings.Contains(filename, "cache"):
		tags = append(tags, "cache")
	case strings.Contains(filename, "devops"):
		tags = append(tags, "devops")
	}
	
	// 从标题推断类型
	title = strings.ToLower(title)
	switch {
	case strings.Contains(title, "安全"):
		tags = append(tags, "security")
	case strings.Contains(title, "性能"):
		tags = append(tags, "performance")
	case strings.Contains(title, "代码"):
		tags = append(tags, "code_style")
	case strings.Contains(title, "测试"):
		tags = append(tags, "testing")
	}
	
	return tags
}

// getProjectTechStacks 获取项目技术栈信息
func (l *FileLoader) getProjectTechStacks() []string {
	techStackPath := filepath.Join(l.basePath, "project", "tech_stack.yaml")

	if _, err := os.Stat(techStackPath); os.IsNotExist(err) {
		return []string{}
	}

	techStackData, err := os.ReadFile(techStackPath)
	if err != nil {
		return []string{}
	}

	var techStack map[string]interface{}
	if err := yaml.Unmarshal(techStackData, &techStack); err != nil {
		return []string{}
	}

	return getStringSlice(techStack, "tech_stacks")
}

// generateTechSpecificRules 根据技术栈生成相应的规避规则
func (l *FileLoader) generateTechSpecificRules(techStacks []string) []Rule {
	var rules []Rule

	for _, tech := range techStacks {
		switch {
		case strings.Contains(tech, "PHP"):
			rules = append(rules, l.generatePHPRules(tech)...)
		case strings.Contains(tech, "Go"):
			rules = append(rules, l.generateGoRules(tech)...)
		case strings.Contains(tech, "Java"):
			rules = append(rules, l.generateJavaRules(tech)...)
		case strings.Contains(tech, "Python"):
			rules = append(rules, l.generatePythonRules(tech)...)
		case strings.Contains(tech, "Node.js"):
			rules = append(rules, l.generateNodeRules(tech)...)
		case strings.Contains(tech, "React") || strings.Contains(tech, "Vue"):
			rules = append(rules, l.generateFrontendRules(tech)...)
		case strings.Contains(tech, "MySQL") || strings.Contains(tech, "PostgreSQL"):
			rules = append(rules, l.generateDatabaseRules(tech)...)
		case strings.Contains(tech, "Redis") || strings.Contains(tech, "Memcached"):
			rules = append(rules, l.generateCacheRules(tech)...)
		case strings.Contains(tech, "Docker") || strings.Contains(tech, "Kubernetes"):
			rules = append(rules, l.generateDevOpsRules(tech)...)
		}
	}

	return rules
}

// generatePHPRules 生成PHP相关的规避规则
func (l *FileLoader) generatePHPRules(tech string) []Rule {
	rules := []Rule{
		{
			Title:       "PHP安全规范",
			Description: "PHP开发中的安全注意事项和规避规则",
			Type:        "security",
			Content:     "使用 PDO 预处理语句防止 SQL 注入，验证所有用户输入，使用 password_hash() 加密密码",
			Priority:    5,
			Enabled:     true,
			Tags:        []string{"security", "php", "sql-injection"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "PHP性能优化",
			Description: "PHP性能优化的关键规则",
			Type:        "performance",
			Content:     "启用 OPcache，使用 Composer 自动加载，避免在循环中执行数据库查询",
			Priority:    4,
			Enabled:     true,
			Tags:        []string{"performance", "php", "optimization"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// 根据具体框架添加特定规则
	if strings.Contains(tech, "Laravel") {
		rules = append(rules, Rule{
			Title:       "Laravel最佳实践",
			Description: "Laravel框架开发的最佳实践",
			Type:        "framework",
			Content:     "使用 Eloquent ORM，遵循 MVC 模式，使用 Artisan 命令，启用 CSRF 保护",
			Priority:    4,
			Enabled:     true,
			Tags:        []string{"framework", "laravel", "best-practices"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return rules
}

// generateGoRules 生成Go相关的规避规则
func (l *FileLoader) generateGoRules(tech string) []Rule {
	return []Rule{
		{
			Title:       "Go代码规范",
			Description: "Go语言开发的标准规范",
			Type:        "code_style",
			Content:     "使用 gofmt 格式化代码，遵循 Go 官方命名约定，使用 go mod 管理依赖",
			Priority:    4,
			Enabled:     true,
			Tags:        []string{"code_style", "go", "golang"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Go错误处理",
			Description: "Go语言错误处理的最佳实践",
			Type:        "error_handling",
			Content:     "始终检查错误返回值，使用 errors.Wrap 包装错误，避免忽略错误",
			Priority:    5,
			Enabled:     true,
			Tags:        []string{"error_handling", "go", "best-practices"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// generateJavaRules 生成Java相关的规避规则
func (l *FileLoader) generateJavaRules(tech string) []Rule {
	return []Rule{
		{
			Title:       "Java代码规范",
			Description: "Java开发的标准规范",
			Type:        "code_style",
			Content:     "遵循 Java 命名约定，使用 Lombok 减少样板代码，启用代码检查工具",
			Priority:    4,
			Enabled:     true,
			Tags:        []string{"code_style", "java", "spring"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// generatePythonRules 生成Python相关的规避规则
func (l *FileLoader) generatePythonRules(tech string) []Rule {
	return []Rule{
		{
			Title:       "Python代码规范",
			Description: "Python开发的标准规范",
			Type:        "code_style",
			Content:     "遵循 PEP 8 规范，使用类型提示，使用虚拟环境管理依赖",
			Priority:    4,
			Enabled:     true,
			Tags:        []string{"code_style", "python", "pep8"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// generateNodeRules 生成Node.js相关的规避规则
func (l *FileLoader) generateNodeRules(tech string) []Rule {
	return []Rule{
		{
			Title:       "Node.js安全规范",
			Description: "Node.js开发中的安全注意事项",
			Type:        "security",
			Content:     "使用 helmet 中间件，验证所有输入，使用 bcrypt 加密密码，定期更新依赖",
			Priority:    5,
			Enabled:     true,
			Tags:        []string{"security", "nodejs", "express"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// generateFrontendRules 生成前端相关的规避规则
func (l *FileLoader) generateFrontendRules(tech string) []Rule {
	return []Rule{
		{
			Title:       "前端安全规范",
			Description: "前端开发中的安全注意事项",
			Type:        "security",
			Content:     "使用 HTTPS，验证用户输入，防止 XSS 攻击，使用 CSP 策略",
			Priority:    5,
			Enabled:     true,
			Tags:        []string{"security", "frontend", "xss"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// generateDatabaseRules 生成数据库相关的规避规则
func (l *FileLoader) generateDatabaseRules(tech string) []Rule {
	return []Rule{
		{
			Title:       "数据库安全规范",
			Description: "数据库操作的安全注意事项",
			Type:        "security",
			Content:     "使用参数化查询防止 SQL 注入，限制数据库用户权限，定期备份数据",
			Priority:    5,
			Enabled:     true,
			Tags:        []string{"security", "database", "sql-injection"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// generateCacheRules 生成缓存相关的规避规则
func (l *FileLoader) generateCacheRules(tech string) []Rule {
	return []Rule{
		{
			Title:       "缓存使用规范",
			Description: "缓存系统使用的最佳实践",
			Type:        "performance",
			Content:     "设置合理的过期时间，避免缓存穿透，使用缓存预热，监控缓存命中率",
			Priority:    4,
			Enabled:     true,
			Tags:        []string{"performance", "cache", "redis"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// generateDevOpsRules 生成DevOps相关的规避规则
func (l *FileLoader) generateDevOpsRules(tech string) []Rule {
	return []Rule{
		{
			Title:       "容器安全规范",
			Description: "容器化部署的安全注意事项",
			Type:        "security",
			Content:     "使用非 root 用户运行容器，定期更新基础镜像，扫描镜像漏洞，限制容器权限",
			Priority:    5,
			Enabled:     true,
			Tags:        []string{"security", "docker", "kubernetes"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// LoadTemplateRules 加载模板规则
func (l *FileLoader) LoadTemplateRules() ([]Rule, error) {
	templatesDir := filepath.Join(l.basePath, "templates")

	// 检查模板目录是否存在
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return []Rule{}, nil
	}

	// 目前返回空规则，用户可以根据需要添加自定义模板
	return []Rule{}, nil
}

// LoadMetadata 加载元数据
func (l *FileLoader) LoadMetadata() (*Metadata, error) {
	projectDir := filepath.Join(l.basePath, "project")
	techStackPath := filepath.Join(projectDir, "tech_stack.yaml")

	// 读取技术栈信息
	techStackData, err := os.ReadFile(techStackPath)
	if err != nil {
		return nil, fmt.Errorf("读取技术栈文件失败: %w", err)
	}

	var techStack map[string]interface{}
	if err := yaml.Unmarshal(techStackData, &techStack); err != nil {
		return nil, fmt.Errorf("解析技术栈文件失败: %w", err)
	}

	// 读取配置文件
	configPath := filepath.Join(l.basePath, "config.yaml")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	metadata := &Metadata{
		ProjectName:   getString(techStack, "project_name", "未知项目"),
		TechStacks:    getStringSlice(techStack, "tech_stacks"),
		AIEditors:     getStringSlice(techStack, "ai_editors"),
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Version:       "1.0.0",
	}

	return metadata, nil
}

// LoadAllRules 加载所有规则
func (l *FileLoader) LoadAllRules() (*RuleSet, error) {
	projectRules, err := l.LoadProjectRules()
	if err != nil {
		return nil, fmt.Errorf("加载项目规则失败: %w", err)
	}

	globalRules, err := l.LoadGlobalRules()
	if err != nil {
		return nil, fmt.Errorf("加载全局规则失败: %w", err)
	}

	templateRules, err := l.LoadTemplateRules()
	if err != nil {
		return nil, fmt.Errorf("加载模板规则失败: %w", err)
	}

	metadata, err := l.LoadMetadata()
	if err != nil {
		return nil, fmt.Errorf("加载元数据失败: %w", err)
	}

	ruleSet := &RuleSet{
		ProjectRules:  projectRules,
		GlobalRules:   globalRules,
		TemplateRules: templateRules,
		Metadata:      *metadata,
	}

	return ruleSet, nil
}

// 辅助函数
func getString(data map[string]interface{}, key, defaultValue string) string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getStringSlice(data map[string]interface{}, key string) []string {
	if value, exists := data[key]; exists {
		if slice, ok := value.([]interface{}); ok {
			result := make([]string, len(slice))
			for i, item := range slice {
				if str, ok := item.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}
	return []string{}
}
