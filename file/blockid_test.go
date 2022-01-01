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
			t.Fatalf("HashCode is failed: %v", err.Error())
		}
		for _, tt2 := range tests {
			blk2 := NewBlockId(tt2.filename, tt2.blknum)
			hash2, err := blk2.HashCode()
			if err != nil {
				t.Fatalf("HashCode is failed: %v", err.Error())
			}
			eq := tt1.filename == tt2.filename && tt1.blknum == tt2.blknum
			if hash1 == hash2 && !eq {
				t.Fatalf("Expected same hash! tt1=%v, tt2=%v", blk1.String(), blk2.String())
			}
			if hash1 != hash2 && eq {
				t.Fatalf("Expected different hash! tt1=%v, tt2=%v", blk1.String(), blk2.String())
			}
		}
	}
}
