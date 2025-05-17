package cmd

import (
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/installer"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// installCmd 表示安装命令
var installCmd = &cobra.Command{
	Use:   "install",
	Short: i18n.T("cmd.install.short"),
	Long:  i18n.T("cmd.install.long"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doInstall()
	},
}

// doInstall 执行安装操作
func doInstall() error {
	logger.Info(i18n.T("cmd.install.start"))
	if err := installer.InstallToPath(); err != nil {
		logger.Error(i18n.Tf("cmd.install.failed", err))
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(installCmd)
}
