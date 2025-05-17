package cmd

import (
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// 版本信息
var (
	Version   = "0.1.0"
	BuildDate = "未知"
	GitCommit = "未知"
)

// versionCmd 显示当前版本信息和项目信息
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: i18n.T("cmd.version.short"),
	Long:  i18n.T("cmd.version.long"),
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info(i18n.T("version.title"))
		logger.Info("==============")
		logger.Info(i18n.Tf("version.version", Version))
		logger.Info(i18n.Tf("version.build_date", BuildDate))
		logger.Info(i18n.Tf("version.git_commit", GitCommit))
		logger.Info("")
		logger.Info(i18n.T("version.description_1"))
		logger.Info(i18n.T("version.description_2"))
		logger.Info("")
		logger.Info(i18n.T("version.author_title"))
		logger.Info("--------")
		logger.Info(i18n.T("version.author_name"))
		logger.Info(i18n.T("version.author_github"))
		logger.Info("")
		logger.Info(i18n.T("version.feedback_title"))
		logger.Info("--------")
		logger.Info(i18n.T("version.feedback_text"))
		logger.Info(i18n.T("version.feedback_url"))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
