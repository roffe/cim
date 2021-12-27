package cim

type Sync struct {
	Data     [][]byte `bin:"len:5,[len:4]"`
	Checksum uint16   `bin:"le,len:2"`
} // 22 bytes

func (s *Sync) SetData(no uint8, data []byte) {
	s.Data[no] = data
}

//func (s *Sync)
