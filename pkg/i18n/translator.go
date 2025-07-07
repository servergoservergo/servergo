package i18n

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Init 初始化国际化支持系统
//
// 该函数完成以下工作:
// 1. 创建翻译Bundle并注册TOML解析器
// 2. 加载嵌入的翻译文件
// 3. 创建指定语言的本地化器
//
// 语言优先级:
// 1. 用户指定的语言参数
// 2. 环境变量中的语言设置 (SERVERGO_LANGUAGE)
// 3. 配置文件中的语言设置
// 4. 操作系统检测的语言
// 5. 最后才使用默认语言(英语)
//
// 参数:
//   - lang: 要初始化的语言代码，如"en"或"zh-CN"，如为空则按优先级自动选择语言
//
// 返回:
//   - 成功时返回nil，失败时返回对应的错误
//
// 示例:
//
//	err := i18n.Init("zh-CN") // 初始化为中文
//	err := i18n.Init("") // 按优先级自动选择语言
func Init(lang string) error {
	// 防止并发初始化
	initMutex.Lock()
	defer initMutex.Unlock()

	// 语言选择优先级:
	// 1. 用户指定的语言参数
	// 2. 环境变量中的语言设置 (SERVERGO_LANGUAGE)
	// 3. 配置文件中的语言设置
	// 4. 操作系统检测的语言
	if lang == "" {
		// 检查环境变量
		envLang := os.Getenv(EnvLanguage)
		if envLang != "" && IsSupportedLanguage(envLang) {
			lang = envLang
		} else {
			// 从配置文件获取语言设置
			configLang := getConfigLanguage()
			if configLang != "" && IsSupportedLanguage(configLang) {
				// 使用配置文件中的语言
				lang = configLang
			} else {
				// 配置文件中没有有效的语言设置，使用系统语言
				lang = DetectOSLanguage()
			}
		}
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
