package utils

import (
	"fmt"
	"net"
	"strconv"
)

// 端口范围常量
const (
	MinPort = 1024  // 避免使用特权端口
	MaxPort = 65535 // 最大端口号
)

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
// 如果首选端口不可用或未指定，会从MinPort开始依次尝试
func FindAvailablePort(preferredPort int) (int, error) {
	// 如果指定了首选端口并且可用，直接返回
	if preferredPort > 0 && IsPortAvailable(preferredPort) {
		return preferredPort, nil
	}

	// 否则从MinPort开始寻找可用端口
	for port := MinPort; port <= MaxPort; port++ {
		if IsPortAvailable(port) {
			return port, nil
		}
	}

	return 0, fmt.Errorf("无法找到可用端口")
}
