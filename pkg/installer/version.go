package installer

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/CC11001100/servergo/pkg/logger"
)

// VersionInfo 存储安装的版本信息
type VersionInfo struct {
	Version     string    `json:"version"`
	InstallPath string    `json:"install_path"`
	InstallTime time.Time `json:"install_time"`
	BackupPath  string    `json:"backup_path,omitempty"`
}

// 获取版本信息文件的路径
func getVersionFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %v", err)
	}

	// 存储在用户目录下的.servergo/version.json
	versionDir := filepath.Join(homeDir, ".servergo")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return "", fmt.Errorf("无法创建版本信息目录: %v", err)
	}

	return filepath.Join(versionDir, "version.json"), nil
}

// SaveVersionInfo 保存当前版本信息
func SaveVersionInfo(version, installPath string) error {
	versionFilePath, err := getVersionFilePath()
	if err != nil {
		return err
	}

	// 检查是否存在旧版本信息，如果存在，保留为备份
	var oldInfo VersionInfo
	oldExists := false
	if data, err := os.ReadFile(versionFilePath); err == nil {
		if err := json.Unmarshal(data, &oldInfo); err == nil {
			oldExists = true
		}
	}

	info := VersionInfo{
		Version:     version,
		InstallPath: installPath,
		InstallTime: time.Now(),
	}

	// 如果有旧版本并且不同于当前版本，保存备份路径
	if oldExists && oldInfo.Version != version {
		backupPath := filepath.Join(filepath.Dir(versionFilePath), "backup", oldInfo.Version)
		if err := os.MkdirAll(backupPath, 0755); err == nil {
			oldExecName := filepath.Base(oldInfo.InstallPath)
			backupExePath := filepath.Join(backupPath, oldExecName)

			// 复制旧版本到备份目录（仅在Windows下，Linux/Mac使用的是符号链接）
			if stat, err := os.Stat(oldInfo.InstallPath); err == nil && !stat.Mode().IsDir() {
				if data, err := os.ReadFile(oldInfo.InstallPath); err == nil {
					if err := os.WriteFile(backupExePath, data, fs.FileMode(0755)); err == nil {
						info.BackupPath = backupExePath
						logger.Info("旧版本 %s 已备份到 %s", oldInfo.Version, backupExePath)
					}
				}
			}
		}
	}

	// 将版本信息写入文件
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("无法序列化版本信息: %v", err)
	}

	if err := os.WriteFile(versionFilePath, data, 0644); err != nil {
		return fmt.Errorf("无法写入版本信息: %v", err)
	}

	return nil
}

// GetInstalledVersion 获取已安装的版本信息
func GetInstalledVersion() (*VersionInfo, error) {
	versionFilePath, err := getVersionFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(versionFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // 没有安装记录
		}
		return nil, fmt.Errorf("无法读取版本信息: %v", err)
	}

	var info VersionInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("无法解析版本信息: %v", err)
	}

	return &info, nil
}

// CompareVersions 比较版本号，返回:
// -1: oldVersion < newVersion
//
//	0: oldVersion == newVersion
//	1: oldVersion > newVersion
func CompareVersions(oldVersion, newVersion string) int {
	// 简单实现，实际可能需要更复杂的版本比较逻辑
	if oldVersion == newVersion {
		return 0
	}

	// 此处可以添加更复杂的版本比较逻辑
	// 例如，将版本号拆分为主要、次要和补丁版本，并逐一比较

	// 简单起见，这里只进行字符串比较
	if oldVersion < newVersion {
		return -1
	}
	return 1
}
