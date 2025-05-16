package auth

import (
	"embed"
	"io/fs"
	"net/http"
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

// GetLoginHTMLContent 获取登录页面的HTML内容
func GetLoginHTMLContent() (string, error) {
	return GetFileContent("login.html")
}
