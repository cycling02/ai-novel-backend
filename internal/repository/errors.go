package repository

import (
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("资源不存在")
	ErrDuplicate    = errors.New("数据重复")
	ErrInvalidInput = errors.New("无效输入")
	ErrUnauthorized = errors.New("未授权")
	ErrInternal     = errors.New("内部错误")
)

// 辅助函数：类型转换
func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	default:
		return ""
	}
}

func getInt(v interface{}) int {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}

func getBool(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func getTime(v interface{}) time.Time {
	if v == nil {
		return time.Time{}
	}
	if t, ok := v.(time.Time); ok {
		return t
	}
	return time.Time{}
}
