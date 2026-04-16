package formatter

import (
	"context"
	"encoding/json"
	"time"
)

// JSONFormatter JSON格式化器
type JSONFormatter struct {
	// 是否美化输出
	Pretty bool
	// 自定义字段映射
	FieldMap map[string]string
}

// NewJSONFormatter 创建JSON格式化器
func NewJSONFormatter(pretty bool) *JSONFormatter {
	return &JSONFormatter{
		Pretty: pretty,
	}
}

// Format 实现JSONFormatter的格式化方法
func (f *JSONFormatter) Format(ctx context.Context, details *Details) string {
	// 创建JSON结构
	jsonData := make(map[string]interface{})

	// 映射字段
	if f.FieldMap != nil {
		// 使用自定义字段映射
		for key, field := range f.FieldMap {
			switch field {
			case "level":
				jsonData[key] = details.Level.String()
			case "time":
				jsonData[key] = details.Time.Format(time.RFC3339)
			case "path":
				jsonData[key] = details.Path
			case "message":
				jsonData[key] = details.Message
			case "value":
				jsonData[key] = details.Value
			}
		}
	} else {
		// 使用默认字段名
		jsonData["level"] = details.Level.String()
		jsonData["time"] = details.Time.Format(time.RFC3339)
		jsonData["path"] = details.Path
		jsonData["message"] = details.Message
		// 如果Value非空，也包含value字段
		if len(details.Value) > 0 {
			jsonData["value"] = details.Value
		}
	}

	// 序列化为JSON
	var output []byte
	var err error
	if f.Pretty {
		output, err = json.MarshalIndent(jsonData, "", "  ")
	} else {
		output, err = json.Marshal(jsonData)
	}

	if err != nil {
		// JSON序列化失败，返回空字符串
		return ""
	}

	return string(output) + "\n"
}

// JSONFormatterFactory 创建JSON格式化器
func JSONFormatterFactory(pretty bool) Formatter {
	return NewJSONFormatter(pretty)
}
