package log

import (
	"sync"

	"github.com/tyabu12/simpledb/file"
)

type Manager struct {
	mu         sync.Mutex
	fileMgr    *file.Manager
	logFile    string
	logPage    *file.Page
	currentBlk *file.BlockId
	// LSN = Log Sequence Number
	latestLSN    int
	lastSavedLSN int
}

func NewManager(fileMgr *file.Manager, logFile string) (*Manager, error) {
	mgr := &Manager{
		fileMgr:      fileMgr,
		logFile:      logFile,
		logPage:      file.NewPageByBlockSize(fileMgr.BlockSize()),
		latestLSN:    0,
		lastSavedLSN: 0,
	}

	// set currentBlk
	logSize, err := fileMgr.Length(logFile)
	if err != nil {
		return nil, err
	}
	if logSize == 0 {
		mgr.currentBlk, err = mgr.appendNewBlock()
		if err != nil {
			return nil, err
		}
	} else {
		mgr.currentBlk = file.NewBlockId(logFile, logSize-1)
		if err = fileMgr.Read(mgr.currentBlk, mgr.logPage); err != nil {
			return nil, err
		}
	}

	return mgr, nil
}

func (mgr *Manager) Flush(lsn int) error {
	if lsn >= mgr.lastSavedLSN {
		return mgr.flush()
	}
	return nil
}

func (mgr *Manager) Iterator() (*Iterator, error) {
	if err := mgr.flush(); err != nil {
		return nil, err
	}
	it, err := NewIterator(mgr.fileMgr, mgr.currentBlk)
	if err != nil {
		return nil, err
	}
	return it, nil
}

func (mgr *Manager) Append(logRec []byte) (int, error) {
	mgr.lock()
	defer mgr.unlock()

	boundary, _, err := mgr.logPage.GetInt(0)
	if err != nil {
		return 0, err
	}
	recSize := len(logRec)
	bytesNeeded := recSize + file.SizeOfInt

	if boundary-bytesNeeded < file.SizeOfInt {
		if err := mgr.flush(); err != nil {
			return 0, err
		}
		mgr.currentBlk, err = mgr.appendNewBlock()
		if err != nil {
			return 0, err
		}
		boundary, _, err = mgr.logPage.GetInt(0)
		if err != nil {
			return 0, err
		}
	}

	recPos := boundary - bytesNeeded
	if _, err := mgr.logPage.SetBytes(recPos, logRec); err != nil {
		return 0, err
	}
	if _, err := mgr.logPage.SetInt(0, recPos); err != nil {
		return 0, err
	}

	mgr.latestLSN++

	return mgr.latestLSN, nil
}

func (mgr *Manager) appendNewBlock() (*file.BlockId, error) {
	blk, err := mgr.fileMgr.Append(mgr.logFile)
	if err != nil {
		return nil, err
	}
	if _, err := mgr.logPage.SetInt(0, mgr.fileMgr.BlockSize()); err != nil {
		return nil, err
	}
	if err := mgr.fileMgr.Write(blk, mgr.logPage); err != nil {
		return nil, err
	}
	return blk, nil
}

func (mgr *Manager) flush() error {
	err := mgr.fileMgr.Write(mgr.currentBlk, mgr.logPage)
	if err != nil {
		return err
	}
	mgr.lastSavedLSN = mgr.latestLSN
	return nil
}

func (mgr *Manager) lock() {
	mgr.mu.Lock()
}

func (mgr *Manager) unlock() {
	mgr.mu.Unlock()
}
