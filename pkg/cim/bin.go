package cim

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/albenik/bcd"
	"github.com/ghostiam/binstruct"
)

const isoDate = "2006-01-02"

func Load(filename string) (*Bin, error) {
	//var fw Bin
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

	if err := fw.Pin.validate(); err != nil {
		log.Fatal(err)
	}

	return &fw, nil
}

type Bin struct {
	filename               string    `bin:"-"`
	md5                    string    `bin:"-"`
	crc32                  string    `bin:"-"`
	MagicByte              uint8     `bin:"len:1"`         // 0x20
	ProgrammingDate        time.Time `bin:"BCDDate,len:3"` // BCD Binary-Coded Decimal yy-mm-dd
	SasOption              uint8     `bin:"len:1"`
	UnknownBytes1          []byte    `bin:"len:6"`
	PnSAAB1                uint32    `bin:"len:4"`
	PnSAAB1_2              string    `bin:"len:2"`
	ConfigurationVersion   uint32    `bin:"len:4"`
	PnBase1                uint32    `bin:"len:4"`
	PnBase1_2              string    `bin:"len:2"`
	Vin                    Vin       // 30 bytes
	ProgrammingID          []string  `bin:"len:3,[len:10]"`
	UnknownData3           []byte    `bin:"len:88"`
	Pin                    Pin       // 20 bytes
	UnknownData4           []byte    `bin:"len:4"`
	UnknownData1           []byte    `bin:"len:44"`
	Const1                 []byte    `bin:"len:10"`
	Keys                   Keys      // 74 bytes
	UnknownData5           []byte    `bin:"len:25"`
	RSync                  RSync     // 22 bytes
	UnknownData6           []byte    `bin:"len:44"`
	UnknownData7           []byte    `bin:"len:14"`
	UnknownData8           []byte    `bin:"len:8"`
	UnknownData9           []byte    `bin:"len:7"`
	UnknownData2           []byte    `bin:"len:14"`
	SnSticker              uint64    `bin:"ReadSN,len:5"`   // BCD
	ProgrammingFactoryDate time.Time `bin:"BCDDateR,len:3"` // Reversed BCD date dd-mm-yy
	UnknownBytes2          []byte    `bin:"len:3"`
	PnDelphi               uint32    `bin:"Uint32l,len:4"`
	UnknownBytes3          []byte    `bin:"len:2"`
	PnSAAB2                uint32    `bin:"Uint32l,len:4"`
	UnknownBytes4          []byte    `bin:"len:3"`
	PSK                    PSK
	UnknownData10          []byte `bin:"len:12"`
	EOF                    byte   `bin:"len:1"` //0x00
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

func (*Bin) ReadSN(r binstruct.Reader) (uint64, error) {
	_, b, err := r.ReadBytes(5)
	if err != nil {
		return 0, nil
	}
	return bcd.ToUint64(b), nil
}

// Uint32l returns a uint32 read as little endian
func (*Bin) Uint32l(r binstruct.Reader) (uint32, error) {
	var out uint32
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
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

func (fw *Bin) Dump() {
	fmt.Println("Bin file:", filepath.Base(fw.filename))
	fmt.Println("MD5:", fw.MD5())
	fmt.Println("CRC32:", fw.CRC32())
	fmt.Println("")

	fmt.Println("VIN:", fw.Vin.Data)
	fmt.Printf("PIN: %q / %q\n", fw.Pin.Pin1, fw.Pin.Pin2)
	fmt.Println("")

	fmt.Printf("Model Year: %02s\n", fw.Vin.Data[9:10])
	fmt.Printf("Steering Angle Sensor: %d\n", fw.SasOption)
	fmt.Println("")

	fmt.Println("Programmed keys:", fw.Keys.KeysKeysCount1)
	for i, k := range fw.Keys.Keys1 {
		fmt.Printf("Key %d: %X / %X\n", i+1, k, fw.Keys.Keys2[i])
	}
	fmt.Printf("ISK High: %X / %X\n", fw.Keys.IskHI1, fw.Keys.IskHI2)
	fmt.Printf("ISK Low: %X / %X\n", fw.Keys.IskLO1, fw.Keys.IskLO2)
	fmt.Println()

	fmt.Println("Remotes:")
	fmt.Printf("PSK High: %X\n", fw.PSK.High)
	fmt.Printf("PSK Low:  %X\n", fw.PSK.Low)
	fmt.Printf("PCF: %s\n", "TODO")
	fmt.Printf("Sync: ")
	for _, v := range fw.RSync.Data {
		fmt.Printf("%X ", v)
	}
	fmt.Println()
	fmt.Println()

	fmt.Println("Programming history:")
	fmt.Printf("- Last programming date: %s\n", fw.ProgrammingDate.Format(isoDate))
	if fw.Vin.SpsCount == 0 {
		fmt.Println("- Factory programming only")
	} else {
		fmt.Printf("- SPS Counter: %d\n", fw.Vin.SpsCount)
		for i := 0; i < int(fw.Vin.SpsCount); i++ {
			fmt.Printf("- Workshop %d ID: %s\n", i+1, strings.TrimRight(fw.ProgrammingID[i], " "))
		}
	}
	fmt.Println()

	fmt.Printf("Serial sticker: %d\n", fw.SnSticker)
	fmt.Printf("Factory programming date: %s\n", fw.ProgrammingFactoryDate.Format(isoDate))
	fmt.Println()

	fmt.Println("Part numbers:")
	fmt.Printf("- End model (HW+SW): %d%s\n", fw.PnSAAB1, fw.PnSAAB1_2)
	fmt.Printf("- Base model (HW+boot): %d%s\n", fw.PnBase1, fw.PnBase1_2)
	fmt.Printf("- Delphi part number: %d\n", fw.PnDelphi)
	fmt.Printf("- SAAB part number (factory?): %d\n", fw.PnSAAB2)
	fmt.Printf("- Configuration Version: %d\n", fw.ConfigurationVersion)
	fmt.Println()
}
