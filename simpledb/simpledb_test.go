package simpledb

import "testing"

func TestNewSimpleDB(t *testing.T) {
	db, err := New(t.TempDir(), 400)
	if err != nil {
		t.Fatal("SimpleDB.New is failed: " + err.Error())
	}
	if db.FileMgr() == nil {
		t.Fatal("FileMgr() is nil")
	}
	if db.LogMgr() == nil {
		t.Fatal("LogMgr() is nil")
	}
}
