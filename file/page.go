package file

import (
	"encoding/binary"
	"errors"
)

const (
	SizeOfInt = 8
)

type Page struct {
	bb []byte
}

// A constructor for creating data buffers
func NewPageByBlockSize(blockSize int) *Page {
	return &Page{bb: make([]byte, blockSize)}
}

// A constructor for creating log page
func NewPageByBytes(b []byte) *Page {
	return &Page{bb: b}
}

func (p *Page) GetInt(offset int) (int, int, error) {
	if offset+SizeOfInt > len(p.bb) {
		return 0, 0, errors.New("out of offset")
	}
	return int(binary.LittleEndian.Uint64(p.bb[offset : offset+SizeOfInt])), offset + SizeOfInt, nil
}

func (p *Page) GetBytes(offset int) ([]byte, int, error) {
	length, offset, err := p.GetInt(offset)
	if err != nil {
		return nil, 0, err
	}
	if offset > len(p.bb) {
		return nil, 0, errors.New("out of offset")
	}
	return p.bb[offset : offset+int(length)], offset + int(length), nil
}

func (p *Page) GetString(offset int) (string, int, error) {
	b, offset, err := p.GetBytes(offset)
	if err != nil {
		return "", 0, err
	}
	return string(b), offset, nil
}

func (p *Page) SetInt(offset int, val int) (int, error) {
	if offset+SizeOfInt > len(p.bb) {
		return 0, errors.New("out of offset")
	}
	b := make([]byte, SizeOfInt)
	binary.LittleEndian.PutUint64(b, uint64(val))
	for i := 0; i < SizeOfInt; i++ {
		p.bb[offset+i] = b[i]
	}
	return offset + SizeOfInt, nil
}

func (p *Page) SetBytes(offset int, val []byte) (int, error) {
	offset, err := p.SetInt(offset, len(val))
	if err != nil {
		return 0, err
	}
	if offset+len(val) > len(p.bb) {
		return 0, errors.New("out of page")
	}
	for i, v := range val {
		p.bb[offset+i] = v
	}
	return offset + len(val), nil
}

func (p *Page) SetString(offset int, val string) (int, error) {
	return p.SetBytes(offset, []byte(val))
}

func PageMaxLength(strlen int) int {
	return SizeOfInt + strlen
}

func (p *Page) contens() []byte {
	return p.bb
}
