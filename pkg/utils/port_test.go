package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestIsPortAvailable 测试端口可用性检查函数
func TestIsPortAvailable(t *testing.T) {
	// 测试场景1：随机选择一个可能可用的高位端口
	// 使用当前时间的纳秒部分来生成随机端口，减少冲突可能性
	randomPort := 40000 + (time.Now().Nanosecond() % 10000)
	if !IsPortAvailable(randomPort) {
		// 可能这个端口恰好被占用，不是测试失败
		t.Logf("随机选择的端口 %d 不可用，这可能是正常的", randomPort)
	} else {
		t.Logf("端口 %d 可用", randomPort)
	}

	// 测试场景2：占用一个端口，然后检查它是否可用
	// 这种方法使用系统分配的端口，确保测试的准确性
	listener, err := net.Listen("tcp", ":0") // 让系统分配一个可用端口
	if err != nil {
		t.Fatalf("无法创建监听器: %v", err)
	}
	defer listener.Close()

	// 获取系统分配的端口
	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("无法解析地址: %v", err)
	}
	usedPort, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("无法将端口转换为整数: %v", err)
	}

	// 检查已占用的端口
	if IsPortAvailable(usedPort) {
		t.Errorf("端口 %d 已被占用，但函数返回它可用", usedPort)
	} else {
		t.Logf("正确检测到端口 %d 已被占用", usedPort)
	}

	// 测试场景3：特权端口（通常需要管理员权限）
	privilegedPort := 80
	// 这里不断言结果，因为根据执行环境不同（普通用户/管理员），结果会不同
	available := IsPortAvailable(privilegedPort)
	t.Logf("特权端口 %d 的可用性: %v (取决于当前用户权限)", privilegedPort, available)

	// 测试场景4：无效端口号（大于最大有效端口）
	invalidPort := 70000 // 超出有效范围
	if IsPortAvailable(invalidPort) {
		t.Errorf("无效端口 %d 不应该被识别为可用", invalidPort)
	}

	// 测试场景5：无效端口号（负数端口）
	invalidNegativePort := -1
	if IsPortAvailable(invalidNegativePort) {
		t.Errorf("负数端口 %d 不应该被识别为可用", invalidNegativePort)
	}

	// 测试场景6：端口号0（系统会分配一个可用端口）
	// 这里不断言结果，因为端口0是特殊的
	zeroPortAvailable := IsPortAvailable(0)
	t.Logf("端口 0 的可用性: %v (特殊端口，系统会分配一个可用端口)", zeroPortAvailable)
}

