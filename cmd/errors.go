package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
)

var (
	// 错误消息模式正则表达式
	flagNeedsArgRegex     = regexp.MustCompile(`flag needs an argument: '([^']+)' in -([^\s]+)`)
	longFlagNeedsArgRegex = regexp.MustCompile(`flag needs an argument: --([^\s]+)`)
	unknownFlagRegex      = regexp.MustCompile(`unknown flag: -([^\s]+)`)
	unknownShorthandRegex = regexp.MustCompile(`unknown shorthand flag: '([^']+)' in -([^\s]+)`)
	invalidArgumentRegex  = regexp.MustCompile(`invalid argument "([^"]+)" for "([^"]+)"`)
)

// 为特定标志准备的友好错误消息映射
var flagDescriptions = map[string]struct {
	DescriptionKey   string // 标志描述的i18n键
	ExpectedValueKey string // 期望值类型的i18n键
	OptionsKey       string // 选项列表的i18n键或直接值
	HasOptions       bool   // 是否有固定选项
}{
	"p":          {"flag.port", "flag.number", "", false},
	"port":       {"flag.port", "flag.number", "", false},
	"d":          {"flag.dir", "flag.directory_path", "", false},
	"dir":        {"flag.dir", "flag.directory_path", "", false},
	"o":          {"flag.auto_open", "flag.bool", "flag.bool_options", true},
	"open":       {"flag.auto_open", "flag.bool", "flag.bool_options", true},
	"a":          {"flag.auth_type", "flag.auth_method", "none, basic, token, form", true},
	"auth":       {"flag.auth_type", "flag.auth_method", "none, basic, token, form", true},
	"u":          {"flag.username", "flag.string", "", false},
	"username":   {"flag.username", "flag.string", "", false},
	"w":          {"flag.password", "flag.string", "", false},
	"password":   {"flag.password", "flag.string", "", false},
	"t":          {"flag.token", "flag.string", "", false},
	"token":      {"flag.token", "flag.string", "", false},
	"l":          {"flag.login_page", "flag.bool", "flag.bool_options", true},
	"login-page": {"flag.login_page", "flag.bool", "flag.bool_options", true},
	"i":          {"flag.dir_list", "flag.bool", "flag.bool_options", true},
	"dir-list":   {"flag.dir_list", "flag.bool", "flag.bool_options", true},
	"m":          {"flag.theme", "flag.theme_name", "", true},
	"theme":      {"flag.theme", "flag.theme_name", "", true},
}

