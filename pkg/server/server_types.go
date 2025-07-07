package server

import (
	"github.com/gin-gonic/gin"

	"github.com/CC11001100/servergo/pkg/auth"
	"github.com/CC11001100/servergo/pkg/dirlist"
)

// Config 保存文件服务器的配置
// 该结构体包含了服务器启动和运行所需的所有配置参数
type Config struct {
	Port int    // 服务器监听的端口，例如: 8080
	Dir  string // 提供服务的目录路径，例如: "/home/user/files"

	// 认证相关配置
	AuthType        auth.AuthType // 认证类型，可选值: auth.NoAuth, auth.BasicAuth, auth.TokenAuth, auth.FormAuth
	Username        string        // 用户名，用于BasicAuth和FormAuth，例如: "admin"
	Password        string        // 密码，用于BasicAuth和FormAuth，例如: "password123"
	Token           string        // 令牌，用于TokenAuth，例如: "abcdef123456"
	EnableLoginPage bool          // 是否启用登录页面，用于FormAuth，例如: true表示启用

	// 目录浏览相关配置
	EnableDirListing bool   // 是否启用目录列表功能，例如: true表示启用
	Theme            string // 目录列表主题，可选值: "default", "bootstrap", "material" 等
}

// FileServer 表示一个文件服务器实例
// 该结构体封装了服务器的所有状态和功能
type FileServer struct {
	config        Config                   // 服务器配置信息
	absDir        string                   // 服务目录的绝对路径，例如: "/home/user/files"
	engine        *gin.Engine              // Gin引擎实例，用于处理HTTP请求
	authenticator auth.Authenticator       // 认证器实例，用于处理用户认证
	dirTemplate   *dirlist.DirListTemplate // 目录列表模板，用于渲染目录页面
}

// GetAbsDir 获取文件服务器的绝对路径
//
// 返回值:
//   - string: 服务器提供服务的目录绝对路径
//
// 使用示例:
// ```
// absPath := srv.GetAbsDir()
// fmt.Printf("服务器的绝对路径: %s\n", absPath)
// ```
func (fs *FileServer) GetAbsDir() string {
	return fs.absDir
}
