package main

import (
	"fmt"

	"github.com/tyabu12/simpledb/file"
	"github.com/tyabu12/simpledb/log"
	"github.com/tyabu12/simpledb/simpledb"
)

func main() {
	db, err := simpledb.New("test", 400)
	if err != nil {
		panic(err)
	}

	logMgr := db.LogMgr()
	createLogRecords(logMgr, 1, 35)
	printLogRecords(logMgr, "The log file now has these records:")
	createLogRecords(logMgr, 36, 70)
	logMgr.Flush(65)
	printLogRecords(logMgr, "The log file now has these records:")
}

func printLogRecords(logMgr *log.Manager, msg string) {
	fmt.Println(msg)

	it, err := logMgr.Iterator()
	if err != nil {
		panic(err)
	}

	for it.HasNext() {
		rec, err := it.Next()
		if err != nil {
			panic(err)
		}
		page := file.NewPageByBytes(rec)
		s, _, err := page.GetString(0)
		if err != nil {
			panic(err)
		}
		nPos := file.PageMaxLength(len(s))
		val, _, err := page.GetInt(nPos)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[%s, %v]\n", s, val)
	}
	fmt.Println()
}

func createLogRecords(logMgr *log.Manager, start int, end int) {
	fmt.Print("Creating records: ")
	for i := start; i <= end; i++ {
		rec := createLogRecord(logMgr, fmt.Sprintf("record%d", i), i+100)
		lsn, err := logMgr.Append(rec)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v ", lsn)
	}
	fmt.Println()
}

func createLogRecord(logMgr *log.Manager, s string, n int) []byte {
	nPos := file.PageMaxLength(len(s))
	b := make([]byte, nPos+file.SizeOfInt)
	page := file.NewPageByBytes(b)
	if _, err := page.SetString(0, s); err != nil {
		panic(err)
	}
	if _, err := page.SetInt(nPos, n); err != nil {
		panic(err)
	}
	return b
}
