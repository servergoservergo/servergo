package i18n

import (
	"os"
	"strings"
)

// DetectOSLanguage 检测操作系统的语言设置
//
// 该函数通过检查环境变量(LANG, LANGUAGE, LC_ALL)来确定系统语言
// 环境变量示例:
// - LANG=zh_CN.UTF-8
// - LANGUAGE=zh_CN:en_US (冒号分隔多种语言，按优先级排序)
// - LC_ALL=zh_CN.UTF-8
//
// 检测顺序:
// 1. 先检查 LANG 环境变量
// 2. 然后检查 LANGUAGE 环境变量(可能包含多个语言，以冒号分隔)
// 3. 最后检查 LC_ALL 环境变量
// 4. 如果都检测不到可识别的语言，返回默认语言
//
// 返回:
//   - 检测到的语言代码，如"en"或"zh-CN"
//   - 如果无法检测或不支持，则返回默认语言代码
//
// 示例:
//
//	lang := i18n.DetectOSLanguage()  // 可能返回"en"或"zh-CN"
func DetectOSLanguage() string {
	// 定义检查环境变量的函数
	checkEnvVar := func(varName string) string {
		value := os.Getenv(varName)
		if value == "" {
			return ""
		}

		// 处理可能包含多个语言定义的情况 (如 LANGUAGE=zh_CN:en_US)
		// 对于LANGUAGE，冒号前的语言优先级更高
		languageCodes := strings.Split(value, ":")
		for _, langCode := range languageCodes {
			if langCode == "" {
				continue
			}

			// 解析语言代码，通常格式为 "en_US.UTF-8" 或 "zh_CN"
			// 先处理可能存在的字符编码部分
			langParts := strings.Split(langCode, ".")
			if len(langParts) == 0 {
				continue
			}

			// 将 en_US 格式转换为 en-US 格式
			normalizedCode := strings.Replace(langParts[0], "_", "-", -1)

			// 提取主要语言部分
			majorLang := strings.ToLower(strings.Split(normalizedCode, "-")[0])

			// 匹配支持的语言
			switch majorLang {
			case "zh":
				return LangZH
			case "en":
				return LangEN
			}

			// 如果有精确匹配，也直接返回
			if normalizedCode == LangZH || normalizedCode == LangEN {
				return normalizedCode
			}
		}
		return ""
	}

	// 按优先级依次检查环境变量
	// 1. LANG 环境变量 (最常用)
	if lang := checkEnvVar("LANG"); lang != "" {
		return lang
	}

	// 2. LANGUAGE 环境变量 (可能包含多个值)
	if lang := checkEnvVar("LANGUAGE"); lang != "" {
		return lang
	}

	// 3. LC_ALL 环境变量 (通常用于覆盖所有本地化设置)
	if lang := checkEnvVar("LC_ALL"); lang != "" {
		return lang
	}

	// 找不到有效语言设置或不支持时返回默认语言
	return DefaultLanguage
}
