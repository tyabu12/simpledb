package buffer

import (
	"testing"

	"github.com/tyabu12/simpledb/file"
	"github.com/tyabu12/simpledb/log"
)

func TestBufferManager(t *testing.T) {
	fileMgr, err := file.NewManager(t.TempDir(), 400)
	if err != nil {
		t.Fatal("new FileManager is failed: " + err.Error())
	}
	logMgr, err := log.NewManager(fileMgr, "temp.log")
	if err != nil {
		t.Fatal("new LogManager is failed: " + err.Error())
	}
	mgr := NewManager(fileMgr, logMgr, 3)

	buf1, err := mgr.Pin(file.NewBlockId("testfile", 1))
	if err != nil {
		t.Fatal("Pin is failed: " + err.Error())
	}
	p := buf1.Contents()
	n, _, err := p.GetInt(80)
	if err != nil {
		t.Fatal("Page.GetInt is failed: " + err.Error())
	}
	if _, err := p.SetInt(80, n+1); err != nil {
		t.Fatal("Page.SetInt is failed: " + err.Error())
	}
	buf1.SetModified(1, 0)
	if err := mgr.Unpin(buf1); err != nil {
		t.Fatal("Unpin is failed: " + err.Error())
	}

	buf2, err := mgr.Pin(file.NewBlockId("testfile", 2))
	if err != nil {
		t.Fatal("Pin is failed: " + err.Error())
	}
	if _, err := mgr.Pin(file.NewBlockId("testfile", 3)); err != nil {
		t.Fatal("Pin is failed: " + err.Error())
	}
	if _, err := mgr.Pin(file.NewBlockId("testfile", 4)); err != nil {
		t.Fatal("Pin is failed: " + err.Error())
	}

	if err := mgr.Unpin(buf2); err != nil {
		t.Fatal("Unpin is failed: " + err.Error())
	}
	buf2, err = mgr.Pin(file.NewBlockId("testfile", 1))
	if err != nil {
		t.Fatal("Pin is failed: " + err.Error())
	}
	p2 := buf2.Contents()
	if _, err := p2.SetInt(80, 9999); err != nil {
		t.Fatal("Page.SetInt is failed: " + err.Error())
	}
	buf2.SetModified(1, 0)
	if err := mgr.Unpin(buf2); err != nil {
		t.Fatal("Unpin is failed: " + err.Error())
	}
}

// func TestBufferManager2(t *testing.T) {
// 	fileMgr, err := file.NewManager(t.TempDir(), 400)
// 	if err != nil {
// 		t.Fatal("new FileManager is failed: " + err.Error())
// 	}
// 	logMgr, err := log.NewManager(fileMgr, "temp.log")
// 	if err != nil {
// 		t.Fatal("new LogManager is failed: " + err.Error())
// 	}
// 	mgr := NewManager(fileMgr, logMgr, 3)

// 	buff := make([]*Buffer, 6)
// 	for i := 0; i < 3; i++ {
// 		buff[i], err = mgr.Pin(file.NewBlockId("testfile", int64(i)))
// 		if err != nil {
// 			t.Fatal("Pin is failed: " + err.Error())
// 		}
// 	}
// 	if err := mgr.Unpin(buff[1]); err != nil {
// 		t.Fatal("Unpin is failed: " + err.Error())
// 	}
// 	buff[1] = nil

// 	buff[3], err = mgr.Pin(file.NewBlockId("testfile", 0))
// 	if err != nil {
// 		t.Fatal("Pin is failed: " + err.Error())
// 	}
// 	buff[4], err = mgr.Pin(file.NewBlockId("testfile", 1))
// 	if err != nil {
// 		t.Fatal("Pin is failed: " + err.Error())
// 	}
// 	buff[5], err = mgr.Pin(file.NewBlockId("testfile", 3))
// 	if err != ErrNoAvailableBuffer {
// 		t.Fatalf("Error expected: %s, actual: %v", ErrNoAvailableBuffer, err)
// 	}

// 	if err := mgr.Unpin(buff[2]); err != nil {
// 		t.Fatal("Unpin is failed: " + err.Error())
// 	}
// 	buff[2] = nil

// 	buff[5], err = mgr.Pin(file.NewBlockId("testfile", 3))
// 	if err != nil {
// 		t.Fatal("Pin is failed: " + err.Error())
// 	}
// }
