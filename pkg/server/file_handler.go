package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
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

	// 调试日志
	logger.Info("[DEBUG] Original reqPath: %s", reqPath)
	logger.Info("[DEBUG] RequestURI: %s", c.Request.RequestURI)

	// 如果请求的是内部静态资源，跳过处理
	if strings.HasPrefix(reqPath, "/_servergo_assets/") {
		c.Next()
		return
	}

	// 处理路径中可能的URL编码
	// 例如: 将 "%20" 转换为空格
	reqPath = strings.Replace(reqPath, "%20", " ", -1)

	logger.Info("[DEBUG] After URL decode reqPath: %s", reqPath)

	// 确保路径不会超出根目录
	// 使用 filepath.Clean 清理路径，移除多余的 . 和 .. 元素
	cleanPath := filepath.Clean(reqPath)
	logger.Info("[DEBUG] Clean path: %s", cleanPath)
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	// 构建完整的文件路径
	fullPath := filepath.Join(fs.absDir, cleanPath)
	logger.Info("[DEBUG] Full path: %s", fullPath)

	// 安全检查：确保解析后的路径仍在根目录内
	// 使用 filepath.Rel 检查完整路径相对于根目录的位置
	relPath, err := filepath.Rel(fs.absDir, fullPath)
	logger.Info("[DEBUG] Relative path: %s, error: %v", relPath, err)
	if err != nil || strings.HasPrefix(relPath, "..") || strings.Contains(relPath, "/..") {
		logger.Info("[DEBUG] Path traversal attack detected! Relative path: %s", relPath)
		c.String(http.StatusForbidden, i18n.T("http.403"))
		return
	}

	// 主要安全检查已通过，移除有问题的第二个检查
	logger.Info("[DEBUG] Path security check passed for: %s", relPath)

	// 获取文件状态
	fileInfo, err := os.Lstat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, i18n.T("http.404"), reqPath)
			return
		}
		c.String(http.StatusInternalServerError, i18n.T("http.500"))
		return
	}

	// 安全检查：防止符号链接路径遍历攻击
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		// 解析符号链接的真实路径
		realPath, err := filepath.EvalSymlinks(fullPath)
		if err != nil {
			logger.Info("[DEBUG] Failed to resolve symlink: %s, error: %v", fullPath, err)
			c.String(http.StatusForbidden, i18n.T("http.403"))
			return
		}
		logger.Info("[DEBUG] Symlink detected. Real path: %s", realPath)

		// 检查符号链接指向的真实路径是否在根目录内
		realRelPath, err := filepath.Rel(fs.absDir, realPath)
		if err != nil || strings.HasPrefix(realRelPath, "..") || strings.Contains(realRelPath, "/..") {
			logger.Info("[DEBUG] Symlink path traversal detected! Real relative path: %s", realRelPath)
			c.String(http.StatusForbidden, i18n.T("http.403"))
			return
		}
		logger.Info("[DEBUG] Symlink security check passed. Real relative path: %s", realRelPath)
	}

	// 重新获取文件状态（如果是符号链接，这次会获取目标文件的状态）
	fileInfo, err = os.Stat(fullPath)
	if err != nil {
		// 文件不存在，返回404
		c.String(http.StatusNotFound, i18n.Tf("http.404", reqPath))
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
		c.String(http.StatusForbidden, i18n.T("http.403"))
		return
	}

	// 如果是文件，则提供该文件
	c.File(fullPath)
}
