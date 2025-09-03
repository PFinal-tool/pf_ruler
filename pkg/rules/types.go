package rules

import (
	"time"
)

// RuleSet 统一规则集合
type RuleSet struct {
	// 项目规则（最高优先级）
	ProjectRules []Rule `yaml:"project_rules" json:"project_rules"`

	// 全局规则（次优先级）
	GlobalRules []Rule `yaml:"global_rules" json:"global_rules"`

	// 模板规则（可选，需用户配置）
	TemplateRules []Rule `yaml:"template_rules" json:"template_rules"`

	// 元数据
	Metadata Metadata `yaml:"metadata" json:"metadata"`
}

// Rule 单条规则
type Rule struct {
	// 规则标题
	Title string `yaml:"title" json:"title"`

	// 规则描述
	Description string `yaml:"description" json:"description"`

	// 规则类型（如：naming, security, performance, style）
	Type string `yaml:"type" json:"type"`

	// 规则内容
	Content string `yaml:"content" json:"content"`

	// 规则优先级（1-5，5为最高）
	Priority int `yaml:"priority" json:"priority"`

	// 是否启用
	Enabled bool `yaml:"enabled" json:"enabled"`

	// 标签（用于分类和搜索）
	Tags []string `yaml:"tags" json:"tags"`

	// 创建时间
	CreatedAt time.Time `yaml:"created_at" json:"created_at"`

	// 更新时间
	UpdatedAt time.Time `yaml:"updated_at" json:"updated_at"`
}

// Metadata 规则元数据
type Metadata struct {
	// 项目名称
	ProjectName string `yaml:"project_name" json:"project_name"`

	// 技术栈
	TechStacks []string `yaml:"tech_stacks" json:"tech_stacks"`

	// 目标AI编辑器
	AIEditors []string `yaml:"ai_editors" json:"ai_editors"`

	// 创建时间
	CreatedAt time.Time `yaml:"created_at" json:"created_at"`

	// 最后更新时间
	LastUpdatedAt time.Time `yaml:"last_updated_at" json:"last_updated_at"`

	// 版本
	Version string `yaml:"version" json:"version"`
}

// Loader 规则加载器接口
type Loader interface {
	// LoadProjectRules 加载项目规则
	LoadProjectRules() ([]Rule, error)

	// LoadGlobalRules 加载全局规则
	LoadGlobalRules() ([]Rule, error)

	// LoadTemplateRules 加载模板规则
	LoadTemplateRules() ([]Rule, error)

	// LoadMetadata 加载元数据
	LoadMetadata() (*Metadata, error)
}

// FileLoader 基于文件的规则加载器
type FileLoader struct {
	basePath string
}

// NewFileLoader 创建新的文件加载器
func NewFileLoader(basePath string) *FileLoader {
	return &FileLoader{
		basePath: basePath,
	}
}
