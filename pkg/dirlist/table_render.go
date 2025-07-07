package dirlist

import (
	"fmt"
	"strings"
)

// renderTableTheme 渲染表格主题格式的目录列表
// 该函数生成一个简洁的文本表格，适合在终端或简单文本界面中显示
func renderTableTheme(data TemplateData) (string, error) {
	if len(data.Items) == 0 {
		return "目录为空", nil
	}

	// 计算各列的最大宽度，确保表格整齐对齐
	maxNameWidth := 4 // "名称" 列的最小宽度
	maxSizeWidth := 4 // "大小" 列的最小宽度
	maxTimeWidth := 8 // "修改时间" 列的最小宽度

	// 遍历所有项目，找出每列的最大宽度
	for _, item := range data.Items {
		// 计算名称列宽度（考虑中文字符）
		nameWidth := getDisplayWidth(item.Name)
		if nameWidth > maxNameWidth {
			maxNameWidth = nameWidth
		}

		// 计算大小列宽度
		sizeWidth := len(item.Size)
		if sizeWidth > maxSizeWidth {
			maxSizeWidth = sizeWidth
		}

		// 计算时间列宽度
		timeWidth := len(item.LastModified)
		if timeWidth > maxTimeWidth {
			maxTimeWidth = timeWidth
		}
	}

	// 构建表格
	var result strings.Builder

	// 写入表头
	headerFormat := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%s\n",
		maxNameWidth, maxSizeWidth, maxTimeWidth)
	result.WriteString(fmt.Sprintf(headerFormat, "名称", "大小", "修改时间", "类型"))

	// 写入分割线
	separator := strings.Repeat("-", maxNameWidth) + "  " +
		strings.Repeat("-", maxSizeWidth) + "  " +
		strings.Repeat("-", maxTimeWidth) + "  " +
		strings.Repeat("-", 4) + "\n"
	result.WriteString(separator)

	// 写入数据行
	for _, item := range data.Items {
		itemType := "文件"
		if item.IsDir {
			itemType = "目录"
		}

		// 确保名称列正确对齐（处理中文字符）
		nameField := padString(item.Name, maxNameWidth)

		rowFormat := fmt.Sprintf("%%s  %%-%ds  %%-%ds  %%s\n",
			maxSizeWidth, maxTimeWidth)
		result.WriteString(fmt.Sprintf(rowFormat,
			nameField, item.Size, item.LastModified, itemType))
	}

	return result.String(), nil
}

// getDisplayWidth 计算字符串的显示宽度
// 中文字符占用2个显示位置，英文字符占用1个显示位置
func getDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r > 127 {
			// 非ASCII字符（如中文）占用2个位置
			width += 2
		} else {
			// ASCII字符占用1个位置
			width += 1
		}
	}
	return width
}

// padString 在字符串右侧填充空格，使其达到指定的显示宽度
// 这个函数处理中文字符的对齐问题
func padString(s string, width int) string {
	currentWidth := getDisplayWidth(s)
	if currentWidth >= width {
		return s
	}

	// 计算需要填充的空格数
	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}
