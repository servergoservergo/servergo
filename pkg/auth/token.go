package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TokenAuthenticator 实现了基于令牌的认证
type TokenAuthenticator struct {
	token string
}

// NewTokenAuth 创建一个TokenAuth认证器
func NewTokenAuth(config Config) *TokenAuthenticator {
	return &TokenAuthenticator{
		token: config.Token,
	}
}

// Middleware 返回一个检查令牌的中间件
func (a *TokenAuthenticator) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从查询参数中检查token
		token := c.Query("token")
		if token == "" {
			// 如果查询参数中没有token，从Header中检查
			token = c.GetHeader("Authorization")
		}

		// 验证token
		if token != a.token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未授权访问，请提供有效的token",
			})
			return
		}

		c.Next()
	}
}

// AuthType 返回认证类型
func (a *TokenAuthenticator) AuthType() AuthType {
	return TokenAuth
}

// LoginPageEnabled 返回是否启用了登录页
func (a *TokenAuthenticator) LoginPageEnabled() bool {
	return false
}
