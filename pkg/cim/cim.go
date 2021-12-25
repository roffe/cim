package cim

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"time"

	"github.com/albenik/bcd"
	"github.com/ghostiam/binstruct"
	"github.com/jedib0t/go-pretty/table"
)

var tableTheme = table.StyleColoredDark

const isoDate = "2006-01-02"

func Load(filename string) (*Bin, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fw := Bin{
		filename: filename,
		md5:      fmt.Sprintf("%x", md5.Sum(b)),
		crc32:    fmt.Sprintf("%08x", crc32.Checksum(b, crc32.MakeTable(crc32.IEEE))),
	}

	for i, bb := range b {
		b[i] = bb ^ 0xFF //xor byte with 0xFF
	}
	if err := binstruct.UnmarshalBE(b, &fw); err != nil {
		return nil, err
	}

	return &fw, nil
}

type Bin struct {
	filename             string    `bin:"-"`
	md5                  string    `bin:"-"`
	crc32                string    `bin:"-"`
	MagicByte            uint8     `bin:"len:1"`         // 0x20
	ProgrammingDate      time.Time `bin:"BCDDate,len:3"` // BCD Binary-Coded Decimal yy-mm-dd
	SasOption            bool      `bin:"Sasopt,len:1"`  // Steering Angle Sensor 0x03 = true
	UnknownBytes1        []byte    `bin:"len:6"`
	PnSAAB1              uint32    `bin:"len:4"`
	PnSAAB1_2            string    `bin:"len:2"`
	ConfigurationVersion uint32    `bin:"len:4"`
	PnBase1              uint32    `bin:"len:4"`
	PnBase1_2            string    `bin:"len:2"`
	Vin                  struct {
		Data     string `bin:"len:17"` // Vin as ASCII
		Unknown  []byte `bin:"len:10"`
		SpsCount uint8
		Checksum uint16 `bin:"len:2"` // CRC16 MCRF4XX
	} // 30 bytes
	ProgrammingID []string `bin:"len:3,[len:10]"` // 3 last sps progrmming ids's groups of 10 characters each. 30 bytes
	UnknownData3  []byte   `bin:"len:88"`
	Pin           Pin      // 20 bytes
	UnknownData4  []byte   `bin:"len:4"`
	UnknownData1  []byte   `bin:"len:44"`
	Const1        []byte   `bin:"len:10"`
	Keys          struct {
		IskHI1         []byte   `bin:"len:4"`
		IskLO1         []byte   `bin:"len:2"`
		Keys1          [][]byte `bin:"len:5,[len:4]"`
		KeysKeysCount1 uint8    `bin:"len:1"`
		KeysUnknown1   []byte   `bin:"len:7"`
		KeyErrors1     uint8    `bin:"len:1"`
		Checksum1      uint16   `bin:"len:2"`
		IskHI2         []byte   `bin:"len:4"`
		IskLO2         []byte   `bin:"len:2"`
		Keys2          [][]byte `bin:"len:5,[len:4]"`
		KeysKeysCount2 uint8    `bin:"len:1"`
		KeysUnknown2   []byte   `bin:"len:7"`
		KeyErrors2     uint8    `bin:"len:1"`
		Checksum2      uint16   `bin:"len:2"` // CRC16 MCRF4XX
	} // 74 bytes
	UnknownData5 []byte `bin:"len:25"`
	Sync         struct {
		Data     [][]byte `bin:"len:5,[len:4]"`
		Checksum uint16   `bin:"len:2"`
	} // 22 bytes
	UnknownData6           []byte    `bin:"len:44"`
	UnknownData7           []byte    `bin:"len:14"`
	UnknownData8           []byte    `bin:"len:8"`
	UnknownData9           []byte    `bin:"len:7"`
	UnknownData2           []byte    `bin:"len:14"`
	SnSticker              uint64    `bin:"ReadSN,len:5"`   // BCD
	ProgrammingFactoryDate time.Time `bin:"BCDDateR,len:3"` // Reversed BCD date dd-mm-yy
	UnknownBytes2          []byte    `bin:"len:3"`
	PnDelphi               uint32    `bin:"Uint32l,len:4"` // Little endian
	UnknownBytes3          []byte    `bin:"len:2"`
	PnSAAB2                uint32    `bin:"Uint32l,len:4"` // Little endian
	UnknownBytes4          []byte    `bin:"len:3"`
	PSK                    struct {
		Low      []byte `bin:"len:4"`
		High     []byte `bin:"len:4"`
		Ide      []byte `bin:"len:4"`
		Checksum uint16 `bin:"len:2"` // CRC16 MCRF4XX
	} // 14 bytes
	UnknownData10 []byte `bin:"len:12"`
	EOF           byte   `bin:"len:1"` //0x00
}

type Pin struct {
	Pin1         []byte `bin:"len:4"` // LE
	Pin1Unknown  []byte `bin:"len:4"`
	Pin1Checksum uint16 `bin:"Uint16l,len:2"`
	Pin2         []byte `bin:"len:4"` // LE
	Pin2Unknown  []byte `bin:"len:4"`
	Pin2Checksum uint16 `bin:"Uint16l,len:2"`
}

// Uint16l returns a uint16 read as little endian
func (*Pin) Uint16l(r binstruct.Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

// Uint32l returns a uint32 read as little endian
func (*Bin) Uint32l(r binstruct.Reader) (uint32, error) {
	var out uint32
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

func (*Bin) Sasopt(r binstruct.Reader) (bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	if b == 0x03 {
		return true, nil
	}
	return false, nil
}

func (*Bin) BCDDate(r binstruct.Reader) (time.Time, error) {
	_, b, err := r.ReadBytes(3)
	if err != nil {
		return time.Time{}, err
	}
	t := fmt.Sprintf("%02X-%02X-%02X", b[0], b[1:2], b[2:3])
	out, err := time.Parse("06-01-02", t) // yy-mm-dd
	if err != nil {
		return time.Time{}, err
	}
	return out, nil
}

func (*Bin) BCDDateR(r binstruct.Reader) (time.Time, error) {
	_, b, err := r.ReadBytes(3)
	if err != nil {
		return time.Time{}, err
	}
	t := fmt.Sprintf("%02X-%02X-%02X", b[0], b[1:2], b[2:3])
	out, err := time.Parse("02-01-06", t) // dd-mm-yy
	if err != nil {
		return time.Time{}, err
	}
	return out, nil
}

func (bin *Bin) Filename() string {
	return bin.filename
}

func (bin *Bin) MD5() string {
	return bin.md5
}

func (bin *Bin) CRC32() string {
	return bin.crc32
}

// Return Serial sticker as uint64, stored as 5byte Binary-Coded Decimal
func (*Bin) ReadSN(r binstruct.Reader) (uint64, error) {
	_, b, err := r.ReadBytes(5)
	if err != nil {
		return 0, nil
	}
	return bcd.ToUint64(b), nil
}

func (bin *Bin) Validate() error {
	tests := []func() error{
		bin.validateVin,
		bin.validatePin,
		bin.validateKeys,
	}

	for _, t := range tests {
		if err := t(); err != nil {
			return err
		}
	}
	return nil
}
