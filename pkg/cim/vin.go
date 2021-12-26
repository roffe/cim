package cim

import (
	"bytes"
	"fmt"
	"log"

	"github.com/roffe/cim/pkg/crc16"
)

func (bin *Bin) validateVin() error {
	v := bin.Vin
	b := bytes.NewBuffer(nil)
	b.Write([]byte(v.Data))
	b.WriteByte(v.Value)
	b.Write(v.Unknown)
	b.WriteByte(v.SpsCount)
	//b = append(b, []byte(v.Data)...)
	//b = append(b, byte(v.Value))
	//b = append(b, v.Unknown...)
	//b = append(b, v.SpsCount)
	c := crc16.Calc(b.Bytes())
	if v.Checksum != c {
		log.Printf("%d || %d", v.Checksum, c)
		return fmt.Errorf("vin cheksum validation failed")

	}
	return nil
}
