package i18n

import (
	"fmt"
)

// loadTranslations 加载翻译文件
//
// 内部函数，负责加载所有翻译文件
// 当前实现只加载嵌入式翻译文件
//
// 返回:
//   - 成功时返回nil，失败时返回错误
func loadTranslations() error {
	// 只加载嵌入式翻译
	return loadEmbeddedTranslations()
}

// loadEmbeddedTranslations 从嵌入式文件系统加载翻译
//
// 内部函数，从嵌入的FS中读取并解析翻译文件
// 加载的文件包括英文(en.toml)和中文(zh-CN.toml)
//
// 返回:
//   - 成功时返回nil，失败时返回错误
//   - 如果任一翻译文件加载失败，则返回相应错误
func loadEmbeddedTranslations() error {
	// 加载英文翻译
	enData, err := embeddedTranslations.ReadFile("translations/en.toml")
	if err != nil {
		return fmt.Errorf("failed to read embedded English translations: %v", err)
	}
	_, err = bundle.ParseMessageFileBytes(enData, "en.toml")
	if err != nil {
		return fmt.Errorf("failed to parse embedded English translations: %v", err)
	}

	// 加载中文翻译
	zhData, err := embeddedTranslations.ReadFile("translations/zh-CN.toml")
	if err != nil {
		return fmt.Errorf("failed to read embedded Chinese translations: %v", err)
	}
	_, err = bundle.ParseMessageFileBytes(zhData, "zh-CN.toml")
	if err != nil {
		return fmt.Errorf("failed to parse embedded Chinese translations: %v", err)
	}

	return nil
}
