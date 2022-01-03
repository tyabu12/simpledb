package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type Manager struct {
	mu          sync.Mutex
	dbDirectory string
	blockSize   int
	isNew       bool
	openFiles   map[string]*os.File
}

func NewManager(dbDirectory string, blockSize int) (*Manager, error) {
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
	return &Manager{
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

func (mgr *Manager) Read(blk *BlockId, p *Page) error {
	mgr.lock()
	defer mgr.unlock()

	f, err := mgr.getFile(blk.filename)
	if err != nil {
		return errors.Wrap(err, "cannot read block "+blk.String())
	}
	if _, err = f.Seek(blk.Number()*int64(mgr.BlockSize()), io.SeekStart); err != nil {
		return errors.Wrap(err, "cannot read block "+blk.String())
	}
	if _, err = f.Read(p.contens()); err != nil {
		return errors.Wrap(err, "cannot read block "+blk.String())
	}
	return nil
}

func (mgr *Manager) Write(blk *BlockId, p *Page) error {
	mgr.lock()
	defer mgr.unlock()

	f, err := mgr.getFile(blk.filename)
	if err != nil {
		return errors.Wrap(err, "cannot write block "+blk.String())
	}
	if _, err = f.Seek(blk.Number()*int64(mgr.BlockSize()), io.SeekStart); err != nil {
		return errors.Wrap(err, "cannot write block "+blk.String())
	}
	if _, err = f.Write(p.contens()); err != nil {
		return errors.Wrap(err, "cannot write block "+blk.String())
	}
	return nil
}

func (mgr *Manager) Append(filename string) (*BlockId, error) {
	mgr.lock()
	defer mgr.unlock()

	newBlockNum, err := mgr.Length(filename)
	if err != nil {
		return nil, errors.Wrap(err, "cannot append block")
	}
	blk := NewBlockId(filename, newBlockNum)

	b := make([]byte, mgr.BlockSize())
	f, err := mgr.getFile(blk.filename)
	if _, err = f.Seek(blk.Number()*int64(mgr.BlockSize()), io.SeekStart); err != nil {
		return nil, errors.Wrap(err, "cannot append block "+blk.String())
	}
	if _, err = f.Write(b); err != nil {
		return nil, errors.Wrap(err, "cannot append block "+blk.String())
	}
	return blk, nil
}

func (mgr *Manager) Length(filename string) (int64, error) {
	f, err := mgr.getFile(filename)
	if err != nil {
		return 0, err
	}
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size() / int64(mgr.BlockSize()), nil
}

func (mgr *Manager) IsNew() bool {
	return mgr.isNew
}

func (mgr *Manager) BlockSize() int {
	return mgr.blockSize
}

func (mgr *Manager) getFile(filename string) (*os.File, error) {
	f, ok := mgr.openFiles[filename]
	if !ok {
		var err error
		name := filepath.Join(mgr.dbDirectory, filename)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, errors.Wrap(err, "cannot open file: "+filename)
		}
		mgr.openFiles[filename] = f
	}
	return f, nil
}

func (mgr *Manager) lock() {
	mgr.mu.Lock()
}

func (mgr *Manager) unlock() {
	mgr.mu.Unlock()
}
