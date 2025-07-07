package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/spf13/viper"
)

// 支持的配置键列表
var validConfigKeys = []string{
	"auto-open",              // 是否自动打开浏览器
	"enable-dir-listing",     // 是否启用目录列表功能
	"theme",                  // 目录列表主题
	"language",               // 界面语言
	"enable-log-persistence", // 是否启用日志持久化
	"start-port",             // 从哪个端口开始递增寻找空闲端口
	// 在这里添加其他支持的配置键
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
	case "auto-open", "enable-dir-listing", "enable-log-persistence":
		// 将输入转换为布尔值
		boolValue, err := parseBoolValue(value)
		if err != nil {
			return err
		}
		viper.Set(key, boolValue)

	case "start-port":
		// 将输入转换为整数
		portValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf(i18n.Tf("error.invalid_port_number", value))
		}
		// 验证端口范围
		if portValue < 0 || portValue > 65535 {
			return fmt.Errorf(i18n.Tf("error.port_out_of_range", portValue))
		}
		viper.Set(key, portValue)

	case "theme":
		// 验证主题名称是否有效
		isValid := false
		for _, theme := range dirlist.GetSupportedThemes() {
			if value == theme {
				isValid = true
				break
			}
		}
		if !isValid {
			// 传入所有支持的主题列表
			themesStr := strings.Join(dirlist.GetSupportedThemes(), ", ")
			return fmt.Errorf(i18n.Tf("error.invalid_theme", value, themesStr))
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
