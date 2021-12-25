package cim

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

func (fw *Bin) Pretty() {
	t := s("CIM Dump analyser: " + filepath.Base(fw.filename))
	t.AppendRows([]table.Row{
		{"MD5", fw.MD5(), "Crc32", fw.CRC32()},
		{"VIN", fw.Vin.Data, "Model Year", fmt.Sprintf("%02s", fw.Vin.Data[9:10])},
		{"PIN", fmt.Sprintf("%s", fw.Pin.Pin1), "Steering Angle Sensor", fw.SasOption},
		{},
	})
	t.Render()

	keys := s(fmt.Sprintf("Programmed keys: %d", fw.Keys.KeysKeysCount1))
	keys.AppendHeader(table.Row{
		"No", "Bank1", "Bank2",
	})
	for i, k := range fw.Keys.Keys1 {
		keys.AppendRow(table.Row{
			i + 1, fmt.Sprintf("%X", k), fmt.Sprintf("%X", fw.Keys.Keys2[i]),
		})
		//fmt.Printf("Key %d: %X / %X\n", i+1, k, fw.Keys.Keys2[i])
	}
	keys.Render()

	isk := s("ISK")
	isk.AppendHeader(table.Row{"ISK", "Bank1", "Bank2"})
	isk.AppendRows([]table.Row{
		{"High", fmt.Sprintf("%X", fw.Keys.IskHI1), fmt.Sprintf("%X", fw.Keys.IskHI2)},
		{"Low", fmt.Sprintf("%X", fw.Keys.IskLO1), fmt.Sprintf("%X", fw.Keys.IskLO2)},
	})
	isk.Render()

	t4 := s("Remotes")
	t4.AppendRows([]table.Row{
		{"PSK High", fmt.Sprintf("%X", fw.PSK.High)},
		{"PSK Low", fmt.Sprintf("%X", fw.PSK.Low)},
		{"PCF", "TODO"},
		{"RSync", fmt.Sprintf("%X", fw.RSync.Data)},
	})
	t4.Render()

	ph := s("Programming history")
	ph.AppendRow(table.Row{"Serial sticker", fw.SnSticker})
	ph.AppendRow(table.Row{"Factory programming date", fw.ProgrammingFactoryDate.Format(isoDate)})
	ph.AppendRow(table.Row{"Last programming date", fw.ProgrammingDate.Format(isoDate)})
	if fw.Vin.SpsCount == 0 {
		ph.AppendRow(table.Row{"Factory programming only"})
	} else {
		ph.AppendRow(table.Row{"SPS Counter", fw.Vin.SpsCount})
		for i := 0; i < int(fw.Vin.SpsCount); i++ {
			ph.AppendRow(table.Row{fmt.Sprintf("Workshop %d ID", i+1), strings.TrimRight(fw.ProgrammingID[i], " ")})
		}
	}

	ph.Render()

	pn := s("Part numbers")
	pn.AppendRows([]table.Row{
		{"End model (HW+SW)", fmt.Sprintf("%d%s", fw.PnSAAB1, fw.PnSAAB1_2)},
		{"Base model (HW+boot)", fmt.Sprintf("%d%s", fw.PnBase1, fw.PnBase1_2)},
		{"Delphi part number", fmt.Sprintf("%d", fw.PnDelphi)},
		{"SAAB part number (factory?)", fmt.Sprintf("%d", fw.PnSAAB2)},
		{"Configuration Version:", fmt.Sprintf("%d", fw.ConfigurationVersion)},
	})
	pn.Render()

}

func s(title string) table.Writer {
	t := table.NewWriter()
	t.SetStyle(tableTheme)
	t.SetOutputMirror(os.Stdout)
	t.Style().Title.Align = text.AlignCenter
	t.SetTitle(title)
	return t
}
