package server

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"reflect"
	"strings"

	"github.com/roffe/cim/pkg/cim"
)

type Section struct {
	ID     string `json:"id"`
	Start  int    `json:"start"`
	Length int    `json:"length"`
	//Confirmed bool
	//Checksum  bool
	Type string `json:"type"`
}

func (s *Section) String() string {
	return fmt.Sprintf(`{id: "%s", start: 0x%02X, length: %d, type: "%s"}`, s.ID, s.Start, s.Length, s.Type)
}

func generateStyles(sections []Section) template.CSS {
	var css strings.Builder
	for _, s := range sections {
		m := md5.Sum([]byte(s.ID))
		col := fmt.Sprintf("%XFF", m[:3])
		b, err := hex.DecodeString(col)
		if err != nil {
			log.Fatal(err)
		}
		for i, bb := range b {
			b[i] = bb << 0x3
		}

		css.WriteString(fmt.Sprintf("\t.section-%s {\n\t\tbackground: #%X;\n\t}\n", strings.ToUpper(s.ID), b))
	}
	return template.CSS(css.String())
}

func generateSections(fw *cim.Bin) []Section {
	var sections []Section
	offset := 0x00
	var itr func(string, reflect.Value)
	itr = func(prefix string, valueOf reflect.Value) {
		for i := 0; i < valueOf.NumField(); i++ {
			tag := valueOf.Type().Field(i).Tag.Get("bin")
			if tag == "-" || tag == "" {
				continue
			}
			field := valueOf.Field(i)
			tname := valueOf.Type().Field(i).Type.String()
			switch tname {
			case "time.Time":
			default:
				if field.Kind() == reflect.Struct {
					itr(valueOf.Type().Field(i).Name, field)
					continue
				}
			}
			parts := strings.Split(tag, ",")
			var previous int
			for _, p := range parts {
				var length int
				var nlength int
				fmt.Sscanf(p, "[len:%d]", &nlength)
				if nlength > 0 {
					asd := (previous*nlength - previous)
					offset += asd
					sections[len(sections)-1].Length = previous * nlength
					previous = 0
					continue
				}
				fmt.Sscanf(p, "len:%d", &length)
				if length == 0 {
					continue
				}
				previous = length

				var fname string
				if prefix == "" {
					fname = valueOf.Type().Field(i).Name

				} else {
					fname = fmt.Sprintf("%s_%s", prefix, valueOf.Type().Field(i).Name)
					fname = strings.TrimSuffix(fname, "1")
					fname = strings.TrimSuffix(fname, "2")
				}
				sections = append(sections, Section{
					ID:     strings.ToUpper(fname),
					Start:  offset, //fmt.Sprintf("0x%X", offset)
					Length: length,
					Type:   tname,
				})
				offset += length
			}
		}
	}

	itr("", reflect.ValueOf(*fw))

	return sections
}
