package file

import "testing"

func TestHashCode(t *testing.T) {
	tests := []struct {
		filename string
		blknum   int64
	}{
		{filename: "abc", blknum: 0},
		{filename: "abc", blknum: 1},
		{filename: "abc", blknum: 2},
		{filename: "abc", blknum: 3},
		{filename: "abc", blknum: 23},
		{filename: "abcd", blknum: 0},
		{filename: "abcd", blknum: 1},
		{filename: "abcd", blknum: 3},
		{filename: "fheoafi", blknum: 10},
	}

	for _, tt1 := range tests {
		blk1 := NewBlockId(tt1.filename, tt1.blknum)
		hash1, err := blk1.HashCode()
		if err != nil {
			t.Fatal("HashCode is failed: " + err.Error())
		}
		for _, tt2 := range tests {
			blk2 := NewBlockId(tt2.filename, tt2.blknum)
			hash2, err := blk2.HashCode()
			if err != nil {
				t.Fatal("HashCode is failed: " + err.Error())
			}
			eq := tt1.filename == tt2.filename && tt1.blknum == tt2.blknum
			if eq && hash1 != hash2 {
				t.Fatalf("Expected same hash, but got different hash! block1=%v, block2=%v", blk1.String(), blk2.String())
			}
			if !eq && hash1 == hash2 {
				t.Fatalf("Expected different hash, but got same hash! block1=%v, block2=%v", blk1.String(), blk2.String())
			}
		}
	}
}
