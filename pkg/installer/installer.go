package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/CC11001100/servergo/pkg/logger"
)

// 当前版本号 - 这个应该从项目的某个地方导入
// 这里临时硬编码，实际应该从项目版本配置中获取
const CurrentVersion = "1.0.0"

// Installer 接口定义了安装程序的方法
type Installer interface {
	// Install 将程序安装到系统PATH中
	Install(executablePath string) error
	// Uninstall 从系统PATH中移除程序
	Uninstall() error
}

// GetInstaller 根据当前操作系统返回对应的安装器
func GetInstaller() (Installer, error) {
	switch runtime.GOOS {
	case "darwin":
		return NewMacInstaller(), nil
	case "windows":
		return NewWindowsInstaller(), nil
	case "linux":
		return NewLinuxInstaller(), nil
	default:
		return nil, fmt.Errorf("暂不支持的操作系统: %s", runtime.GOOS)
	}
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

	logger.Info("当前可执行文件路径: %s", execPath)
	logger.Info("当前版本: %s", CurrentVersion)

	// 检查是否已安装
	installedVer, err := GetInstalledVersion()
	if err != nil {
		logger.Warning("检查已安装版本时出错: %v", err)
	}

	// 获取对应操作系统的安装器
	installer, err := GetInstaller()
	if err != nil {
		return err
	}

	// 如果已经安装了相同版本，询问用户是否继续
	if installedVer != nil {
		compareResult := CompareVersions(installedVer.Version, CurrentVersion)
		switch compareResult {
		case 0:
			logger.Info("检测到已安装相同版本 %s，将重新安装...", CurrentVersion)
		case -1:
			logger.Info("检测到已安装旧版本 %s，将升级到新版本 %s...", installedVer.Version, CurrentVersion)
		case 1:
			logger.Warning("警告: 当前版本 %s 低于已安装版本 %s，降级可能导致兼容性问题", CurrentVersion, installedVer.Version)
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
	if err := SaveVersionInfo(CurrentVersion, execPath); err != nil {
		logger.Warning("保存版本信息失败: %v", err)
	}

	return nil
}

// UninstallFromPath 执行卸载过程
func UninstallFromPath() error {
	// 获取对应操作系统的安装器
	installer, err := GetInstaller()
	if err != nil {
		return err
	}

	// 获取已安装版本信息
	installedVer, _ := GetInstalledVersion()
	if installedVer != nil {
		logger.Info("正在卸载版本 %s...", installedVer.Version)
	}

	// 执行卸载
	return installer.Uninstall()
}
