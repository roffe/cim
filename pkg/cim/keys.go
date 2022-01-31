package cim

import (
	"bytes"
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

type Keys struct {
	IskHI1    []byte   `bin:"len:4" json:"isk_hi1"`
	IskLO1    []byte   `bin:"len:2" json:"isk_lo1"`
	Data1     [][]byte `bin:"len:5,[len:4]" json:"data1"`
	Count1    uint8    `bin:"len:1" json:"count1"`
	Constant1 []byte   `bin:"len:7" json:"constant1"`
	Errors1   uint8    `bin:"len:1" json:"errors1"`
	Checksum1 uint16   `bin:"le,len:2" json:"checksum1"`
	IskHI2    []byte   `bin:"len:4" json:"isk_hi2"`
	IskLO2    []byte   `bin:"len:2" json:"isk_lo2"`
	Data2     [][]byte `bin:"len:5,[len:4]" json:"data2"`
	Count2    uint8    `bin:"len:1" json:"count2"`
	Constant2 []byte   `bin:"len:7" json:"constant2"`
	Errors2   uint8    `bin:"len:1" json:"errors2"`
	Checksum2 uint16   `bin:"le,len:2" json:"checksum2"` // CRC16 MCRF4XX
} // 74 bytes

// Set key count
func (k *Keys) Count(no uint8) error {
	if no > 5 {
		return fmt.Errorf("max 5 keys")
	}
	k.Count1 = no
	k.Count2 = no
	k.updateChecksum()
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
	k.Data1[keyno], k.Data2[keyno] = value, value
	k.updateChecksum()
	return nil
}

func (k *Keys) SetErrorCount(value uint8) error {
	k.Errors1, k.Errors2 = value, value
	k.updateChecksum()
	return nil
}

func (k *Keys) SetKeyCount(keys uint8) error {
	if keys > 5 {
		return fmt.Errorf("maximum number of keys is 5")
	}
	k.Count1, k.Count2 = keys, keys
	return nil
}

func (k *Keys) SetIsk(high, low []byte) error {
	if len(high) != 4 {
		return fmt.Errorf("invalid isk high value length")
	}
	if len(low) != 2 {
		return fmt.Errorf("invalid isk low value length")
	}
	k.IskHI1, k.IskHI2 = high, high
	k.IskLO1, k.IskLO2 = low, low
	k.updateChecksum()
	return nil
}

func (k *Keys) validate() error {
	if k.Checksum1 != k.Checksum2 {
		return fmt.Errorf("key 0 checksums missmatch in bin")
	}

	if k.Count1 != k.Count2 {
		return fmt.Errorf("key count missmatch in data banks %d | %d", k.Count1, k.Count2)
	}

	if k.Errors1 != k.Errors2 {
		return fmt.Errorf("key errors missmatch in data banks %d | %d", k.Errors1, k.Errors2)
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

	c1, c2 := k.Crc16()

	if c1 != k.Checksum1 {
		return fmt.Errorf("keys calculated bank 1 checksum %X does not match stored %X", c1, k.Checksum1)
	}

	if c2 != k.Checksum1 {
		return fmt.Errorf("keys calculated bank 1 checksum %X does not match stored %X", c1, k.Checksum1)
	}

	return nil
}

func (k *Keys) Crc16() (uint16, uint16) {
	d1 := bytes.NewBuffer([]byte{})
	d1.Write(k.IskHI1)
	d1.Write(k.IskLO1)
	for _, k := range k.Data1 {
		d1.Write(k)
	}
	d1.WriteByte(k.Count1)
	d1.Write(k.Constant1)
	d1.WriteByte(k.Errors1)

	d2 := bytes.NewBuffer([]byte{})
	d2.Write(k.IskHI2)
	d2.Write(k.IskLO2)
	for _, k := range k.Data2 {
		d2.Write(k)
	}
	d2.WriteByte(k.Count2)
	d2.Write(k.Constant2)
	d2.WriteByte(k.Errors2)

	return crc16.Calc(d1.Bytes()), crc16.Calc(d2.Bytes())
}

func (k *Keys) updateChecksum() {
	k.Checksum1, k.Checksum2 = k.Crc16()
}
