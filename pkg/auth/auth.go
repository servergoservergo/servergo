package auth

import (
	"github.com/gin-gonic/gin"
)

// AuthType 表示认证类型
type AuthType string

const (
	// NoAuth 表示不使用认证
	NoAuth AuthType = "none"
	// BasicAuth 表示使用HTTP Basic认证
	BasicAuth AuthType = "basic"
	// TokenAuth 表示使用令牌认证
	TokenAuth AuthType = "token"
	// FormAuth 表示使用表单认证（登录页面）
	FormAuth AuthType = "form"
)

// Authenticator 接口定义了认证器的方法
type Authenticator interface {
	// Middleware 返回一个Gin中间件，用于处理认证
	Middleware() gin.HandlerFunc
	// AuthType 返回认证类型
	AuthType() AuthType
	// LoginPageEnabled 返回是否启用了登录页
	LoginPageEnabled() bool
}

// Config 保存认证配置
type Config struct {
	// AuthType 认证类型
	Type AuthType
	// Username 用户名，用于BasicAuth和FormAuth
	Username string
	// Password 密码，用于BasicAuth和FormAuth
	Password string
	// Token 令牌，用于TokenAuth
	Token string
	// EnableLoginPage 是否启用登录页面
	EnableLoginPage bool
	// Realm 认证域，用于BasicAuth
	Realm string
}

// NewAuthenticator 根据配置创建一个认证器
func NewAuthenticator(config Config) Authenticator {
	switch config.Type {
	case BasicAuth:
		return NewBasicAuth(config)
	case TokenAuth:
		return NewTokenAuth(config)
	case FormAuth:
		return NewFormAuth(config)
	default:
		return NewNoAuth()
	}
}

// NoAuthenticator 实现了一个不进行认证的认证器
type NoAuthenticator struct{}

// NewNoAuth 创建一个不进行认证的认证器
func NewNoAuth() *NoAuthenticator {
	return &NoAuthenticator{}
}

// Middleware 返回一个不进行认证的中间件
func (a *NoAuthenticator) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// AuthType 返回认证类型
func (a *NoAuthenticator) AuthType() AuthType {
	return NoAuth
}

// LoginPageEnabled 返回是否启用了登录页
func (a *NoAuthenticator) LoginPageEnabled() bool {
	return false
}
