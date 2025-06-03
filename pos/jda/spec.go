package jda

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimeNano 是一个自定义类型，用于存储纳秒级时间戳并支持与 RFC3339Nano 字符串的转换
type TimeNano int64

// MarshalJSON 实现自定义的 JSON 序列化：int64 时间戳 -> 字符串
func (t TimeNano) MarshalJSON() ([]byte, error) {
	// 将纳秒时间戳转换为 time.Time
	tm := time.Unix(0, int64(t))
	// 格式化为 RFC3339Nano
	s := fmt.Sprintf(`"%s"`, tm.Format(time.RFC3339Nano))
	return []byte(s), nil
}

// UnmarshalJSON 实现自定义的 JSON 反序列化：字符串 -> int64 时间戳
func (t *TimeNano) UnmarshalJSON(data []byte) error {
	// 去除引号
	s := string(data)
	s = strings.Trim(s, `"`)

	// 处理空字符串
	if s == "" || s == "null" {
		*t = 0
		return nil
	}

	// 尝试将字符串解析为 time.Time
	tm, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}

	// 设置纳秒时间戳
	*t = TimeNano(tm.UnixNano())
	return nil
}

// NumericString 是一个自定义类型，用于存储数值并支持与字符串的转换
// 包括特殊格式如 "10K", "1M" 等
type NumericString float64

// MarshalJSON 实现自定义的 JSON 序列化：float64 -> 字符串
func (n NumericString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%g"`, float64(n))), nil
}

// UnmarshalJSON 实现自定义的 JSON 反序列化：字符串 -> float64
func (n *NumericString) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, `"`)

	// 处理空字符串
	if s == "" || s == "null" {
		*n = 0
		return nil
	}

	// 检查是否有特殊单位 (K, M, etc.)
	re := regexp.MustCompile(`^([0-9]+\.?[0-9]*|\.[0-9]+)([kKmM])$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) == 3 {
		// 有单位的情况
		baseValue, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return err
		}

		// 根据单位调整值
		unit := matches[2]
		switch strings.ToUpper(unit) {
		case "K":
			baseValue *= 1000
		case "M":
			baseValue *= 1000000
		}

		*n = NumericString(baseValue)
		return nil
	}

	// 尝试直接解析为 float64
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*n = NumericString(value)
	return nil
}

// F64 返回底层的 float64 值
func (n NumericString) F64() float64 {
	return float64(n)
}
