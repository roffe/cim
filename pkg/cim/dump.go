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
	fmt.Printf("PIN: %q / %q\n", fw.Pin.Pin1, fw.Pin.Pin2)
	fmt.Println("")

	fmt.Printf("Model Year: %02s\n", fw.Vin.Data[9:10])
	fmt.Printf("Steering Angle Sensor: %t\n", fw.SasOption)
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
