package main

import (
	"fmt"

	"github.com/tyabu12/simpledb/file"
	"github.com/tyabu12/simpledb/simpledb"
)

func main() {
	db, err := simpledb.New("filetest", 400)
	if err != nil {
		panic(err)
	}
	fm := db.FileMgr()

	blk, err := fm.Append("tempTest")
	if err != nil {
		panic(err)
	}

	p1 := file.NewPageByBlockSize(fm.BlockSize())
	pos1 := 88
	s := "abcdefghijklm"
	if _, err := p1.SetString(pos1, s); err != nil {
		panic(err)
	}
	size := file.PageMaxLength(len(s))
	pos2 := pos1 + size
	if _, err := p1.SetInt(pos2, 345); err != nil {
		panic(err)
	}
	if err := fm.Write(blk, p1); err != nil {
		panic(err)
	}

	p2 := file.NewPageByBlockSize(fm.BlockSize())
	if err := fm.Read(blk, p2); err != nil {
		panic(err)
	}
	val2, _, err := p2.GetInt(pos2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("offset %v contains %v\n", pos2, val2)
	val1, _, err := p2.GetString(pos1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("offset %v contains %v\n", pos1, val1)
}
