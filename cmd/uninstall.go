package cmd

import (
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/installer"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// uninstallCmd 表示卸载命令
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: i18n.T("cmd.uninstall.short"),
	Long:  i18n.T("cmd.uninstall.long"),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info(i18n.T("cmd.uninstall.start"))

		// 调用installer包中的卸载函数
		if err := installer.UninstallFromPath(); err != nil {
			logger.Error(i18n.Tf("cmd.uninstall.failed", err))
			return err
		}

		logger.Info(i18n.T("cmd.uninstall.complete"))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(uninstallCmd)
}
