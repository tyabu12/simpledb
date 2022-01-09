package file

import (
	"bytes"
	"testing"
)

func TestSetGetInt(t *testing.T) {
	tests := []struct {
		blockSize int
		offset    int
		val       int
	}{
		{10, 0, 0},
		{100, 0, 100},
		{100, 10, 110},
		{100, 30, 0},
		{100, 50, 30},
		{100, 50, -30},
	}

	for _, tt := range tests {
		p := NewPageByBlockSize(tt.blockSize)

		nextOffset, err := p.SetInt(tt.offset, tt.val)
		if err != nil {
			t.Fatal("SetInt is failed: " + err.Error())
		}
		if nextOffset != tt.offset+SizeOfInt {
			t.Fatalf("offset of SetInt is invalid, expected=%v, got=%v", tt.offset+SizeOfInt, nextOffset)
		}

		val, nextOffset, err := p.GetInt(tt.offset)
		if err != nil {
			t.Fatal("GetInt is failed: " + err.Error())
		}
		if nextOffset != tt.offset+SizeOfInt {
			t.Fatalf("offset of GetInt is invalid, expected=%v, got=%v", tt.offset+SizeOfInt, nextOffset)
		}
		if val != tt.val {
			t.Fatalf("value of GetInt is invalid, expected=%v, got=%v", tt.val, val)
		}
	}
}

func TestSetGetBytes(t *testing.T) {
	tests := []struct {
		blockSize int
		offset    int
		val       []byte
	}{
		{100, 0, []byte{}},
		{100, 0, []byte{0x01, 0x02, 0x03}},
		{100, 10, []byte{0x33, 0x44, 0x2B}},
		{100, 30, []byte{0xAF, 0x88, 0x33, 0x11, 0x9A}},
		{100, 50, []byte("abcefghigjklmn")},
	}

	for _, tt := range tests {
		p := NewPageByBlockSize(tt.blockSize)
		sizeofBytes := len(tt.val)

		nextOffset, err := p.SetBytes(tt.offset, tt.val)
		if err != nil {
			t.Fatal("SetBytes is failed: " + err.Error())
		}
		if nextOffset != tt.offset+SizeOfInt+sizeofBytes {
			t.Fatalf("offset of SetBytes is invalid, expected=%v, got=%v", tt.offset+SizeOfInt+sizeofBytes, nextOffset)
		}

		val, nextOffset, err := p.GetBytes(tt.offset)
		if err != nil {
			t.Fatal("GetInt is failed: " + err.Error())
		}
		if nextOffset != tt.offset+SizeOfInt+sizeofBytes {
			t.Fatalf("offset of GetBytes is invalid, expected=%v, got=%v", tt.offset+SizeOfInt+sizeofBytes, nextOffset)
		}
		if bytes.Compare(val, tt.val) != 0 {
			t.Fatalf("value of GetBytes is invalid, expected=%v, got=%v", tt.val, val)
		}
	}
}

func TestSetGetString(t *testing.T) {
	tests := []struct {
		blockSize int
		offset    int
		val       string
	}{
		{100, 0, "abc"},
		{100, 10, "jwaojfieawo024"},
		{100, 30, ""},
		{100, 50, "abaca"},
	}

	for _, tt := range tests {
		p := NewPageByBlockSize(tt.blockSize)
		sizeofString := len([]byte(tt.val))

		nextOffset, err := p.SetString(tt.offset, tt.val)
		if err != nil {
			t.Fatal("SetString is failed: " + err.Error())
		}
		if nextOffset != tt.offset+SizeOfInt+sizeofString {
			t.Fatalf("offset of SetString is invalid, expected=%v, got=%v", tt.offset+SizeOfInt+sizeofString, nextOffset)
		}

		val, nextOffset, err := p.GetString(tt.offset)
		if err != nil {
			t.Fatal("GetInt is failed: " + err.Error())
		}
		if nextOffset != tt.offset+SizeOfInt+sizeofString {
			t.Fatalf("offset of GetString is invalid, expected=%v, got=%v", tt.offset+SizeOfInt+sizeofString, nextOffset)
		}
		if val != tt.val {
			t.Fatalf("value of GetString is invalid, expected=%v, got=%v", tt.val, val)
		}
	}
}
