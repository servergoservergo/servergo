package i18n

import (
	"embed"
	"encoding/json"
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
	LangEnglish = "en"
	LangChinese = "zh-CN"
	LangDefault = LangEnglish // 默认回退语言
)

//go:embed translations/*.toml
var embeddedTranslations embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	mutex     sync.RWMutex
	loaded    bool
)

// Init 初始化i18n包
func Init(lang string) error {
	mutex.Lock()
	defer mutex.Unlock()

	// 如果未指定语言，尝试检测系统语言
	if lang == "" {
		lang = DetectOSLanguage()
	}

	// 只支持特定语言，不支持的语言回退到默认语言
	if !IsSupportedLanguage(lang) {
		lang = LangDefault
	}

	// 创建一个新的Bundle并设置默认语言
	bundle = i18n.NewBundle(language.English)

	// 注册解析函数，根据文件扩展名选择相应的解析器
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载翻译文件
	err := loadTranslations(lang)
	if err != nil {
		return fmt.Errorf("failed to load translations: %v", err)
	}

	// 创建本地化器
	localizer = i18n.NewLocalizer(bundle, lang)

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
			currentLanguage = LangChinese
			break
		} else if strings.HasPrefix(langStr, "en") {
			currentLanguage = LangEnglish
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

	return nil
}

// IsSupportedLanguage 检查指定语言是否受支持
func IsSupportedLanguage(lang string) bool {
	switch lang {
	case LangEnglish, LangChinese:
		return true
	default:
		return false
	}
}

// GetSupportedLanguages 获取所有支持的语言列表
func GetSupportedLanguages() []string {
	return []string{LangEnglish, LangChinese}
}

// GetLanguageDisplayName 获取语言的显示名称
func GetLanguageDisplayName(lang string) string {
	switch lang {
	case LangEnglish:
		return "English"
	case LangChinese:
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
					return LangChinese
				}
				if strings.HasPrefix(langCode, "en") {
					return LangEnglish
				}
			}
		}
	}

	// 找不到有效语言设置时返回默认语言
	return LangDefault
}

// loadTranslations 加载翻译文件
func loadTranslations(lang string) error {
	// 从嵌入式文件系统加载翻译数据
	err := loadEmbeddedTranslations()
	if err != nil {
		return fmt.Errorf("failed to load embedded translations: %v", err)
	}

	// 尝试从外部文件加载（如果存在）
	// 获取可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		// 忽略错误，只使用内嵌翻译
		return nil
	}

	// 翻译文件目录（与可执行文件同级的translations目录）
	transDir := filepath.Join(filepath.Dir(exePath), "translations")

	// 如果外部翻译目录存在，尝试加载外部翻译文件
	if _, err := os.Stat(transDir); err == nil {
		enFile := filepath.Join(transDir, "en.toml")
		zhFile := filepath.Join(transDir, "zh-CN.toml")

		// 尝试加载英文翻译文件（如果存在）
		if _, err := os.Stat(enFile); err == nil {
			// LoadMessageFile返回两个值：消息文件和错误
			_, loadErr := bundle.LoadMessageFile(enFile)
			if loadErr != nil {
				logger.Info("加载英文翻译文件失败: %v", loadErr)
			} else {
				logger.Info("加载外部英文翻译文件成功: %s", enFile)
			}
		}

		// 尝试加载中文翻译文件（如果存在）
		if _, err := os.Stat(zhFile); err == nil {
			// LoadMessageFile返回两个值：消息文件和错误
			_, loadErr := bundle.LoadMessageFile(zhFile)
			if loadErr != nil {
				logger.Info("加载中文翻译文件失败: %v", loadErr)
			} else {
				logger.Info("加载外部中文翻译文件成功: %s", zhFile)
			}
		}
	}

	return nil
}

// loadEmbeddedTranslations 加载内嵌的翻译数据
func loadEmbeddedTranslations() error {
	// 加载嵌入式英文翻译
	enData, err := embeddedTranslations.ReadFile("translations/en.toml")
	if err != nil {
		return fmt.Errorf("failed to read embedded English translation: %v", err)
	}

	// 加载嵌入式中文翻译
	zhData, err := embeddedTranslations.ReadFile("translations/zh-CN.toml")
	if err != nil {
		return fmt.Errorf("failed to read embedded Chinese translation: %v", err)
	}

	// 解析英文翻译
	_, err = bundle.ParseMessageFileBytes(enData, "en.toml")
	if err != nil {
		return fmt.Errorf("failed to parse English translation: %v", err)
	}

	// 解析中文翻译
	_, err = bundle.ParseMessageFileBytes(zhData, "zh-CN.toml")
	if err != nil {
		return fmt.Errorf("failed to parse Chinese translation: %v", err)
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
