package cim

import (
	"bytes"
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type Keys struct {
	IskHI1     []byte   `bin:"len:4"`
	IskLO1     []byte   `bin:"len:2"`
	Data1      [][]byte `bin:"len:5,[len:4]"`
	Count1     uint8    `bin:"len:1"`
	Constant1  []byte   `bin:"len:7"`
	KeyErrors1 uint8    `bin:"len:1"`
	Checksum1  uint16   `bin:"le,len:2"`
	IskHI2     []byte   `bin:"len:4"`
	IskLO2     []byte   `bin:"len:2"`
	Data2      [][]byte `bin:"len:5,[len:4]"`
	Count2     uint8    `bin:"len:1"`
	Constant2  []byte   `bin:"len:7"`
	Errors2    uint8    `bin:"len:1"`
	Checksum2  uint16   `bin:"le,len:2"` // CRC16 MCRF4XX
} // 74 bytes

// Set key count
func (k *Keys) Count(no uint8) error {
	if no > 5 {
		return fmt.Errorf("max 5 keys")
	}
	k.Count1 = no
	k.Count2 = no
	k.updateKeysChecksum()
	return nil
}

// SetKey set key value,0 is the first key
func (k *Keys) SetKey(keyno uint8, value []byte) error {
	if keyno > 5 {
		return fmt.Errorf("invalid key position")
	}
	if len(value) != 4 {
		return fmt.Errorf("invalid key size")
	}
	k.Data1[keyno] = value
	k.Data2[keyno] = value
	k.updateKeysChecksum()
	return nil
}

func (k *Keys) SetIskHi(value []byte) {
	k.IskHI1 = value
	k.IskHI2 = value
	k.updateKeysChecksum()
}

func (k *Keys) SetIskLo(value []byte) {
	k.IskLO1 = value
	k.IskLO2 = value
	k.updateKeysChecksum()
}

func (k *Keys) validate() error {

	if k.Checksum1 != k.Checksum2 {
		return fmt.Errorf("key 0 checksums missmatch in bin")
	}

	if k.Count1 != k.Count2 {
		return fmt.Errorf("key count missmatch in data banks %d | %d", k.Count1, k.Count2)
	}

	if k.KeyErrors1 != k.Errors2 {
		return fmt.Errorf("key errors missmatch in data banks %d | %d", k.KeyErrors1, k.Errors2)
	}

	if !bytes.Equal(k.Data1[0], k.Data2[0]) {
		return fmt.Errorf("key 1, bank 1 and 2 does not match, corrupt memory?")
	}

	if !bytes.Equal(k.Data1[1], k.Data2[1]) {
		return fmt.Errorf("key 2, bank 1 and 2 does not match, corrupt memory?")
	}

	if !bytes.Equal(k.Data1[2], k.Data2[2]) {
		return fmt.Errorf("key 3, bank 1 and 2 does not match, corrupt memory?")
	}

	if !bytes.Equal(k.Data1[3], k.Data2[3]) {
		return fmt.Errorf("key 4, bank 1 and 2 does not match, corrupt memory?")
	}

	if !bytes.Equal(k.Data1[4], k.Data2[4]) {
		return fmt.Errorf("key 5, bank 1 and 2 does not match, corrupt memory?")
	}

	c1, c2 := k.Checksum()

	if c1 != k.Checksum1 {
		return fmt.Errorf("calculated bank 1 checksum %X does not match stored %X", c1, k.Checksum1)
	}

	if c2 != k.Checksum1 {
		return fmt.Errorf("calculated bank 1 checksum %X does not match stored %X", c1, k.Checksum1)
	}

	return nil
}

func (k *Keys) Checksum() (uint16, uint16) {
	d1 := bytes.NewBuffer([]byte{})
	d1.Write(k.IskHI1)
	d1.Write(k.IskLO1)
	for _, k := range k.Data1 {
		d1.Write(k)
	}
	d1.WriteByte(k.Count1)
	d1.Write(k.Constant1)
	d1.WriteByte(k.KeyErrors1)
	k1crc := crc16.Calc(d1.Bytes())

	d2 := bytes.NewBuffer([]byte{})
	d2.Write(k.IskHI2)
	d2.Write(k.IskLO2)
	for _, k := range k.Data2 {
		d2.Write(k)
	}
	d2.WriteByte(k.Count2)
	d2.Write(k.Constant2)
	d2.WriteByte(k.Errors2)
	k2crc := crc16.Calc(d1.Bytes())
	return k1crc, k2crc
}

func (k *Keys) updateKeysChecksum() {
	k.Checksum1, k.Checksum2 = k.Checksum()
}
