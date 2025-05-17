package dirlist

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	"github.com/CC11001100/servergo/pkg/i18n"
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
	JsonTheme    = "json"  // 新增JSON主题，支持API调用
	TableTheme   = "table" // 新增彩色表格主题
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
	SizeBytes    int64  // 原始大小（字节）
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
		JsonTheme,
		TableTheme,
	}
}

// NewDirListTemplate 创建新的目录列表模板
func NewDirListTemplate(theme string) (*DirListTemplate, error) {
	// 如果没有指定主题或主题无效，使用默认主题
	if theme == "" || !isThemeValid(theme) {
		theme = DefaultTheme
	}

	// 对于特殊主题，不需要加载HTML模板文件
	if theme == JsonTheme || theme == TableTheme {
		return &DirListTemplate{
			theme:    theme,
			template: nil, // 对于特殊主题，不需要template对象
		}, nil
	}

	// 加载模板
	templatePath := fmt.Sprintf("templates/%s/index.html", theme)

	// 创建模板函数映射
	funcMap := template.FuncMap{
		"subtract": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"last": func(x int, a []FileItem) bool {
			return x == len(a)-1
		},
	}

	// 创建一个新的模板
	tmpl := template.New("").Funcs(funcMap)

	// 解析模板文件
	tmpl, err := tmpl.ParseFS(templatesFS, templatePath)
	if err != nil {
		// 如果指定主题加载失败，尝试回退到默认主题
		if theme != DefaultTheme {
			return NewDirListTemplate(DefaultTheme)
		}
		return nil, fmt.Errorf(i18n.Tf("dirlist.theme_load_error", err))
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
	// 对于特殊主题，使用特殊处理方式
	switch t.theme {
	case JsonTheme:
		return t.renderJSON(data)
	case TableTheme:
		return renderTableTheme(data)
	}

	// 对于HTML主题，使用常规模板渲染
	// 如果模板对象为nil，返回错误
	if t.template == nil {
		return "", fmt.Errorf(i18n.Tf("dirlist.template_not_initialized", t.theme))
	}

	result := &strings.Builder{}
	if err := t.template.Execute(result, data); err != nil {
		return "", fmt.Errorf(i18n.Tf("dirlist.template_render_error", err))
	}

	return result.String(), nil
}

// renderJSON 直接渲染JSON数据，避免HTML模板的限制
func (t *DirListTemplate) renderJSON(data TemplateData) (string, error) {
	// 创建特殊的JSON结构
	type jsonItem struct {
		Name          string `json:"name"`
		IsDirectory   bool   `json:"is_directory"`
		Size          int64  `json:"size"`           // 改为数值类型的字节大小
		SizeFormatted string `json:"size_formatted"` // 保留原格式化大小
		LastModified  string `json:"last_modified"`
		Path          string `json:"path"`
		URL           string `json:"url"`
	}

	type jsonData struct {
		Path            string     `json:"path"`
		Timestamp       string     `json:"timestamp"`
		ParentDirectory string     `json:"parent_directory"`
		Contents        []jsonItem `json:"contents"`
	}

	// 转换数据
	contents := make([]jsonItem, len(data.Items))
	for i, item := range data.Items {
		url := item.Path
		if item.IsDir {
			url += "/"
		}

		contents[i] = jsonItem{
			Name:          item.Name,
			IsDirectory:   item.IsDir,
			Size:          item.SizeBytes,
			SizeFormatted: item.Size,
			LastModified:  item.LastModified,
			Path:          item.Path,
			URL:           url,
		}
	}

	jsonResult := jsonData{
		Path:            data.DirPath,
		Timestamp:       data.CurrentTime,
		ParentDirectory: data.ParentDir,
		Contents:        contents,
	}

	// 序列化为JSON字符串
	jsonBytes, err := json.MarshalIndent(jsonResult, "", "    ")
	if err != nil {
		return "", fmt.Errorf(i18n.Tf("dirlist.json_marshal_error", err))
	}

	return string(jsonBytes), nil
}

// GetTheme 获取当前使用的主题
func (t *DirListTemplate) GetTheme() string {
	return t.theme
}

// GetContentType 获取输出内容类型
func (t *DirListTemplate) GetContentType() string {
	switch t.theme {
	case JsonTheme:
		return "application/json"
	case TableTheme:
		return "text/plain; charset=utf-8"
	default:
		return "text/html; charset=utf-8"
	}
}

// GetStaticAssets 获取静态资源文件系统
func GetStaticAssets() fs.FS {
	// 从嵌入式文件系统中获取静态资源
	subFS, _ := fs.Sub(templatesFS, "templates")
	return subFS
}
