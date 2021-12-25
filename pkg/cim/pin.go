package cim

import (
	"encoding/binary"
	"fmt"

	"github.com/ghostiam/binstruct"
	"github.com/roffe/cim/pkg/crc16"
)

type Pin struct {
	Pin1         []byte `bin:"len:4"` // LE
	Pin1Unknown  []byte `bin:"len:4"`
	Pin1Checksum uint16 `bin:"Uint16l,len:2"`
	Pin2         []byte `bin:"len:4"` // LE
	Pin2Unknown  []byte `bin:"len:4"`
	Pin2Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*Pin) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

func (p *Pin) validate() error {
	p1ccrc := crc16.Calc(append(p.Pin1[:], p.Pin1Unknown[:]...))
	p2ccrc := crc16.Calc(append(p.Pin2[:], p.Pin2Unknown[:]...))

	if p.Pin1Checksum != p1ccrc {
		return fmt.Errorf("pin 1 checksum in binary does not match calculated checksum")
	}

	if p.Pin2Checksum != p2ccrc {
		return fmt.Errorf("pin 2 checksum in binary does not match calculated checksum")
	}

	if p1ccrc != p2ccrc {
		return fmt.Errorf("calculated checksums does not match")
	}

	if p.Pin1Checksum != p.Pin2Checksum {
		return fmt.Errorf("stored pin checksums in binary does not match")
	}

	return nil
}
