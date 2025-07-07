package cmd

import (
	"fmt"
	"strings"

	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
)

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
				// 使用全局定义的有效主题列表
				themesStr := strings.Join(dirlist.GetSupportedThemes(), ", ")
				msg.WriteString(i18n.T("cmd.theme.options") + themesStr + "\n")
			} else if args[0] == "auto-open" || args[0] == "enable-dir-listing" || args[0] == "enable-log-persistence" {
				msg.WriteString(i18n.T("cmd.bool.options") + "\n")
			} else if args[0] == "language" {
				// 使用语言模块提供的支持语言列表
				supportedLangs := strings.Join(i18n.GetSupportedLanguages(), ", ")
				msg.WriteString(i18n.T("cmd.language.options") + supportedLangs + "\n")
			} else if args[0] == "start-port" {
				msg.WriteString(i18n.T("cmd.start_port.options") + "\n")
			}
		}
	}

	return msg.String()
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
	msg.WriteString("  - " + i18n.T("error.enable_log_persistence_desc") + "\n")
	msg.WriteString("  - " + i18n.T("error.start_port_desc") + "\n")

	return fmt.Errorf(msg.String())
}
