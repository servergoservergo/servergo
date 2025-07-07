package dirlist

import (
	"embed"
	"html/template"
	"io/fs"
)

// 使用嵌入文件系统导入模板文件
//
//go:embed templates
var templatesFS embed.FS

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
	RepoURL     string     // GitHub仓库URL
	Stars       int        // GitHub Star数量
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

// GetStaticAssets 获取静态资源文件系统
func GetStaticAssets() fs.FS {
	// 从嵌入式文件系统中获取静态资源
	subFS, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		return templatesFS
	}
	return subFS
}
