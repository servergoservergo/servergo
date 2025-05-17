package cmd

import (
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/CC11001100/servergo/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd 显示当前版本信息和项目信息
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: i18n.T("cmd.version.short"),
	Long:  i18n.T("cmd.version.long"),
	Run: func(cmd *cobra.Command, args []string) {
		// 获取本地化的未知字符串
		localizedBuildTime := version.BuildTime
		localizedGitCommit := version.GitCommit
		if version.BuildTime == "unknown" {
			localizedBuildTime = i18n.T("version.unknown_date")
		}
		if version.GitCommit == "unknown" {
			localizedGitCommit = i18n.T("version.unknown_commit")
		}

		logger.Info(i18n.T("version.title"))
		logger.Info("==============")
		logger.Info(i18n.Tf("version.version", version.Version))
		logger.Info(i18n.Tf("version.build_date", localizedBuildTime))
		logger.Info(i18n.Tf("version.git_commit", localizedGitCommit))
		if version.GitRef != "unknown" {
			logger.Info(i18n.Tf("version.git_ref", version.GitRef))
		}
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
