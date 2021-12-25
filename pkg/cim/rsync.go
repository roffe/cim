package cim

type RSync struct {
	Data     [][]byte `bin:"len:5,[len:4]"`
	Checksum [2]byte
}
