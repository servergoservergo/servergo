package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/CC11001100/servergo/pkg/logger"
)

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

	// 获取对应操作系统的安装器
	installer, err := GetInstaller()
	if err != nil {
		return err
	}

	// 执行安装
	return installer.Install(execPath)
}

// UninstallFromPath 执行卸载过程
func UninstallFromPath() error {
	// 获取对应操作系统的安装器
	installer, err := GetInstaller()
	if err != nil {
		return err
	}

	// 执行卸载
	return installer.Uninstall()
}
