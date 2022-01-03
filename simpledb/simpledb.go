package simpledb

import (
	"github.com/tyabu12/simpledb/file"
	"github.com/tyabu12/simpledb/log"
)

const (
	logFileName = "simpledb.log"
)

type SimpleDB struct {
	filename  string
	blockSize int
	fileMgr   *file.Manager
	logMgr    *log.Manager
}

func New(filename string, blockSize int) (*SimpleDB, error) {
	fileMgr, err := file.NewManager(filename, blockSize)
	if err != nil {
		return nil, err
	}
	logMgr, err := log.NewManager(fileMgr, logFileName)
	if err != nil {
		return nil, err
	}
	return &SimpleDB{
		filename:  filename,
		blockSize: blockSize,
		fileMgr:   fileMgr,
		logMgr:    logMgr,
	}, nil
}

func (db *SimpleDB) FileMgr() *file.Manager {
	return db.fileMgr
}

func (db *SimpleDB) LogMgr() *log.Manager {
	return db.logMgr
}
