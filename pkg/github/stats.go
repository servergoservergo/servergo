package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	// GitHub API URL
	apiURL = "https://api.github.com/repos/CC11001100/servergo"
	// 缓存过期时间（1小时）
	cacheExpiration = 1 * time.Hour
)

// RepoStats 存储仓库统计信息
type RepoStats struct {
	Stars      int       `json:"stargazers_count"`
	UpdateTime time.Time `json:"-"`
}

var (
	stats     *RepoStats
	statsMux  sync.RWMutex
	lastFetch time.Time
)

// GetStats 获取仓库统计信息（带缓存）
func GetStats() (*RepoStats, error) {
	statsMux.RLock()
	if stats != nil && time.Since(lastFetch) < cacheExpiration {
		defer statsMux.RUnlock()
		return stats, nil
	}
	statsMux.RUnlock()

	// 需要更新缓存
	return updateStats()
}

// updateStats 从GitHub API更新统计信息
func updateStats() (*RepoStats, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub stats: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var newStats RepoStats
	if err := json.Unmarshal(body, &newStats); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	newStats.UpdateTime = time.Now()

	statsMux.Lock()
	stats = &newStats
	lastFetch = time.Now()
	statsMux.Unlock()

	return &newStats, nil
}

// GetRepoURL 返回仓库URL
func GetRepoURL() string {
	return "https://github.com/CC11001100/servergo"
}
