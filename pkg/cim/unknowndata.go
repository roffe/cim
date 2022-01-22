package cim

import (
	"bytes"
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type UnknownData1 struct {
	Data1     []byte `bin:"len:20"`
	Checksum1 uint16 `bin:"le,len:2"`
	Data2     []byte `bin:"len:20"`
	Checksum2 uint16 `bin:"le,len:2"`
}

const (
	errChecksum = "checksum does not match calculated checksum"
)

func (u *UnknownData1) validate() error {
	if !bytes.Equal(u.Data1, u.Data2) {
		return fmt.Errorf("UnknownData1 data1 & data2 does not match, corrupt memory?")
	}
	c1, c2 := u.Crc16()
	if u.Checksum1 != c1 {
		return fmt.Errorf("UnknownData1 data1 " + errChecksum)
	}
	if u.Checksum2 != c2 {
		return fmt.Errorf("UnknownData1 data2 " + errChecksum)
	}
	return nil
}

func (u *UnknownData1) Crc16() (uint16, uint16) {
	return crc16.Calc(u.Data1), crc16.Calc(u.Data2)
}

type UnknownData2 struct {
	Data1     []byte `bin:"len:5"`
	Checksum1 uint16 `bin:"le,len:2"`
	Data2     []byte `bin:"len:5"`
	Checksum2 uint16 `bin:"le,len:2"`
}

func (u *UnknownData2) validate() error {
	if !bytes.Equal(u.Data1, u.Data2) {
		return fmt.Errorf("UnknownData2 data1 & data2 does not match, corrupt memory?")
	}
	c1, c2 := u.Crc16()
	if u.Checksum1 != c1 {
		return fmt.Errorf("UnknownData2 data1 " + errChecksum)
	}
	if u.Checksum2 != c2 {
		return fmt.Errorf("UnknownData2 data2 " + errChecksum)
	}
	return nil
}

func (u *UnknownData2) Crc16() (uint16, uint16) {
	return crc16.Calc(u.Data1), crc16.Calc(u.Data2)
}

type UnknownData3 struct {
	Data1     []byte `bin:"len:42" json:"data1"`
	Checksum1 uint16 `bin:"le,len:2" json:"checksum1" `
	Data2     []byte `bin:"len:42" json:"data2"`
	Checksum2 uint16 `bin:"le,len:2" json:"checksum2"`
}

func (u *UnknownData3) validate() error {
	if !bytes.Equal(u.Data1, u.Data2) {
		return fmt.Errorf("UnknownData3 data1 & data2 does not match, corrupt memory?")
	}
	c1, c2 := u.Crc16()
	if u.Checksum1 != c1 {
		return fmt.Errorf("UnknownData3 data1 " + errChecksum)
	}
	if u.Checksum2 != c2 {
		return fmt.Errorf("UnknownData3 data2 " + errChecksum)
	}
	return nil
}

func (u *UnknownData3) Crc16() (uint16, uint16) {
	return crc16.Calc(u.Data1), crc16.Calc(u.Data2)
}

type UnknownData4 struct {
	Data     []byte `bin:"len:2" json:"data"`
	Checksum uint16 `bin:"le,len:2" json:"checksum"`
}

func (u *UnknownData4) validate() error {
	if u.Checksum != u.Crc16() {
		return fmt.Errorf("UnknownData4 data " + errChecksum)
	}
	return nil
}

func (u *UnknownData4) Crc16() uint16 {
	return crc16.Calc(u.Data)
}

type UnknownData5 struct {
	Data     []byte `bin:"len:23" json:"data"`
	Checksum uint16 `bin:"le,len:2" json:"checksum"`
}

func (u *UnknownData5) validate() error {
	if u.Checksum != u.Crc16() {
		return fmt.Errorf("UnknownData5 calculated:%04d, actuall: %04d data %v", u.Crc16(), u.Checksum, errChecksum)
	}
	return nil
}

func (u *UnknownData5) Crc16() uint16 {
	return crc16.Calc(u.Data)
}

type UnknownData6 struct {
	Data1     []byte `bin:"len:20" json:"data1"`
	Checksum1 uint16 `bin:"le,len:2" json:"checksum1"`
	Data2     []byte `bin:"len:20" json:"data2"`
	Checksum2 uint16 `bin:"le,len:2" json:"checksum2"`
}

func (u *UnknownData6) validate() error {
	c1, c2 := u.Crc16()
	if !bytes.Equal(u.Data1, u.Data2) {
		return fmt.Errorf("UnknownData6 data1 & data2 does not match, corrupt memory?")
	}
	if u.Checksum1 != c1 {
		return fmt.Errorf("UnknownData6 data1 " + errChecksum)
	}
	if u.Checksum2 != c2 {
		return fmt.Errorf("UnknownData6 data2 " + errChecksum)
	}
	return nil
}

func (u *UnknownData6) Crc16() (uint16, uint16) {
	return crc16.Calc(u.Data1), crc16.Calc(u.Data2)
}

type UnknownData7 struct {
	Data1     []byte `bin:"len:5" json:"data1"`
	Checksum1 uint16 `bin:"le,len:2" json:"checksum1"`
	Data2     []byte `bin:"len:5" json:"data2"`
	Checksum2 uint16 `bin:"le,len:2" json:"checksum2"`
}

func (u *UnknownData7) validate() error {
	c1, c2 := u.Crc16()
	if !bytes.Equal(u.Data1, u.Data2) {
		return fmt.Errorf("UnknownData7 data1 & data2 does not match, corrupt memory?")
	}
	if u.Checksum1 != c1 {
		return fmt.Errorf("UnknownData7 data1 " + errChecksum)
	}
	if u.Checksum2 != c2 {
		return fmt.Errorf("UnknownData7 data2 " + errChecksum)
	}
	return nil
}

func (u *UnknownData7) Crc16() (uint16, uint16) {
	return crc16.Calc(u.Data1), crc16.Calc(u.Data2)
}

type UnknownData8 struct {
	Data     []byte `bin:"len:6" json:"data"`
	Checksum uint16 `bin:"le,len:2" json:"checksum"`
}

func (u *UnknownData8) validate() error {
	if u.Checksum != u.Crc16() {
		return fmt.Errorf("UnknownData5 data " + errChecksum)
	}
	return nil
}

func (u *UnknownData8) Crc16() uint16 {
	return crc16.Calc(u.Data)
}

type UnknownData9 struct {
	Data     []byte `bin:"len:5" json:"data"`
	Checksum uint16 `bin:"le,len:2" json:"checksum"`
}

func (u *UnknownData9) validate() error {
	if u.Checksum != u.Crc16() {
		return fmt.Errorf("UnknownData5 data " + errChecksum)
	}
	return nil
}

func (u *UnknownData9) Crc16() uint16 {
	return crc16.Calc(u.Data)
}

type UnknownData10 struct {
	Data1     []byte `bin:"len:4" json:"data1"`
	Checksum1 uint16 `bin:"le,len:2" json:"checksum1"`
	Data2     []byte `bin:"len:4" json:"data2"`
	Checksum2 uint16 `bin:"le,len:2" json:"checksum2"`
}

func (u *UnknownData10) validate() error {
	if bytes.Equal(u.Data1, []byte{0x00, 0x00, 0x00, 0x00}) && bytes.Equal(u.Data2, []byte{0x00, 0x00, 0x00, 0x00}) {
		return nil
	}

	if !bytes.Equal(u.Data1, u.Data2) {
		return fmt.Errorf("UnknownData10 data1 & data2 does not match, corrupt memory?")
	}

	c1, c2 := u.Crc16()
	if u.Checksum1 != c1 {
		return fmt.Errorf("UnknownData10 data1 "+errChecksum+" %X, %X, %X", u.Data1, u.Checksum1, c1)
	}
	if u.Checksum2 != c2 {
		return fmt.Errorf("UnknownData10 data2 "+errChecksum+" %X, %X, %X", u.Data1, u.Checksum2, c2)
	}
	return nil
}

func (u *UnknownData10) Crc16() (uint16, uint16) {
	return crc16.Calc(u.Data1), crc16.Calc(u.Data2)
}
