package simpledb

import "github.com/tyabu12/simpledb/file"

type SimpleDB struct {
	filename  string
	blockSize int
	fileMgr   *file.FileManager
}

func New(filename string, blockSize int) (*SimpleDB, error) {
	fileMgr, err := file.NewFileManager(filename, blockSize)
	if err != nil {
		return nil, err
	}
	return &SimpleDB{
		filename:  filename,
		blockSize: blockSize,
		fileMgr:   fileMgr,
	}, nil
}

func (db *SimpleDB) FileMgr() *file.FileManager {
	return db.fileMgr
}
