package cim

import (
	"encoding/binary"

	"github.com/ghostiam/binstruct"
)

type Const1 struct {
	Data     []byte `bin:"len:8"`
	Checksum uint16 `bin:"Uint16l,len:2"` // CRC16 MCRF4XX
}

// Uint16l returns a uint16 read as little endian
func (*Const1) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}
