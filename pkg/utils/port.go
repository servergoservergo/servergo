package utils

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// 端口范围常量
const (
	MinPort = 1024  // 避免使用特权端口
	MaxPort = 65535 // 最大端口号
)

// 确保随机数生成器被初始化
func init() {
	// 使用当前时间作为随机种子
	rand.Seed(time.Now().UnixNano())
}

// IsPortAvailable 检查指定端口是否可用
func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// FindAvailablePort 寻找一个可用的端口
// 如果指定了首选端口(preferredPort > 0)，会先检查它是否可用
// 如果首选端口不可用或未指定，会从随机端口开始探测
func FindAvailablePort(preferredPort int) (int, error) {
	// 如果指定了首选端口并且可用，直接返回
	if preferredPort > 0 && IsPortAvailable(preferredPort) {
		return preferredPort, nil
	}

	// 随机选择一个起始端口
	portRange := MaxPort - MinPort + 1
	startPort := MinPort + rand.Intn(portRange)

	// 从随机端口开始，向上循环查找
	for port := startPort; port <= MaxPort; port++ {
		if IsPortAvailable(port) {
			return port, nil
		}
	}

	// 如果到达最大端口仍未找到，从最小端口到随机端口再次尝试
	for port := MinPort; port < startPort; port++ {
		if IsPortAvailable(port) {
			return port, nil
		}
	}

	return 0, fmt.Errorf("无法找到可用端口")
}
