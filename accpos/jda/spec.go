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

// 定义UTC时间格式，没有时区信息
const timeFormat = "2006-01-02T15:04:05.999999999"

// MarshalJSON 实现自定义的 JSON 序列化：int64 时间戳 -> 字符串
func (t TimeNano) MarshalJSON() ([]byte, error) {
	// 将纳秒时间戳转换为 time.Time
	tm := time.Unix(0, int64(t)).UTC()
	// 格式化为没有时区的时间格式
	s := fmt.Sprintf(`"%s"`, tm.Format(timeFormat))
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

	// 解析固定格式的时间字符串
	tm, err := time.Parse(timeFormat, s)
	if err != nil {
		// parse as long
		var i int64
		i, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid time format: %s", s)
		}
		// 如果解析为整数成功，则将其视为纳秒时间戳
		*t = TimeNano(i)
		return nil
	}

	// 设置纳秒时间戳（UTC 时间）
	*t = TimeNano(tm.UnixNano())
	return nil
}

// String 实现 Stringer 接口，返回原始格式的时间字符串
func (t TimeNano) String() string {
	// 特殊处理零值
	if t == 0 {
		return ""
	}

	// 将纳秒时间戳转换为 time.Time
	tm := time.Unix(0, int64(t)).UTC()

	// 格式化为原始格式
	return tm.Format(timeFormat)
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

// StringInt64 是一个自定义类型，用于在 JSON 中表示为字符串的 int64 值
type StringInt64 int64

// MarshalJSON 实现自定义的 JSON 序列化：int64 -> 字符串
func (si StringInt64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, int64(si))), nil
}

// UnmarshalJSON 实现自定义的 JSON 反序列化：字符串 -> int64
func (si *StringInt64) UnmarshalJSON(data []byte) error {
	// 去除引号
	s := string(data)
	s = strings.Trim(s, `"`)

	// 处理空字符串或 null
	if s == "" || s == "null" {
		*si = 0
		return nil
	}

	// 解析为 int64
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*si = StringInt64(i)
	return nil
}