// FriendlyErrorMessage 将Cobra的标准错误消息转换为更友好的用户提示，支持国际化
func FriendlyErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	// 提取原始错误信息
	errMsg := err.Error()

	// 匹配：短标志需要参数 (例如 -m)
	if matches := flagNeedsArgRegex.FindStringSubmatch(errMsg); len(matches) == 3 {
		flagChar := matches[1] // 标志字符
		flagName := matches[2] // 完整标志名称

		// 查找标志描述
		if desc, ok := flagDescriptions[flagChar]; ok {
			msg := i18n.Tf("errors.flag_needs_value", i18n.T(desc.DescriptionKey), i18n.T(desc.ExpectedValueKey))

			if desc.HasOptions {
				var options string
				if desc.DescriptionKey == "flag.theme" {
					// 主题选项特殊处理
					options = strings.Join(dirlist.GetSupportedThemes(), ", ")
				} else if desc.OptionsKey == "flag.bool_options" {
					// 布尔值选项特殊处理
					options = i18n.T(desc.OptionsKey)
				} else {
					// 直接使用选项字符串
					options = desc.OptionsKey
				}
				msg += i18n.Tf("errors.available_options", options)
			}

			// 添加示例
			exampleFlag := "-" + flagName
			if len(flagChar) > 1 {
				exampleFlag = "--" + flagChar
			}

			if desc.HasOptions {
				var firstOption string
				if desc.DescriptionKey == "flag.theme" {
					// 使用第一个主题作为示例
					firstOption = dirlist.GetSupportedThemes()[0]
				} else if desc.OptionsKey == "flag.bool_options" {
					// 使用true作为布尔示例
					firstOption = "true"
				} else {
					// 从选项中取第一个作为示例
					firstOption = strings.Split(desc.OptionsKey, ", ")[0]
				}
				msg += i18n.Tf("errors.example", fmt.Sprintf("%s %s", exampleFlag, firstOption))
			} else {
				msg += i18n.Tf("errors.example", fmt.Sprintf("%s %s", exampleFlag, i18n.T("errors.value_placeholder")))
			}

			return msg
		}

		// 如果没有找到具体描述，返回一个通用的友好消息
		return i18n.Tf("errors.general_flag_needs_value", flagChar, flagName)
	}

	// 匹配：长标志需要参数 (例如 --theme)
	if matches := longFlagNeedsArgRegex.FindStringSubmatch(errMsg); len(matches) == 2 {
		flagName := matches[1] // 标志名称

		// 查找标志描述
		if desc, ok := flagDescriptions[flagName]; ok {
			msg := i18n.Tf("errors.flag_needs_value", i18n.T(desc.DescriptionKey), i18n.T(desc.ExpectedValueKey))

			if desc.HasOptions {
				var options string
				if desc.DescriptionKey == "flag.theme" {
					// 主题选项特殊处理
					options = strings.Join(dirlist.GetSupportedThemes(), ", ")
				} else if desc.OptionsKey == "flag.bool_options" {
					// 布尔值选项特殊处理
					options = i18n.T(desc.OptionsKey)
				} else {
					// 直接使用选项字符串
					options = desc.OptionsKey
				}
				msg += i18n.Tf("errors.available_options", options)
			}

			// 添加示例
			exampleFlag := "--" + flagName

			if desc.HasOptions {
				var firstOption string
				if desc.DescriptionKey == "flag.theme" {
					// 使用第一个主题作为示例
					firstOption = dirlist.GetSupportedThemes()[0]
				} else if desc.OptionsKey == "flag.bool_options" {
					// 使用true作为布尔示例
					firstOption = "true"
				} else {
					// 从选项中取第一个作为示例
					firstOption = strings.Split(desc.OptionsKey, ", ")[0]
				}
				msg += i18n.Tf("errors.example", fmt.Sprintf("%s %s", exampleFlag, firstOption))
			} else {
				msg += i18n.Tf("errors.example", fmt.Sprintf("%s %s", exampleFlag, i18n.T("errors.value_placeholder")))
			}

			return msg
		}

		// 如果没有找到具体描述，返回一个通用的友好消息
		return i18n.Tf("errors.general_flag_needs_value_long", flagName)
	}

	// 匹配：未知标志
	if matches := unknownFlagRegex.FindStringSubmatch(errMsg); len(matches) == 2 {
		flagName := matches[1]
		return i18n.Tf("errors.unknown_flag", flagName)
	}

	// 匹配：无效的参数
	if matches := invalidArgumentRegex.FindStringSubmatch(errMsg); len(matches) == 3 {
		arg := matches[1]
		flag := matches[2]

		// 查找标志描述
		if desc, ok := flagDescriptions[flag]; ok && desc.HasOptions {
			var options string
			if desc.DescriptionKey == "flag.theme" {
				// 主题选项特殊处理
				options = strings.Join(dirlist.GetSupportedThemes(), ", ")
			} else if desc.OptionsKey == "flag.bool_options" {
				// 布尔值选项特殊处理
				options = i18n.T(desc.OptionsKey)
			} else {
				// 直接使用选项字符串
				options = desc.OptionsKey
			}
			return i18n.Tf("errors.invalid_flag_value_with_options", flag, arg, options)
		}

		return i18n.Tf("errors.invalid_flag_value", flag, arg)
	}

	// 对于其他错误，保持原样
	return errMsg
}

// GetFriendlyThemeErrorMessage 返回一个关于主题参数的友好错误消息
func GetFriendlyThemeErrorMessage() string {
	themeOptions := strings.Join(dirlist.GetSupportedThemes(), ", ")
	return i18n.Tf("errors.theme_help", themeOptions)
}
