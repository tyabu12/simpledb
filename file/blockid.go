package file

import (
	"crypto/sha1"
	"fmt"
)

type BlockId struct {
	filename string
	blknum   int64
}

func NewBlockId(filename string, blknum int64) *BlockId {
	return &BlockId{
		filename: filename,
		blknum:   blknum,
	}
}

func (blk *BlockId) Filename() string {
	return blk.filename
}

func (blk *BlockId) Number() int64 {
	return blk.blknum
}

func (blk *BlockId) Equals(blk2 *BlockId) bool {
	return blk.filename == blk2.filename && blk.blknum == blk2.blknum
}

func (blk *BlockId) String() string {
	return fmt.Sprintf("[file=%v, block=%v]", blk.filename, blk.blknum)
}

func (blk *BlockId) HashCode() (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(blk.String()))
	if err != nil {
		return "", err
	}
	return string(hash.Sum(nil)), nil
}
