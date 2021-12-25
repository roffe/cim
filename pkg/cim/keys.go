package cim

import (
	"fmt"
)

func (bin *Bin) validateKeys() error {
	//fmt.Printf("Bank 1: %x %x %x %x %x %x %x %x %x %x %x\n", k.IskHI1, k.IskLO1, k.Key1, k.Key2, k.Key3, k.Key4, k.Key5, k.KeysKeysCount1, k.KeysUnknown1, k.KeyErrors1, k.Checksum1)
	//fmt.Printf("Bank 2: %x %x %x %x %x %x %x %x %x %x %x\n", k.IskHI2, k.IskLO2, k.Key1_2, k.Key2_2, k.Key3_2, k.Key4_2, k.Key5_2, k.KeysKeysCount2, k.KeysUnknown2, k.KeyErrors2, k.Checksum2)
	k := bin.Keys
	if k.Checksum1 != k.Checksum2 {
		return fmt.Errorf("key checksums missmatch in bin")
	}

	//if !bytes.Equal(k.Keys1[:], k.Keys2[:]) {
	//	return fmt.Errorf("key data bank 1 and 2 does not match, corrupt memory?")
	//}

	return nil
}
