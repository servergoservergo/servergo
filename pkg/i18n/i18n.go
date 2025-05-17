package i18n

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// 支持的语言常量
const (
	// LangEN 英语
	LangEN = "en"
	// LangZH 简体中文
	LangZH = "zh-CN"
	// LangDefault 默认语言
	LangDefault = LangEN
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
		return fmt.Errorf("unsupported language: %s", lang)
	}

	// 创建一个新的Bundle，使用英语作为回退语言
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载翻译文件
	err := loadTranslations()
	if err != nil {
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

	if localizer == nil {
		return LangDefault
	}

	// 我们没有直接的方法获取 localizer 的当前语言
	// 所以简单返回当前已经保存的语言值
	// 创建时保存语言设置
	currentLanguage := LangDefault
	for _, lang := range bundle.LanguageTags() {
		langStr := lang.String()
		if strings.HasPrefix(langStr, "zh") {
			currentLanguage = LangZH
			break
		} else if strings.HasPrefix(langStr, "en") {
			currentLanguage = LangEN
			break
		}
	}

	return currentLanguage
}

// SetLanguage 设置当前语言
func SetLanguage(lang string) error {
	// 检查语言是否受支持
	if !IsSupportedLanguage(lang) {
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
	return LangDefault
}

// loadTranslations 加载翻译文件
func loadTranslations() error {
	// 先尝试加载嵌入式翻译
	loadEmbeddedTranslations()

	// 再尝试从外部文件加载翻译
	// 尝试加载英文翻译
	enFile := filepath.Join("translations", "en.toml")
	if _, err := os.Stat(enFile); err == nil {
		_, loadErr := bundle.LoadMessageFile(enFile)
		if loadErr != nil {
			// 只有在出错时才记录日志
			logger.Info("Failed to load English translation file: %v", loadErr)
		}
	}

	// 尝试加载中文翻译
	zhFile := filepath.Join("translations", "zh-CN.toml")
	if _, err := os.Stat(zhFile); err == nil {
		_, loadErr := bundle.LoadMessageFile(zhFile)
		if loadErr != nil {
			// 只有在出错时才记录日志
			logger.Info("Failed to load Chinese translation file: %v", loadErr)
		}
	}

	return nil
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
		return messageID
	}

	return translation
}

// Tf 格式化翻译文本，支持变量替换
func Tf(messageID string, args ...interface{}) string {
	// 获取基本翻译文本
	translated := T(messageID)

	// 如果有替换参数，进行格式化
	if len(args) > 0 {
		return fmt.Sprintf(translated, args...)
	}

	return translated
}
