package auth

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/CC11001100/servergo/pkg/i18n"
)

//go:embed web
var authFS embed.FS

// 提取auth子目录
var authFiles fs.FS

func init() {
	var err error
	// 提取web子目录
	authFiles, err = fs.Sub(authFS, "web")
	if err != nil {
		panic("无法加载认证页面资源: " + err.Error())
	}
}

// GetAuthFileSystem 返回认证页面的文件系统
func GetAuthFileSystem() http.FileSystem {
	return http.FS(authFiles)
}

// GetFileContent 从嵌入式文件系统中读取文件内容
func GetFileContent(filename string) (string, error) {
	content, err := fs.ReadFile(authFiles, filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// LoginPageData 登录页面的数据结构
type LoginPageData struct {
	Lang             string
	Title            string
	Header           string
	Subheader        string
	UsernameLabel    string
	PasswordLabel    string
	ButtonText       string
	Footer           string
	ErrorEmptyFields string
	ErrorCredentials string
}

// GetLoginHTMLContent 获取登录页面的HTML内容
func GetLoginHTMLContent() (string, error) {
	// 获取当前语言
	currentLang := i18n.GetCurrentLanguage()

	// 读取登录页面模板
	htmlContent, err := GetFileContent("login.html")
	if err != nil {
		return "", err
	}

	// 创建模板
	tmpl, err := template.New("login").Parse(htmlContent)
	if err != nil {
		return "", err
	}

	// 准备模板数据
	data := LoginPageData{
		Lang:             currentLang,
		Title:            i18n.T("login.title"),
		Header:           i18n.T("login.header"),
		Subheader:        i18n.T("login.subheader"),
		UsernameLabel:    i18n.T("login.username"),
		PasswordLabel:    i18n.T("login.password"),
		ButtonText:       i18n.T("login.button"),
		Footer:           i18n.T("login.footer"),
		ErrorEmptyFields: i18n.T("login.error.empty_fields"),
		ErrorCredentials: i18n.T("login.error.credentials"),
	}

	// 渲染模板
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
