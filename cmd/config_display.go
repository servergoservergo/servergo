package cmd

import (
	"fmt"
	"os"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// formatBoolValue 格式化布尔值为更友好的显示文本
func formatBoolValue(value bool) string {
	if value {
		return i18n.T("config.enabled")
	}
	return i18n.T("config.disabled")
}

// 格式化语言值，显示友好的语言名称
func formatLanguageValue(lang string) string {
	return i18n.GetLanguageDisplayName(lang)
}

// UpdateConfigTableHeaders 更新配置表格的表头为当前语言
func UpdateConfigTableHeaders(t table.Writer) {
	t.ResetHeaders()

	// 强制重新获取当前语言的翻译
	itemHeader := i18n.T("config.item")
	valueHeader := i18n.T("config.current_value")
	descHeader := i18n.T("config.description")

	t.AppendHeader(table.Row{itemHeader, valueHeader, descHeader})
}

// displayConfigTable 显示配置表格
func displayConfigTable(cfg config.Config) {
	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表格样式
	t.SetStyle(table.StyleColoredBright)

	// 设置表头 - 确保使用当前语言的翻译
	UpdateConfigTableHeaders(t)

	// 添加配置信息行 - 确保使用当前语言的翻译
	t.AppendRows([]table.Row{
		{"auto-open", formatBoolValue(cfg.AutoOpen), i18n.T("config.auto_open_desc")},
		{"enable-dir-listing", formatBoolValue(cfg.EnableDirListing), i18n.T("config.enable_dir_listing_desc")},
		{"theme", cfg.Theme, i18n.T("config.theme_desc")},
		{"language", formatLanguageValue(cfg.Language), i18n.T("config.language_desc")},
		{"enable-log-persistence", formatBoolValue(cfg.EnableLogPersistence), i18n.T("config.enable_log_persistence_desc")},
		{"start-port", fmt.Sprintf("%d", cfg.StartPort), i18n.T("config.start_port_desc")},
	})

	// 设置列对齐方式
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMax: 30},
		{Number: 2, Align: text.AlignCenter, AlignHeader: text.AlignCenter, WidthMax: 20},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignCenter},
	})

	// 输出表格
	t.Render()
}
