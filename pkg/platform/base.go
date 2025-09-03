package platform

import (
	"github/pfinal/pf_ruler/pkg/rules"
)

// PlatformAdapter 平台适配器接口
// 支持通过"适配器接口"新增 AI 编辑器平台，无需修改核心逻辑
type PlatformAdapter interface {
	// Name 返回平台名称（如 "trae"、"cursor"）
	Name() string
	
	// DefaultOutputPath 返回平台规则默认输出路径
	// 如 ".trae/rules/project_rules.md"、".cursor/rules.json"
	DefaultOutputPath() string
	
	// Convert 将统一规则转换为平台格式
	Convert(ruleSet *rules.RuleSet) ([]byte, error)
}

// PlatformRegistry 平台注册表
type PlatformRegistry struct {
	adapters map[string]PlatformAdapter
}

// NewPlatformRegistry 创建新的平台注册表
func NewPlatformRegistry() *PlatformRegistry {
	return &PlatformRegistry{
		adapters: make(map[string]PlatformAdapter),
	}
}

// Register 注册平台适配器
func (r *PlatformRegistry) Register(adapter PlatformAdapter) {
	r.adapters[adapter.Name()] = adapter
}

// Get 获取平台适配器
func (r *PlatformRegistry) Get(name string) (PlatformAdapter, bool) {
	adapter, exists := r.adapters[name]
	return adapter, exists
}

// ListSupported 列出所有支持的平台
func (r *PlatformRegistry) ListSupported() []string {
	platforms := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		platforms = append(platforms, name)
	}
	return platforms
}
