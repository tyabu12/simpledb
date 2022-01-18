package buffer

import (
	"github.com/tyabu12/simpledb/file"
	"github.com/tyabu12/simpledb/log"
)

type Buffer struct {
	fileMgr  *file.Manager
	logMgr   *log.Manager
	contents *file.Page
	blk      *file.BlockId
	pins     int
	txNum    int
	lsn      int
}

func NewBuffer(fileMgr *file.Manager, logMgr *log.Manager) *Buffer {
	return &Buffer{
		fileMgr:  fileMgr,
		logMgr:   logMgr,
		contents: file.NewPageByBlockSize(fileMgr.BlockSize()),
		blk:      nil,
		pins:     0,
		txNum:    -1,
		lsn:      -1,
	}
}

func (buff *Buffer) Contents() *file.Page {
	return buff.contents
}

func (buff *Buffer) Block() *file.BlockId {
	return buff.blk
}

func (buff *Buffer) IsPinned() bool {
	return buff.pins > 0
}

func (buff *Buffer) SetModified(txNum int, lsn int) {
	buff.txNum = txNum
	if lsn >= 0 {
		buff.lsn = lsn
	}
}

func (buff *Buffer) ModifyingTx() int {
	return buff.txNum
}

func (buff *Buffer) assignToBlock(blk *file.BlockId) error {
	if err := buff.flush(); err != nil {
		return err
	}
	buff.blk = blk
	if err := buff.fileMgr.Read(buff.blk, buff.contents); err != nil {
		return err
	}
	buff.pins = 0
	return nil
}

func (buff *Buffer) flush() error {
	if buff.txNum < 0 {
		return nil
	}
	if err := buff.logMgr.Flush(buff.lsn); err != nil {
		return err
	}
	if err := buff.fileMgr.Write(buff.blk, buff.contents); err != nil {
		return err
	}
	buff.txNum = -1
	return nil
}

func (buff *Buffer) pin() {
	buff.pins++
}

func (buff *Buffer) unpin() {
	buff.pins--
}
