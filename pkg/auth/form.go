package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// FormAuthenticator 实现了基于表单的认证
type FormAuthenticator struct {
	username        string
	password        string
	enableLoginPage bool
	loginTemplate   string
}

// NewFormAuth 创建一个FormAuth认证器
func NewFormAuth(config Config) *FormAuthenticator {
	return &FormAuthenticator{
		username:        config.Username,
		password:        config.Password,
		enableLoginPage: config.EnableLoginPage,
		loginTemplate:   loginPageHTML, // 默认的登录页面HTML
	}
}

// Middleware 返回表单认证中间件
func (a *FormAuthenticator) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 静态登录页面的路由
		if c.Request.URL.Path == "/login" {
			if c.Request.Method == "GET" {
				c.Header("Content-Type", "text/html")
				c.String(http.StatusOK, a.loginTemplate)
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
					// 验证失败，返回登录页面并显示错误
					c.Header("Content-Type", "text/html")
					errorHTML := a.loginTemplate
					// 简单替换一个错误提示
					errorHTML = strings.Replace(errorHTML, "<!--ERROR_MESSAGE-->",
						"<div class='error'>用户名或密码不正确</div>", 1)
					c.String(http.StatusUnauthorized, errorHTML)
					c.Abort()
					return
				}
			}
		}

		// 检查Cookie认证状态
		auth, _ := c.Cookie("servergo_auth")
		if auth != "true" {
			// 未认证，重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 认证通过，继续请求
		c.Next()
	}
}

// SetupRoutes 设置表单认证的路由
func (a *FormAuthenticator) SetupRoutes(router *gin.Engine) {
	router.GET("/login", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, a.loginTemplate)
	})

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == a.username && password == a.password {
			// 设置cookie表示已登录
			c.SetCookie("servergo_auth", "true", 3600, "/", "", false, true)

			// 重定向到根目录
			c.Redirect(http.StatusFound, "/")
		} else {
			// 验证失败，返回登录页面并显示错误
			c.Header("Content-Type", "text/html")
			errorHTML := a.loginTemplate
			// 简单替换一个错误提示
			errorHTML = strings.Replace(errorHTML, "<!--ERROR_MESSAGE-->",
				"<div class='error'>用户名或密码不正确</div>", 1)
			c.String(http.StatusUnauthorized, errorHTML)
		}
	})

	router.GET("/logout", func(c *gin.Context) {
		c.SetCookie("servergo_auth", "", -1, "/", "", false, true)
		c.Redirect(http.StatusFound, "/login")
	})
}

// AuthType 返回认证类型
func (a *FormAuthenticator) AuthType() AuthType {
	return FormAuth
}

// LoginPageEnabled 返回是否启用了登录页
func (a *FormAuthenticator) LoginPageEnabled() bool {
	return a.enableLoginPage
}

// 默认的登录页面HTML
var loginPageHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ServerGo - 登录</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }
        .login-container {
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            max-width: 400px;
            width: 100%;
        }
        h1 {
            color: #333;
            margin-bottom: 20px;
            text-align: center;
        }
        form {
            display: flex;
            flex-direction: column;
        }
        label {
            margin-bottom: 8px;
            color: #555;
        }
        input {
            padding: 10px;
            margin-bottom: 15px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
        }
        button {
            padding: 12px;
            background-color: #4caf50;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #45a049;
        }
        .error {
            color: #ff3860;
            margin-bottom: 15px;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <h1>ServerGo 文件服务器</h1>
        <!--ERROR_MESSAGE-->
        <form method="post" action="/login">
            <label for="username">用户名</label>
            <input type="text" id="username" name="username" required>
            
            <label for="password">密码</label>
            <input type="password" id="password" name="password" required>
            
            <button type="submit">登录</button>
        </form>
    </div>
</body>
</html>
`
