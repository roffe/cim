package cim

type PSK struct {
	Low      []byte `bin:"len:4"`
	High     []byte `bin:"len:4"`
	Ide      []byte `bin:"len:4"`
	Checksum uint16 `bin:"len:2"`
}
