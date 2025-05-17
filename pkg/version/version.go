// Package version 提供应用程序版本信息
// 版本信息通过编译时的 -ldflags 参数注入
package version

// 版本信息 - 这些变量的值将在编译时通过 -ldflags 参数设置
var (
	// Version 表示程序版本号，格式为 vX.Y.Z
	Version = "dev"
	// BuildTime 构建时间，格式为 YYYY-MM-DD HH:MM:SS
	BuildTime = "unknown"
	// GitCommit Git提交哈希的短版本（前8位）
	GitCommit = "unknown"
	// GitRef Git分支或标签
	GitRef = "unknown"
)

// GetVersion 返回完整的版本信息字符串
func GetVersion() string {
	return Version
}

// GetVersionInfo 返回包含所有版本信息的map
func GetVersionInfo() map[string]string {
	return map[string]string{
		"version":   Version,
		"buildTime": BuildTime,
		"gitCommit": GitCommit,
		"gitRef":    GitRef,
	}
}

// GetBuildInfo 返回构建信息字符串
func GetBuildInfo() string {
	return BuildTime + " (" + GitCommit + ")"
}
