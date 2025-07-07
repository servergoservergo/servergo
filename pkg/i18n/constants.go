// Package i18n 提供国际化支持功能
// 本包实现了应用程序的多语言翻译机制，支持中英文切换
package i18n

import (
	"embed"
)

// 环境变量名称，用于设置语言
const (
	// EnvLanguage 环境变量名称，用于指定应用程序语言
	// 例如：SERVERGO_LANGUAGE=zh-CN
	EnvLanguage = "SERVERGO_LANGUAGE"
)

// 支持的语言常量定义
// 系统目前支持英文(en)和简体中文(zh-CN)两种语言
const (
	// LangEN 代表英语
	// 值为"en"，符合ISO 639-1标准
	LangEN = "en"

	// LangZH 代表简体中文
	// 值为"zh-CN"，遵循语言代码-国家/地区代码格式
	LangZH = "zh-CN"

	// DefaultLanguage 默认语言设置
	// 当未指定语言或检测不到系统语言时使用
	DefaultLanguage = LangEN
)

// 嵌入式翻译文件
// 使用go:embed指令将translations目录下的所有.toml文件嵌入到程序中
//
//go:embed translations/*.toml
var embeddedTranslations embed.FS

// IsSupportedLanguage 检查指定语言是否受支持
//
// 参数:
//   - lang: 要检查的语言代码，如"en"或"zh-CN"
//
// 返回:
//   - 如果语言受支持返回true，否则返回false
//
// 示例:
//
//	if i18n.IsSupportedLanguage("fr") { /* 法语不支持，不会执行 */ }
//	if i18n.IsSupportedLanguage("zh-CN") { /* 中文支持，会执行 */ }
func IsSupportedLanguage(lang string) bool {
	switch lang {
	case LangEN, LangZH:
		return true
	default:
		return false
	}
}

// GetSupportedLanguages 获取所有支持的语言列表
//
// 返回:
//   - 包含所有支持的语言代码的字符串切片
//
// 示例:
//
//	langs := i18n.GetSupportedLanguages() // 返回 []string{"en", "zh-CN"}
func GetSupportedLanguages() []string {
	return []string{LangEN, LangZH}
}

// GetLanguageDisplayName 获取语言的可显示名称
//
// 将语言代码转换为用户友好的语言名称
//
// 参数:
//   - lang: 语言代码，如"en"或"zh-CN"
//
// 返回:
//   - 语言的显示名称，如"English"或"中文(简体)"
//   - 对于不支持的语言，返回原始代码
//
// 示例:
//
//	name := i18n.GetLanguageDisplayName("en") // 返回 "English"
//	name := i18n.GetLanguageDisplayName("zh-CN") // 返回 "中文(简体)"
func GetLanguageDisplayName(lang string) string {
	switch lang {
	case LangEN:
		return "English"
	case LangZH:
		return "中文(简体)"
	default:
		return lang
	}
}
