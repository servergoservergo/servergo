package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/CC11001100/servergo/pkg/auth"
)

// setupTestServer 设置测试服务器和测试环境
func setupTestServer(t *testing.T) (*FileServer, string, func()) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建临时目录用于测试
	tempDir, err := ioutil.TempDir("", "servergo-test-")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := ioutil.WriteFile(testFile, []byte("Test File Content"), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建测试子目录
	testSubDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(testSubDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("创建测试子目录失败: %v", err)
	}

	// 创建测试子目录中的文件
	testSubFile := filepath.Join(testSubDir, "subfile.txt")
	if err := ioutil.WriteFile(testSubFile, []byte("Sub File Content"), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("创建测试子目录文件失败: %v", err)
	}

	// 创建index.html文件
	indexFile := filepath.Join(tempDir, "index.html")
	if err := ioutil.WriteFile(indexFile, []byte("<html><body>Index Page</body></html>"), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("创建index.html文件失败: %v", err)
	}

	// 创建带有index.html的测试子目录
	testIndexDir := filepath.Join(tempDir, "indexdir")
	if err := os.Mkdir(testIndexDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("创建测试子目录失败: %v", err)
	}
	testIndexFile := filepath.Join(testIndexDir, "index.html")
	if err := ioutil.WriteFile(testIndexFile, []byte("<html><body>Index in Subdir</body></html>"), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("创建子目录index.html文件失败: %v", err)
	}

	// 创建服务器配置
	config := Config{
		Port:             8080,
		Dir:              tempDir,
		AuthType:         auth.NoAuth,
		EnableDirListing: true,
		Theme:            "default",
	}

	// 创建服务器实例
	srv, err := New(config)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("创建服务器失败: %v", err)
	}

	// 返回清理函数
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return srv, tempDir, cleanup
}

// TestHandleFileRequest 测试handleFileRequest函数
func TestHandleFileRequest(t *testing.T) {
	srv, _, cleanup := setupTestServer(t)
	defer cleanup()

	// 创建Gin测试上下文
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(srv.handleFileRequest)

	// 测试用例
	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedBody string
		setupFunc    func() // 可选的额外设置函数
	}{
		{
			name:         "访问存在的文件",
			path:         "/test.txt",
			expectedCode: http.StatusOK,
			expectedBody: "Test File Content",
		},
		{
			name:         "访问存在的子目录文件",
			path:         "/subdir/subfile.txt",
			expectedCode: http.StatusOK,
			expectedBody: "Sub File Content",
		},
		{
			name:         "访问不存在的文件",
			path:         "/notexist.txt",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "访问根目录时找到index.html",
			path:         "/",
			expectedCode: http.StatusOK,
			expectedBody: "<html><body>Index Page</body></html>",
		},
		{
			name:         "访问带有index.html的子目录",
			path:         "/indexdir",
			expectedCode: http.StatusOK,
			expectedBody: "<html><body>Index in Subdir</body></html>",
		},
		{
			name:         "访问URL编码的路径",
			path:         "/test%20.txt",      // 测试空格编码
			expectedCode: http.StatusNotFound, // 因为我们没有创建名为"test .txt"的文件
		},
		{
			name:         "禁用目录列表时访问子目录",
			path:         "/subdir",
			expectedCode: http.StatusForbidden,
			setupFunc: func() {
				srv.config.EnableDirListing = false
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 运行设置函数（如果有）
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			// 创建测试请求
			req, _ := http.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			// 处理请求
			router.ServeHTTP(w, req)

			// 检查状态码
			if w.Code != tt.expectedCode {
				t.Errorf("状态码 = %d, 期望 %d", w.Code, tt.expectedCode)
			}

			// 如果期望有特定的响应体，则检查
			if tt.expectedBody != "" && tt.expectedCode == http.StatusOK {
				if w.Body.String() != tt.expectedBody {
					t.Errorf("响应体 = %s, 期望 %s", w.Body.String(), tt.expectedBody)
				}
			}
		})
	}
}

// BenchmarkHandleFileRequest 对handleFileRequest进行基准测试
func BenchmarkHandleFileRequest(b *testing.B) {
	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 创建临时目录和文件
	tempDir, err := ioutil.TempDir("", "servergo-bench-")
	if err != nil {
		b.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	content := "This is a test file content for benchmarking."
	if err := ioutil.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建服务器
	config := Config{
		Port:             8080,
		Dir:              tempDir,
		AuthType:         auth.NoAuth,
		EnableDirListing: true,
		Theme:            "default",
	}

	srv, err := New(config)
	if err != nil {
		b.Fatalf("创建服务器失败: %v", err)
	}

	// 创建路由器
	router := gin.New()
	router.NoRoute(srv.handleFileRequest)

	// 重置基准计时器
	b.ResetTimer()

	// 运行基准测试：访问文件
	b.Run("访问文件", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", "/test.txt", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}
