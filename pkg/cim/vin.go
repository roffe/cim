package cim

type Vin struct {
	Data     string `bin:"len:17"`
	Unknown  []byte `bin:"len:10"`
	SpsCount uint8
	Checksum uint16 `bin:"len:2"`
}
