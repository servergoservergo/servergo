package cmd

// 命令行标志
var (
	// 是否自动打开浏览器（命令行标志）
	autoOpen bool

	// 认证相关标志
	authType        string // 认证类型：none, basic, token, form
	username        string // 用户名
	password        string // 密码
	token           string // 令牌
	enableLoginPage bool   // 是否启用登录页面

	// 目录浏览相关标志
	enableDirListing bool   // 是否启用目录列表功能
	theme            string // 目录列表主题

	// 日志相关标志
	logLevel             string // 日志级别
	enableLogPersistence bool   // 是否启用日志持久化
)

// 别名列表 - 预留位置供后续扩展
var startCmdAliases = []string{
	"run",
	"serve",
	"launch",
	// 这里可以继续添加更多别名
}
