package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/CC11001100/servergo/pkg/logger"
)

// LinuxInstaller 实现了Linux系统下的安装逻辑
type LinuxInstaller struct {
	// Linux通常将命令行工具安装到/usr/local/bin目录
	targetDir string
}

// NewLinuxInstaller 创建一个新的Linux安装器
func NewLinuxInstaller() *LinuxInstaller {
	return &LinuxInstaller{
		targetDir: "/usr/local/bin",
	}
}

// Install 在Linux系统中安装程序到PATH
func (l *LinuxInstaller) Install(executablePath string) error {
	// 检查目标目录是否存在，如果不存在则尝试创建
	if _, err := os.Stat(l.targetDir); os.IsNotExist(err) {
		logger.Info("目标目录 %s 不存在，尝试创建...", l.targetDir)

		// 使用sudo创建目录（因为/usr/local/bin通常需要root权限）
		cmd := exec.Command("sudo", "mkdir", "-p", l.targetDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("无法创建目录 %s: %v", l.targetDir, err)
		}
	}

	// 目标路径
	targetPath := filepath.Join(l.targetDir, "servergo")

	// 检查是否已经安装
	if _, err := os.Stat(targetPath); err == nil {
		logger.Warning("servergo 已经安装在 %s，将覆盖现有安装", targetPath)

		// 使用sudo删除现有文件
		cmd := exec.Command("sudo", "rm", "-f", targetPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("无法删除现有安装: %v", err)
		}
	}

	// 创建符号链接
	logger.Info("创建符号链接: %s -> %s", executablePath, targetPath)

	// 使用sudo创建符号链接
	cmd := exec.Command("sudo", "ln", "-sf", executablePath, targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("使用sudo创建符号链接失败: %v", err)
	}

	logger.Info("servergo 已成功安装到 %s", targetPath)
	logger.Info("现在您可以在任何目录下使用 'servergo' 命令")

	return nil
}

// Uninstall 在Linux系统中从PATH移除程序
func (l *LinuxInstaller) Uninstall() error {
	targetPath := filepath.Join(l.targetDir, "servergo")

	// 检查是否已安装
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		logger.Warning("未找到安装的servergo: %s", targetPath)
		return nil
	}

	// 删除符号链接
	logger.Info("移除 %s", targetPath)

	// 使用sudo删除符号链接
	cmd := exec.Command("sudo", "rm", "-f", targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("使用sudo移除失败: %v", err)
	}

	logger.Info("servergo 已成功从 %s 移除", targetPath)
	return nil
}
