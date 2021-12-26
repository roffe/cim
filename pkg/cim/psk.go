package cim

import (
	"encoding/binary"

	"github.com/ghostiam/binstruct"
)

type PSK struct {
	Low      []byte `bin:"len:4"`
	High     []byte `bin:"len:2"`
	Constant []byte `bin:"len:4"`
	Unknown  []byte `bin:"len:2"`
	Checksum uint16 `bin:"Uint16l,len:2"` // CRC16 MCRF4XX
}

// Uint16l returns a uint16 read as little endian
func (*PSK) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}
