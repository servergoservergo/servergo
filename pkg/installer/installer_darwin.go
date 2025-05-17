//go:build darwin
// +build darwin

package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/CC11001100/servergo/pkg/logger"
)

// MacInstaller 实现了macOS系统下的安装逻辑
type MacInstaller struct {
	// macOS通常将命令行工具安装到/usr/local/bin目录
	targetDir string
}

// NewInstaller 创建一个适用于macOS的安装器
func NewInstaller() Installer {
	return &MacInstaller{
		targetDir: "/usr/local/bin",
	}
}

// Install 在macOS系统中安装程序到PATH
func (m *MacInstaller) Install(executablePath string) error {
	logger.Info("当前可执行文件路径: %s", executablePath)

	// 检查目标目录是否存在，如果不存在则尝试创建
	if _, err := os.Stat(m.targetDir); os.IsNotExist(err) {
		logger.Info("目标目录 %s 不存在，尝试创建...", m.targetDir)
		if err := os.MkdirAll(m.targetDir, 0755); err != nil {
			return fmt.Errorf("无法创建目录 %s: %v", m.targetDir, err)
		}
	}

	// 目标路径
	targetPath := filepath.Join(m.targetDir, "servergo")

	// 检查是否已经安装
	existingInstall := false
	if _, err := os.Stat(targetPath); err == nil {
		existingInstall = true
	}

	// 如果已安装，先检测是符号链接还是实际文件
	if existingInstall {
		fileInfo, err := os.Lstat(targetPath)
		if err == nil && fileInfo.Mode()&os.ModeSymlink != 0 {
			// 是符号链接，读取链接指向
			if linkTarget, err := os.Readlink(targetPath); err == nil {
				logger.Info("检测到现有安装，符号链接指向: %s", linkTarget)
				if linkTarget == executablePath {
					logger.Info("重新安装相同位置的程序，无需更新符号链接")
					return nil
				}
			}
		}

		logger.Warning("servergo 已经安装在 %s，将更新为新版本", targetPath)
		if err := os.Remove(targetPath); err != nil {
			return fmt.Errorf("无法删除现有安装: %v", err)
		}
	}

	// 创建符号链接
	logger.Info("创建符号链接: %s -> %s", executablePath, targetPath)
	if err := os.Symlink(executablePath, targetPath); err != nil {
		// 如果创建符号链接失败（可能是权限问题），尝试使用sudo
		logger.Warning("创建符号链接失败: %v", err)
		logger.Info("尝试使用sudo创建符号链接...")

		cmd := exec.Command("sudo", "ln", "-sf", executablePath, targetPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("使用sudo创建符号链接失败: %v", err)
		}
	}

	if existingInstall {
		logger.Info("servergo 已成功更新到 %s", targetPath)
	} else {
		logger.Info("servergo 已成功安装到 %s", targetPath)
	}
	logger.Info("现在您可以在任何目录下使用 'servergo' 命令")

	return nil
}

// Uninstall 在macOS系统中从PATH移除程序
func (m *MacInstaller) Uninstall() error {
	targetPath := filepath.Join(m.targetDir, "servergo")

	// 检查是否已安装
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		logger.Warning("未找到安装的servergo: %s", targetPath)
		return nil
	}

	// 删除符号链接
	logger.Info("移除 %s", targetPath)
	if err := os.Remove(targetPath); err != nil {
		// 如果删除失败（可能是权限问题），尝试使用sudo
		logger.Warning("移除失败: %v", err)
		logger.Info("尝试使用sudo移除...")

		cmd := exec.Command("sudo", "rm", "-f", targetPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("使用sudo移除失败: %v", err)
		}
	}

	logger.Info("servergo 已成功从 %s 移除", targetPath)
	return nil
}
