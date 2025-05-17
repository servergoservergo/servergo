package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/gin-gonic/gin"
)

// FormAuthenticator 实现了基于表单的认证
type FormAuthenticator struct {
	username        string
	password        string
	enableLoginPage bool
}

// NewFormAuth 创建一个FormAuth认证器
func NewFormAuth(config Config) *FormAuthenticator {
	return &FormAuthenticator{
		username:        config.Username,
		password:        config.Password,
		enableLoginPage: config.EnableLoginPage,
	}
}

// Middleware 返回表单认证中间件
func (a *FormAuthenticator) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 静态资源路径不需要认证
		if strings.HasPrefix(c.Request.URL.Path, "/auth/") &&
			(strings.HasSuffix(c.Request.URL.Path, ".css") ||
				strings.HasSuffix(c.Request.URL.Path, ".js")) {
			c.Next()
			return
		}

		// 处理登录页面请求
		if c.Request.URL.Path == "/auth/login" {
			if c.Request.Method == "GET" {
				// 获取登录页面HTML
				loginHTML, err := GetLoginHTMLContent()
				if err != nil {
					c.String(http.StatusInternalServerError, i18n.T("login.error.server"))
					c.Abort()
					return
				}

				c.Header("Content-Type", "text/html; charset=utf-8")
				c.String(http.StatusOK, loginHTML)
				c.Abort()
				return
			} else if c.Request.Method == "POST" {
				// 处理登录表单提交
				username := c.PostForm("username")
				password := c.PostForm("password")

				if username == a.username && password == a.password {
					// 设置cookie表示已登录
					c.SetCookie("servergo_auth", "true", 3600, "/", "", false, true)

					// 重定向到根目录
					c.Redirect(http.StatusFound, "/")
					c.Abort()
					return
				} else {
					// 验证失败，重定向到登录页面并显示错误
					errorMsg := i18n.T("login.error.credentials")
					c.Redirect(http.StatusFound, fmt.Sprintf("/auth/login?error=%s", errorMsg))
					c.Abort()
					return
				}
			}
		}

		// 处理登出请求
		if c.Request.URL.Path == "/auth/logout" {
			c.SetCookie("servergo_auth", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		// 检查Cookie认证状态
		auth, _ := c.Cookie("servergo_auth")
		if auth != "true" {
			// 未认证，重定向到登录页面
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		// 认证通过，继续请求
		c.Next()
	}
}

// SetupRoutes 设置表单认证的路由
func (a *FormAuthenticator) SetupRoutes(router *gin.Engine) {
	// 提供静态资源
	router.StaticFS("/auth", GetAuthFileSystem())

	// 登录处理在中间件中已经实现，这里不需要额外的路由
}

// AuthType 返回认证类型
func (a *FormAuthenticator) AuthType() AuthType {
	return FormAuth
}

// LoginPageEnabled 返回是否启用了登录页
func (a *FormAuthenticator) LoginPageEnabled() bool {
	return a.enableLoginPage
}
