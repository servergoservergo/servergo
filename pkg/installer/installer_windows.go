//go:build windows
// +build windows

package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/CC11001100/servergo/pkg/logger"
)

// WindowsInstaller 实现了Windows系统下的安装逻辑
type WindowsInstaller struct {
	// Windows通常使用环境变量PATH
	// 也可以安装到已在PATH中的目录，如C:\Windows\System32
	targetDir string
}

// NewInstaller 创建一个适用于Windows的安装器
func NewInstaller() Installer {
	// 默认安装到用户目录下的.servergo\bin目录
	// 这样不需要管理员权限
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果获取用户目录失败，使用当前目录
		homeDir, _ = os.Getwd()
	}

	targetDir := filepath.Join(homeDir, ".servergo", "bin")
	logger.Info("将安装到 %s 目录", targetDir)

	return &WindowsInstaller{
		targetDir: targetDir,
	}
}

// Install 在Windows系统中安装程序到PATH
func (w *WindowsInstaller) Install(executablePath string) error {
	// 检查目标目录是否存在，如果不存在则尝试创建
	if _, err := os.Stat(w.targetDir); os.IsNotExist(err) {
		logger.Info("目标目录 %s 不存在，创建中...", w.targetDir)
		if err := os.MkdirAll(w.targetDir, 0755); err != nil {
			return fmt.Errorf("无法创建目录 %s: %v", w.targetDir, err)
		}
	}

	// 目标路径
	targetPath := filepath.Join(w.targetDir, "servergo.exe")

	// 检查是否已经安装
	existingInstall := false
	backupNeeded := false
	if _, err := os.Stat(targetPath); err == nil {
		existingInstall = true

		// 检查是否与当前文件相同
		if isSameFile(executablePath, targetPath) {
			logger.Info("检测到安装的文件与当前文件相同，跳过复制")
			backupNeeded = false
		} else {
			logger.Warning("servergo 已经安装在 %s，将更新为新版本", targetPath)
			backupNeeded = true
		}
	}

	// 如果需要备份之前的版本
	if backupNeeded {
		// 创建备份目录
		backupDir := filepath.Join(w.targetDir, "backup", time.Now().Format("20060102-150405"))
		if err := os.MkdirAll(backupDir, 0755); err == nil {
			backupPath := filepath.Join(backupDir, "servergo.exe")
			logger.Info("备份现有安装到 %s", backupPath)

			// 复制文件到备份目录
			if input, err := os.ReadFile(targetPath); err == nil {
				if err = os.WriteFile(backupPath, input, 0755); err != nil {
					logger.Warning("备份失败: %v", err)
				}
			} else {
				logger.Warning("无法读取现有安装进行备份: %v", err)
			}
		} else {
			logger.Warning("无法创建备份目录: %v", err)
		}

		// 删除现有安装
		if err := os.Remove(targetPath); err != nil {
			return fmt.Errorf("无法删除现有安装: %v", err)
		}
	}

	// 如果是新安装或者需要更新
	if !existingInstall || backupNeeded {
		// 复制可执行文件
		logger.Info("复制可执行文件: %s -> %s", executablePath, targetPath)
		input, err := os.ReadFile(executablePath)
		if err != nil {
			return fmt.Errorf("无法读取源文件: %v", err)
		}
		if err = os.WriteFile(targetPath, input, 0755); err != nil {
			return fmt.Errorf("无法写入目标文件: %v", err)
		}
	}

	// 将目标目录添加到用户PATH环境变量
	if err := w.addToPath(); err != nil {
		return err
	}

	if existingInstall && backupNeeded {
		logger.Info("servergo 已成功更新到 %s", targetPath)
	} else if existingInstall {
		logger.Info("servergo 已经是最新版本")
	} else {
		logger.Info("servergo 已成功安装到 %s", targetPath)
	}
	logger.Info("请重新打开命令提示符或PowerShell窗口，然后尝试使用 'servergo' 命令")

	return nil
}

// Uninstall 在Windows系统中从PATH移除程序
func (w *WindowsInstaller) Uninstall() error {
	targetPath := filepath.Join(w.targetDir, "servergo.exe")

	// 检查是否已安装
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		logger.Warning("未找到安装的servergo: %s", targetPath)
		return nil
	}

	// 删除可执行文件
	logger.Info("移除 %s", targetPath)
	if err := os.Remove(targetPath); err != nil {
		return fmt.Errorf("无法移除文件: %v", err)
	}

	// 从PATH中移除目标目录
	if err := w.removeFromPath(); err != nil {
		return err
	}

	logger.Info("servergo 已成功卸载")
	logger.Info("请重新打开命令提示符或PowerShell窗口使更改生效")

	return nil
}

// addToPath 将目标目录添加到用户的PATH环境变量
func (w *WindowsInstaller) addToPath() error {
	// 获取当前的PATH环境变量
	pathEnv, exists := os.LookupEnv("PATH")
	if !exists {
		return fmt.Errorf("无法获取PATH环境变量")
	}

	// 检查目标目录是否已经在PATH中
	paths := strings.Split(pathEnv, string(os.PathListSeparator))
	for _, path := range paths {
		if strings.EqualFold(path, w.targetDir) {
			logger.Info("目标目录已在PATH中")
			return nil
		}
	}

	// 将目标目录添加到PATH中
	logger.Info("将 %s 添加到PATH环境变量", w.targetDir)

	// 使用setx命令更新用户PATH环境变量
	// 注意：setx有长度限制，如果PATH太长可能会被截断
	cmd := exec.Command("setx", "PATH", pathEnv+string(os.PathListSeparator)+w.targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("无法更新PATH环境变量: %v", err)
	}

	return nil
}

// removeFromPath 从用户PATH环境变量中移除目标目录
func (w *WindowsInstaller) removeFromPath() error {
	// 获取当前的PATH环境变量
	pathEnv, exists := os.LookupEnv("PATH")
	if !exists {
		return fmt.Errorf("无法获取PATH环境变量")
	}

	// 分割PATH环境变量
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	// 移除目标目录
	var newPaths []string
	for _, path := range paths {
		if !strings.EqualFold(path, w.targetDir) {
			newPaths = append(newPaths, path)
		}
	}

	// 如果目录不在PATH中，则无需修改
	if len(paths) == len(newPaths) {
		logger.Info("目标目录不在PATH中，无需修改")
		return nil
	}

	// 使用setx命令更新用户PATH环境变量
	newPathEnv := strings.Join(newPaths, string(os.PathListSeparator))
	logger.Info("从PATH环境变量移除 %s", w.targetDir)

	cmd := exec.Command("setx", "PATH", newPathEnv)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("无法更新PATH环境变量: %v", err)
	}

	return nil
}
