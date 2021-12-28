package cim

import (
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type Sync struct {
	Data     [][]byte `bin:"len:5,[len:4]" json:"data"`
	Checksum uint16   `bin:"le,len:2" json:"checksum"`
} // 22 bytes

func (s *Sync) SetData(no uint8, data []byte) {
	s.Data[no] = data
}

func (s *Sync) validate() error {
	if s.Checksum != s.Crc16() {
		if s.Checksum != s.Crc16() {
			return fmt.Errorf("sync data checksum does not match calculated checksum")
		}
	}
	return nil
}

func (p *Sync) Crc16() uint16 {
	var data []byte
	for _, b := range p.Data {
		data = append(data, b...)
	}
	return crc16.Calc(data)
}
