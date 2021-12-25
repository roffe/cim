package cim

import (
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

func (bin *Bin) validateVin() error {
	v := bin.Vin
	var b []byte
	b = append(b, []byte(v.Data)...)
	b = append(b, v.Unknown...)
	b = append(b, v.SpsCount)

	c := crc16.Calc(b)
	if v.Checksum != c {
		return fmt.Errorf("vin cheksum validation failed")
	}
	return nil
}
