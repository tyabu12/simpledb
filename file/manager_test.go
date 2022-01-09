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
		fileMgr, err := NewManager(t.TempDir(), 400)
		if err != nil {
			t.Fatal("NewManager is failed: " + err.Error())
		}

		blk, err := fileMgr.Append("tempTest")
		if err != nil {
			t.Fatal("Apeend new Block is failed: " + err.Error())
		}

		p1 := NewPageByBlockSize(fileMgr.BlockSize())
		pos1 := tt.startOffset
		if _, err := p1.SetString(pos1, tt.stringVal); err != nil {
			t.Fatal("SetString to page1 is failed: " + err.Error())
		}
		size := PageMaxLength(len(tt.stringVal))
		pos2 := pos1 + size
		if _, err := p1.SetInt(pos2, tt.intVal); err != nil {
			t.Fatal("SetInt to page1 is failed: " + err.Error())
		}
		if err := fileMgr.Write(blk, p1); err != nil {
			t.Fatal("Writing page is failed: " + err.Error())
		}

		p2 := NewPageByBlockSize(fileMgr.BlockSize())
		if err := fileMgr.Read(blk, p2); err != nil {
			t.Fatal("Reading page is failed: " + err.Error())
		}
		intVal, _, err := p2.GetInt(pos2)
		if err != nil {
			t.Fatal("GetInt is failed: " + err.Error())
		}
		if intVal != tt.intVal {
			t.Fatalf("value of GetInt is invalid, expected=%v, got=%v", tt.intVal, intVal)
		}

		stringVal, _, err := p2.GetString(pos1)
		if err != nil {
			t.Fatal("GetString is failed: " + err.Error())
		}
		if stringVal != tt.stringVal {
			t.Fatalf("value of GetString is invalid, expected=%v, got=%v", tt.intVal, intVal)
		}
	}
}
