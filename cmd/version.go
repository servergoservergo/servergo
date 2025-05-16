package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 版本信息
var (
	Version   = "0.1.0"
	BuildDate = "未知"
	GitCommit = "未知"
)

// versionCmd 显示当前版本信息和项目信息
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示ServerGo的版本信息",
	Long:  `显示ServerGo的详细版本信息、构建信息、作者信息以及问题反馈渠道等。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ServerGo 版本信息")
		fmt.Println("==============")
		fmt.Printf("版本: %s\n", Version)
		fmt.Printf("构建日期: %s\n", BuildDate)
		fmt.Printf("Git提交哈希: %s\n", GitCommit)
		fmt.Println()
		fmt.Println("ServerGo 是一个快速的HTTP文件服务器工具")
		fmt.Println("类似于Python的http.server模块，但提供更好的性能和更多功能")
		fmt.Println()
		fmt.Println("作者信息")
		fmt.Println("--------")
		fmt.Println("作者: CC11001100")
		fmt.Println("GitHub: https://github.com/CC11001100")
		fmt.Println()
		fmt.Println("反馈问题")
		fmt.Println("--------")
		fmt.Println("如果您发现任何问题或有功能请求，请访问:")
		fmt.Println("https://github.com/CC11001100/servergo/issues")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
