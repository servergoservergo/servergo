package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CC11001100/servergo/pkg/auth"
)

// TestNew 测试服务器实例的创建
func TestNew(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := ioutil.TempDir("", "servergo-test-")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试用例
	tests := []struct {
		name      string
		config    Config
		expectErr bool
	}{
		{
			name: "有效配置",
			config: Config{
				Port:             8080,
				Dir:              tempDir,
				AuthType:         auth.NoAuth,
				EnableDirListing: true,
				Theme:            "default",
			},
			expectErr: false,
		},
		{
			name: "无效目录",
			config: Config{
				Port:             8080,
				Dir:              "/not/exist/dir",
				AuthType:         auth.NoAuth,
				EnableDirListing: true,
			},
			expectErr: true,
		},
		{
			name: "非目录路径",
			config: Config{
				Port:             8080,
				Dir:              "/etc/hosts", // 假设这个文件存在
				AuthType:         auth.NoAuth,
				EnableDirListing: true,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, err := New(tt.config)
			if (err != nil) != tt.expectErr {
				t.Errorf("New() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && srv == nil {
				t.Errorf("New() returned nil server but no error")
			}

			// 如果成功创建了服务器，验证一些基本属性
			if err == nil {
				if srv.absDir == "" {
					t.Errorf("New() created server with empty absDir")
				}
				if srv.engine == nil {
					t.Errorf("New() created server with nil engine")
				}
				if srv.authenticator == nil {
					t.Errorf("New() created server with nil authenticator")
				}
				if srv.dirTemplate == nil {
					t.Errorf("New() created server with nil dirTemplate")
				}

				// 测试GetAbsDir方法
				absDir := srv.GetAbsDir()
				if absDir != srv.absDir {
					t.Errorf("GetAbsDir() = %v, want %v", absDir, srv.absDir)
				}
			}
		})
	}
}

// TestStartServer 测试服务器的启动
// 这是一个集成测试，会实际启动服务器和发送HTTP请求
func TestStartServer(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := ioutil.TempDir("", "servergo-test-")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 在临时目录中创建一个测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ServerGo!"
	err = ioutil.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 使用一个随机端口创建服务器
	port := 8090 // 假设这个端口可用
	config := Config{
		Port:             port,
		Dir:              tempDir,
		AuthType:         auth.NoAuth,
		EnableDirListing: true,
		Theme:            "default",
	}

	srv, err := New(config)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}

	// 在后台启动服务器
	go func() {
		err := srv.Start()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("启动服务器失败: %v", err)
		}
	}()

	// 给服务器一点启动时间
	time.Sleep(100 * time.Millisecond)

	// 测试访问文件
	url := "http://localhost:8090/test.txt"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("请求文件失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("文件请求返回状态码 %d, 期望 %d", resp.StatusCode, http.StatusOK)
	}

	// 检查文件内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("读取响应失败: %v", err)
	}

	if string(body) != testContent {
		t.Errorf("文件内容 = %q, 期望 %q", string(body), testContent)
	}

	// 注意：这个测试不会实际停止服务器，因为没有好的方法在测试中停止它
	// 在实际应用中，我们通常会使用 context.Context 和 server.Shutdown() 进行优雅关闭
}
