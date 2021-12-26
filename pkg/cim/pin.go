package cim

import (
	"fmt"

	"github.com/roffe/cim/pkg/crc16"
)

func (bin *Bin) validatePin() error {
	p := bin.Pin
	p1ccrc := crc16.Calc(append(p.Pin1[:], p.Pin1Unknown[:]...))
	p2ccrc := crc16.Calc(append(p.Pin2[:], p.Pin2Unknown[:]...))

	if p.Pin1Checksum != p1ccrc {
		return fmt.Errorf("pin 1 checksum in binary does not match calculated checksum")
	}

	if p.Pin2Checksum != p2ccrc {
		return fmt.Errorf("pin 2 checksum in binary does not match calculated checksum")
	}

	if p1ccrc != p2ccrc {
		return fmt.Errorf("calculated checksums does not match")
	}

	if p.Pin1Checksum != p.Pin2Checksum {
		return fmt.Errorf("stored pin checksums in binary does not match")
	}

	return nil
}
