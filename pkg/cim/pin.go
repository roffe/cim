package cim

import (
	"encoding/hex"
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type Pin struct {
	Data1     []byte `bin:"len:4"` // LE
	Unknown1  []byte `bin:"len:4"`
	Checksum1 uint16 `bin:"le,len:2"`
	Data2     []byte `bin:"len:4"` // LE
	Unknown2  []byte `bin:"len:4"`
	Checksum2 uint16 `bin:"le,len:2"`
} // 20 bytes

func (p *Pin) validate() error {

	p1crc, p2crc := p.Crc16()

	if p.Checksum1 != p1crc {
		return fmt.Errorf("pin 1 checksum in binary does not match calculated checksum")
	}

	if p.Checksum2 != p2crc {
		return fmt.Errorf("pin 2 checksum in binary does not match calculated checksum")
	}

	if p1crc != p2crc {
		return fmt.Errorf("calculated checksums does not match")
	}

	if p.Checksum1 != p.Checksum2 {
		return fmt.Errorf("stored pin checksums in binary does not match")
	}

	return nil
}

func (p *Pin) Crc16() (uint16, uint16) {
	p1 := crc16.Calc(append(p.Data1[:], p.Unknown1[:]...))
	p2 := crc16.Calc(append(p.Data2[:], p.Unknown2[:]...))
	return p1, p2

}

func (p *Pin) Set(pin string) error {
	if len(pin) != 8 {
		return fmt.Errorf("invalid pin length")
	}

	b, err := hex.DecodeString(pin)
	if err != nil {
		return fmt.Errorf("invalid pin: %v", err)
	}

	p.Data1 = b
	p.Data2 = b
	p.Checksum1, p.Checksum2 = p.Crc16()

	return nil
}
