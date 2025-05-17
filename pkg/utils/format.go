package utils

import (
	"fmt"

	"github.com/CC11001100/servergo/pkg/i18n"
)

// FormatSize 格式化文件大小为人类可读格式
//
// 将字节大小转换为适当的单位（B, KB, MB, GB, TB）
// 对于大于等于1KB的大小，保留一位小数
//
// 示例:
//   - FormatSize(1024) => "1.0 KB"
//   - FormatSize(1048576) => "1.0 MB"
//   - FormatSize(800) => "800 B"
func FormatSize(size int64) string {
	if size < 0 {
		return fmt.Sprintf("%d %s", size, i18n.T("format.size_byte"))
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d %s", size, i18n.T("format.size_byte"))
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{
		i18n.T("format.size_kb"),
		i18n.T("format.size_mb"),
		i18n.T("format.size_gb"),
		i18n.T("format.size_tb"),
	}
	return fmt.Sprintf("%.1f %s", float64(size)/float64(div), units[exp])
}
