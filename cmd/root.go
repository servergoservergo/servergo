package cmd

import (
	"fmt"
	"os"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// 全局配置变量
var (
	port int
	dir  string
)

// RootCmd 表示没有调用子命令时的基础命令
var RootCmd = &cobra.Command{
	Use:   "servergo",
	Short: "快速启动HTTP文件服务器",
	Long: `ServerGo 是一个简单的命令行工具，用于快速启动HTTP文件服务器，
类似于Python的http.server模块，但使用Go实现，提供更好的性能。`,
}

func init() {
	// 在这里我们可以初始化一些全局标志或配置
	// 初始化配置（忽略错误，因为首次运行可能没有配置文件）
	_ = config.InitConfig()

	// 设置默认值
	config.SetDefaults()

	// 初始化国际化支持
	cfg := config.GetConfig()
	if err := i18n.Init(cfg.Language); err != nil {
		logger.Warning(i18n.Tf("i18n.init_failed", err))
	} else {
		// 更新命令的描述为当前语言
		RootCmd.Short = i18n.T("cmd.root.short")
		RootCmd.Long = i18n.T("cmd.root.long")
	}

	// 初始化日志系统
	// 可以根据需要设置日志级别
	logger.Default.SetLevel(logger.INFO)
}

// Execute 添加所有子命令到根命令并执行
func Execute() {
	// 如果没有提供子命令，则修改参数列表，使其包含 "start" 子命令
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "start")
		logger.Debug("未指定子命令，默认运行start命令")
	}

	// 更新所有命令的描述为当前语言
	updateCommandDescriptions()

	// 阻止Cobra默认的错误打印和使用说明打印
	RootCmd.SilenceErrors = true
	RootCmd.SilenceUsage = true

	// 执行命令
	err := RootCmd.Execute()
	if err != nil {
		// 使用fmt直接打印错误到标准错误输出，不使用logger避免日期前缀等
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// updateCommandDescriptions 更新所有命令的描述为当前语言
func updateCommandDescriptions() {
	// 更新根命令
	RootCmd.Short = i18n.T("cmd.root.short")
	RootCmd.Long = i18n.T("cmd.root.long")

	// 更新子命令
	for _, cmd := range RootCmd.Commands() {
		switch cmd.Name() {
		case "config":
			cmd.Short = i18n.T("cmd.config.short")
			cmd.Long = i18n.T("cmd.config.long")
			// 更新config的子命令
			for _, subcmd := range cmd.Commands() {
				switch subcmd.Name() {
				case "list":
					subcmd.Short = i18n.T("cmd.config.list.short")
					subcmd.Long = i18n.T("cmd.config.list.long")
				case "get":
					subcmd.Short = i18n.T("cmd.config.get.short")
					subcmd.Long = i18n.T("cmd.config.get.long")
				case "set":
					subcmd.Short = i18n.T("cmd.config.set.short")
					subcmd.Long = i18n.T("cmd.config.set.long")
				}
			}
		case "start":
			cmd.Short = i18n.T("cmd.start.short")
			cmd.Long = i18n.T("cmd.start.long")
		case "install":
			cmd.Short = i18n.T("cmd.install.short")
			cmd.Long = i18n.T("cmd.install.long")
		case "uninstall":
			cmd.Short = i18n.T("cmd.uninstall.short")
			cmd.Long = i18n.T("cmd.uninstall.long")
		case "version":
			cmd.Short = i18n.T("cmd.version.short")
			cmd.Long = i18n.T("cmd.version.long")
		}
	}
}
