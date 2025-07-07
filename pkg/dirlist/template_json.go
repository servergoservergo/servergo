package dirlist

import (
	"encoding/json"
	"fmt"

	"github.com/CC11001100/servergo/pkg/i18n"
)

// renderJSON 直接渲染JSON数据，避免HTML模板的限制
func (t *DirListTemplate) renderJSON(data TemplateData) (string, error) {
	// 创建特殊的JSON结构
	type jsonItem struct {
		Name          string `json:"name"`
		IsDirectory   bool   `json:"is_directory"`
		Size          int64  `json:"size"`           // 改为数值类型的字节大小
		SizeFormatted string `json:"size_formatted"` // 保留原格式化大小
		LastModified  string `json:"last_modified"`
		Path          string `json:"path"`
		URL           string `json:"url"`
	}

	type jsonData struct {
		Path            string     `json:"path"`
		Timestamp       string     `json:"timestamp"`
		ParentDirectory string     `json:"parent_directory"`
		Contents        []jsonItem `json:"contents"`
	}

	// 转换数据
	contents := make([]jsonItem, len(data.Items))
	for i, item := range data.Items {
		url := item.Path
		if item.IsDir {
			url += "/"
		}

		contents[i] = jsonItem{
			Name:          item.Name,
			IsDirectory:   item.IsDir,
			Size:          item.SizeBytes,
			SizeFormatted: item.Size,
			LastModified:  item.LastModified,
			Path:          item.Path,
			URL:           url,
		}
	}

	jsonResult := jsonData{
		Path:            data.DirPath,
		Timestamp:       data.CurrentTime,
		ParentDirectory: data.ParentDir,
		Contents:        contents,
	}

	// 序列化为JSON字符串
	jsonBytes, err := json.MarshalIndent(jsonResult, "", "    ")
	if err != nil {
		return "", fmt.Errorf(i18n.Tf("dirlist.json_marshal_error", err))
	}

	return string(jsonBytes), nil
}