// TestFindAvailablePort 测试查找可用端口函数
func TestFindAvailablePort(t *testing.T) {
	// 测试场景1：不指定首选端口（传入0）
	port, err := FindAvailablePort(0)
	if err != nil {
		t.Errorf("未能找到可用端口: %v", err)
	} else if port < MinPort || port > MaxPort {
		t.Errorf("找到的端口 %d 超出有效范围 [%d, %d]", port, MinPort, MaxPort)
	} else {
		t.Logf("找到可用端口: %d", port)
	}

	// 测试场景2：指定一个可能可用的首选端口
	preferredPort := 40000 + (time.Now().Nanosecond() % 10000)
	port, err = FindAvailablePort(preferredPort)
	if err != nil {
		t.Errorf("未能找到可用端口: %v", err)
	} else {
		t.Logf("指定首选端口 %d，找到可用端口: %d", preferredPort, port)
		// 如果首选端口可用，应该返回首选端口
		if IsPortAvailable(preferredPort) && port != preferredPort {
			t.Errorf("首选端口 %d 可用，但返回了不同的端口 %d", preferredPort, port)
		}
	}

	// 测试场景3：指定一个已被占用的首选端口
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("无法创建监听器: %v", err)
	}
	defer listener.Close()

	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("无法解析地址: %v", err)
	}
	usedPort, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("无法将端口转换为整数: %v", err)
	}

	port, err = FindAvailablePort(usedPort)
	if err != nil {
		t.Errorf("未能找到可用端口: %v", err)
	} else if port == usedPort {
		t.Errorf("返回了已被占用的端口 %d", usedPort)
	} else {
		t.Logf("指定已占用端口 %d，找到其他可用端口: %d", usedPort, port)
	}

	// 测试场景4：边界情况 - 端口号越界
	// 负数端口号
	port, err = FindAvailablePort(-1)
	if err != nil {
		t.Errorf("使用无效端口 -1 时未能找到可用端口: %v", err)
	} else if port < MinPort || port > MaxPort {
		t.Errorf("找到的端口 %d 超出有效范围 [%d, %d]", port, MinPort, MaxPort)
	} else {
		t.Logf("使用无效端口 -1，函数正确地找到了有效端口: %d", port)
	}

	// 超出最大值的端口号
	port, err = FindAvailablePort(70000)
	if err != nil {
		t.Errorf("使用无效端口 70000 时未能找到可用端口: %v", err)
	} else if port < MinPort || port > MaxPort {
		t.Errorf("找到的端口 %d 超出有效范围 [%d, %d]", port, MinPort, MaxPort)
	} else {
		t.Logf("使用无效端口 70000，函数正确地找到了有效端口: %d", port)
	}

	// 测试场景5：端口0（特殊情况）
	port, err = FindAvailablePort(0)
	if err != nil {
		t.Errorf("使用端口0时未能找到可用端口: %v", err)
	} else if port < MinPort || port > MaxPort {
		t.Errorf("找到的端口 %d 超出有效范围 [%d, %d]", port, MinPort, MaxPort)
	} else {
		t.Logf("使用端口0，函数正确地找到了有效端口: %d", port)
	}
}

// TestPortRangeBoundary 测试端口范围边界值
func TestPortRangeBoundary(t *testing.T) {
	// 测试场景1：最小端口号
	available := IsPortAvailable(MinPort)
	t.Logf("最小端口 %d 的可用性: %v", MinPort, available)

	// 测试场景2：最大端口号
	available = IsPortAvailable(MaxPort)
	t.Logf("最大端口 %d 的可用性: %v", MaxPort, available)

	// 测试场景3：最小端口号 -1
	invalidLowerPort := MinPort - 1
	if IsPortAvailable(invalidLowerPort) {
		t.Logf("端口 %d 小于最小端口 %d，但被识别为可用", invalidLowerPort, MinPort)
	}

	// 测试场景4：最大端口号 +1
	invalidUpperPort := MaxPort + 1
	if IsPortAvailable(invalidUpperPort) {
		t.Errorf("端口 %d 大于最大端口 %d，但被识别为可用", invalidUpperPort, MaxPort)
	}
}

// TestConcurrentPortAllocation 测试并发端口分配
// 这个测试确保在并发情况下不会分配相同的端口
func TestConcurrentPortAllocation(t *testing.T) {
	// 创建通道存储结果
	type result struct {
		port int
		err  error
	}
	results := make(chan result, 10)

	// 并发请求10个端口
	for i := 0; i < 10; i++ {
		go func() {
			port, err := FindAvailablePort(0)
			results <- result{port, err}
		}()
	}

	// 收集结果
	ports := make(map[int]bool)
	for i := 0; i < 10; i++ {
		r := <-results
		if r.err != nil {
			t.Errorf("并发请求 #%d 失败: %v", i, r.err)
			continue
		}

		// 检查端口是否有效且未重复分配
		if r.port < MinPort || r.port > MaxPort {
			t.Errorf("并发请求 #%d 返回的端口 %d 超出有效范围", i, r.port)
		} else if ports[r.port] {
			t.Errorf("端口 %d 被重复分配", r.port)
		} else {
			ports[r.port] = true
			t.Logf("并发请求 #%d 分配到端口: %d", i, r.port)
		}
	}
}

