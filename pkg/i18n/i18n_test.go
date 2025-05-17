package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDetectOSLanguage 测试操作系统语言检测功能
func TestDetectOSLanguage(t *testing.T) {
	tests := []struct {
		name            string
		envVars         map[string]string
		expectedLang    string
		restoreOriginal map[string]string
	}{
		{
			name: "测试中文LANG环境变量",
			envVars: map[string]string{
				"LANG":     "zh_CN.UTF-8",
				"LANGUAGE": "",
				"LC_ALL":   "",
			},
			expectedLang: LangZH,
		},
		{
			name: "测试英文LANG环境变量",
			envVars: map[string]string{
				"LANG":     "en_US.UTF-8",
				"LANGUAGE": "",
				"LC_ALL":   "",
			},
			expectedLang: LangEN,
		},
		{
			name: "测试LANGUAGE环境变量优先级",
			envVars: map[string]string{
				"LANG":     "", // 确保LANG为空，这样才会检查LANGUAGE
				"LANGUAGE": "zh_CN:en_US",
				"LC_ALL":   "",
			},
			expectedLang: LangZH,
		},
		{
			name: "测试LC_ALL环境变量",
			envVars: map[string]string{
				"LANG":     "",
				"LANGUAGE": "",
				"LC_ALL":   "zh_CN.UTF-8",
			},
			expectedLang: LangZH,
		},
		{
			name: "测试无效语言回退到默认语言",
			envVars: map[string]string{
				"LANG":     "fr_FR.UTF-8",
				"LANGUAGE": "fr_FR:de_DE",
				"LC_ALL":   "",
			},
			expectedLang: DefaultLanguage,
		},
		{
			name: "测试无环境变量情况",
			envVars: map[string]string{
				"LANG":     "",
				"LANGUAGE": "",
				"LC_ALL":   "",
			},
			expectedLang: DefaultLanguage,
		},
	}

	// 保存原始环境变量以便后续恢复
	origLang := os.Getenv("LANG")
	origLanguage := os.Getenv("LANGUAGE")
	origLcAll := os.Getenv("LC_ALL")
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LANGUAGE", origLanguage)
		os.Setenv("LC_ALL", origLcAll)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清除可能影响测试的环境变量
			os.Unsetenv("LANG")
			os.Unsetenv("LANGUAGE")
			os.Unsetenv("LC_ALL")

			// 设置测试环境变量
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// 执行测试
			detectedLang := DetectOSLanguage()
			if detectedLang != tt.expectedLang {
				t.Errorf("DetectOSLanguage() = %v, 期望 %v", detectedLang, tt.expectedLang)
			}
		})
	}
}

