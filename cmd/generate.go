package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github/pfinal/pf_ruler/pkg/platform"
	"github/pfinal/pf_ruler/pkg/rules"
)

var (
	// 命令标志
	platformFlag string
	forceFlag    bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成跨平台规则文件",
	Long: `加载 .ruler 目录中的统一规则，根据用户指定的 AI 编辑器平台，
自动转换为该平台的原生规则格式，并输出到对应目录。

支持平台：
  - trae: 生成 .trae/rules/project_rules.md 文件
  - cursor: 生成 .cursor/rules.json 文件

示例：
  pf_ruler generate                    # 生成默认平台规则
  pf_ruler generate --platform=cursor  # 生成指定平台规则
  pf_ruler generate --platform=cursor --force  # 强制覆盖现有文件
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 平台参数校验
		if err := validatePlatform(); err != nil {
			redBold("❌ 平台参数错误：", err)
			os.Exit(1)
		}

		// 2. 加载统一规则
		ruleSet, err := loadUnifiedRules()
		if err != nil {
			redBold("❌ 加载规则失败：", err)
			os.Exit(1)
		}

		// 3. 跨平台规则转换
		if err := convertAndOutput(ruleSet); err != nil {
			redBold("❌ 规则转换失败：", err)
			os.Exit(1)
		}

		greenBold("✅ 规则生成完成！")
	},
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(generateCmd)

	// 添加标志
	generateCmd.Flags().StringVarP(&platformFlag, "platform", "p", "", "目标平台 (trae, cursor)")
	generateCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "强制覆盖现有文件")
}

// validatePlatform 验证平台参数
func validatePlatform() error {
	// 如果没有指定平台，使用配置文件中的默认平台
	if platformFlag == "" {
		configPath := filepath.Join(".ruler", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			// 这里可以读取配置文件获取默认平台
			// 暂时使用 trae 作为默认平台
			platformFlag = "trae"
		} else {
			platformFlag = "trae" // 默认使用 trae
		}
	}

	// 检查平台是否支持
	supportedPlatforms := []string{"trae", "cursor"}
	for _, platform := range supportedPlatforms {
		if platformFlag == platform {
			return nil
		}
	}

	return fmt.Errorf("不支持的平台 \"%s\"，当前支持：%v", platformFlag, supportedPlatforms)
}

// loadUnifiedRules 加载统一规则
func loadUnifiedRules() (*rules.RuleSet, error) {
	// 检查 .ruler 目录是否存在
	if _, err := os.Stat(".ruler"); os.IsNotExist(err) {
		return nil, fmt.Errorf(".ruler 目录不存在，请先运行 pf_ruler init 命令")
	}

	// 创建规则加载器
	loader := rules.NewFileLoader(".ruler")

	// 加载所有规则
	ruleSet, err := loader.LoadAllRules()
	if err != nil {
		return nil, fmt.Errorf("加载规则失败: %w", err)
	}

	// 统计规则数量
	totalRules := len(ruleSet.ProjectRules) + len(ruleSet.GlobalRules) + len(ruleSet.TemplateRules)
	
	if totalRules == 0 {
		yellowBold("⚠️  未检测到任何规则文件，将使用默认模板")
	} else {
		greenBold(fmt.Sprintf("✅ 已加载规则（项目规则 %d 条 + 全局规则 %d 条）", 
			len(ruleSet.ProjectRules), len(ruleSet.GlobalRules)))
	}

	return ruleSet, nil
}

// convertAndOutput 转换并输出规则
func convertAndOutput(ruleSet *rules.RuleSet) error {
	// 创建平台注册表
	registry := platform.NewPlatformRegistry()

	// 注册支持的平台适配器
	registry.Register(platform.NewTraeAdapter())
	registry.Register(platform.NewCursorAdapter())

	// 获取目标平台适配器
	adapter, exists := registry.Get(platformFlag)
	if !exists {
		return fmt.Errorf("平台适配器不存在: %s", platformFlag)
	}

	// 转换规则
	outputData, err := adapter.Convert(ruleSet)
	if err != nil {
		return fmt.Errorf("规则转换失败: %w", err)
	}

	greenBold(fmt.Sprintf("✅ 已完成 %s 规则格式转换", platformFlag))

	// 确保输出目录存在
	if err := ensureOutputDirectory(adapter); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 输出文件管理
	outputPath := adapter.DefaultOutputPath()
	if err := writeOutputFile(outputPath, outputData); err != nil {
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	return nil
}

// ensureOutputDirectory 确保输出目录存在
func ensureOutputDirectory(adapter platform.PlatformAdapter) error {
	outputPath := adapter.DefaultOutputPath()
	outputDir := filepath.Dir(outputPath)
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	
	return nil
}

// writeOutputFile 写入输出文件
func writeOutputFile(outputPath string, data []byte) error {
	// 检查文件是否已存在
	if _, err := os.Stat(outputPath); err == nil {
		if !forceFlag {
			// 询问用户是否覆盖
			yellowBold(fmt.Sprintf("⚠️  文件已存在: %s", outputPath))
			yellowBold("使用 --force 标志强制覆盖，或手动删除后重试")
			return fmt.Errorf("文件已存在，请使用 --force 标志覆盖")
		}
		
		// 强制覆盖
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			return fmt.Errorf("覆盖文件失败: %w", err)
		}
		
		greenBold(fmt.Sprintf("✅ 已覆盖现有文件: %s", outputPath))
	} else {
		// 新文件创建
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			return fmt.Errorf("创建文件失败: %w", err)
		}
		
		greenBold(fmt.Sprintf("✅ %s 规则已生成: %s", platformFlag, outputPath))
	}
	
	return nil
}
