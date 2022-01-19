package buffer

import (
	"errors"
	"sync"
	"time"

	"github.com/tyabu12/simpledb/file"
	"github.com/tyabu12/simpledb/log"
)

const (
	timeoutMilliSeconds = 10000
)

var (
	ErrNoAvailableBuffer = errors.New("No available buffers")
)

type Manager struct {
	mu           sync.Mutex
	cond         *sync.Cond
	bufferPool   []*Buffer
	numAvailable int
}

func NewManager(fileMgr *file.Manager, logMgr *log.Manager, numBuffers int) *Manager {
	bufferPool := []*Buffer{}
	for i := 0; i < numBuffers; i++ {
		bufferPool = append(bufferPool, NewBuffer(fileMgr, logMgr))
	}
	mgr := &Manager{
		bufferPool:   bufferPool,
		numAvailable: numBuffers,
	}
	mgr.cond = sync.NewCond(&mgr.mu)
	return mgr
}

func (mgr *Manager) Pin(blk *file.BlockId) (*Buffer, error) {
	mgr.lock()
	defer mgr.unlock()

	startTime := time.Now()
	for {
		buff, err := mgr.tryToPin(blk)
		if buff != nil {
			return buff, nil
		}
		if err != nil {
			return nil, err
		}
		mgr.cond.Wait()
		if mgr.waitingTooLong(startTime) {
			return nil, ErrNoAvailableBuffer
		}
	}
}

func (mgr *Manager) Unpin(buff *Buffer) error {
	mgr.lock()
	defer mgr.unlock()

	buff.unpin()
	if !buff.IsPinned() {
		mgr.numAvailable++
		mgr.cond.Broadcast()
	}
	return nil
}

func (mgr *Manager) Available() int {
	mgr.lock()
	defer mgr.unlock()
	return mgr.numAvailable
}

func (mgr *Manager) FlushAll(txNum int) error {
	mgr.lock()
	defer mgr.unlock()

	for _, buff := range mgr.bufferPool {
		if buff.ModifyingTx() == txNum {
			if err := buff.flush(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mgr *Manager) waitingTooLong(startTime time.Time) bool {
	return time.Now().Sub(startTime).Milliseconds() > timeoutMilliSeconds
}

func (mgr *Manager) tryToPin(blk *file.BlockId) (*Buffer, error) {
	var buff *Buffer
	buff = mgr.findExistingBuffer(blk)
	if buff == nil {
		buff = mgr.chooseUnpinnedBuffer()
		if buff == nil {
			return nil, nil
		}
		if err := buff.assignToBlock(blk); err != nil {
			return nil, err
		}
	}
	if !buff.IsPinned() {
		mgr.numAvailable--
	}
	buff.pin()
	return buff, nil
}

func (mgr *Manager) findExistingBuffer(blk *file.BlockId) *Buffer {
	for _, buff := range mgr.bufferPool {
		b := buff.Block()
		if b != nil && b.Equals(blk) {
			return buff
		}
	}
	return nil
}

func (mgr *Manager) chooseUnpinnedBuffer() *Buffer {
	for _, buff := range mgr.bufferPool {
		if !buff.IsPinned() {
			return buff
		}
	}
	return nil
}

func (mgr *Manager) lock() {
	mgr.mu.Lock()
}

func (mgr *Manager) unlock() {
	mgr.mu.Unlock()
}
