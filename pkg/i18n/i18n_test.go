package i18n

import (
	"os"
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
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LANGUAGE", origLanguage)
		os.Setenv("LC_ALL", origLcAll)
	}()

	// 清除可能影响测试的环境变量
	os.Unsetenv("LANG")
	os.Unsetenv("LANGUAGE")
	os.Unsetenv("LC_ALL")

	// 测试场景：中文环境
	os.Setenv("LANG", "zh_CN.UTF-8")

	// 测试初始化
	err := Init("")
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 验证是否使用了系统语言
	if GetCurrentLanguage() != LangZH {
		t.Errorf("Init() with empty lang didn't use system language, got %v, want %v",
			GetCurrentLanguage(), LangZH)
	}
}
