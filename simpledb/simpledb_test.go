package simpledb

import "testing"

func TestNewSimpleDB(t *testing.T) {
	db, err := New(t.TempDir(), 400)
	if err != nil {
		t.Fatalf("SimpleDB.New is failed: " + err.Error())
	}
	if db.FileMgr() == nil {
		t.Fatalf("FileMgr() is nil")
	}
	if db.LogMgr() == nil {
		t.Fatalf("LogMgr() is nil")
	}
}
