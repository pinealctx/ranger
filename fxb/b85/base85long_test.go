package b85

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestBase85LongConverter_AsString(t *testing.T) {
	c := NewBase85LongConverter()
	if c.AsString(0) != "0" {
		t.Errorf("AsString(0) = %s; want 0", c.AsString(0))
	}
	if c.AsString(1) != "1" {
		t.Errorf("AsString(1) = %s; want 1", c.AsString(1))
	}
	if c.AsString(85) != "10" {
		t.Errorf("AsString(85) = %s; want 10", c.AsString(85))
	}
	if c.AsString(86) != "11" {
		t.Errorf("AsString(86) = %s; want 11", c.AsString(86))
	}
	if c.AsString(85*85) != "100" {
		t.Errorf("AsString(85*85) = %s; want 100", c.AsString(85*85))
	}
	fmt.Println(c.AsString(205055))
	fmt.Println(c.AsString(-3197852304731898214))
}

func TestBase85LongConverter_Parse(t *testing.T) {
	c := NewBase85LongConverterS()
	if c.AsString(0) != "0" {
		t.Errorf("AsString(0) = %s; want 0", c.AsString(0))
	}
	if c.AsString(1) != "1" {
		t.Errorf("AsString(1) = %s; want 1", c.AsString(1))
	}
	if c.AsString(85) != "10" {
		t.Errorf("AsString(85) = %s; want 10", c.AsString(85))
	}
	if c.AsString(86) != "11" {
		t.Errorf("AsString(86) = %s; want 11", c.AsString(86))
	}
	if c.AsString(85*85) != "100" {
		t.Errorf("AsString(85*85) = %s; want 100", c.AsString(85*85))
	}
	fmt.Println(c.AsString(205055))
	fmt.Println(c.AsString(-3197852304731898214))
}

func TestBase85LongCompare(t *testing.T) {
	c1 := NewBase85LongConverter()
	c2 := NewBase85LongConverterS()
	for i := 0; i < 10000; i++ {
		// random a int64 value
		r, err := randomReadBytesAsInt64()
		if err != nil {
			t.Error(err)
		}
		// convert to base85 string
		fmt.Println("r:", r)
		s1 := c1.AsString(r)
		s2 := c2.AsString(r)
		if s1 != s2 {
			t.Errorf("AsString(%d) = %s; want %s", r, s1, s2)
		}
	}
}

func randomReadBytesAsInt64() (int64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	return int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 |
		int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56, nil
}
