package cim

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ghostiam/binstruct"
)

type Keys struct {
	IskHI1         []byte   `bin:"len:4"`
	IskLO1         []byte   `bin:"len:2"`
	Keys1          [][]byte `bin:"len:5,[len:4]"`
	KeysKeysCount1 uint8    `bin:"len:1"`
	KeysConstant1  []byte   `bin:"len:7"`
	KeyErrors1     uint8    `bin:"len:1"`
	Checksum1      uint16   `bin:"Uint16l,len:2"`
	IskHI2         []byte   `bin:"len:4"`
	IskLO2         []byte   `bin:"len:2"`
	Keys2          [][]byte `bin:"len:5,[len:4]"`
	KeysKeysCount2 uint8    `bin:"len:1"`
	KeysConstant2  []byte   `bin:"len:7"`
	KeyErrors2     uint8    `bin:"len:1"`
	Checksum2      uint16   `bin:"Uint16l,len:2"` // CRC16 MCRF4XX
}

// Uint16l returns a uint16 read as little endian
func (*Keys) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}
func (bin *Bin) validateKeys() error {
	k := bin.Keys
	if k.Checksum1 != k.Checksum2 {
		return fmt.Errorf("key 0 checksums missmatch in bin")
	}

	if k.KeyErrors1 != k.KeyErrors2 {
		return fmt.Errorf("key errors missmatch in data banks %d | %d", k.KeyErrors1, k.KeyErrors2)
	}

	if !bytes.Equal(k.Keys1[0], k.Keys2[0]) {
		return fmt.Errorf("key 1 data bank 1 and 2 does not match, corrupt memory?")
	}

	if !bytes.Equal(k.Keys1[1], k.Keys2[1]) {
		return fmt.Errorf("key 2 data bank 1 and 2 does not match, corrupt memory?")
	}

	if !bytes.Equal(k.Keys1[2], k.Keys2[2]) {
		return fmt.Errorf("key 3 data bank 1 and 2 does not match, corrupt memory?")
	}

	if !bytes.Equal(k.Keys1[3], k.Keys2[3]) {
		return fmt.Errorf("key 4 data bank 1 and 2 does not match, corrupt memory?")
	}

	if !bytes.Equal(k.Keys1[4], k.Keys2[4]) {
		return fmt.Errorf("key 5 data bank 1 and 2 does not match, corrupt memory?")
	}

	return nil
}
