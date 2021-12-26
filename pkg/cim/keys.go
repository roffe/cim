package cim

import (
	"bytes"
	"fmt"
)

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
