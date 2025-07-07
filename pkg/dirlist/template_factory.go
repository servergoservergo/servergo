package dirlist

import (
	"fmt"
	"html/template"

	"github.com/CC11001100/servergo/pkg/i18n"
)

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
