package log

import (
	"github.com/tyabu12/simpledb/file"
)

type Iterator struct {
	fileMgr    *file.Manager
	blk        *file.BlockId
	page       *file.Page
	currentPos int
}

func NewIterator(fileMgr *file.Manager, blk *file.BlockId) (*Iterator, error) {
	it := &Iterator{
		fileMgr: fileMgr,
		page:    file.NewPageByBlockSize(fileMgr.BlockSize()),
	}
	if err := it.moveToBlock(blk); err != nil {
		return nil, err
	}
	return it, nil
}

func (it *Iterator) HasNext() bool {
	return it.currentPos < it.fileMgr.BlockSize() || it.blk.Number() > 0
}

func (it *Iterator) Next() ([]byte, error) {
	if it.currentPos == it.fileMgr.BlockSize() {
		blk := file.NewBlockId(it.blk.Filename(), it.blk.Number()-1)
		if err := it.moveToBlock(blk); err != nil {
			return nil, err
		}
	}
	var rec []byte
	var err error
	rec, it.currentPos, err = it.page.GetBytes(it.currentPos)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (it *Iterator) moveToBlock(blk *file.BlockId) error {
	var err error
	it.blk = blk
	err = it.fileMgr.Read(it.blk, it.page)
	if err != nil {
		return err
	}
	it.currentPos, _, err = it.page.GetInt(0)
	if err != nil {
		return err
	}
	return nil
}
