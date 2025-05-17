package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// 用于测试的临时配置目录
var testConfigDir string

// 模拟配置函数，返回测试目录
func mockGetConfigDir() (string, error) {
	return testConfigDir, nil
}

// 测试前的设置函数
func setupTestConfig(t *testing.T) {
	// 创建临时目录作为测试用的配置目录
	tempDir, err := os.MkdirTemp("", "servergo-test-config")
	if err != nil {
		t.Fatalf("无法创建测试临时目录: %v", err)
	}
	testConfigDir = tempDir

	// 清除可能存在的viper配置
	viper.Reset()

	// 记录清理函数
	t.Cleanup(func() {
		// 清理临时目录
		os.RemoveAll(testConfigDir)

		// 重置viper
		viper.Reset()
	})
}

// 测试获取配置文件路径
func TestGetConfigFilePath(t *testing.T) {
	setupTestConfig(t)

	// 因为我们不能直接修改getConfigDir，只能测试整体功能
	// 或者使用临时的环境变量来影响其行为

	// 检查是否能获取用户主目录
	_, err := os.UserHomeDir()
	if err != nil {
		t.Skip("无法获取用户主目录，跳过测试")
	}

	// 检查配置文件路径是否包含必要元素
	configPath, err := GetConfigFilePath()
	assert.NoError(t, err)
	assert.Contains(t, configPath, ConfigFileName+"."+ConfigFileType)
	assert.Contains(t, configPath, "."+AppName)
}

// 测试初始化配置
func TestInitConfig(t *testing.T) {
	setupTestConfig(t)

	// 初始化配置
	err := InitConfig()
	assert.NoError(t, err)

	// 由于我们无法直接访问viper的内部状态，
	// 这里只能验证函数不返回错误
	// 如果需要进一步验证，可以通过检查后续函数的行为间接验证
}

// 测试设置默认配置
func TestSetDefaults(t *testing.T) {
	// 重置viper以确保清除任何现有配置
	viper.Reset()

	// 设置默认值
	SetDefaults()

	// 验证默认值是否正确设置
	assert.True(t, viper.GetBool("auto-open"))
	assert.True(t, viper.GetBool("enable-dir-listing"))
	assert.Equal(t, "default", viper.GetString("theme"))
}

// 测试获取配置
func TestGetConfig(t *testing.T) {
	setupTestConfig(t)

	// 设置一些配置值
	viper.Set("auto-open", false)
	viper.Set("enable-dir-listing", false)
	viper.Set("theme", "dark")

	// 获取配置
	cfg := GetConfig()

	// 验证配置值
	assert.False(t, cfg.AutoOpen)
	assert.False(t, cfg.EnableDirListing)
	assert.Equal(t, "dark", cfg.Theme)
}

// 测试保存配置
func TestSaveConfig(t *testing.T) {
	// 这个测试无法在不修改原始函数的情况下进行完整测试
	// 因为SaveConfig内部使用了getConfigDir，我们无法直接修改它
	// 这里我们只能测试一些基本功能，或者使用环境变量来影响它

	// 获取临时目录用于测试
	tempDir, err := os.MkdirTemp("", "servergo-test-config")
	if err != nil {
		t.Fatalf("无法创建测试临时目录: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 使用临时环境变量来影响HOME目录
	// 注意：这只在getConfigDir使用os.UserHomeDir时有效
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// 创建测试配置
	testCfg := Config{
		AutoOpen:         false,
		EnableDirListing: true,
		Theme:            "test-theme",
	}

	// 在真实环境中测试保存配置
	// 注意：这可能会修改真实的配置文件，具体取决于getConfigDir的实现
	t.Skip("跳过实际写入文件的测试，以防意外修改实际配置文件")

	// 以下代码在Skip之后不会执行，但保留以供参考
	err = SaveConfig(testCfg)
	assert.NoError(t, err)
}
