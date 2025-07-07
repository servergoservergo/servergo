package dirlist

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/CC11001100/servergo/pkg/utils"
)

// sortFileItems 对文件项目进行排序
// 排序规则：目录在前，文件在后；同类型按名称字母顺序排序
func sortFileItems(items []FileItem) {
	sort.Slice(items, func(i, j int) bool {
		// 如果一个是目录，一个是文件，目录排在前面
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		// 同类型按名称排序（忽略大小写）
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})
}

// formatFileItem 格式化单个文件项目信息
// 将文件的原始信息转换为适合显示的格式
func formatFileItem(name string, isDir bool, size int64, modTime time.Time, basePath string) FileItem {
	// 格式化文件大小
	var sizeStr string
	if isDir {
		sizeStr = "-" // 目录不显示大小
	} else {
		sizeStr = utils.FormatSize(size)
	}

	// 格式化修改时间
	timeStr := modTime.Format("2006-01-02 15:04")

	// 构建文件的相对路径
	relativePath := filepath.Join(basePath, name)
	// 在Windows上也使用正斜杠，保持一致性
	relativePath = filepath.ToSlash(relativePath)

	return FileItem{
		Name:         name,
		IsDir:        isDir,
		Size:         sizeStr,
		SizeBytes:    size,
		LastModified: timeStr,
		Path:         relativePath,
	}
}

// buildBreadcrumb 构建面包屑导航路径
// 将目录路径转换为面包屑形式的导航链接
func buildBreadcrumb(currentPath string) string {
	if currentPath == "" || currentPath == "/" {
		return "/"
	}

	// 标准化路径分隔符
	currentPath = filepath.ToSlash(currentPath)
	if strings.HasPrefix(currentPath, "/") {
		currentPath = currentPath[1:]
	}

	parts := strings.Split(currentPath, "/")
	var breadcrumbs []string

	// 根目录链接
	breadcrumbs = append(breadcrumbs, `<a href="/">/</a>`)

	// 构建每一级的链接
	currentDir := ""
	for i, part := range parts {
		if part == "" {
			continue
		}

		currentDir = filepath.Join(currentDir, part)
		currentDir = filepath.ToSlash(currentDir)

		if i == len(parts)-1 {
			// 最后一级不加链接，只显示文本
			breadcrumbs = append(breadcrumbs, part)
		} else {
			// 中间级别加上链接
			breadcrumbs = append(breadcrumbs, fmt.Sprintf(`<a href="/%s/">%s</a>`, currentDir, part))
		}
	}

	return strings.Join(breadcrumbs, " / ")
}

// calculateTotalSize 计算目录中所有文件的总大小
// 不包括子目录的大小，只计算当前目录下的文件
func calculateTotalSize(items []FileItem) (int64, int) {
	var totalSize int64
	fileCount := 0

	for _, item := range items {
		if !item.IsDir {
			totalSize += item.SizeBytes
			fileCount++
		}
	}

	return totalSize, fileCount
}

// generateSummaryInfo 生成目录摘要信息
// 返回目录中文件和子目录的统计信息
func generateSummaryInfo(items []FileItem) string {
	if len(items) == 0 {
		return "空目录"
	}

	dirCount := 0
	fileCount := 0
	totalSize, _ := calculateTotalSize(items)

	for _, item := range items {
		if item.IsDir {
			dirCount++
		} else {
			fileCount++
		}
	}

	var parts []string

	if dirCount > 0 {
		parts = append(parts, fmt.Sprintf("%d个目录", dirCount))
	}

	if fileCount > 0 {
		parts = append(parts, fmt.Sprintf("%d个文件", fileCount))
		if totalSize > 0 {
			parts = append(parts, fmt.Sprintf("总大小: %s", utils.FormatSize(totalSize)))
		}
	}

	return strings.Join(parts, ", ")
}
