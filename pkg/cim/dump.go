package cim

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (fw *Bin) Dump() {
	fmt.Println("Bin file:", filepath.Base(fw.filename))
	fmt.Println("MD5:", fw.MD5())
	fmt.Println("CRC32:", fw.CRC32())
	fmt.Println("")

	fmt.Println("VIN:", fw.Vin.Data)
	fmt.Printf("PIN: %q / %q\n", fw.Pin.Data1, fw.Pin.Data2)
	fmt.Println("")

	fmt.Printf("Model Year: %02s\n", fw.Vin.Data[9:10])
	fmt.Printf("Steering Angle Sensor: %t\n", fw.SasOpt())
	fmt.Println("")

	fmt.Println("Programmed keys:", fw.Keys.Count1)
	for i, k := range fw.Keys.Data1 {
		fmt.Printf("Key %d: %X / %X\n", i+1, k, fw.Keys.Data2[i])
	}
	fmt.Printf("ISK High: %X / %X\n", fw.Keys.IskHI1, fw.Keys.IskHI2)
	fmt.Printf("ISK Low: %X / %X\n", fw.Keys.IskLO1, fw.Keys.IskLO2)
	fmt.Println()

	fmt.Println("Remotes:")
	fmt.Printf("PSK High: %X\n", fw.PSK.High)
	fmt.Printf("PSK Low:  %X\n", fw.PSK.Low)
	fmt.Printf("PCF: %s\n", "6732F2C5")
	fmt.Printf("Sync: ")
	for _, v := range fw.Sync.Data {
		fmt.Printf("%X ", v)
	}
	fmt.Println()
	fmt.Println()

	fmt.Println("Programming history:")
	fmt.Printf("- Last programming date: %s\n", fw.ProgrammingDate.Format(IsoDate))
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
	fmt.Printf("Factory programming date: %s\n", fw.ProgrammingFactoryDate.Format(IsoDate))
	fmt.Println()

	fmt.Println("Part numbers:")
	fmt.Printf("- End model (HW+SW): %d%s\n", fw.PartNo1, fw.PartNo1Suffix)
	fmt.Printf("- Base model (HW+boot): %d%s\n", fw.PnBase1, fw.PnBase1Suffix)
	fmt.Printf("- Delphi part number: %d\n", fw.DelphiPN)
	fmt.Printf("- SAAB part number: %d\n", fw.PartNo)
	fmt.Printf("- Configuration Version: %d\n", fw.ConfigurationVersion)
	fmt.Println()
}
