package dirlist

import (
	"fmt"
	"strings"

	"github.com/CC11001100/servergo/pkg/github"
	"github.com/CC11001100/servergo/pkg/i18n"
)

// Render 渲染目录列表模板
func (t *DirListTemplate) Render(data TemplateData) (string, error) {
	// 获取GitHub统计信息
	if stats, err := github.GetStats(); err == nil {
		data.Stars = stats.Stars
	}
	data.RepoURL = github.GetRepoURL()

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
