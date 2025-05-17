package cmd

import (
	"fmt"
	"os"
	"strings"

	"errors"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 支持的配置键列表
var validConfigKeys = []string{
	"auto-open",          // 是否自动打开浏览器
	"enable-dir-listing", // 是否启用目录列表功能
	"theme",              // 目录列表主题
	"language",           // 界面语言
	// 在这里添加其他支持的配置键
}

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

		// 创建表格
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		// 设置表格样式
		t.SetStyle(table.StyleColoredBright)

		// 设置表头 - 确保使用当前语言的翻译
		UpdateConfigTableHeaders(t)

		// 添加配置信息行 - 确保使用当前语言的翻译
		t.AppendRows([]table.Row{
			{"auto-open", formatBoolValue(cfg.AutoOpen), i18n.T("config.auto_open_desc")},
			{"enable-dir-listing", formatBoolValue(cfg.EnableDirListing), i18n.T("config.enable_dir_listing_desc")},
			{"theme", cfg.Theme, i18n.T("config.theme_desc")},
			{"language", formatLanguageValue(cfg.Language), i18n.T("config.language_desc")},
		})

		// 设置列对齐方式
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMax: 30},
			{Number: 2, Align: text.AlignCenter, AlignHeader: text.AlignCenter, WidthMax: 20},
			{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignCenter},
		})

		// 输出表格
		t.Render()

		return nil
	},
}

// 生成配置命令缺少参数的友好错误信息
func generateConfigCommandHelp(cmdName string, args []string) string {
	var msg strings.Builder

	if cmdName == "get" {
		msg.WriteString(i18n.T("cmd.get.missing_arg") + "\n\n")
	} else if cmdName == "set" {
		if len(args) == 0 {
			msg.WriteString(i18n.T("cmd.set.missing_key_value") + "\n\n")
		} else {
			msg.WriteString(i18n.T("cmd.set.missing_value") + "\n\n")
		}
	}

	msg.WriteString(i18n.T("cmd.available_items") + "\n")
	for _, key := range validConfigKeys {
		fmt.Fprintf(&msg, "  - %s\n", key)
	}

	msg.WriteString("\n" + i18n.T("cmd.usage") + "\n")
	if cmdName == "get" {
		msg.WriteString("  " + i18n.T("cmd.get.usage") + "\n\n")
		msg.WriteString(i18n.T("cmd.examples") + "\n")
		msg.WriteString("  " + i18n.T("cmd.get.example1") + "\n")
		msg.WriteString("  " + i18n.T("cmd.get.example2") + "\n")
		msg.WriteString("  " + i18n.T("cmd.get.example3") + "\n")
		msg.WriteString("  " + i18n.T("cmd.get.example4") + "\n")
	} else if cmdName == "set" {
		msg.WriteString("  " + i18n.T("cmd.set.usage") + "\n\n")
		msg.WriteString(i18n.T("cmd.examples") + "\n")
		msg.WriteString("  " + i18n.T("cmd.set.example1") + "\n")
		msg.WriteString("  " + i18n.T("cmd.set.example2") + "\n")
		msg.WriteString("  " + i18n.T("cmd.set.example3") + "\n")
		msg.WriteString("  " + i18n.T("cmd.set.example4") + "\n")

		if len(args) == 1 {
			msg.WriteString("\n" + i18n.T("cmd.provided_item") + args[0] + "\n")
			if args[0] == "theme" {
				msg.WriteString(i18n.T("cmd.theme.options") + "\n")
			} else if args[0] == "auto-open" || args[0] == "enable-dir-listing" {
				msg.WriteString(i18n.T("cmd.bool.options") + "\n")
			} else if args[0] == "language" {
				msg.WriteString(i18n.T("cmd.language.options") + "\n")
			}
		}
	}

	return msg.String()
}

