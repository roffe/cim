package cim

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"time"

	"github.com/albenik/bcd"
	"github.com/ghostiam/binstruct"
	"github.com/jedib0t/go-pretty/table"
)

var tableTheme = table.StyleColoredDark

const IsoDate = "2006-01-02"

// Load a file from disk
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

	//xor bytes if magicByte is not 0x20
	if b[0] != 0x20 {
		for i, bb := range b {
			b[i] = bb ^ 0xFF
		}
	}

	if err := binstruct.UnmarshalBE(b, &fw); err != nil {
		return nil, err
	}

	return &fw, nil
}

// Load a byte slice as a named binary
func LoadBytes(filename string, b []byte) (*Bin, error) {
	fw := Bin{
		filename: filename,
		md5:      fmt.Sprintf("%x", md5.Sum(b)),
		crc32:    fmt.Sprintf("%08x", crc32.Checksum(b, crc32.MakeTable(crc32.IEEE))),
	}

	//xor bytes if magicByte is not 0x20
	if b[0] != 0x20 {
		for i, bb := range b {
			b[i] = bb ^ 0xFF
		}
	}
	// Unpack bytes into struct
	if err := binstruct.UnmarshalBE(b, &fw); err != nil {
		return nil, err
	}
	return &fw, nil
}

type Bin struct {
	filename               string        `bin:"-"`
	md5                    string        `bin:"-"`
	crc32                  string        `bin:"-"`
	MagicByte              byte          `bin:"len:1"`         // 0x20
	ProgrammingDate        time.Time     `bin:"BCDDate,len:3"` // BCD Binary-Coded Decimal yy-mm-dd
	SasOption              uint8         `bin:"len:1"`         // Steering Angle Sensor 0x03 = true
	UnknownBytes1          []byte        `bin:"len:6"`
	PartNo1                uint32        `bin:"len:4"` // End model (HW+SW)
	PartNo1Suffix          string        `bin:"len:2"`
	ConfigurationVersion   uint32        `bin:"len:4"`
	PnBase1                uint32        `bin:"len:4"` // Base model (HW+boot)
	PnBase1Suffix          string        `bin:"len:2"`
	Vin                    Vin           `bin:"len:30"`         //30 bytes
	ProgrammingID          []string      `bin:"len:3,[len:10]"` // 3 last sps progrmming ids's groups of 10 characters each. 30 bytes
	UnknownData3           UnknownData3  `bin:"len:88"`
	Pin                    Pin           `bin:"len:20"` // 20 bytes
	UnknownData4           UnknownData4  `bin:"len:4"`
	UnknownData1           UnknownData1  `bin:"len:44"`
	Const1                 Const1        `bin:"len:10"` // 10 bytes
	Keys                   Keys          `bin:"len:74"` // 74 bytes
	UnknownData5           UnknownData5  `bin:"len:25"`
	Sync                   Sync          `bin:"len:22"` // 22 bytes
	UnknownData6           UnknownData6  `bin:"len:44"`
	UnknownData7           UnknownData7  `bin:"len:14"`
	UnknownData8           UnknownData8  `bin:"len:8"`
	UnknownData9           UnknownData9  `bin:"len:7"`
	UnknownData2           UnknownData2  `bin:"len:14"`
	SnSticker              uint64        `bin:"ReadSN,len:5"`   // BCD
	ProgrammingFactoryDate time.Time     `bin:"BCDDateR,len:3"` // Reversed BCD date dd-mm-yy
	UnknownBytes2          []byte        `bin:"len:3"`
	DelphiPN               uint32        `bin:"le,len:4"` // Little endian, Delphi part number
	UnknownBytes3          []byte        `bin:"len:2"`
	PartNo                 uint32        `bin:"le,len:4"` // Little endian, SAAB part number (factory?)
	UnknownData14          []byte        `bin:"len:3"`
	PSK                    PSK           `bin:"len:14"` // 14 byte
	UnknownData10          UnknownData10 `bin:"len:12"`
	EOF                    byte          `bin:"len:1"` //0x00
}

// Validate all checksums and known tests to ensure a healthy bin
func (bin *Bin) Validate() error {
	tests := []func() error{
		bin.Vin.validate,
		bin.Pin.validate,
		bin.Keys.validate,
		bin.Const1.validate,
		bin.UnknownData1.validate,
		bin.UnknownData2.validate,
		bin.UnknownData3.validate,
		bin.UnknownData4.validate,
		bin.UnknownData5.validate,
		bin.UnknownData6.validate,
		bin.UnknownData7.validate,
		bin.UnknownData8.validate,
		bin.UnknownData9.validate,
		bin.UnknownData10.validate,
		bin.PSK.validate,
	}

	for _, v := range tests {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func (bin *Bin) SasOpt() bool {
	return bin.SasOption == 0x03
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

// Return Serial sticker as uint64, stored as 5byte Binary-Coded Decimal (BCD)
func (*Bin) ReadSN(r binstruct.Reader) (uint64, error) {
	_, b, err := r.ReadBytes(5)
	if err != nil {
		return 0, nil
	}
	return bcd.ToUint64(b), nil
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

// Return model year from VIN
func (bin *Bin) ModelYear() string {
	return fmt.Sprintf("%02s", bin.Vin.Data[9:10])
}
