package server

import (
	"fmt"
)

// formatSize 格式化文件大小
// 将字节大小转换为人类易读的格式
//
// 参数:
//   - size: 文件大小，单位为字节
//
// 返回值:
//   - string: 格式化后的大小字符串
//
// 示例:
//   - formatSize(1024) => "1.0 KB"
//   - formatSize(1048576) => "1.0 MB"
//   - formatSize(800) => "800 B"
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
