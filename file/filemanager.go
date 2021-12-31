package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type FileManager struct {
	mu          sync.Mutex
	dbDirectory string
	blockSize   int64
	isNew       bool
	openFiles   map[string]*os.File
}

func NewFileManager(dbDirectory string, blockSize int64) (*FileManager, error) {
	isNew := false
	if !fileExists(dbDirectory) {
		if err := os.MkdirAll(dbDirectory, 0777); err != nil {
			return nil, err
		}
		isNew = true
	}
	if err := removeTemporaryTables(dbDirectory); err != nil {
		return nil, err
	}
	return &FileManager{
		dbDirectory: dbDirectory,
		blockSize:   blockSize,
		isNew:       isNew,
		openFiles:   make(map[string]*os.File),
	}, nil
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func removeTemporaryTables(dbDirectory string) error {
	return filepath.WalkDir(dbDirectory, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed filepath.WalkDir")
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasPrefix(entry.Name(), "temp") {
			return nil
		}
		if err := os.Remove(path); err != nil {
			return err
		}
		return nil
	})
}

func (fm *FileManager) Read(blk *BlockId, p *Page) error {
	fm.lock()
	defer fm.unlock()

	f, err := fm.getFile(blk.filename)
	if err != nil {
		return errors.Wrap(err, "cannot read block "+blk.String())
	}
	if _, err = f.Seek(blk.Number()*fm.BlockSize(), io.SeekStart); err != nil {
		return errors.Wrap(err, "cannot read block "+blk.String())
	}
	if _, err = f.Read(p.contens()); err != nil {
		return errors.Wrap(err, "cannot read block "+blk.String())
	}
	return nil
}

func (fm *FileManager) Write(blk *BlockId, p *Page) error {
	fm.lock()
	defer fm.unlock()

	f, err := fm.getFile(blk.filename)
	if err != nil {
		return errors.Wrap(err, "cannot write block "+blk.String())
	}
	if _, err = f.Seek(blk.Number()*fm.BlockSize(), io.SeekStart); err != nil {
		return errors.Wrap(err, "cannot write block "+blk.String())
	}
	if _, err = f.Write(p.contens()); err != nil {
		return errors.Wrap(err, "cannot write block "+blk.String())
	}
	return nil
}

func (fm *FileManager) Append(filename string) (*BlockId, error) {
	fm.lock()
	defer fm.unlock()

	newBlockNum, err := fm.Length(filename)
	if err != nil {
		return nil, errors.Wrap(err, "cannot append block")
	}
	blk := NewBlockId(filename, newBlockNum)

	b := make([]byte, fm.BlockSize())
	f, err := fm.getFile(blk.filename)
	if _, err = f.Seek(blk.Number()*fm.BlockSize(), io.SeekStart); err != nil {
		return nil, errors.Wrap(err, "cannot append block "+blk.String())
	}
	if _, err = f.Write(b); err != nil {
		return nil, errors.Wrap(err, "cannot append block "+blk.String())
	}
	return blk, nil
}

func (fm *FileManager) Length(filename string) (int64, error) {
	f, err := fm.getFile(filename)
	if err != nil {
		return 0, err
	}
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size() / fm.BlockSize(), nil
}

func (fm *FileManager) IsNew() bool {
	return fm.isNew
}

func (fm *FileManager) BlockSize() int64 {
	return fm.blockSize
}

func (fm *FileManager) getFile(filename string) (*os.File, error) {
	f, ok := fm.openFiles[filename]
	if !ok {
		var err error
		name := filepath.Join(fm.dbDirectory, filename)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, errors.Wrap(err, "cannot open file: "+filename)
		}
		fm.openFiles[filename] = f
	}
	return f, nil
}

func (fm *FileManager) lock() {
	fm.mu.Lock()
}

func (fm *FileManager) unlock() {
	fm.mu.Unlock()
}
