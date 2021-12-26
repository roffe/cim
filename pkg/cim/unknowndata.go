package cim

import (
	"encoding/binary"

	"github.com/ghostiam/binstruct"
)

type UnknownData3 struct {
	Data1         []byte `bin:"len:42"`
	Data1Checksum uint16 `bin:"Uint16l,len:2"`
	Data2         []byte `bin:"len:42"`
	Data2Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData3) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData4 struct {
	Data     []byte `bin:"len:2"`
	Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData4) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData1 struct {
	Data1         []byte `bin:"len:20"`
	Data1Checksum uint16 `bin:"Uint16l,len:2"`
	Data2         []byte `bin:"len:20"`
	Data2Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData1) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData5 struct {
	Data     []byte `bin:"len:23"`
	Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData5) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData6 struct {
	Data1         []byte `bin:"len:20"`
	Data1Checksum uint16 `bin:"Uint16l,len:2"`
	Data2         []byte `bin:"len:20"`
	Data2Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData6) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData7 struct {
	Data1         []byte `bin:"len:5"`
	Data1Checksum uint16 `bin:"Uint16l,len:2"`
	Data2         []byte `bin:"len:5"`
	Data2Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData7) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData8 struct {
	Data     []byte `bin:"len:6"`
	Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData8) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData9 struct {
	Data     []byte `bin:"len:5"`
	Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData9) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData2 struct {
	Data1         []byte `bin:"len:5"`
	Data1Checksum uint16 `bin:"Uint16l,len:2"`
	Data2         []byte `bin:"len:5"`
	Data2Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData2) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

type UnknownData10 struct {
	Data1         []byte `bin:"len:4"`
	Data1Checksum uint16 `bin:"Uint16l,len:2"`
	Data2         []byte `bin:"len:4"`
	Data2Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*UnknownData10) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}
