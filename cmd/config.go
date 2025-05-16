package cmd

import (
	"fmt"
	"strings"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 支持的配置键列表
var validConfigKeys = []string{
	"auto-open",          // 是否自动打开浏览器
	"enable-dir-listing", // 是否启用目录列表功能
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

		// 显示配置信息
		logger.Info("当前配置:")
		logger.Info("====================")
		logger.Info("配置文件路径: %s", cfgPath)
		logger.Info("---")
		logger.Info("auto-open = %t", cfg.AutoOpen)
		logger.Info("enable-dir-listing = %t", cfg.EnableDirListing)
		// 其他配置项...

		return nil
	},
}

// configGetCmd 获取指定配置
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "获取指定配置",
	Long:  `获取指定配置项的当前值。`,
	Args:  cobra.ExactArgs(1),
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

		// 打印配置值
		logger.Info("%v", value)
		return nil
	},
}

// configSetCmd 设置指定配置
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "设置指定配置",
	Long:  `设置指定配置项的值。`,
	Args:  cobra.ExactArgs(2),
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

		// 保存配置
		if err := viper.WriteConfig(); err != nil {
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

	default:
		// 这里不应该到达，因为已经在前面验证了key的有效性
		return fmt.Errorf("未知的配置项: %s", key)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(configCmd)

	// 添加子命令
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
