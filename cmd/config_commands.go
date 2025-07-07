package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd 表示配置相关的命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: i18n.T("cmd.config.short"),
	Long:  i18n.T("cmd.config.long"),
}

// configListCmd 列出所有配置
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("cmd.config.list.short"),
	Long:  i18n.T("cmd.config.list.long"),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化配置
		if err := config.InitConfig(); err != nil {
			return err
		}

		// 获取当前配置
		cfg := config.GetConfig()

		// 确保使用当前语言
		if err := i18n.Init(cfg.Language); err != nil {
			logger.Warning(fmt.Sprintf("Failed to initialize i18n: %v", err))
		}

		// 获取配置文件路径
		cfgPath, err := config.GetConfigFilePath()
		if err != nil {
			return err
		}

		// 检查配置文件是否实际存在
		fileExists := false
		if _, err := os.Stat(cfgPath); err == nil {
			fileExists = true
		}

		// 先单独显示配置文件路径信息
		logger.Info(i18n.T("config.current_config"))
		if fileExists {
			logger.Info(i18n.Tf("config.file_path", cfgPath))
		} else {
			logger.Info(i18n.T("config.file_not_created"))
			logger.Info(i18n.Tf("config.default_path", cfgPath))
		}
		logger.Info("")

		// 显示配置表格
		displayConfigTable(cfg)

		return nil
	},
}

// configGetCmd 获取指定配置
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: i18n.T("cmd.config.get.short"),
	Long:  i18n.T("cmd.config.get.long"),
	// 不使用标准参数验证，改用自定义验证
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// 没有提供配置项，列出所有可用的配置项
			fmt.Println(i18n.T("cmd.available_items"))
			for _, key := range validConfigKeys {
				fmt.Println("  -", key)
			}
			os.Exit(0)
		} else if len(args) > 1 {
			// 提供了过多的参数
			return errors.New(generateConfigCommandHelp("get", args))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化配置
		if err := config.InitConfig(); err != nil {
			return err
		}

		key := args[0]

		// 验证key是否合法
		if !isValidConfigKey(key) {
			return generateInvalidKeyError(key)
		}

		// 获取配置值
		value := viper.Get(key)
		if value == nil {
			return fmt.Errorf(i18n.Tf("error.config_item_not_exist", key))
		}

		// 直接输出原始值，不做格式转换
		fmt.Println(value)
		return nil
	},
}

// configSetCmd 设置指定配置
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: i18n.T("cmd.config.set.short"),
	Long:  i18n.T("cmd.config.set.long"),
	// 不使用标准参数验证，改用自定义验证
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// 没有提供任何参数，显示常规帮助信息
			return errors.New(generateConfigCommandHelp("set", args))
		} else if len(args) == 1 {
			// 只提供了配置项名称，但没有提供值
			key := args[0]
			if !isValidConfigKey(key) {
				return generateInvalidKeyError(key)
			}

			// 根据配置项类型，直接显示可选的值并退出程序
			switch key {
			case "theme":
				// 直接显示所有可用主题
				fmt.Println(i18n.T("cmd.theme.options"))
				for _, theme := range dirlist.GetSupportedThemes() {
					fmt.Println("  -", theme)
				}
				os.Exit(0)
			case "auto-open", "enable-dir-listing", "enable-log-persistence":
				fmt.Println(i18n.T("cmd.bool.options"))
				fmt.Println("  - true, yes, y, 1, on")
				fmt.Println("  - false, no, n, 0, off")
				os.Exit(0)
			case "language":
				fmt.Println(i18n.T("cmd.language.options"))
				langs := i18n.GetSupportedLanguages()
				for _, lang := range langs {
					// 尝试获取该语言的显示名称
					displayName := i18n.GetLanguageDisplayName(lang)
					fmt.Printf("  - %s (%s)\n", lang, displayName)
				}
				os.Exit(0)
			}
		} else if len(args) > 2 {
			// 参数太多
			return errors.New(generateConfigCommandHelp("set", []string{args[0]}))
		}

		// 正确提供了两个参数
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化配置
		if err := config.InitConfig(); err != nil {
			return err
		}

		key := args[0]
		value := args[1]

		// 验证key是否合法
		if !isValidConfigKey(key) {
			return generateInvalidKeyError(key)
		}

		// 设置配置值
		if err := setConfigValue(key, value); err != nil {
			return err
		}

		// 获取完整的配置对象
		cfg := config.GetConfig()

		// 保存配置
		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf(i18n.Tf("error.cannot_save_config", err))
		}

		// 特殊处理语言变化的消息
		if key == "language" {
			if err := i18n.Init(value); err != nil {
				logger.Warning(i18n.Tf("config.language_init_error", err))
			}

			// 语言改变后，更新所有命令描述
			UpdateCommandDescriptions()

			languageDisplayName := i18n.GetLanguageDisplayName(value)
			logger.Info(i18n.Tf("config.language_changed", languageDisplayName))
		} else {
			logger.Info(i18n.Tf("config.item_set", key, value))
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	// 添加子命令
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
