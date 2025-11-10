/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// 定义彩色输出函数
var (
	green   = color.New(color.FgGreen).PrintlnFunc()
	cyan    = color.New(color.FgCyan).PrintlnFunc()
	yellow  = color.New(color.FgYellow).PrintlnFunc()
	red     = color.New(color.FgRed).PrintlnFunc()
	magenta = color.New(color.FgMagenta).PrintlnFunc()

	greenBold  = color.New(color.FgGreen, color.Bold).PrintlnFunc()
	cyanBold   = color.New(color.FgCyan, color.Bold).PrintlnFunc()
	yellowBold = color.New(color.FgYellow, color.Bold).PrintlnFunc()
	redBold    = color.New(color.FgRed, color.Bold).PrintlnFunc()
)

// 显示 pfinalclub logo
func showLogo() {
	logo := `
    ____  __  __  ___   __  _______  _   _ 
   |  _ \|  \/  |/ _ \ / / |__   __|| \ | |
   | |_) | \  / | | | / /     | |   |  \| |
   |  __/| |\/| | |_| \ \     | |   | . ` + "`" + ` |
   |_|   |_|  |_|\___/  \_\    |_|   |_|\_|
  =======================================
  `
	color.New(color.FgMagenta, color.Bold).Print(logo)
	color.New(color.FgBlue, color.Bold).Println("       AI编辑器规则统一管理工具 v1.1       ")
	color.New(color.FgBlue, color.Bold).Println("    https://github.com/pfinal/pf_ruler    ")
	color.New(color.FgBlue, color.Bold).Println("  =======================================")
	color.New(color.Reset).Println()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pf_ruler",
	Short: "pf_ruler - AI编辑器规则统一管理工具",
	Long: `pf_ruler 是一款基于 Go 语言的命令行工具，用于统一管理各 AI 编辑器的规则配置。

` + color.GreenString("核心功能：") + `
  - 一键初始化规则目录结构
  - 交互式收集项目需求
  - 跨平台规则生成（支持 Trae、Cursor 等）

` + color.CyanString("使用示例：") + `
  ` + color.YellowString("pf_ruler init") + `       # 初始化规则目录
  ` + color.YellowString("pf_ruler generate") + `   # 生成默认平台规则
  ` + color.YellowString("pf_ruler generate --platform=cursor") + `  # 生成指定平台规则
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		// 只有在不是帮助命令时才显示logo
		if len(args) == 0 || args[0] != "help" {
			showLogo()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		redBold("❌ 执行命令时发生错误：", err)
		os.Exit(1)
	}
}

func init() {
	// 配置根命令
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
