package dirlist

// 支持的主题常量
const (
	DefaultTheme = "default" // 默认主题
	DarkTheme    = "dark"    // 深色主题
	BlueTheme    = "blue"    // 蓝色主题
	GreenTheme   = "green"   // 绿色主题
	RetroTheme   = "retro"   // 复古主题
	JsonTheme    = "json"    // JSON格式（非HTML主题）
	TableTheme   = "table"   // 表格格式（非HTML主题）
)

// ValidThemes 是所有支持的目录列表主题
var ValidThemes = []string{
	DefaultTheme,
	DarkTheme,
	BlueTheme,
	GreenTheme,
	RetroTheme,
	JsonTheme,
	TableTheme,
}

// IsValidTheme 检查提供的主题是否有效
func IsValidTheme(theme string) bool {
	for _, validTheme := range ValidThemes {
		if theme == validTheme {
			return true
		}
	}
	return false
}