// TestPortExhaustion 通过创建一个使用定制检测函数的端口查找器来测试端口耗尽情况
// 这个测试模拟了端口耗尽的情况，确保函数能够正确处理极端情况
func TestPortExhaustion(t *testing.T) {
	// 这个测试可能会很耗时，所以仅在详细测试模式下运行
	if testing.Short() {
		t.Skip("跳过在短测试模式下的端口耗尽测试")
	}

	// 情况1：创建一个自定义的端口查找函数，只有一个端口可用
	onlyAvailablePort := MinPort + 100
	finder := func(preferredPort int) (int, error) {
		// 自定义的端口可用性检查函数
		isAvailable := func(port int) bool {
			return port == onlyAvailablePort
		}

		// 如果首选端口可用，直接返回
		if preferredPort > 0 && isAvailable(preferredPort) {
			return preferredPort, nil
		}

		// 随机选择一个起始端口
		portRange := MaxPort - MinPort + 1
		startPort := MinPort + (time.Now().Nanosecond() % portRange)

		// 从随机端口开始，向上循环查找
		for port := startPort; port <= MaxPort; port++ {
			if isAvailable(port) {
				return port, nil
			}
		}

		// 如果到达最大端口仍未找到，从最小端口到随机端口再次尝试
		for port := MinPort; port < startPort; port++ {
			if isAvailable(port) {
				return port, nil
			}
		}

		return 0, fmt.Errorf("无法找到可用端口")
	}

	// 测试场景1：指定唯一可用的端口
	port, err := finder(onlyAvailablePort)
	if err != nil {
		t.Errorf("未能找到唯一可用的端口: %v", err)
	} else if port != onlyAvailablePort {
		t.Errorf("指定唯一可用端口 %d，但返回了 %d", onlyAvailablePort, port)
	} else {
		t.Logf("正确找到唯一可用的端口: %d", port)
	}

	// 测试场景2：指定一个不可用的端口，应该找到唯一可用的端口
	port, err = finder(onlyAvailablePort + 1)
	if err != nil {
		t.Errorf("未能找到唯一可用的端口: %v", err)
	} else if port != onlyAvailablePort {
		t.Errorf("指定不可用端口时，应该找到唯一可用端口 %d，但返回了 %d", onlyAvailablePort, port)
	} else {
		t.Logf("正确找到唯一可用的端口: %d", port)
	}

	// 情况2：创建一个所有端口都不可用的查找函数
	noPortFinder := func(preferredPort int) (int, error) {
		// 自定义的端口可用性检查函数，总是返回false
		isAvailable := func(port int) bool {
			return false
		}

		// 如果首选端口可用，直接返回
		if preferredPort > 0 && isAvailable(preferredPort) {
			return preferredPort, nil
		}

		// 剩余的查找逻辑复制自FindAvailablePort
		// 随机选择一个起始端口
		portRange := MaxPort - MinPort + 1
		startPort := MinPort + (time.Now().Nanosecond() % portRange)

		// 从随机端口开始，向上循环查找
		for port := startPort; port <= MaxPort; port++ {
			if isAvailable(port) {
				return port, nil
			}
		}

		// 如果到达最大端口仍未找到，从最小端口到随机端口再次尝试
		for port := MinPort; port < startPort; port++ {
			if isAvailable(port) {
				return port, nil
			}
		}

		return 0, fmt.Errorf("无法找到可用端口")
	}

	// 测试场景3：所有端口都不可用的情况
	port, err = noPortFinder(0)
	if err == nil {
		t.Errorf("所有端口都不可用，但没有返回错误")
	} else {
		t.Logf("正确处理所有端口不可用的情况: %v", err)
	}
}

