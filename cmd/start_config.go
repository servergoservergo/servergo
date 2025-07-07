package cmd

import (
	"fmt"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// processStartConfig 处理start命令的配置参数
func processStartConfig(cmd *cobra.Command) error {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		return err
	}

	// 读取配置中的默认值
	cfg := config.GetConfig()

	// 对所有配置项，如果命令行未指定则使用配置文件中的值
	if !cmd.Flags().Changed("open") {
		autoOpen = cfg.AutoOpen
	}
	if !cmd.Flags().Changed("dir-list") {
		enableDirListing = cfg.EnableDirListing
	}
	if !cmd.Flags().Changed("theme") {
		theme = cfg.Theme
	}
	if !cmd.Flags().Changed("language") {
		if err := i18n.Init(cfg.Language); err != nil {
			logger.Warning(fmt.Sprintf("初始化语言失败: %v", err))
		}
	}
	if !cmd.Flags().Changed("enable-log-persistence") {
		enableLogPersistence = cfg.EnableLogPersistence
	}

	return nil
}

// processLogConfig 处理日志相关配置
func processLogConfig(cmd *cobra.Command) error {
	// 设置日志级别
	if cmd.Flags().Changed("log-level") {
		switch logLevel {
		case "debug":
			logger.Default.SetLevel(logger.DEBUG)
		case "info":
			logger.Default.SetLevel(logger.INFO)
		case "warn", "warning":
			logger.Default.SetLevel(logger.WARNING)
		case "error":
			logger.Default.SetLevel(logger.ERROR)
		default:
			logger.Warning(i18n.Tf("error.invalid_log_level", logLevel))
			logger.Default.SetLevel(logger.INFO) // 默认使用INFO级别
		}
	}

	// 处理日志持久化设置
	if cmd.Flags().Changed("enable-log-persistence") {
		// 如果命令行指定了日志持久化设置，则更新配置
		cfg := config.GetConfig()
		cfg.EnableLogPersistence = enableLogPersistence

		// 保存更新后的配置
		if err := config.SaveConfig(cfg); err != nil {
			logger.Warning(i18n.Tf("error.save_config_failed", err))
		}

		// 重新初始化日志系统
		logConfig := logger.LogConfig{
			Level:         logger.Default.GetLevel(),
			EnableFileLog: enableLogPersistence,
			Filename:      "servergo.log",
		}

		newLogger, err := logger.New(logConfig)
		if err != nil {
			logger.Warning(i18n.Tf("error.logger_init_failed", err))
		} else {
			logger.Default = newLogger
			if enableLogPersistence {
				logger.Info(i18n.T("logger.persistence_enabled"))
			} else {
				logger.Info(i18n.T("logger.persistence_disabled"))
			}
		}
	}

	return nil
}
