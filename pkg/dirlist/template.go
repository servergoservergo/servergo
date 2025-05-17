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

// 支持的主题常量在 themes.go 中定义

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
	return ValidThemes
}

// 检查模板文件是否存在
func templateFileExists(path string) bool {
	// 尝试使用嵌入文件系统检查文件是否存在
	_, err := templatesFS.ReadFile(path)
	return err == nil
}

// NewDirListTemplate 创建新的目录列表模板
func NewDirListTemplate(theme string) (*DirListTemplate, error) {
	// 如果没有指定主题或主题无效，使用默认主题
	if theme == "" || !IsValidTheme(theme) {
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

	// 首先检查模板文件是否存在
	if !templateFileExists(templatePath) {
		errMsg := fmt.Sprintf("主题模板文件不存在: %s", templatePath)
		if theme != DefaultTheme {
			fmt.Printf("警告: %s，尝试回退到默认主题\n", errMsg)
			return NewDirListTemplate(DefaultTheme)
		}
		return nil, fmt.Errorf("默认主题模板文件不存在: %s", templatePath)
	}

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

	// 使用主题名称作为模板名称，这样错误信息会更明确
	tmplName := fmt.Sprintf("%s_theme", theme)

	// 创建一个新的模板
	tmpl := template.New(tmplName).Funcs(funcMap)

	// 读取模板文件内容
	content, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		errDetail := fmt.Sprintf("无法读取主题模板文件(%s): %v", templatePath, err)
		if theme != DefaultTheme {
			fmt.Printf("警告: %s，尝试回退到默认主题\n", errDetail)
			return NewDirListTemplate(DefaultTheme)
		}
		return nil, fmt.Errorf("无法读取默认主题模板文件: %s", err)
	}

	// 解析模板内容
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		errDetail := fmt.Sprintf("无法解析主题模板(%s): %v", templatePath, err)
		if theme != DefaultTheme {
			fmt.Printf("警告: %s，尝试回退到默认主题\n", errDetail)
			return NewDirListTemplate(DefaultTheme)
		}
		return nil, fmt.Errorf(i18n.Tf("dirlist.theme_load_error", errDetail))
	}

	// 额外检查是否真的成功解析了模板
	if tmpl == nil || tmpl.Templates() == nil || len(tmpl.Templates()) == 0 {
		return nil, fmt.Errorf("成功解析模板，但模板集合为空 (%s)", templatePath)
	}

	return &DirListTemplate{
		theme:    theme,
		template: tmpl,
	}, nil
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
		// 提供更详细的错误信息
		return "", fmt.Errorf(i18n.Tf("dirlist.template_not_initialized", t.theme) +
			fmt.Sprintf(" (theme=%s, tmpl=%v)", t.theme, t.template))
	}

	// 输出模板名称，以便更好地诊断问题
	templateName := "unnamed"
	if t.template != nil && t.template.Name() != "" {
		templateName = t.template.Name()
	}

	result := &strings.Builder{}
	if err := t.template.Execute(result, data); err != nil {
		return "", fmt.Errorf(i18n.Tf("dirlist.template_render_error", err) +
			fmt.Sprintf(" (theme=%s, template_name=%s)", t.theme, templateName))
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