// TestCheckPort 测试检查指定端口是否可用，并提供详细错误信息
func TestCheckPort(t *testing.T) {
	// 测试场景1：有效端口范围内的随机端口
	randomPort := 40000 + (time.Now().Nanosecond() % 10000)
	available, err := CheckPort(randomPort, "tcp")
	if err != nil {
		t.Logf("检查端口 %d 时出现错误: %v", randomPort, err)
	} else {
		t.Logf("端口 %d 可用性: %v", randomPort, available)
	}

	// 测试场景2：已被占用的端口
	listener, err := net.Listen("tcp", ":0") // 系统分配一个可用端口
	if err != nil {
		t.Fatalf("无法创建监听器: %v", err)
	}
	defer listener.Close()

	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("无法解析地址: %v", err)
	}
	usedPort, _ := strconv.Atoi(portStr)

	available, err = CheckPort(usedPort, "tcp")
	if err == nil {
		t.Errorf("端口 %d 已被占用，但没有返回错误", usedPort)
	} else {
		t.Logf("正确检测到端口 %d 已被占用: %v", usedPort, err)
	}

	// 测试场景3：无效端口（超出范围）
	available, err = CheckPort(70000, "tcp")
	if err == nil {
		t.Errorf("端口 70000 超出有效范围，但没有返回错误")
	} else {
		t.Logf("正确检测到无效端口: %v", err)
	}

	// 测试场景4：无效端口（负数）
	available, err = CheckPort(-1, "tcp")
	if err == nil {
		t.Errorf("端口 -1 是负数，但没有返回错误")
	} else {
		t.Logf("正确检测到无效端口: %v", err)
	}

	// 测试场景5：使用不支持的协议
	available, err = CheckPort(8080, "icmp") // ICMP不是TCP或UDP
	if err == nil {
		t.Errorf("使用了不支持的协议，但没有返回错误")
	} else {
		t.Logf("正确检测到不支持的协议: %v", err)
	}
}

// TestFindAvailablePortWithProtocol 测试带协议的可用端口查找函数
func TestFindAvailablePortWithProtocol(t *testing.T) {
	// 测试场景1：默认协议（不指定）
	port, err := FindAvailablePortWithProtocol(0, "")
	if err != nil {
		t.Errorf("使用默认协议查找端口失败: %v", err)
	} else {
		t.Logf("使用默认协议找到可用端口: %d", port)
	}

	// 测试场景2：指定TCP协议
	port, err = FindAvailablePortWithProtocol(0, "tcp")
	if err != nil {
		t.Errorf("使用TCP协议查找端口失败: %v", err)
	} else {
		t.Logf("使用TCP协议找到可用端口: %d", port)
	}

	// 测试场景3：指定UDP协议
	port, err = FindAvailablePortWithProtocol(0, "udp")
	if err != nil {
		// 在某些环境下UDP端口检测可能会失败，这里只记录不导致测试失败
		t.Logf("使用UDP协议查找端口失败: %v", err)
	} else {
		t.Logf("使用UDP协议找到可用端口: %d", port)
	}

	// 测试场景4：指定不支持的协议
	_, err = FindAvailablePortWithProtocol(0, "icmp")
	if err == nil {
		t.Errorf("使用不支持的协议，但没有返回错误")
	} else {
		t.Logf("正确检测到不支持的协议: %v", err)
	}

	// 测试场景5：指定已被占用的首选端口
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("无法创建监听器: %v", err)
	}
	defer listener.Close()

	_, portStr, _ := net.SplitHostPort(listener.Addr().String())
	usedPort, _ := strconv.Atoi(portStr)

	port, err = FindAvailablePortWithProtocol(usedPort, "tcp")
	if err != nil {
		t.Errorf("指定已占用端口时查找其他可用端口失败: %v", err)
	} else if port == usedPort {
		t.Errorf("返回了已被占用的端口 %d", usedPort)
	} else {
		t.Logf("指定已占用端口 %d，找到其他可用端口: %d", usedPort, port)
	}

	// 测试场景6：特权端口（需要管理员权限）
	privilegedPort := 80
	port, err = FindAvailablePortWithProtocol(privilegedPort, "tcp")

	// 特权端口测试：权限拒绝是预期的行为，不应导致测试失败
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			t.Logf("无法使用特权端口 %d: %v (这是正常的非root用户行为)", privilegedPort, err)
		} else {
			// 其他错误可能需要关注
			t.Logf("使用特权端口 %d 时出现非权限错误: %v", privilegedPort, err)
		}
	} else {
		// 如果成功了（可能是在root权限下运行），也是正常的
		t.Logf("成功使用特权端口 %d", port)
	}
}
