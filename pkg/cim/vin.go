package cim

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/ghostiam/binstruct"
	"github.com/roffe/cim/pkg/crc16"
)

type Vin struct {
	Data     string `bin:"len:17"` // Vin as ASCII
	Value    uint8  `bit:"len:1"`  // unknown value, seems to be a counter?
	Unknown  []byte `bin:"len:9"`
	SpsCount uint8
	Checksum uint16 `bin:"Uint16l,len:2"` // CRC16 MCRF4XX
} // 30 bytes

// Uint16l returns a uint16 read as little endian
func (*Vin) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

func (bin *Bin) validateVin() error {
	v := bin.Vin
	var b []byte
	b = append(b, []byte(v.Data)...)
	b = append(b, byte(v.Value))
	b = append(b, v.Unknown...)
	b = append(b, v.SpsCount)
	c := crc16.Calc(b)
	if v.Checksum != c {
		log.Printf("%d || %d", v.Checksum, c)
		return fmt.Errorf("vin cheksum validation failed")

	}
	return nil
}
