package cim

type UnknownData3 struct {
	Data1         []byte `bin:"len:42"`
	Data1Checksum uint16 `bin:"le,len:2"`
	Data2         []byte `bin:"len:42"`
	Data2Checksum uint16 `bin:"le,len:2"`
}

type UnknownData4 struct {
	Data     []byte `bin:"len:2"`
	Checksum uint16 `bin:"le,len:2"`
}

type UnknownData1 struct {
	Data1         []byte `bin:"len:20"`
	Data1Checksum uint16 `bin:"le,len:2"`
	Data2         []byte `bin:"len:20"`
	Data2Checksum uint16 `bin:"le,len:2"`
}

type UnknownData5 struct {
	Data     []byte `bin:"len:23"`
	Checksum uint16 `bin:"le,len:2"`
}

type UnknownData6 struct {
	Data1         []byte `bin:"len:20"`
	Data1Checksum uint16 `bin:"le,len:2"`
	Data2         []byte `bin:"len:20"`
	Data2Checksum uint16 `bin:"le,len:2"`
}

type UnknownData7 struct {
	Data1         []byte `bin:"len:5"`
	Data1Checksum uint16 `bin:"le,len:2"`
	Data2         []byte `bin:"len:5"`
	Data2Checksum uint16 `bin:"le,len:2"`
}

type UnknownData8 struct {
	Data     []byte `bin:"len:6"`
	Checksum uint16 `bin:"le,len:2"`
}

type UnknownData9 struct {
	Data     []byte `bin:"len:5"`
	Checksum uint16 `bin:"le,len:2"`
}

type UnknownData2 struct {
	Data1         []byte `bin:"len:5"`
	Data1Checksum uint16 `bin:"le,len:2"`
	Data2         []byte `bin:"len:5"`
	Data2Checksum uint16 `bin:"le,len:2"`
}

type UnknownData10 struct {
	Data1         []byte `bin:"len:4"`
	Data1Checksum uint16 `bin:"le,len:2"`
	Data2         []byte `bin:"len:4"`
	Data2Checksum uint16 `bin:"le,len:2"`
}
