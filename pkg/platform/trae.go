package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github/pfinal/pf_ruler/pkg/rules"
)

// TraeAdapter Trae平台适配器
type TraeAdapter struct{}

// NewTraeAdapter 创建新的Trae适配器
func NewTraeAdapter() *TraeAdapter {
	return &TraeAdapter{}
}

// Name 返回平台名称
func (t *TraeAdapter) Name() string {
	return "trae"
}

// DefaultOutputPath 返回Trae规则默认输出路径
func (t *TraeAdapter) DefaultOutputPath() string {
	return ".trae/rules/project_rules.md"
}

// Convert 将统一规则转换为Trae格式
func (t *TraeAdapter) Convert(ruleSet *rules.RuleSet) ([]byte, error) {
	var content strings.Builder
	
	// 写入标题
	content.WriteString(fmt.Sprintf("# %s 项目规则集\n\n", ruleSet.Metadata.ProjectName))
	
	// 写入元数据
	content.WriteString("## 项目信息\n\n")
	content.WriteString(fmt.Sprintf("- **项目名称**: %s\n", ruleSet.Metadata.ProjectName))
	content.WriteString(fmt.Sprintf("- **技术栈**: %s\n", strings.Join(ruleSet.Metadata.TechStacks, ", ")))
	content.WriteString(fmt.Sprintf("- **目标AI编辑器**: %s\n", strings.Join(ruleSet.Metadata.AIEditors, ", ")))
	content.WriteString(fmt.Sprintf("- **生成时间**: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("- **版本**: %s\n\n", ruleSet.Metadata.Version))
	
	// 写入项目规则（最高优先级）
	if len(ruleSet.ProjectRules) > 0 {
		content.WriteString("## 项目特定规则\n\n")
		content.WriteString("*这些规则具有最高优先级，适用于当前项目*\n\n")
		
		for _, rule := range ruleSet.ProjectRules {
			if !rule.Enabled {
				continue
			}
			
			content.WriteString(fmt.Sprintf("### %s\n\n", rule.Title))
			content.WriteString(fmt.Sprintf("**类型**: %s  |  **优先级**: %d  |  **标签**: %s\n\n", 
				rule.Type, rule.Priority, strings.Join(rule.Tags, ", ")))
			content.WriteString(fmt.Sprintf("%s\n\n", rule.Description))
			content.WriteString(fmt.Sprintf("**规则内容**:\n%s\n\n", rule.Content))
		}
	}
	
	// 写入全局规则（次优先级）
	if len(ruleSet.GlobalRules) > 0 {
		content.WriteString("## 全局通用规则\n\n")
		content.WriteString("*这些规则适用于所有项目，具有中等优先级*\n\n")
		
		for _, rule := range ruleSet.GlobalRules {
			if !rule.Enabled {
				continue
			}
			
			content.WriteString(fmt.Sprintf("### %s\n\n", rule.Title))
			content.WriteString(fmt.Sprintf("**类型**: %s  |  **优先级**: %d  |  **标签**: %s\n\n", 
				rule.Type, rule.Priority, strings.Join(rule.Tags, ", ")))
			content.WriteString(fmt.Sprintf("%s\n\n", rule.Description))
			content.WriteString(fmt.Sprintf("**规则内容**:\n%s\n\n", rule.Content))
		}
	}
	
	// 写入模板规则（可选）
	if len(ruleSet.TemplateRules) > 0 {
		content.WriteString("## 自定义模板规则\n\n")
		content.WriteString("*这些规则来自用户自定义模板*\n\n")
		
		for _, rule := range ruleSet.TemplateRules {
			if !rule.Enabled {
				continue
			}
			
			content.WriteString(fmt.Sprintf("### %s\n\n", rule.Title))
			content.WriteString(fmt.Sprintf("**类型**: %s  |  **优先级**: %d  |  **标签**: %s\n\n", 
				rule.Type, rule.Priority, strings.Join(rule.Tags, ", ")))
			content.WriteString(fmt.Sprintf("%s\n\n", rule.Description))
			content.WriteString(fmt.Sprintf("**规则内容**:\n%s\n\n", rule.Content))
		}
	}
	
	// 写入使用说明
	content.WriteString("## 使用说明\n\n")
	content.WriteString("本规则集由 pf_ruler 工具自动生成，用于指导 AI 编辑器生成符合项目规范的代码。\n\n")
	content.WriteString("### 规则优先级\n\n")
	content.WriteString("1. **项目特定规则** - 最高优先级，覆盖其他规则\n")
	content.WriteString("2. **全局通用规则** - 中等优先级，适用于所有项目\n")
	content.WriteString("3. **自定义模板规则** - 可选，来自用户配置\n\n")
	content.WriteString("### 更新规则\n\n")
	content.WriteString("如需更新规则，请修改 `.ruler` 目录下的相应文件，然后重新运行 `pf_ruler generate --platform=trae` 命令。\n")
	
	return []byte(content.String()), nil
}

// EnsureOutputDirectory 确保输出目录存在
func (t *TraeAdapter) EnsureOutputDirectory() error {
	outputPath := t.DefaultOutputPath()
	outputDir := filepath.Dir(outputPath)
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	
	return nil
}
