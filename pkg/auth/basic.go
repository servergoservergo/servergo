package auth

import (
	"github.com/gin-gonic/gin"
)

// BasicAuthenticator 实现了HTTP Basic认证
type BasicAuthenticator struct {
	username string
	password string
	realm    string
}

// NewBasicAuth 创建一个BasicAuth认证器
func NewBasicAuth(config Config) *BasicAuthenticator {
	realm := "ServerGo Protected Area"
	if config.Realm != "" {
		realm = config.Realm
	}

	return &BasicAuthenticator{
		username: config.Username,
		password: config.Password,
		realm:    realm,
	}
}

// Middleware 返回Gin的Basic认证中间件
func (a *BasicAuthenticator) Middleware() gin.HandlerFunc {
	accounts := gin.Accounts{
		a.username: a.password,
	}
	return gin.BasicAuth(accounts)
}

// AuthType 返回认证类型
func (a *BasicAuthenticator) AuthType() AuthType {
	return BasicAuth
}

// LoginPageEnabled 返回是否启用了登录页
func (a *BasicAuthenticator) LoginPageEnabled() bool {
	return false // Basic Auth使用浏览器自带的认证对话框
}
