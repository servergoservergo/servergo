package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/CC11001100/servergo/pkg/version"
)

// Installer 接口定义了安装程序的方法
type Installer interface {
	// Install 将程序安装到系统PATH中
	Install(executablePath string) error
	// Uninstall 从系统PATH中移除程序
	Uninstall() error
}

// GetExecutablePath 获取当前可执行文件的路径
func GetExecutablePath() (string, error) {
	// 获取当前可执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("无法获取可执行文件路径: %v", err)
	}

	// 确保使用绝对路径
	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return "", fmt.Errorf("无法获取绝对路径: %v", err)
	}

	// 解析符号链接（如果有的话）
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return "", fmt.Errorf("无法解析符号链接: %v", err)
	}

	return realPath, nil
}

// InstallToPath 执行安装过程
func InstallToPath() error {
	// 获取当前可执行文件的路径
	execPath, err := GetExecutablePath()
	if err != nil {
		return err
	}

	currentVersion := version.Version
	logger.Info("当前可执行文件路径: %s", execPath)
	logger.Info("当前版本: %s", currentVersion)
	logger.Info("构建信息: %s", version.GetBuildInfo())

	// 获取对应操作系统的安装器
	installer := NewInstaller()

	// 检查是否已安装
	installedVer, err := GetInstalledVersion()
	if err == nil && installedVer != nil {
		compareResult := CompareVersions(installedVer.Version, currentVersion)
		switch compareResult {
		case 0:
			logger.Info("检测到已安装相同版本 %s，将重新安装...", currentVersion)
		case -1:
			logger.Info("检测到已安装旧版本 %s，将升级到新版本 %s...", installedVer.Version, currentVersion)
		case 1:
			logger.Warning("警告: 当前版本 %s 低于已安装版本 %s，降级可能导致兼容性问题", currentVersion, installedVer.Version)
			logger.Info("继续安装，将替换为当前版本...")
		}
	} else {
		logger.Info("未检测到已安装版本，进行首次安装...")
	}

	// 执行安装
	if err := installer.Install(execPath); err != nil {
		return err
	}

	// 保存版本信息
	if err := SaveVersionInfo(currentVersion, execPath); err != nil {
		logger.Warning("保存版本信息失败: %v", err)
	}

	return nil
}

// UninstallFromPath 执行卸载过程
func UninstallFromPath() error {
	// 获取对应操作系统的安装器
	installer := NewInstaller()

	// 获取已安装版本信息
	installedVer, _ := GetInstalledVersion()
	if installedVer != nil {
		logger.Info("正在卸载版本 %s...", installedVer.Version)
	}

	// 执行卸载
	return installer.Uninstall()
}

// isSameFile 检查两个文件是否相同 (通用辅助函数)
func isSameFile(file1, file2 string) bool {
	// 检查文件大小
	stat1, err1 := os.Stat(file1)
	stat2, err2 := os.Stat(file2)

	if err1 != nil || err2 != nil {
		return false
	}

	// 如果文件大小不同，则肯定不是同一个文件
	if stat1.Size() != stat2.Size() {
		return false
	}

	// 读取两个文件内容进行比较
	content1, err1 := os.ReadFile(file1)
	content2, err2 := os.ReadFile(file2)

	if err1 != nil || err2 != nil {
		return false
	}

	// 内容完全相同则认为是同一个文件
	return string(content1) == string(content2)
}
