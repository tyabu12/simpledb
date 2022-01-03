package file

import (
	"testing"
)

func TestManager(t *testing.T) {
	tests := []struct {
		stringVal   string
		intVal      int
		startOffset int
	}{
		{
			stringVal:   "abcdefghijklm",
			intVal:      345,
			startOffset: 88,
		},
	}

	for _, tt := range tests {
		fileMgr, err := NewManager("filetest", 400)
		if err != nil {
			t.Fatalf("NewManager is failed: %v", err.Error())
		}

		blk, err := fileMgr.Append("tempTest")
		if err != nil {
			t.Fatalf("Apeend new Block is failed: %v", err.Error())
		}

		p1 := NewPageByBlockSize(fileMgr.BlockSize())
		pos1 := tt.startOffset
		if _, err := p1.SetString(pos1, tt.stringVal); err != nil {
			t.Fatalf("SetString to page1 is failed: %v", err.Error())
		}
		size := PageMaxLength(len(tt.stringVal))
		pos2 := pos1 + size
		if _, err := p1.SetInt(pos2, tt.intVal); err != nil {
			t.Fatalf("SetInt to page1 is failed: %v", err.Error())
		}
		if err := fileMgr.Write(blk, p1); err != nil {
			t.Fatalf("Writing page is failed: %v", err.Error())
		}

		p2 := NewPageByBlockSize(fileMgr.BlockSize())
		if err := fileMgr.Read(blk, p2); err != nil {
			t.Fatalf("Reading page is failed: %v", err.Error())
		}
		intVal, _, err := p2.GetInt(pos2)
		if err != nil {
			t.Fatalf("GetInt is failed: %s", err.Error())
		}
		if intVal != tt.intVal {
			t.Fatalf("value of GetInt is invalid, expected=%v, got=%v", tt.intVal, intVal)
		}

		stringVal, _, err := p2.GetString(pos1)
		if err != nil {
			t.Fatalf("GetString is failed: %s", err.Error())
		}
		if stringVal != tt.stringVal {
			t.Fatalf("value of GetString is invalid, expected=%v, got=%v", tt.intVal, intVal)
		}
	}
}
