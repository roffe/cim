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
		{"MD5", fw.MD5()},
		{"Crc32", fw.CRC32()},
		{"VIN", fw.Vin.Data},
		{"Model Year", fw.Vin.Data[9:10]},
		{"Steering Angle Sensor", fw.SasOption},
	})
	t.Render()

	pin := s("Pin")
	pin.AppendHeader(table.Row{"#", "Bank 1", "Bank 2"})
	pin.AppendRow(table.Row{0, fmt.Sprintf("%X", fw.Pin.Pin1), fmt.Sprintf("%X", fw.Pin.Pin2)})
	pin.Render()

	keys := s(fmt.Sprintf("Programmed keys: %d", fw.Keys.KeysKeysCount1))
	keys.AppendHeader(table.Row{
		"#", "Bank 1", "Bank 2",
	})
	for i, k := range fw.Keys.Keys1 {
		keys.AppendRow(table.Row{
			i + 1, fmt.Sprintf("%X", k), fmt.Sprintf("%X", fw.Keys.Keys2[i]),
		})
	}
	keys.Render()

	isk := s("ISK")
	isk.AppendHeader(table.Row{"ISK", "Bank1", "Bank2"})
	isk.AppendRows([]table.Row{
		{"High", fmt.Sprintf("%X", fw.Keys.IskHI1), fmt.Sprintf("%X", fw.Keys.IskHI2)},
		{"Low", fmt.Sprintf("%X", fw.Keys.IskLO1), fmt.Sprintf("%X", fw.Keys.IskLO2)},
	})
	isk.Render()

	r := s("Remotes")
	r.AppendRows([]table.Row{
		{"PSK High", fmt.Sprintf("%X", fw.PSK.High)},
		{"PSK Low", fmt.Sprintf("%X", fw.PSK.Low)},
		{"PCF", "TODO"},
		{"Sync", fmt.Sprintf("%X", fw.Sync.Data)},
	})
	r.Render()

	ph := s("Programming history")
	ph.AppendRow(table.Row{"Serial sticker", fw.SnSticker})
	ph.AppendRow(table.Row{"Factory programming date", fw.ProgrammingFactoryDate.Format(IsoDate)})
	ph.AppendRow(table.Row{"Last programming date", fw.ProgrammingDate.Format(IsoDate)})
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
		{"End model (HW+SW)", fmt.Sprintf("%d%s", fw.PartNo1, fw.PartNo1Suffix)},
		{"Base model (HW+boot)", fmt.Sprintf("%d%s", fw.PnBase1, fw.PnBase1Suffix)},
		{"Delphi part number", fw.DelphiPN},
		{"SAAB part number", fw.PartNo},
		{"Configuration Version:", fw.ConfigurationVersion},
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
