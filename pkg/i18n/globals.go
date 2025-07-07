package i18n

import (
	"fmt"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

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
