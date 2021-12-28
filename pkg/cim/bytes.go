package cim

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"

	"github.com/albenik/bcd"
)

type writeOp struct {
	w io.Writer        // writer
	o binary.ByteOrder // byte-order
	v interface{}      // data
}

func (bin *Bin) XORBytes() ([]byte, error) {
	b, err := bin.Bytes()
	if err != nil {
		return nil, err
	}

	for i, bb := range b {
		b[i] = bb ^ 0xFF
	}
	return b, nil
}

// Return the byte representation of the binary
func (bin *Bin) Bytes() ([]byte, error) {
	o := bytes.NewBuffer([]byte{})
	ops := []writeOp{
		{o, binary.BigEndian, bin.MagicByte},

		{o, binary.LittleEndian, bin.programmingDate()[1:4]},

		{o, binary.LittleEndian, bin.SasOption},

		{o, binary.BigEndian, bin.UnknownBytes1},

		{o, binary.BigEndian, bin.PartNo1},
		{o, binary.BigEndian, []byte(bin.PartNo1Rev)},

		{o, binary.BigEndian, bin.ConfigurationVersion},

		{o, binary.BigEndian, bin.PnBase1},
		{o, binary.BigEndian, []byte(bin.PnBase1Rev)},

		{o, binary.BigEndian, []byte(bin.Vin.Data)},
		{o, binary.BigEndian, bin.Vin.Value},
		{o, binary.BigEndian, bin.Vin.Unknown},
		{o, binary.LittleEndian, bin.Vin.SpsCount},
		{o, binary.LittleEndian, bin.Vin.Checksum},

		{o, binary.BigEndian, []byte(strings.Join(bin.ProgrammingID, ""))},

		{o, binary.BigEndian, bin.UnknownData3.Data1},
		{o, binary.LittleEndian, bin.UnknownData3.Checksum1},
		{o, binary.BigEndian, bin.UnknownData3.Data2},
		{o, binary.LittleEndian, bin.UnknownData3.Checksum2},

		{o, binary.BigEndian, bin.Pin.Data1},
		{o, binary.BigEndian, bin.Pin.Unknown1},
		{o, binary.LittleEndian, bin.Pin.Checksum1},
		{o, binary.BigEndian, bin.Pin.Data2},
		{o, binary.BigEndian, bin.Pin.Unknown2},
		{o, binary.LittleEndian, bin.Pin.Checksum2},

		{o, binary.BigEndian, bin.UnknownData4.Data},
		{o, binary.LittleEndian, bin.UnknownData4.Checksum},

		{o, binary.BigEndian, bin.UnknownData1.Data1},
		{o, binary.LittleEndian, bin.UnknownData1.Checksum1},
		{o, binary.BigEndian, bin.UnknownData1.Data2},
		{o, binary.LittleEndian, bin.UnknownData1.Checksum2},

		{o, binary.BigEndian, bin.Const1.Data},
		{o, binary.LittleEndian, bin.Const1.Checksum},

		{o, binary.BigEndian, bin.Keys.IskHI1},
		{o, binary.BigEndian, bin.Keys.IskLO1},
		{o, binary.BigEndian, bin.Keys.Data1[0]},
		{o, binary.BigEndian, bin.Keys.Data1[1]},
		{o, binary.BigEndian, bin.Keys.Data1[2]},
		{o, binary.BigEndian, bin.Keys.Data1[3]},
		{o, binary.BigEndian, bin.Keys.Data1[4]},
		{o, binary.LittleEndian, bin.Keys.Count1},
		{o, binary.BigEndian, bin.Keys.Constant1},
		{o, binary.LittleEndian, bin.Keys.Errors1},
		{o, binary.LittleEndian, bin.Keys.Checksum1},
		{o, binary.BigEndian, bin.Keys.IskHI2},
		{o, binary.BigEndian, bin.Keys.IskLO2},
		{o, binary.BigEndian, bin.Keys.Data2[0]},
		{o, binary.BigEndian, bin.Keys.Data2[1]},
		{o, binary.BigEndian, bin.Keys.Data2[2]},
		{o, binary.BigEndian, bin.Keys.Data2[3]},
		{o, binary.BigEndian, bin.Keys.Data2[4]},
		{o, binary.LittleEndian, bin.Keys.Count2},
		{o, binary.BigEndian, bin.Keys.Constant2},
		{o, binary.LittleEndian, bin.Keys.Errors2},
		{o, binary.LittleEndian, bin.Keys.Checksum2},

		{o, binary.BigEndian, bin.UnknownData5.Data},
		{o, binary.LittleEndian, bin.UnknownData5.Checksum},

		{o, binary.BigEndian, bin.Sync.Data[0]},
		{o, binary.BigEndian, bin.Sync.Data[1]},
		{o, binary.BigEndian, bin.Sync.Data[2]},
		{o, binary.BigEndian, bin.Sync.Data[3]},
		{o, binary.BigEndian, bin.Sync.Data[4]},

		{o, binary.LittleEndian, bin.Sync.Checksum},

		{o, binary.BigEndian, bin.UnknownData6.Data1},
		{o, binary.LittleEndian, bin.UnknownData6.Checksum1},
		{o, binary.BigEndian, bin.UnknownData6.Data2},
		{o, binary.LittleEndian, bin.UnknownData6.Checksum2},

		{o, binary.BigEndian, bin.UnknownData7.Data1},
		{o, binary.LittleEndian, bin.UnknownData7.Checksum1},
		{o, binary.BigEndian, bin.UnknownData7.Data2},
		{o, binary.LittleEndian, bin.UnknownData7.Checksum2},

		{o, binary.BigEndian, bin.UnknownData8.Data},
		{o, binary.LittleEndian, bin.UnknownData8.Checksum},

		{o, binary.BigEndian, bin.UnknownData9.Data},
		{o, binary.LittleEndian, bin.UnknownData9.Checksum},

		{o, binary.BigEndian, bin.UnknownData2.Data1},
		{o, binary.LittleEndian, bin.UnknownData2.Checksum1},
		{o, binary.BigEndian, bin.UnknownData2.Data2},
		{o, binary.LittleEndian, bin.UnknownData2.Checksum2},
		{o, binary.BigEndian, bcd.FromUint64(bin.SnSticker)[3:8]},

		{o, binary.LittleEndian, bin.programmingFactoryDate()[1:4]},
		{o, binary.LittleEndian, bin.UnknownBytes2},

		{o, binary.LittleEndian, bin.DelphiPN},

		{o, binary.BigEndian, bin.UnknownBytes3},

		{o, binary.LittleEndian, bin.PartNo},

		{o, binary.LittleEndian, bin.UnknownData14},

		{o, binary.LittleEndian, bin.PSK.Low},
		{o, binary.LittleEndian, bin.PSK.High},
		{o, binary.LittleEndian, bin.PSK.Constant},
		{o, binary.LittleEndian, bin.PSK.Unknown},
		{o, binary.LittleEndian, bin.PSK.Checksum},

		{o, binary.BigEndian, bin.UnknownData10.Data1},
		{o, binary.LittleEndian, bin.UnknownData10.Checksum1},
		{o, binary.BigEndian, bin.UnknownData10.Data2},
		{o, binary.LittleEndian, bin.UnknownData10.Checksum2},

		{o, binary.BigEndian, bin.EOF},
	}

	for _, o := range ops {
		if err := binary.Write(o.w, o.o, o.v); err != nil {
			return nil, err
		}
	}

	return o.Bytes(), nil
}
