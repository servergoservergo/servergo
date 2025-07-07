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

	// 新增现代设计系列
	ModernTheme   = "modern"   // 现代简约主题
	MaterialTheme = "material" // Material Design主题
	MinimalTheme  = "minimal"  // 极简主题
	GlassTheme    = "glass"    // 玻璃主题

	// 自然主题系列
	OceanTheme  = "ocean"  // 海洋主题
	ForestTheme = "forest" // 森林主题
	SunsetTheme = "sunset" // 日落主题
	AutumnTheme = "autumn" // 秋天主题
	WinterTheme = "winter" // 冬天主题
	SpringTheme = "spring" // 春天主题
	SummerTheme = "summer" // 夏天主题

	// 科技风格系列
	CyberpunkTheme = "cyberpunk" // 赛博朋克主题
	NeonTheme      = "neon"      // 霓虹主题
	MatrixTheme    = "matrix"    // 矩阵主题
	TerminalTheme  = "terminal"  // 终端主题
	SpaceTheme     = "space"     // 太空主题

	// 色彩主题系列
	NeonBlueTheme   = "neon-blue"  // 蓝色霓虹主题
	NeonPinkTheme   = "neon-pink"  // 粉色霓虹主题
	GradientTheme   = "gradient"   // 渐变主题
	MonochromeTheme = "monochrome" // 单色主题

	// 环境主题系列
	ArcticTheme  = "arctic"  // 北极主题
	DesertTheme  = "desert"  // 沙漠主题
	VolcanoTheme = "volcano" // 火山主题
	GalaxyTheme  = "galaxy"  // 银河主题

	// 风格主题系列
	VintageTheme   = "vintage"   // 复古主题
	CorporateTheme = "corporate" // 企业主题
	PaperTheme     = "paper"     // 纸质主题
	BootstrapTheme = "bootstrap" // Bootstrap主题
	NatureTheme    = "nature"    // 自然主题

	// 额外主题
	TechnologyTheme = "technology" // 科技主题
	ElegantTheme    = "elegant"    // 优雅主题
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

	// 新增主题
	ModernTheme,
	MaterialTheme,
	MinimalTheme,
	GlassTheme,
	OceanTheme,
	ForestTheme,
	SunsetTheme,
	AutumnTheme,
	WinterTheme,
	SpringTheme,
	SummerTheme,
	CyberpunkTheme,
	NeonTheme,
	MatrixTheme,
	TerminalTheme,
	SpaceTheme,
	NeonBlueTheme,
	NeonPinkTheme,
	GradientTheme,
	MonochromeTheme,
	ArcticTheme,
	DesertTheme,
	VolcanoTheme,
	GalaxyTheme,
	VintageTheme,
	CorporateTheme,
	PaperTheme,
	BootstrapTheme,
	NatureTheme,
	TechnologyTheme,
	ElegantTheme,
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
