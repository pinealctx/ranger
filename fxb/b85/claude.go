package b85

import (
	"errors"
	"strings"
)

type Base85LongConverter struct {
	chars     string
	charToNum map[rune]int
}

const (
	MaxLength   = 10
	base85Chars = "0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz\"#$%&'()*+,-./ "
)

// NewBase85LongConverter creates a new Base85LongConverter instance
func NewBase85LongConverter() *Base85LongConverter {
	converter := &Base85LongConverter{
		chars:     base85Chars,
		charToNum: make(map[rune]int),
	}

	// Initialize the reverse lookup map
	for i, ch := range base85Chars {
		converter.charToNum[ch] = i
	}

	return converter
}

// Parse converts a base85 string to int64
func (c *Base85LongConverter) Parse(text string) (int64, error) {
	if len(text) > MaxLength {
		return 0, errors.New("text exceeds maximum length of 10")
	}

	if len(text) == 0 {
		return 0, errors.New("empty string")
	}

	var result int64 = 0
	base := int64(len(c.chars))

	for _, ch := range text {
		val, exists := c.charToNum[ch]
		if !exists {
			return 0, errors.New("invalid character in input")
		}
		result = result*base + int64(val)
	}

	return result, nil
}

// AsString converts int64 to base85 string
func (c *Base85LongConverter) AsString(value int64) string {
	if value == 0 {
		return "0"
	}

	var builder strings.Builder
	base := int64(len(c.chars))
	temp := value

	// Convert to base85 in reverse order
	digits := make([]byte, 0, MaxLength)
	for temp != 0 {
		idx := temp % base
		digits = append(digits, c.chars[idx])
		temp /= base
	}

	// Reverse the digits
	for i := len(digits) - 1; i >= 0; i-- {
		builder.WriteByte(digits[i])
	}

	return builder.String()
}