// TestInitWithEmptyLang 测试Init函数在未提供语言参数时是否能正确检测系统语言
func TestInitWithEmptyLang(t *testing.T) {
	// 保存原始环境变量以便后续恢复
	origLang := os.Getenv("LANG")
	origLanguage := os.Getenv("LANGUAGE")
	origLcAll := os.Getenv("LC_ALL")
	origServergoLang := os.Getenv(EnvLanguage)
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LANGUAGE", origLanguage)
		os.Setenv("LC_ALL", origLcAll)
		os.Setenv(EnvLanguage, origServergoLang)
	}()

	// 清除所有可能影响测试的环境变量
	os.Unsetenv("LANG")
	os.Unsetenv("LANGUAGE")
	os.Unsetenv("LC_ALL")
	os.Unsetenv(EnvLanguage)

	// 重置init状态
	loaded = false
	currentLanguage = ""

	// 创建临时目录确保不使用本地配置文件
	tempDir, err := os.MkdirTemp("", "servergo-i18n-empty-test")
	if err != nil {
		t.Fatalf("无法创建测试临时目录: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 设置HOME目录为临时目录
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", origHome)

	// 测试场景：中文环境
	os.Setenv("LANG", "zh_CN.UTF-8")

	// 测试初始化
	err = Init("")
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 验证是否使用了系统语言
	if GetCurrentLanguage() != LangZH {
		t.Errorf("Init() with empty lang didn't use system language, got %v, want %v",
			GetCurrentLanguage(), LangZH)
	}
}

// TestLanguagePriorityOrder 测试语言选择的优先级顺序
func TestLanguagePriorityOrder(t *testing.T) {
	// 保存原始环境变量以便后续恢复
	origLang := os.Getenv("LANG")
	origLanguage := os.Getenv("LANGUAGE")
	origLcAll := os.Getenv("LC_ALL")
	origServergoLang := os.Getenv(EnvLanguage)
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LANGUAGE", origLanguage)
		os.Setenv("LC_ALL", origLcAll)
		os.Setenv(EnvLanguage, origServergoLang)
	}()

	// 清除可能影响测试的环境变量
	os.Unsetenv("LANG")
	os.Unsetenv("LANGUAGE")
	os.Unsetenv("LC_ALL")
	os.Unsetenv(EnvLanguage)

	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "servergo-i18n-test")
	if err != nil {
		t.Fatalf("无法创建测试临时目录: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建模拟的配置文件
	configDir := filepath.Join(tempDir, ".servergo")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("无法创建配置目录: %v", err)
	}

	// 创建带有中文语言设置的配置文件
	configFile := filepath.Join(configDir, "config.yaml")
	configContent := []byte("language: zh-CN\nauto-open: true\n")
	if err := os.WriteFile(configFile, configContent, 0644); err != nil {
		t.Fatalf("无法创建配置文件: %v", err)
	}

	// 创建原始主目录的备份
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	// 测试场景1: 用户指定语言应该优先
	t.Run("1. 用户指定语言优先", func(t *testing.T) {
		// 重置init状态
		loaded = false
		currentLanguage = ""

		// 设置HOME为测试目录
		os.Setenv("HOME", tempDir)

		// 设置所有环境变量，以确保它们不会干扰测试
		os.Setenv(EnvLanguage, LangEN)   // 环境变量语言为英语
		os.Setenv("LANG", "en_US.UTF-8") // 系统语言为英语

		// 用户明确指定中文，应该使用中文
		err := Init(LangZH)
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}

		if GetCurrentLanguage() != LangZH {
			t.Errorf("用户指定语言未被优先使用, got %v, want %v",
				GetCurrentLanguage(), LangZH)
		}
	})

	// 测试场景1.5: 环境变量优先于配置文件
	t.Run("2. 环境变量优先于配置文件", func(t *testing.T) {
		// 重置init状态
		loaded = false
		currentLanguage = ""

		// 设置HOME为测试目录，使用包含中文配置的配置文件
		os.Setenv("HOME", tempDir)

		// 设置环境变量为英语，配置文件为中文
		os.Setenv(EnvLanguage, LangEN)
		os.Setenv("LANG", "zh_CN.UTF-8") // 系统语言不应影响结果

		// 不指定语言，应该使用环境变量中的英语
		err := Init("")
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}

		if GetCurrentLanguage() != LangEN {
			t.Errorf("环境变量语言未被优先于配置文件使用, got %v, want %v",
				GetCurrentLanguage(), LangEN)
		}
	})

	// 测试场景2: 配置文件语言应该优先于系统语言
	t.Run("3. 配置文件语言优先于系统语言", func(t *testing.T) {
		// 重置init状态
		loaded = false
		currentLanguage = ""

		// 设置HOME为测试目录
		os.Setenv("HOME", tempDir)

		// 设置系统环境变量为英语，清除应用程序环境变量
		os.Setenv("LANG", "en_US.UTF-8")
		os.Unsetenv(EnvLanguage)

		// 不指定语言，应该使用配置文件中的中文
		err := Init("")
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}

		if GetCurrentLanguage() != LangZH {
			t.Errorf("配置文件语言未被正确使用, got %v, want %v",
				GetCurrentLanguage(), LangZH)
		}
	})

	// 测试场景3: 系统语言兜底
	t.Run("4. 系统语言作为兜底", func(t *testing.T) {
		// 重置init状态
		loaded = false
		currentLanguage = ""

		// 创建一个新的临时目录，不包含配置文件
		emptyTempDir, err := os.MkdirTemp("", "servergo-i18n-test-empty")
		if err != nil {
			t.Fatalf("无法创建测试临时目录: %v", err)
		}
		defer os.RemoveAll(emptyTempDir)

		// 设置HOME为空测试目录，确保不使用原配置文件
		os.Setenv("HOME", emptyTempDir)

		// 清除应用程序环境变量，只设置系统语言为中文
		os.Unsetenv(EnvLanguage)
		os.Setenv("LANG", "zh_CN.UTF-8")

		// 不指定语言，应该使用系统检测的中文
		err = Init("")
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}

		if GetCurrentLanguage() != LangZH {
			t.Errorf("系统语言未被正确作为兜底选项, got %v, want %v",
				GetCurrentLanguage(), LangZH)
		}
	})
}
