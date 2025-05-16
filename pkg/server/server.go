package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Config 保存文件服务器的配置
type Config struct {
	Port int    // 服务器监听的端口
	Dir  string // 提供服务的目录路径
}

// FileServer 表示一个文件服务器实例
type FileServer struct {
	config Config
	absDir string
}

// New 创建一个新的文件服务器实例
func New(config Config) (*FileServer, error) {
	// 获取绝对路径
	absDir, err := filepath.Abs(config.Dir)
	if err != nil {
		return nil, fmt.Errorf("无法获取绝对路径: %v", err)
	}

	// 检查目录是否存在
	info, err := os.Stat(absDir)
	if err != nil {
		return nil, fmt.Errorf("无法访问目录 %s: %v", absDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s 不是一个目录", absDir)
	}

	return &FileServer{
		config: config,
		absDir: absDir,
	}, nil
}

// Start 启动文件服务器
func (fs *FileServer) Start() error {
	// 创建文件服务器处理器
	fileServer := http.FileServer(http.Dir(fs.absDir))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		fileServer.ServeHTTP(w, r)
	}))

	// 启动服务器
	hostAddr := ":" + strconv.Itoa(fs.config.Port)
	fmt.Printf("启动文件服务器在 http://localhost:%d\n", fs.config.Port)
	fmt.Printf("提供目录: %s\n", fs.absDir)
	fmt.Println("按 Ctrl+C 停止服务器")

	return http.ListenAndServe(hostAddr, nil)
}

// GetAbsDir 获取文件服务器的绝对路径
func (fs *FileServer) GetAbsDir() string {
	return fs.absDir
}
