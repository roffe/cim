package cim

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/albenik/bcd"
)

type writeOp struct {
	order binary.ByteOrder // Endianess
	data  interface{}
}

// Return the byte representation of the memory dump
func (bin *Bin) Bytes() ([]byte, error) {
	ops := []writeOp{
		{binary.BigEndian, bin.MagicByte},
		{binary.LittleEndian, bin.programmingDate()[1:4]},
		{binary.LittleEndian, bin.SasOption},
		{binary.BigEndian, bin.UnknownBytes1},
		{binary.BigEndian, bin.PartNo1},
		{binary.BigEndian, []byte(bin.PartNo1Rev)},
		{binary.BigEndian, bin.ConfigurationVersion},
		{binary.BigEndian, bin.PnBase1},
		{binary.BigEndian, []byte(bin.PnBase1Rev)},
		{binary.BigEndian, []byte(bin.Vin.Data)},
		{binary.BigEndian, bin.Vin.Value},
		{binary.BigEndian, bin.Vin.Unknown},
		{binary.LittleEndian, bin.Vin.SpsCount},
		{binary.LittleEndian, bin.Vin.Checksum},
		{binary.BigEndian, []byte(strings.Join(bin.ProgrammingID, ""))},
		{binary.BigEndian, bin.UnknownData3.Data1},
		{binary.LittleEndian, bin.UnknownData3.Checksum1},
		{binary.BigEndian, bin.UnknownData3.Data2},
		{binary.LittleEndian, bin.UnknownData3.Checksum2},
		{binary.BigEndian, bin.Pin.Data1},
		{binary.BigEndian, bin.Pin.Unknown1},
		{binary.LittleEndian, bin.Pin.Checksum1},
		{binary.BigEndian, bin.Pin.Data2},
		{binary.BigEndian, bin.Pin.Unknown2},
		{binary.LittleEndian, bin.Pin.Checksum2},
		{binary.BigEndian, bin.UnknownData4.Data},
		{binary.LittleEndian, bin.UnknownData4.Checksum},
		{binary.BigEndian, bin.UnknownData1.Data1},
		{binary.LittleEndian, bin.UnknownData1.Checksum1},
		{binary.BigEndian, bin.UnknownData1.Data2},
		{binary.LittleEndian, bin.UnknownData1.Checksum2},
		{binary.BigEndian, bin.Const1.Data},
		{binary.LittleEndian, bin.Const1.Checksum},
		{binary.BigEndian, bin.Keys.IskHI1},
		{binary.BigEndian, bin.Keys.IskLO1},
		{binary.BigEndian, bin.Keys.Data1[0]},
		{binary.BigEndian, bin.Keys.Data1[1]},
		{binary.BigEndian, bin.Keys.Data1[2]},
		{binary.BigEndian, bin.Keys.Data1[3]},
		{binary.BigEndian, bin.Keys.Data1[4]},
		{binary.LittleEndian, bin.Keys.Count1},
		{binary.BigEndian, bin.Keys.Constant1},
		{binary.LittleEndian, bin.Keys.Errors1},
		{binary.LittleEndian, bin.Keys.Checksum1},
		{binary.BigEndian, bin.Keys.IskHI2},
		{binary.BigEndian, bin.Keys.IskLO2},
		{binary.BigEndian, bin.Keys.Data2[0]},
		{binary.BigEndian, bin.Keys.Data2[1]},
		{binary.BigEndian, bin.Keys.Data2[2]},
		{binary.BigEndian, bin.Keys.Data2[3]},
		{binary.BigEndian, bin.Keys.Data2[4]},
		{binary.LittleEndian, bin.Keys.Count2},
		{binary.BigEndian, bin.Keys.Constant2},
		{binary.LittleEndian, bin.Keys.Errors2},
		{binary.LittleEndian, bin.Keys.Checksum2},
		{binary.BigEndian, bin.UnknownData5.Data},
		{binary.LittleEndian, bin.UnknownData5.Checksum},
		{binary.BigEndian, bin.Sync.Data[0]},
		{binary.BigEndian, bin.Sync.Data[1]},
		{binary.BigEndian, bin.Sync.Data[2]},
		{binary.BigEndian, bin.Sync.Data[3]},
		{binary.BigEndian, bin.Sync.Data[4]},
		{binary.LittleEndian, bin.Sync.Checksum},
		{binary.BigEndian, bin.UnknownData6.Data1},
		{binary.LittleEndian, bin.UnknownData6.Checksum1},
		{binary.BigEndian, bin.UnknownData6.Data2},
		{binary.LittleEndian, bin.UnknownData6.Checksum2},
		{binary.BigEndian, bin.UnknownData7.Data1},
		{binary.LittleEndian, bin.UnknownData7.Checksum1},
		{binary.BigEndian, bin.UnknownData7.Data2},
		{binary.LittleEndian, bin.UnknownData7.Checksum2},
		{binary.BigEndian, bin.UnknownData8.Data},
		{binary.LittleEndian, bin.UnknownData8.Checksum},
		{binary.BigEndian, bin.UnknownData9.Data},
		{binary.LittleEndian, bin.UnknownData9.Checksum},
		{binary.BigEndian, bin.UnknownData2.Data1},
		{binary.LittleEndian, bin.UnknownData2.Checksum1},
		{binary.BigEndian, bin.UnknownData2.Data2},
		{binary.LittleEndian, bin.UnknownData2.Checksum2},
		{binary.BigEndian, bcd.FromUint64(bin.SnSticker)[3:8]},
		{binary.LittleEndian, bin.programmingFactoryDate()[1:4]},
		{binary.LittleEndian, bin.UnknownBytes2},
		{binary.LittleEndian, bin.DelphiPN},
		{binary.BigEndian, bin.UnknownBytes3},
		{binary.LittleEndian, bin.PartNo},
		{binary.LittleEndian, bin.UnknownData14},
		{binary.LittleEndian, bin.PSK.Low},
		{binary.LittleEndian, bin.PSK.High},
		{binary.LittleEndian, bin.PSK.Constant},
		{binary.LittleEndian, bin.PSK.Unknown},
		{binary.LittleEndian, bin.PSK.Checksum},
		{binary.BigEndian, bin.UnknownData10.Data1},
		{binary.LittleEndian, bin.UnknownData10.Checksum1},
		{binary.BigEndian, bin.UnknownData10.Data2},
		{binary.LittleEndian, bin.UnknownData10.Checksum2},
		{binary.BigEndian, bin.EOF},
	}
	buf := bytes.NewBuffer([]byte{})
	for _, o := range ops {
		if err := binary.Write(buf, o.order, o.data); err != nil {
			return nil, fmt.Errorf("failed to write bytes to buffer: %v", err)
		}
	}

	return buf.Bytes(), nil
}

// Return Bytes() Xored for you ready for flashing
func (bin *Bin) XORBytes() ([]byte, error) {
	b, err := bin.Bytes()
	if err != nil {
		return nil, fmt.Errorf("xor bytes failed to read bytes: %v", err)
	}
	for i, bb := range b {
		b[i] = bb ^ 0xFF
	}
	return b, nil
}
