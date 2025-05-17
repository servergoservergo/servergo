package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// handleFileRequest 处理文件请求
// 此方法为内部方法，用作Gin的处理函数
//
// 参数:
//   - c: Gin的上下文，包含请求和响应信息
//
// 功能:
//  1. 处理静态文件请求
//  2. 如果请求的是目录且启用了目录列表，则显示目录内容
//  3. 如果请求的是文件，则直接提供文件下载
//  4. 处理各种错误情况，如文件不存在或无权限
func (fs *FileServer) handleFileRequest(c *gin.Context) {
	// 获取请求路径
	reqPath := c.Request.URL.Path

	// 如果请求的是内部静态资源，跳过处理
	if strings.HasPrefix(reqPath, "/_servergo_assets/") {
		c.Next()
		return
	}

	// 处理路径中可能的URL编码
	// 例如: 将 "%20" 转换为空格
	reqPath = strings.Replace(reqPath, "%20", " ", -1)

	// 确保路径不会超出根目录
	// 例如: "../config" 会被转换为 "/config"
	cleanPath := filepath.Clean(reqPath)
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	// 将请求路径转换为服务器文件系统上的实际路径
	fullPath := filepath.Join(fs.absDir, cleanPath)

	// 获取文件状态
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		// 文件不存在，返回404
		c.String(http.StatusNotFound, "404 Not Found: %s", reqPath)
		return
	}

	// 如果是目录，检查是否启用了目录列表功能
	if fileInfo.IsDir() {
		// 检查该目录下是否有index.html文件
		indexPath := filepath.Join(fullPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			// 如果存在index.html，则提供该文件
			c.File(indexPath)
			return
		}

		// 如果启用了目录列表功能，则显示目录内容
		if fs.config.EnableDirListing {
			fs.renderDirectoryListing(c, fullPath, reqPath)
			return
		}

		// 未启用目录列表功能，返回403禁止访问
		c.String(http.StatusForbidden, "403 Forbidden: Directory listing disabled")
		return
	}

	// 如果是文件，则提供该文件
	c.File(fullPath)
}
