package cim

import (
	"bytes"
	"fmt"
	"log"

	"github.com/roffe/cim/pkg/crc16"
)

type Vin struct {
	Data     string `bin:"len:17"` // Vin as ASCII
	Counter  uint8  `bin:"len:1"`  // unknown value, seems to be a counter?
	Unknown  []byte `bin:"len:9"`
	SpsCount uint8  `bin:"len:1"`
	Checksum uint16 `bin:"le,len:2"` // CRC16 MCRF4XX
} // 30 bytes

func (v *Vin) validate() error {
	c := v.Crc16()
	if v.Checksum != c {
		log.Printf("%d || %d", v.Checksum, c)
		return fmt.Errorf("vin cheksum validation failed")
	}
	return nil
}

func (v *Vin) Crc16() uint16 {
	b := bytes.NewBuffer(nil)
	b.Write([]byte(v.Data))
	b.WriteByte(v.Counter)
	b.Write(v.Unknown)
	b.WriteByte(v.SpsCount)
	return crc16.Calc(b.Bytes())
}

func (v *Vin) Set(vin string) error {
	if len(vin) > 17 {
		return fmt.Errorf("vin is to long")
	}
	v.Data = fmt.Sprintf("%-17s", vin)
	v.Checksum = v.Crc16()

	return nil
}
