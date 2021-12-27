package cim

import (
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type Const1 struct {
	Data     []byte `bin:"len:8"`
	Checksum uint16 `bin:"le,len:2"` // CRC16 MCRF4XX
} // 10 bytes

func (c *Const1) validate() error {
	if c.Checksum != c.Crc16() {
		return fmt.Errorf("const1 checksum does not match calculated checksum")
	}
	return nil
}

func (c *Const1) Crc16() uint16 {
	return crc16.Calc(c.Data)
}
