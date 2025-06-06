package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/spf13/viper"
)

const (
	// 应用名称
	AppName = "servergo"
	// 配置文件名（无扩展名）
	ConfigFileName = "config"
	// 配置文件类型
	ConfigFileType = "yaml"
)

// Config 结构表示应用程序的配置
type Config struct {
	// 启动后是否自动打开浏览器
	AutoOpen bool `mapstructure:"auto-open"`
	// 是否启用目录列表功能
	EnableDirListing bool `mapstructure:"enable-dir-listing"`
	// 目录列表主题
	Theme string `mapstructure:"theme"`
	// 界面语言
	Language string `mapstructure:"language"`
	// 是否启用日志持久化
	EnableLogPersistence bool `mapstructure:"enable-log-persistence"`
	// 从哪个端口开始递增寻找空闲端口(0表示随机选择)
	StartPort int `mapstructure:"start-port"`
	// 认证相关配置
	Username string `json:"username" yaml:"username"` // 默认用户名
	Password string `json:"password" yaml:"password"` // 默认密码
	// 其他配置项可以在这里添加
}

// 获取配置目录路径
func getConfigDir() (string, error) {
	// 获取用户HOME目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	// 配置目录路径
	configDir := filepath.Join(homeDir, "."+AppName)

	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return configDir, nil
}

// 获取配置文件路径
func GetConfigFilePath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, ConfigFileName+"."+ConfigFileType), nil
}

// InitConfig 初始化配置
func InitConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	// 设置配置文件的路径和名称
	viper.AddConfigPath(configDir)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)

	// 读取配置文件，如果不存在则忽略错误
	_ = viper.ReadInConfig()

	// 设置默认值
	SetDefaults()

	return nil
}

// GetConfig 获取当前配置
func GetConfig() Config {
	var cfg Config
	_ = viper.Unmarshal(&cfg)
	return cfg
}

// SaveConfig 保存配置到文件
func SaveConfig(cfg Config) error {
	// 将配置设置到viper
	viper.Set("auto-open", cfg.AutoOpen)
	viper.Set("enable-dir-listing", cfg.EnableDirListing)
	viper.Set("theme", cfg.Theme)
	viper.Set("language", cfg.Language)
	viper.Set("enable-log-persistence", cfg.EnableLogPersistence)
	viper.Set("start-port", cfg.StartPort)
	viper.Set("username", cfg.Username)
	viper.Set("password", cfg.Password)
	// 其他配置项设置...

	// 获取配置目录
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %v", err)
	}

	// 获取配置文件路径
	configPath := filepath.Join(configDir, ConfigFileName+"."+ConfigFileType)

	// 检查配置文件是否存在，如果不存在则创建一个空的配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 创建空配置文件
		file, err := os.Create(configPath)
		if err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
		file.Close()
	}

	// 保存配置到文件
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// SetLanguage 设置语言并保存到配置
func SetLanguage(lang string) error {
	// 验证语言是否支持
	if !i18n.IsSupportedLanguage(lang) {
		// 注意：这里不使用 i18n.T 是为了避免在语言设置出错时可能导致的循环依赖
		// 由于这个错误消息主要是内部使用，保持简单的英文错误即可
		return fmt.Errorf("language '%s' is not supported", lang)
	}

	// 获取当前配置
	cfg := GetConfig()

	// 更新语言设置
	cfg.Language = lang

	// 保存配置
	if err := SaveConfig(cfg); err != nil {
		return err
	}

	// 更新i18n包的当前语言
	return i18n.SetLanguage(lang)
}

// GetLanguage 获取当前语言设置
func GetLanguage() string {
	return GetConfig().Language
}

// 设置默认配置
func SetDefaults() {
	viper.SetDefault("auto-open", true)              // 默认自动打开浏览器
	viper.SetDefault("enable-dir-listing", true)     // 默认启用目录列表功能
	viper.SetDefault("theme", "default")             // 默认使用默认主题
	viper.SetDefault("enable-log-persistence", true) // 默认启用日志持久化
	viper.SetDefault("start-port", 0)                // 默认从0开始递增寻找空闲端口
	viper.SetDefault("username", "admin")            // 默认用户名
	viper.SetDefault("password", "")                 // 默认密码为空，将自动生成

	// 语言默认设置为自动检测
	detectLang := i18n.DetectOSLanguage()
	viper.SetDefault("language", detectLang) // 默认使用系统语言
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		AutoOpen:             true,
		EnableDirListing:     true,
		Theme:                "default",
		Language:             i18n.DetectOSLanguage(),
		EnableLogPersistence: true,
		StartPort:            0,
		Username:             "admin", // 默认用户名
		Password:             "",      // 默认密码为空，将自动生成
	}
}
