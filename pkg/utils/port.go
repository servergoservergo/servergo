package utils

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
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
// 简单版本，直接返回布尔值
func IsPortAvailable(port int) bool {
	// 验证端口范围
	if port < 0 || port > 65535 {
		return false
	}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// CheckPort 检查指定端口是否可用，提供更详细的错误信息
// 这是一个增强版函数，返回端口状态和具体错误
func CheckPort(port int, protocol string) (available bool, err error) {
	// 验证端口号范围
	if port < 0 || port > 65535 {
		return false, fmt.Errorf("port %d is out of valid range (0-65535)", port)
	}

	// 如果未指定协议，默认检查TCP
	if protocol == "" {
		protocol = "tcp"
	}

	// 尝试监听指定端口
	ln, err := net.Listen(protocol, ":"+strconv.Itoa(port))
	if err != nil {
		// 区分不同类型的错误
		if strings.Contains(err.Error(), "permission denied") {
			return false, fmt.Errorf("permission denied to bind port %d: %v", port, err)
		}
		if strings.Contains(err.Error(), "address already in use") {
			return false, fmt.Errorf("port %d is already in use", port)
		}
		return false, fmt.Errorf("error checking port %d: %v", port, err)
	}

	// 记得关闭监听器
	defer ln.Close()
	return true, nil
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

	return 0, fmt.Errorf("could not find an available port")
}

// FindAvailablePortWithProtocol 寻找一个可用的端口，支持指定协议
// 提供更详细的错误信息，同时支持TCP或UDP协议
func FindAvailablePortWithProtocol(preferredPort int, protocol string) (int, error) {
	// 设置默认协议
	if protocol == "" {
		protocol = "tcp"
	}

	// 验证协议类型
	if protocol != "tcp" && protocol != "udp" {
		return 0, fmt.Errorf("unsupported protocol: %s, only tcp and udp are supported", protocol)
	}

	// 如果指定了首选端口，先检查它是否可用
	if preferredPort > 0 {
		available, err := CheckPort(preferredPort, protocol)
		if err != nil {
			// 如果是权限错误或者范围错误，直接返回错误
			if strings.Contains(err.Error(), "permission denied") ||
				strings.Contains(err.Error(), "out of valid range") {
				return 0, err
			}
			// 其他错误，如端口已被占用，继续尝试找其他端口
		} else if available {
			return preferredPort, nil
		}
	}

	// 随机选择一个起始端口
	portRange := MaxPort - MinPort + 1
	startPort := MinPort + rand.Intn(portRange)

	// 从随机端口开始，向上循环查找
	for port := startPort; port <= MaxPort; port++ {
		available, err := CheckPort(port, protocol)
		if err == nil && available {
			return port, nil
		}
	}

	// 如果到达最大端口仍未找到，从最小端口到随机端口再次尝试
	for port := MinPort; port < startPort; port++ {
		available, err := CheckPort(port, protocol)
		if err == nil && available {
			return port, nil
		}
	}

	return 0, fmt.Errorf("exhausted all ports, could not find an available %s port", protocol)
}
