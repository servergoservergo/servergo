package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	// 密码字符集
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars  = "0123456789"
	specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// GenerateRandomPassword 生成一个随机密码
// length: 密码长度
// useSpecial: 是否使用特殊字符
func GenerateRandomPassword(length int, useSpecial bool) string {
	if length < 16 {
		length = 16 // 确保最小长度为16
	}

	// 构建字符集
	charset := lowerChars + upperChars + numberChars
	if useSpecial {
		charset += specialChars
	}

	// 确保至少包含一个小写字母、一个大写字母和一个数字
	password := make([]byte, length)

	// 生成一个小写字母
	password[0] = lowerChars[randomInt(len(lowerChars))]
	// 生成一个大写字母
	password[1] = upperChars[randomInt(len(upperChars))]
	// 生成一个数字
	password[2] = numberChars[randomInt(len(numberChars))]

	// 生成剩余的字符
	for i := 3; i < length; i++ {
		password[i] = charset[randomInt(len(charset))]
	}

	// 打乱密码顺序
	for i := len(password) - 1; i > 0; i-- {
		j := randomInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}

// randomInt 生成一个0到max-1之间的随机数
func randomInt(max int) int {
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(result.Int64())
}
