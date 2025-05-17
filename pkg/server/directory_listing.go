package server

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/utils"
)

// renderDirectoryListing 渲染目录列表页面
// 此方法为内部方法，用于生成目录浏览页面
//
// 参数:
//   - c: Gin的上下文，包含请求和响应信息
//   - fullPath: 目录在文件系统中的完整路径，例如: "/home/user/files/images"
//   - reqPath: 请求的相对路径，例如: "/images"
//
// 功能:
//  1. 读取指定目录的内容
//  2. 将目录内容组织为文件项列表
//  3. 根据模板渲染HTML或JSON格式的目录列表
//  4. 处理错误情况
func (fs *FileServer) renderDirectoryListing(c *gin.Context, fullPath, reqPath string) {
	// 读取目录内容
	files, err := os.ReadDir(fullPath)
	if err != nil {
		c.String(http.StatusInternalServerError, i18n.Tf("http.500_dir_content", err))
		return
	}

	// 创建文件项列表
	items := make([]dirlist.FileItem, 0, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		// 构建文件路径
		itemPath := filepath.Join(reqPath, file.Name())
		if !strings.HasPrefix(itemPath, "/") {
			itemPath = "/" + itemPath
		}

		// 格式化大小
		var sizeStr string
		var sizeBytes int64
		if file.IsDir() {
			sizeStr = "-"
			sizeBytes = 0
		} else {
			sizeBytes = info.Size()
			sizeStr = utils.FormatSize(sizeBytes)
		}

		// 添加到列表
		items = append(items, dirlist.FileItem{
			Name:         file.Name(),                                  // 文件名，例如: "image.jpg"
			IsDir:        file.IsDir(),                                 // 是否是目录，例如: false
			Size:         sizeStr,                                      // 人类可读的大小，例如: "2.5 MB"
			SizeBytes:    sizeBytes,                                    // 原始字节大小，例如: 2621440
			LastModified: info.ModTime().Format("2006-01-02 15:04:05"), // 格式化的修改时间
			Path:         itemPath,                                     // 访问路径，例如: "/images/image.jpg"
		})
	}

	// 按照目录在前，文件在后的方式排序
	sort.Slice(items, func(i, j int) bool {
		// 如果一个是目录一个不是，目录在前
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		// 否则按名称字母顺序排序
		return items[i].Name < items[j].Name
	})

	// 计算上级目录路径
	var parentDir string
	if reqPath != "/" {
		parentDir = filepath.Dir(reqPath)
		if parentDir == "." {
			parentDir = "/"
		}
		if !strings.HasPrefix(parentDir, "/") {
			parentDir = "/" + parentDir
		}
		if parentDir != "/" && strings.HasSuffix(parentDir, "/") {
			parentDir = parentDir[:len(parentDir)-1]
		}
	}

	// 准备模板数据
	data := dirlist.TemplateData{
		DirPath:     reqPath,                                  // 当前目录路径，例如: "/images"
		Items:       items,                                    // 文件列表
		ParentDir:   parentDir,                                // 父目录路径，例如: "/"
		CurrentTime: time.Now().Format("2006-01-02 15:04:05"), // 当前时间，用于显示在页面
	}

	// 渲染模板
	html, err := fs.dirTemplate.Render(data)
	if err != nil {
		c.String(http.StatusInternalServerError, i18n.Tf("http.500_template", err))
		return
	}

	// 使用模板提供的内容类型，支持HTML和JSON格式
	c.Header("Content-Type", fs.dirTemplate.GetContentType())
	c.String(http.StatusOK, html)
}
