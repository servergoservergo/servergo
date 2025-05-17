package dirlist

import (
	"bytes"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 生成彩色表格形式的目录列表
func renderTableTheme(data TemplateData) (string, error) {
	// 创建表格对象
	t := table.NewWriter()
	var buf bytes.Buffer
	t.SetOutputMirror(&buf)

	// 设置表格样式
	t.SetStyle(table.StyleColoredBright)

	// 设置表头
	t.AppendHeader(table.Row{
		"名称",
		"类型",
		"大小",
		"修改时间",
	})

	// 如果有上级目录，添加一行用于返回上级目录
	if data.ParentDir != "" {
		t.AppendRow(table.Row{
			text.Colors{text.FgBlue, text.Bold}.Sprint(".."),
			text.Colors{text.FgBlue, text.Bold}.Sprint("目录"),
			"-",
			"-",
		})
	}

	// 添加文件和目录项
	for _, item := range data.Items {
		var nameColored, typeColored string

		if item.IsDir {
			// 目录用蓝色加粗显示
			nameColored = text.Colors{text.FgBlue, text.Bold}.Sprint(item.Name + "/")
			typeColored = text.Colors{text.FgBlue}.Sprint("目录")
		} else {
			// 根据扩展名设置不同颜色
			switch getFileType(item.Name) {
			case "image":
				nameColored = text.Colors{text.FgMagenta}.Sprint(item.Name)
				typeColored = text.Colors{text.FgMagenta}.Sprint("图片")
			case "video":
				nameColored = text.Colors{text.FgCyan}.Sprint(item.Name)
				typeColored = text.Colors{text.FgCyan}.Sprint("视频")
			case "audio":
				nameColored = text.Colors{text.FgYellow}.Sprint(item.Name)
				typeColored = text.Colors{text.FgYellow}.Sprint("音频")
			case "archive":
				nameColored = text.Colors{text.FgRed}.Sprint(item.Name)
				typeColored = text.Colors{text.FgRed}.Sprint("压缩")
			case "code":
				nameColored = text.Colors{text.FgGreen}.Sprint(item.Name)
				typeColored = text.Colors{text.FgGreen}.Sprint("代码")
			default:
				nameColored = item.Name
				typeColored = "文件"
			}
		}

		t.AppendRow(table.Row{
			nameColored,
			typeColored,
			item.Size,
			item.LastModified,
		})
	}

	// 设置列配置，调整对齐方式
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMax: 40},
		{Number: 2, Align: text.AlignCenter, AlignHeader: text.AlignCenter, WidthMax: 10},
		{Number: 3, Align: text.AlignRight, AlignHeader: text.AlignCenter, WidthMax: 15},
		{Number: 4, Align: text.AlignCenter, AlignHeader: text.AlignCenter, WidthMax: 20},
	})

	// 添加表脚显示当前路径和时间
	t.AppendFooter(table.Row{"当前路径", data.DirPath, "时间", data.CurrentTime})

	// 渲染表格
	t.Render()

	return buf.String(), nil
}

// 根据文件名判断文件类型，返回类型名称
func getFileType(filename string) string {
	// 图片文件
	if isFileType(filename, []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}) {
		return "image"
	}
	// 视频文件
	if isFileType(filename, []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv"}) {
		return "video"
	}
	// 音频文件
	if isFileType(filename, []string{".mp3", ".wav", ".ogg", ".flac", ".aac"}) {
		return "audio"
	}
	// 压缩文件
	if isFileType(filename, []string{".zip", ".rar", ".7z", ".tar", ".gz", ".bz2"}) {
		return "archive"
	}
	// 代码/文本文件
	if isFileType(filename, []string{".go", ".c", ".cpp", ".java", ".py", ".js", ".html", ".css", ".php", ".sh", ".txt", ".md"}) {
		return "code"
	}
	// 默认类型
	return "other"
}

// 判断文件是否属于指定类型
func isFileType(filename string, extensions []string) bool {
	lowName := strings.ToLower(filename)
	for _, ext := range extensions {
		if strings.HasSuffix(lowName, ext) {
			return true
		}
	}
	return false
}

// TableDirList 通过彩色表格显示目录内容
func TableDirList(dirPath string, items []FileItem) {
	// 创建表格对象
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表格样式
	t.SetStyle(table.StyleColoredBright)

	// 设置表头
	t.AppendHeader(table.Row{
		"名称",
		"类型",
		"大小",
		"修改时间",
	})

	// 如果不是根目录，添加返回上级目录行
	if dirPath != "/" {
		t.AppendRow(table.Row{
			text.Colors{text.FgBlue, text.Bold}.Sprint(".."),
			text.Colors{text.FgBlue, text.Bold}.Sprint("目录"),
			"-",
			"-",
		})
	}

	// 添加文件和目录项
	for _, item := range items {
		var nameColored, typeColored string

		if item.IsDir {
			// 目录用蓝色加粗显示
			nameColored = text.Colors{text.FgBlue, text.Bold}.Sprint(item.Name + "/")
			typeColored = text.Colors{text.FgBlue}.Sprint("目录")
		} else {
			// 根据扩展名设置不同颜色
			switch getFileType(item.Name) {
			case "image":
				nameColored = text.Colors{text.FgMagenta}.Sprint(item.Name)
				typeColored = text.Colors{text.FgMagenta}.Sprint("图片")
			case "video":
				nameColored = text.Colors{text.FgCyan}.Sprint(item.Name)
				typeColored = text.Colors{text.FgCyan}.Sprint("视频")
			case "audio":
				nameColored = text.Colors{text.FgYellow}.Sprint(item.Name)
				typeColored = text.Colors{text.FgYellow}.Sprint("音频")
			case "archive":
				nameColored = text.Colors{text.FgRed}.Sprint(item.Name)
				typeColored = text.Colors{text.FgRed}.Sprint("压缩")
			case "code":
				nameColored = text.Colors{text.FgGreen}.Sprint(item.Name)
				typeColored = text.Colors{text.FgGreen}.Sprint("代码")
			default:
				nameColored = item.Name
				typeColored = "文件"
			}
		}

		t.AppendRow(table.Row{
			nameColored,
			typeColored,
			item.Size,
			item.LastModified,
		})
	}

	// 设置列配置
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMax: 40},
		{Number: 2, Align: text.AlignCenter, AlignHeader: text.AlignCenter, WidthMax: 10},
		{Number: 3, Align: text.AlignRight, AlignHeader: text.AlignCenter, WidthMax: 15},
		{Number: 4, Align: text.AlignCenter, AlignHeader: text.AlignCenter, WidthMax: 20},
	})

	// 添加表脚
	t.AppendFooter(table.Row{
		"总计",
		len(items) + func() int {
			if dirPath != "/" && len(items) > 0 {
				return 1
			} else {
				return 0
			}
		}(), // 包含可能的 ".." 条目
		"目录时间",
		time.Now().Format("2006-01-02 15:04:05"),
	})

	// 渲染表格
	t.Render()
}
