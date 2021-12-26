package cim

import (
	"encoding/binary"

	"github.com/ghostiam/binstruct"
)

type Sync struct {
	Data     [][]byte `bin:"len:5,[len:4]"`
	Checksum uint16   `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*Sync) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}
