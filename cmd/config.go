package cmd

import (
	"fmt"
	"strings"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 支持的配置键列表
var validConfigKeys = []string{
	"default_port", // 默认端口
	"default_dir",  // 默认目录
	"auto_open",    // 是否自动打开浏览器
	// 在这里添加其他支持的配置键
}

// configCmd 表示配置相关的命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理ServerGo配置",
	Long: `管理ServerGo的持久化配置。
可以设置默认端口、默认目录等配置，这些配置将保存在用户主目录下的.servergo目录中。
设置的配置将在未指定相应选项时使用。

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
		fmt.Println("当前配置:")
		fmt.Println("====================")
		fmt.Printf("配置文件路径: %s\n", cfgPath)
		fmt.Println("---")
		fmt.Printf("default_port = %d\n", cfg.DefaultPort)
		fmt.Printf("default_dir = %s\n", cfg.DefaultDir)
		fmt.Printf("auto_open = %t\n", cfg.AutoOpen)
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
		fmt.Printf("%v\n", value)
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

		fmt.Printf("配置项 '%s' 已设置为 '%s'\n", key, value)
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
	errorMsg += "  - default_port: 默认端口，0表示自动探测\n"
	errorMsg += "  - default_dir: 默认目录路径\n"
	errorMsg += "  - auto_open: 启动服务器后是否自动打开浏览器，可接受的值: true/false, yes/no, 1/0\n"
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
	case "default_port":
		// 尝试将值转换为整数
		var portValue int
		if _, err := fmt.Sscanf(value, "%d", &portValue); err != nil {
			return fmt.Errorf("端口必须是一个整数: %v", err)
		}
		viper.Set(key, portValue)

	case "default_dir":
		// 目录值直接设置为字符串
		viper.Set(key, value)

	case "auto_open":
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
