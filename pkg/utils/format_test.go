package utils

import (
	"testing"
)

// TestFormatSize 测试FormatSize函数在不同输入大小下的输出
func TestFormatSize(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{"零字节", 0, "0 B"},
		{"小于1KB", 800, "800 B"},
		{"等于1KB", 1024, "1.0 KB"},
		{"1KB到1MB之间", 234567, "229.1 KB"},
		{"等于1MB", 1048576, "1.0 MB"},
		{"1MB到1GB之间", 23456789, "22.4 MB"},
		{"等于1GB", 1073741824, "1.0 GB"},
		{"1GB到1TB之间", 45678901234, "42.5 GB"},
		{"等于1TB", 1099511627776, "1.0 TB"},
		{"大于1TB", 12345678901234, "11.2 TB"},
		{"负值", -1024, "-1024 B"}, // 虽然实际应用中不太可能有负值，但也要测试
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSize(tt.size)
			if result != tt.expected {
				t.Errorf("FormatSize(%d) = %s, 期望值 %s", tt.size, result, tt.expected)
			}
		})
	}
}

// BenchmarkFormatSize 对FormatSize函数进行基准测试
func BenchmarkFormatSize(b *testing.B) {
	// 准备不同范围的大小值进行基准测试
	sizes := []int64{0, 100, 1024, 1048576, 1073741824, 1099511627776}

	for _, size := range sizes {
		b.Run("Size-"+FormatSize(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FormatSize(size)
			}
		})
	}
}
