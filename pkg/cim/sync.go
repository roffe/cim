package cim

import (
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type Sync struct {
	Data     [][]byte `bin:"len:5,[len:4]" json:"data"`
	Checksum uint16   `bin:"le,len:2" json:"checksum"`
} // 22 bytes

func (s *Sync) SetData(no uint8, data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("Sync data %d invalid length %d, should be 4 bytes", no, len(data))
	}
	s.Data[no] = data
	s.updateChecksum()
	return nil
}

func (s *Sync) validate() error {
	if s.Checksum != s.Crc16() {
		return fmt.Errorf("sync data checksum: %X does not match calculated checksum: %X", s.Checksum, s.Crc16())
	}
	return nil
}

func (s *Sync) Crc16() uint16 {
	var data []byte
	for _, b := range s.Data {
		data = append(data, b...)
	}
	return crc16.Calc(data)
}

func (s *Sync) updateChecksum() {
	s.Checksum = s.Crc16()
}
