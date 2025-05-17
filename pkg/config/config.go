package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	// 其他配置项可以在这里添加
}

// 获取配置目录路径
func getConfigDir() (string, error) {
	// 获取用户HOME目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %v", err)
	}

	// 配置目录路径
	configDir := filepath.Join(homeDir, "."+AppName)

	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("无法创建配置目录: %v", err)
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
	// 其他配置项设置...

	// 获取配置目录
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("无法获取配置目录: %v", err)
	}

	// 获取配置文件路径
	configPath := filepath.Join(configDir, ConfigFileName+"."+ConfigFileType)

	// 检查配置文件是否存在，如果不存在则创建一个空的配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 创建空配置文件
		file, err := os.Create(configPath)
		if err != nil {
			return fmt.Errorf("无法创建配置文件: %v", err)
		}
		file.Close()
	}

	// 保存配置到文件
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("无法写入配置文件: %v", err)
	}

	return nil
}

// 设置默认配置
func SetDefaults() {
	viper.SetDefault("auto-open", true)          // 默认自动打开浏览器
	viper.SetDefault("enable-dir-listing", true) // 默认启用目录列表功能
	viper.SetDefault("theme", "default")         // 默认使用默认主题
}
