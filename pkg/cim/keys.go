package cim

import (
	"bytes"
	"fmt"
)

type Keys struct {
	IskHI1         []byte   `bin:"len:4"`
	IskLO1         []byte   `bin:"len:2"`
	Keys1          [][]byte `bin:"len:5,[len:4]"`
	KeysKeysCount1 uint8    `bin:"len:1"`
	KeysUnknown1   []byte   `bin:"len:7"`
	KeyErrors1     uint8    `bin:"len:1"`
	Checksum1      []byte   `bin:"len:2"`
	IskHI2         []byte   `bin:"len:4"`
	IskLO2         []byte   `bin:"len:2"`
	Keys2          [][]byte `bin:"len:5,[len:4]"`
	KeysKeysCount2 uint8    `bin:"len:1"`
	KeysUnknown2   []byte   `bin:"len:7"`
	KeyErrors2     uint8    `bin:"len:1"`
	Checksum2      []byte   `bin:"len:2"`
}

func (k *Keys) validate() error {
	//fmt.Printf("Bank 1: %x %x %x %x %x %x %x %x %x %x %x\n", k.IskHI1, k.IskLO1, k.Key1, k.Key2, k.Key3, k.Key4, k.Key5, k.KeysKeysCount1, k.KeysUnknown1, k.KeyErrors1, k.Checksum1)
	//fmt.Printf("Bank 2: %x %x %x %x %x %x %x %x %x %x %x\n", k.IskHI2, k.IskLO2, k.Key1_2, k.Key2_2, k.Key3_2, k.Key4_2, k.Key5_2, k.KeysKeysCount2, k.KeysUnknown2, k.KeyErrors2, k.Checksum2)

	if !bytes.Equal(k.Checksum1[:], k.Checksum2[:]) {
		return fmt.Errorf("key checksums missmatch in bin")
	}

	//if !bytes.Equal(k.Keys1[:], k.Keys2[:]) {
	//	return fmt.Errorf("key data bank 1 and 2 does not match, corrupt memory?")
	//}

	return nil
}
