package cmd

import (
	"fmt"
	"os"
	"strings"

	"errors"

	"github.com/CC11001100/servergo/pkg/config"
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
	// 在这里添加其他支持的配置键
}

// configCmd 表示配置相关的命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理ServerGo配置",
	Long: `管理ServerGo的持久化配置。
可以设置如是否自动打开浏览器等配置，这些配置将保存在用户主目录下的.servergo目录中。

支持以下子命令:
  list - 列出所有配置
  get  - 获取指定配置的值
  set  - 设置指定配置的值`,
}

// configListCmd 列出所有配置
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置",
	Long:  `列出ServerGo的所有配置项及其当前值。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化配置
		if err := config.InitConfig(); err != nil {
			return err
		}

		// 获取当前配置
		cfg := config.GetConfig()

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
		logger.Info("ServerGo 当前配置")
		if fileExists {
			logger.Info("配置文件路径: %s", cfgPath)
		} else {
			logger.Info("配置文件尚未创建，当前使用默认配置")
			logger.Info("默认配置文件路径将为: %s", cfgPath)
		}
		logger.Info("")

		// 创建表格
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		// 设置表格样式
		t.SetStyle(table.StyleColoredBright)

		// 设置表头
		t.AppendHeader(table.Row{"配置项", "当前值", "说明"})

		// 添加配置信息行
		t.AppendRows([]table.Row{
			{"auto-open", formatBoolValue(cfg.AutoOpen), "启动服务器后是否自动打开浏览器"},
			{"enable-dir-listing", formatBoolValue(cfg.EnableDirListing), "是否启用目录列表功能"},
			{"theme", cfg.Theme, "目录列表主题"},
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
		msg.WriteString("缺少参数：请提供一个配置项名称\n\n")
	} else if cmdName == "set" {
		if len(args) == 0 {
			msg.WriteString("缺少参数：请提供配置项名称和值\n\n")
		} else {
			msg.WriteString("缺少参数：请提供配置项的值\n\n")
		}
	}

	msg.WriteString("可用的配置项:\n")
	for _, key := range validConfigKeys {
		fmt.Fprintf(&msg, "  - %s\n", key)
	}

	msg.WriteString("\n使用说明:\n")
	if cmdName == "get" {
		msg.WriteString("  servergo config get <配置项>\n\n")
		msg.WriteString("示例:\n")
		msg.WriteString("  servergo config get auto-open\n")
		msg.WriteString("  servergo config get theme\n")
		msg.WriteString("  servergo config get enable-dir-listing\n")
	} else if cmdName == "set" {
		msg.WriteString("  servergo config set <配置项> <值>\n\n")
		msg.WriteString("示例:\n")
		msg.WriteString("  servergo config set auto-open false       # 关闭自动打开浏览器\n")
		msg.WriteString("  servergo config set theme dark            # 设置暗色主题\n")
		msg.WriteString("  servergo config set enable-dir-listing true   # 启用目录列表\n")

		if len(args) == 1 {
			msg.WriteString("\n您提供的配置项: " + args[0] + "\n")
			if args[0] == "theme" {
				msg.WriteString("可选主题: default, dark, blue, green, retro, json, table\n")
			} else if args[0] == "auto-open" || args[0] == "enable-dir-listing" {
				msg.WriteString("可接受的布尔值: true/false, yes/no, 1/0\n")
			}
		}
	}

	return msg.String()
}

// configGetCmd 获取指定配置
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "获取指定配置",
	Long:  `获取指定配置项的当前值。`,
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
			return fmt.Errorf("配置项 '%s' 不存在", key)
		}

		// 直接输出原始值，不做格式转换
		fmt.Println(value)
		return nil
	},
}

// configSetCmd 设置指定配置
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "设置指定配置",
	Long:  `设置指定配置项的值。`,
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
			return fmt.Errorf("无法保存配置: %v", err)
		}

		logger.Info("配置项 '%s' 已设置为 '%s'", key, value)
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
	errorMsg := fmt.Sprintf("不支持的配置项: '%s'\n\n支持的配置项有:\n", key)
	for _, validKey := range validConfigKeys {
		errorMsg += fmt.Sprintf("  - %s\n", validKey)
	}

	// 添加配置项说明
	errorMsg += "\n配置项说明:\n"
	errorMsg += "  - auto-open: 启动服务器后是否自动打开浏览器，可接受的值: true/false, yes/no, 1/0\n"
	errorMsg += "  - enable-dir-listing: 是否启用目录列表功能，可接受的值: true/false, yes/no, 1/0\n"
	errorMsg += "  - theme: 目录列表主题，可接受的值: default, dark, blue, green, retro, json, table\n"
	// 添加其他配置项的说明...

	return fmt.Errorf(errorMsg)
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
		return false, fmt.Errorf("无法解析为布尔值，支持的值: true/false, yes/no, y/n, 1/0, on/off")
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
			return fmt.Errorf("无效的主题名称: %s\n支持的主题有: default, dark, blue, green, retro, json, table", value)
		}
		viper.Set(key, value)

	default:
		// 这里不应该到达，因为已经在前面验证了key的有效性
		return fmt.Errorf("未知的配置项: %s", key)
	}

	return nil
}

// 格式化布尔值，使其更易读（添加颜色）
func formatBoolValue(value bool) string {
	if value {
		return text.Colors{text.FgGreen, text.Bold}.Sprint("开启")
	}
	return text.Colors{text.FgRed}.Sprint("关闭")
}

func init() {
	RootCmd.AddCommand(configCmd)

	// 添加子命令
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
