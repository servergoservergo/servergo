package dirlist

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"strings"
)

// 使用嵌入文件系统导入模板文件
//
//go:embed templates
var templatesFS embed.FS

// 支持的主题
const (
	DefaultTheme = "default"
	DarkTheme    = "dark"
	BlueTheme    = "blue"  // 新增蓝色主题
	GreenTheme   = "green" // 新增绿色主题
	RetroTheme   = "retro" // 新增复古主题
)

// DirListTemplate 表示目录列表模板
type DirListTemplate struct {
	theme    string
	template *template.Template
}

// 模板数据
type TemplateData struct {
	DirPath     string     // 当前目录路径
	Items       []FileItem // 文件/目录列表
	ParentDir   string     // 上级目录路径
	CurrentTime string     // 当前时间
}

// 文件或目录项
type FileItem struct {
	Name         string // 文件名称
	IsDir        bool   // 是否是目录
	Size         string // 格式化后的大小
	LastModified string // 修改时间
	Path         string // 文件相对路径
}

// 返回所有支持的主题列表
func GetSupportedThemes() []string {
	return []string{
		DefaultTheme,
		DarkTheme,
		BlueTheme,
		GreenTheme,
		RetroTheme,
	}
}

// NewDirListTemplate 创建新的目录列表模板
func NewDirListTemplate(theme string) (*DirListTemplate, error) {
	// 如果没有指定主题或主题无效，使用默认主题
	if theme == "" || !isThemeValid(theme) {
		theme = DefaultTheme
	}

	// 加载模板
	templatePath := fmt.Sprintf("templates/%s/index.html", theme)

	// 创建一个新的模板
	// 不需要指定模板名称，直接使用ParseFS解析文件
	tmpl, err := template.ParseFS(templatesFS, templatePath)
	if err != nil {
		// 如果指定主题加载失败，尝试回退到默认主题
		if theme != DefaultTheme {
			return NewDirListTemplate(DefaultTheme)
		}
		return nil, fmt.Errorf("无法加载主题模板: %v", err)
	}

	return &DirListTemplate{
		theme:    theme,
		template: tmpl,
	}, nil
}

// 检查主题是否有效
func isThemeValid(theme string) bool {
	for _, validTheme := range GetSupportedThemes() {
		if theme == validTheme {
			return true
		}
	}
	return false
}

// Render 渲染目录列表模板
func (t *DirListTemplate) Render(data TemplateData) (string, error) {
	// 使用模板渲染数据
	result := &strings.Builder{}
	if err := t.template.Execute(result, data); err != nil {
		return "", fmt.Errorf("渲染模板失败: %v", err)
	}

	return result.String(), nil
}

// GetTheme 获取当前使用的主题
func (t *DirListTemplate) GetTheme() string {
	return t.theme
}

// GetStaticAssets 获取静态资源文件系统
func GetStaticAssets() fs.FS {
	// 从嵌入式文件系统中获取静态资源
	subFS, _ := fs.Sub(templatesFS, "templates")
	return subFS
}