// configGetCmd 获取指定配置
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: i18n.T("cmd.config.get.short"),
	Long:  i18n.T("cmd.config.get.long"),
	// 不使用标准参数验证，改用自定义验证
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			// 直接返回原始错误字符串，不需要格式化，避免fmt错误信息重复
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
		if len(args) != 2 {
			// 直接返回原始错误字符串，不需要格式化，避免fmt错误信息重复
			return errors.New(generateConfigCommandHelp("set", args))
		}
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

			// 语言改变后，也需要重新加载表格表头的翻译
			// 这里不需要显式操作，因为表头是在每次configListCmd运行时创建的

			languageDisplayName := i18n.GetLanguageDisplayName(value)
			logger.Info(i18n.Tf("config.language_changed", languageDisplayName))
		} else {
			logger.Info(i18n.Tf("config.item_set", key, value))
		}

		return nil
	},
}

// 检查配置键是否有效
func isValidConfigKey(key string) bool {
	for _, validKey := range validConfigKeys {
		if key == validKey {
			return true
		}
	}
	return false
}

// 生成无效key的友好错误信息
func generateInvalidKeyError(key string) error {
	var msg strings.Builder

	// 不支持的配置项
	fmt.Fprintf(&msg, "%s\n\n", i18n.Tf("error.invalid_config_key", key))

	// 支持的配置项列表
	msg.WriteString(i18n.T("error.available_keys") + "\n")
	for _, validKey := range validConfigKeys {
		fmt.Fprintf(&msg, "  - %s\n", validKey)
	}

	// 添加配置项说明
	msg.WriteString("\n" + i18n.T("error.key_descriptions") + "\n")
	msg.WriteString("  - " + i18n.T("error.auto_open_desc") + "\n")
	msg.WriteString("  - " + i18n.T("error.enable_dir_listing_desc") + "\n")
	msg.WriteString("  - " + i18n.T("error.theme_desc") + "\n")
	msg.WriteString("  - " + i18n.T("error.language_desc") + "\n")

	return fmt.Errorf(msg.String())
}

// 将字符串转换为布尔值，支持多种表示形式
func parseBoolValue(value string) (bool, error) {
	// 转为小写便于比较
	val := strings.ToLower(value)

	switch val {
	case "true", "yes", "y", "1", "on":
		return true, nil
	case "false", "no", "n", "0", "off":
		return false, nil
	default:
		return false, fmt.Errorf(i18n.T("error.invalid_bool"))
	}
}

// 设置配置值（根据类型转换）
func setConfigValue(key, value string) error {
	switch key {
	case "auto-open", "enable-dir-listing":
		// 将输入转换为布尔值
		boolValue, err := parseBoolValue(value)
		if err != nil {
			return err
		}
		viper.Set(key, boolValue)

	case "theme":
		// 验证主题名称是否有效
		validThemes := []string{"default", "dark", "blue", "green", "retro", "json", "table"}
		isValid := false
		for _, theme := range validThemes {
			if value == theme {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf(i18n.Tf("error.invalid_theme", value))
		}
		viper.Set(key, value)

	case "language":
		// 验证语言是否被支持
		if !i18n.IsSupportedLanguage(value) {
			supportedLangs := strings.Join(i18n.GetSupportedLanguages(), ", ")
			return fmt.Errorf(i18n.Tf("error.invalid_language", value, supportedLangs))
		}
		viper.Set(key, value)

		// 语言设置特殊处理：同时更新i18n包的语言设置
		if err := config.SetLanguage(value); err != nil {
			return fmt.Errorf(i18n.Tf("error.cannot_set_language", err))
		}

	default:
		// 这里不应该到达，因为已经在前面验证了key的有效性
		return fmt.Errorf(i18n.Tf("error.unknown_config_item", key))
	}

	return nil
}

// formatBoolValue 格式化布尔值为更友好的显示文本
func formatBoolValue(value bool) string {
	if value {
		return i18n.T("config.enabled")
	}
	return i18n.T("config.disabled")
}

// 格式化语言值，显示友好的语言名称
func formatLanguageValue(lang string) string {
	return i18n.GetLanguageDisplayName(lang)
}

// UpdateConfigTableHeaders 更新配置表格的表头为当前语言
func UpdateConfigTableHeaders(t table.Writer) {
	t.ResetHeaders()

	// 强制重新获取当前语言的翻译
	itemHeader := i18n.T("config.item")
	valueHeader := i18n.T("config.current_value")
	descHeader := i18n.T("config.description")

	t.AppendHeader(table.Row{itemHeader, valueHeader, descHeader})
}

func init() {
	RootCmd.AddCommand(configCmd)

	// 添加子命令
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
