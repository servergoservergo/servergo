package cmd

import (
	"github.com/CC11001100/servergo/pkg/installer"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// 命令行标志
var (
	// 是否卸载（如果为true，则执行卸载操作）
	uninstall bool
)

// installCmd 表示安装命令
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "将servergo安装到系统PATH中",
	Long: `将servergo安装到系统PATH中，以便可以在任何目录下直接使用servergo命令。

支持macOS、Windows和Linux系统。

在macOS和Linux上，默认安装到/usr/local/bin目录（可能需要sudo权限）。
在Windows上，默认安装到用户主目录下的.servergo\bin目录，并将该目录添加到PATH环境变量。

可以使用--uninstall参数卸载servergo。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if uninstall {
			return doUninstall()
		}
		return doInstall()
	},
}

// doInstall 执行安装操作
func doInstall() error {
	logger.Info("开始安装servergo到系统PATH中...")
	if err := installer.InstallToPath(); err != nil {
		logger.Error("安装失败: %v", err)
		return err
	}
	return nil
}

// doUninstall 执行卸载操作
func doUninstall() error {
	logger.Info("开始从系统PATH中卸载servergo...")
	if err := installer.UninstallFromPath(); err != nil {
		logger.Error("卸载失败: %v", err)
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(installCmd)

	// 添加标志
	installCmd.Flags().BoolVarP(&uninstall, "uninstall", "u", false, "卸载servergo（从PATH中移除）")
}
