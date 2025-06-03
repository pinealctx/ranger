package b85

import (
	"errors"
	"fmt"
	"strings"
)

// Base85LongConverterS is a converter for converting int64 values to and from Base85 strings.
type Base85LongConverterS struct {
	chars string // Base85 字符集
}

// NewBase85LongConverterS new a Base85LongConverterS.
func NewBase85LongConverterS() *Base85LongConverterS {
	// Base85 字符集
	chars := "0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz\"#$%&'()*+,-./ "
	return &Base85LongConverterS{chars: chars}
}

// AsString converts an int64 value to a Base85 string.
func (c *Base85LongConverterS) AsString(value int64) string {
	if value == 0 {
		return string(c.chars[0])
	}

	var result strings.Builder
	base := int64(len(c.chars))

	for value > 0 {
		remainder := value % base
		result.WriteByte(c.chars[remainder])
		value /= base
	}

	// 反转字符串
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// Parse converts a Base85 string to an int64 value.
func (c *Base85LongConverterS) Parse(s string) (int64, error) {
	if len(s) > MaxLength {
		return 0, errors.New("input string exceeds maximum length")
	}

	var result int64
	base := int64(len(c.chars))

	for _, char := range s {
		index := strings.IndexRune(c.chars, char)
		if index == -1 {
			return 0, fmt.Errorf("invalid character in input: %c", char)
		}
		result = result*base + int64(index)
	}

	return result, nil
}
