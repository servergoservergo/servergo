package i18n

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// 支持的语言常量
const (
	// LangEN 英语
	LangEN = "en"
	// LangZH 简体中文
	LangZH = "zh-CN"
	// DefaultLanguage 默认语言
	DefaultLanguage = LangEN
)

//go:embed translations/*.toml
var embeddedTranslations embed.FS

var (
	bundle          *i18n.Bundle
	localizer       *i18n.Localizer
	mutex           sync.RWMutex
	loaded          bool
	currentLanguage string
	initMutex       sync.Mutex
)

// Init 初始化国际化支持
func Init(lang string) error {
	// 防止并发初始化
	initMutex.Lock()
	defer initMutex.Unlock()

	// 如果语言未指定，使用默认语言
	if lang == "" {
		lang = DefaultLanguage
	}

	// 验证语言是否受支持
	if !IsSupportedLanguage(lang) {
		// 这里不能使用 i18n.T，因为 i18n 系统还未初始化
		return fmt.Errorf("unsupported language: %s", lang)
	}

	// 创建一个新的Bundle，使用英语作为回退语言
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载翻译文件
	err := loadTranslations()
	if err != nil {
		// 这里不能使用 i18n.T，因为 i18n 系统还未初始化完成
		return fmt.Errorf("failed to load translations: %v", err)
	}

	// 创建本地化器
	localizer = i18n.NewLocalizer(bundle, lang)

	// 保存当前语言
	currentLanguage = lang

	// 标记为已加载
	loaded = true

	return nil
}

// GetCurrentLanguage 返回当前设置的语言
func GetCurrentLanguage() string {
	mutex.RLock()
	defer mutex.RUnlock()

	if currentLanguage == "" {
		return DefaultLanguage
	}

	return currentLanguage
}

// SetLanguage 设置当前语言
func SetLanguage(lang string) error {
	// 检查语言是否受支持
	if !IsSupportedLanguage(lang) {
		// 此时 i18n 已经初始化，但为避免循环引用，不使用 i18n.T
		return fmt.Errorf("language '%s' is not supported", lang)
	}

	mutex.Lock()
	defer mutex.Unlock()

	// 重新创建本地化器
	localizer = i18n.NewLocalizer(bundle, lang)

	// 更新当前语言
	currentLanguage = lang

	// 标记为已加载
	loaded = true

	return nil
}

// IsSupportedLanguage 检查指定语言是否受支持
func IsSupportedLanguage(lang string) bool {
	switch lang {
	case LangEN, LangZH:
		return true
	default:
		return false
	}
}

// GetSupportedLanguages 获取所有支持的语言列表
func GetSupportedLanguages() []string {
	return []string{LangEN, LangZH}
}

// GetLanguageDisplayName 获取语言的显示名称
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

// DetectOSLanguage 检测操作系统的语言设置
func DetectOSLanguage() string {
	// 环境变量中有语言设置时优先使用
	envLangs := []string{
		os.Getenv("LANG"),
		os.Getenv("LANGUAGE"),
		os.Getenv("LC_ALL"),
	}

	for _, envLang := range envLangs {
		if envLang != "" {
			// 解析语言代码，通常格式为 "en_US.UTF-8"
			langParts := strings.Split(envLang, ".")
			if len(langParts) > 0 {
				langCode := strings.Replace(langParts[0], "_", "-", -1)
				// 简单匹配支持的语言
				if strings.HasPrefix(langCode, "zh") {
					return LangZH
				}
				if strings.HasPrefix(langCode, "en") {
					return LangEN
				}
			}
		}
	}

	// 找不到有效语言设置时返回默认语言
	return DefaultLanguage
}

// loadTranslations 加载翻译文件
func loadTranslations() error {
	// 只加载嵌入式翻译
	return loadEmbeddedTranslations()
}

// loadEmbeddedTranslations 从嵌入式文件系统加载翻译
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

// T 是主要的翻译函数，根据当前语言返回对应的翻译
func T(messageID string) string {
	// 确保已初始化，如果尚未初始化，可能会导致nil指针异常
	if !loaded {
		// 如果尚未加载，尝试初始化
		if err := Init(""); err != nil {
			// 初始化失败，直接返回原始ID
			return messageID
		}
	}

	mutex.RLock()
	defer mutex.RUnlock()

	// 如果初始化失败或localizer为nil，直接返回消息ID
	if localizer == nil {
		return messageID
	}

	// 翻译消息
	translation, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})

	if err != nil {
		// 找不到翻译时返回原始ID
		return messageID
	}

	return translation
}

// Tf 是带格式化的翻译函数，类似 fmt.Sprintf(T(messageID), args...)
func Tf(messageID string, args ...interface{}) string {
	return fmt.Sprintf(T(messageID), args...)
}
