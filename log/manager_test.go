package log

import (
	"testing"

	"github.com/tyabu12/simpledb/file"
)

func TestManager(t *testing.T) {
	tests := []struct {
		texts []string
	}{
		{texts: []string{"abc123", "lorem ipsum", "jfeiajifao", "fejaiof"}},
		{texts: []string{"abc123", "lorem ipsum", "jfeiajifao", "fejaiof", "jfioeajifoajoifejajfeoaj awjeofajfj fjaofjeaojf", "439892u98u9afjaoa jaj1joajwaij"}},
	}

	fileMgr, err := file.NewManager(t.TempDir(), 71)
	if err != nil {
		t.Fatal("new FileManager is failed: " + err.Error())
	}
	logMgr, err := NewManager(fileMgr, "temp.log")
	if err != nil {
		t.Fatal("new LogManager is failed: " + err.Error())
	}

	for _, tt := range tests {
		createLogRecord(t, logMgr, tt.texts)
		logMgr.Flush(65)
		checkLogRecords(t, logMgr, tt.texts)
	}
}

func checkLogRecords(t *testing.T, logMgr *Manager, expectedTexts []string) {
	it, err := logMgr.Iterator()
	if err != nil {
		t.Fatal("LogManager.Iterator is failed: " + err.Error())
	}
	for i := len(expectedTexts) - 1; i >= 0; i-- {
		expected := expectedTexts[i]
		if !it.HasNext() {
			t.Fatal("HasNext is expected true")
		}
		rec, err := it.Next()
		if err != nil {
			t.Fatal("Next is failed: " + err.Error())
		}
		page := file.NewPageByBytes(rec)
		text, _, err := page.GetString(0)
		if err != nil {
			t.Fatal("page.GetString(0) is failed: " + err.Error())
		}
		if text != expected {
			t.Fatalf("expected=%v, actual=%v", expected, text)
		}
	}
}

func createLogRecord(t *testing.T, logMgr *Manager, texts []string) {
	for _, text := range texts {
		rec := make([]byte, file.PageMaxLength(len(text)))
		page := file.NewPageByBytes(rec)
		if _, err := page.SetString(0, text); err != nil {
			t.Fatal("page.SetString(0) is failed: " + err.Error())
		}
		if _, err := logMgr.Append(rec); err != nil {
			t.Fatal("LogManager.Append is failed: " + err.Error())
		}
	}
}
