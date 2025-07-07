package i18n

import (
	"os"
	"path/filepath"
	"strings"
)

// getConfigLanguage 尝试从配置中获取语言设置
// 为避免循环依赖，不直接导入config包，而是通过动态获取
func getConfigLanguage() string {
	// 尝试读取配置文件中的语言设置
	configLang := readLanguageFromConfigFile()
	if configLang != "" {
		return configLang
	}

	return ""
}

// readLanguageFromConfigFile 尝试直接读取配置文件中的语言设置
// 这是一个临时解决方案，避免循环依赖
func readLanguageFromConfigFile() string {
	// 尝试获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// 构建配置文件路径 (与config包中的路径构建保持一致)
	configDir := filepath.Join(homeDir, ".servergo")
	configFile := filepath.Join(configDir, "config.yaml")

	// 检查文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return ""
	}

	// 读取并解析配置文件
	// 注意：这里使用了一个简单的方法来解析YAML，实际情况可能需要更复杂的逻辑
	data, err := os.ReadFile(configFile)
	if err != nil {
		return ""
	}

	// 尝试提取language字段
	// 这里使用一个简单的方法，实际可能需要更复杂的解析
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "language:") {
			langValue := strings.TrimSpace(strings.TrimPrefix(line, "language:"))
			return langValue
		}
	}

	return ""
}
