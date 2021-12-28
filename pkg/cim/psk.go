package cim

import (
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type PSK struct {
	Low      []byte `bin:"len:4"`
	High     []byte `bin:"len:2"`
	Constant []byte `bin:"len:4"`
	Unknown  []byte `bin:"len:2"`
	Checksum uint16 `bin:"le,len:2"` // CRC16 MCRF4XX
} // 14 bytes

func (p *PSK) validate() error {
	if p.Checksum != p.Crc16() {
		return fmt.Errorf("psk data checksum does not match calculated checksum")
	}
	return nil
}

func (p *PSK) Crc16() uint16 {
	var b []byte
	b = append(b, p.Low...)
	b = append(b, p.High...)
	b = append(b, p.Constant...)
	b = append(b, p.Unknown...)
	return crc16.Calc(b)
}
