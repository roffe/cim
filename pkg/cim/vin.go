package cim

import (
	"bytes"
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type Vin struct {
	Data     string `bin:"len:17" json:"data"` // Vin as ASCII
	Value    uint8  `bin:"len:1" json:"value"` // unknown value, seems to be a counter?
	Unknown  []byte `bin:"len:9" json:"unknown"`
	SpsCount uint8  `bin:"len:1" json:"spscount"`
	Checksum uint16 `bin:"le,len:2" json:"checksum"` // CRC16 MCRF4XX
} // 30 bytes

func (v *Vin) validate() error {
	c := v.Crc16()
	if v.Checksum != c {
		return fmt.Errorf("vin cheksum validation failed, calculated: %X in bin: %X", c, v.Checksum)
	}
	return nil
}

func (v *Vin) Crc16() uint16 {
	b := bytes.NewBuffer(nil)
	b.Write([]byte(v.Data))
	b.WriteByte(v.Value)
	b.Write(v.Unknown)
	b.WriteByte(v.SpsCount)
	return crc16.Calc(b.Bytes())
}

func (v *Vin) Set(vin string) error {
	if len(vin) > 17 {
		return fmt.Errorf("vin is to long")
	}
	v.Data = fmt.Sprintf("%-17s", vin)
	v.updateChecksum()
	return nil
}

func (v *Vin) SetValue(val uint8) {
	v.Value = val
	v.updateChecksum()
}

func (v *Vin) SetSpsCount(value uint8) {
	v.SpsCount = value
	v.updateChecksum()
}

func (v *Vin) updateChecksum() {
	v.Checksum = v.Crc16()
}
