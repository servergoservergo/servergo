// Package i18n 提供国际化支持功能
// 本包实现了应用程序的多语言翻译机制，支持中英文切换
// 核心功能包括：
// 1. 从嵌入式文件系统加载翻译文件
// 2. 管理和切换当前语言设置
// 3. 提供翻译文本查询接口
// 4. 自动检测操作系统语言设置
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

// 全局变量，用于存储国际化相关的状态和对象
var (
	// bundle 是go-i18n翻译文件的集合
	// 用于存储和管理所有语言的翻译消息
	bundle *i18n.Bundle

	// localizer 是当前语言的本地化器
	// 根据当前语言提供翻译功能
	localizer *i18n.Localizer

	// mutex 是读写锁，用于保护并发访问共享数据
	// 在读取或修改当前语言设置时提供线程安全
	mutex sync.RWMutex

	// loaded 标记i18n系统是否已初始化
	// 初始化完成后置为true
	loaded bool

	// currentLanguage 存储当前设置的语言代码
	// 例如"en"或"zh-CN"
	currentLanguage string

	// initMutex 是用于保护初始化过程的互斥锁
	// 防止并发调用Init函数导致的竞争条件
	initMutex sync.Mutex
)

// Init 初始化国际化支持系统
//
// 该函数完成以下工作:
// 1. 创建翻译Bundle并注册TOML解析器
// 2. 加载嵌入的翻译文件
// 3. 创建指定语言的本地化器
//
// 参数:
//   - lang: 要初始化的语言代码，如"en"或"zh-CN"，如为空则自动检测系统语言
//
// 返回:
//   - 成功时返回nil，失败时返回对应的错误
//
// 示例:
//
//	err := i18n.Init("zh-CN") // 初始化为中文
//	err := i18n.Init("") // 使用系统语言初始化，如系统语言不支持则使用默认语言
func Init(lang string) error {
	// 防止并发初始化
	initMutex.Lock()
	defer initMutex.Unlock()

	// 如果语言未指定，先尝试检测操作系统语言
	if lang == "" {
		lang = DetectOSLanguage()
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

// GetCurrentLanguage 返回当前设置的语言代码
//
// 该函数线程安全，会使用读锁保护对currentLanguage的访问
//
// 返回:
//   - 当前语言代码，如"en"或"zh-CN"
//   - 如果尚未设置语言，返回默认语言代码
//
// 示例:
//
//	lang := i18n.GetCurrentLanguage() // 可能返回"en"或"zh-CN"
func GetCurrentLanguage() string {
	mutex.RLock()
	defer mutex.RUnlock()

	if currentLanguage == "" {
		return DefaultLanguage
	}

	return currentLanguage
}

// SetLanguage 设置当前语言
//
// 该函数会切换当前的语言设置，更新本地化器
// 线程安全，会使用写锁保护对共享数据的修改
//
// 参数:
//   - lang: 要设置的语言代码，如"en"或"zh-CN"
//
// 返回:
//   - 成功时返回nil，如果指定了不支持的语言则返回错误
//
// 示例:
//
//	err := i18n.SetLanguage("zh-CN") // 切换到中文
//	err := i18n.SetLanguage("en") // 切换到英文
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

// T 是主要的翻译函数，根据当前语言返回对应的翻译文本
//
// 该函数会查找messageID对应的翻译，自动处理i18n系统未初始化的情况
// 使用读锁保护对共享数据的访问，线程安全
//
// 参数:
//   - messageID: 翻译项的唯一标识符，如"app.name"或"error.not_found"
//
// 返回:
//   - 当前语言环境下messageID对应的翻译文本
//   - 如果未找到翻译或发生错误，则返回原始messageID
//
// 示例:
//
//	msg := i18n.T("app.name") // 如果当前语言是中文，返回"ServerGo"
//	msg := i18n.T("error.not_found") // 如果找不到翻译，返回"error.not_found"
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

// Tf 是带格式化的翻译函数
//
// 该函数先通过T获取翻译，然后使用fmt.Sprintf进行格式化
// 适用于包含变量的翻译文本，如"Hello, %s"
//
// 参数:
//   - messageID: 翻译项的唯一标识符
//   - args: 格式化参数，用于替换翻译文本中的占位符
//
// 返回:
//   - 格式化后的翻译文本
//
// 示例:
//
//	// 假设"greeting"的中文翻译是"你好，%s!"
//	msg := i18n.Tf("greeting", "张三") // 返回"你好，张三!"
//
//	// 假设"items.count"的英文翻译是"You have %d items"
//	msg := i18n.Tf("items.count", 5) // 返回"You have 5 items"
func Tf(messageID string, args ...interface{}) string {
	return fmt.Sprintf(T(messageID), args...)
}
