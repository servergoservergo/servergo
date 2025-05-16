package cmd

import (
	"github.com/CC11001100/servergo/pkg/installer"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// uninstallCmd 表示卸载命令
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "从系统PATH中卸载servergo",
	Long: `从系统PATH中卸载servergo命令，移除之前通过install命令安装的内容。

该命令会根据当前操作系统执行适当的卸载操作：
- 在macOS和Linux上，会删除/usr/local/bin目录下的servergo符号链接（可能需要sudo权限）
- 在Windows上，会删除用户主目录下.servergo\bin中的servergo.exe文件，并从PATH环境变量中移除该目录

注意：此命令只会移除servergo命令本身，不会删除您创建的任何配置文件或数据。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("开始从系统PATH中卸载servergo...")

		// 调用installer包中的卸载函数
		if err := installer.UninstallFromPath(); err != nil {
			logger.Error("卸载失败: %v", err)
			return err
		}

		logger.Info("卸载完成！servergo已从系统PATH中移除")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(uninstallCmd)
}
