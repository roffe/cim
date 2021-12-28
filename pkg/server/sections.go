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
	Checksum bool
	Type     string `json:"type"`
}

func (s *Section) String() string {
	return fmt.Sprintf(`{id: "%s", start: 0x%02X, length: %d, type: "%s", checksum: %t}`, s.ID, s.Start, s.Length, s.Type, s.Checksum)
}

func generateStyles(sections []Section) template.CSS {
	var css strings.Builder
	for _, s := range sections {
		b, err := hex.DecodeString(
			fmt.Sprintf("%x", md5.Sum([]byte(s.ID)[:3])),
		)
		if err != nil {
			log.Fatal(err)
		}
		// bit-shift value by 2 for some nicer colors
		for i, bb := range b {
			b[i] = bb << 0x2
		}
		css.WriteString(fmt.Sprintf("\t.section-%s {\n\t\tbackground: #%X;\n\t}\n", strings.ToUpper(s.ID), b))
	}
	return template.CSS(css.String())
}

func generateSections(fw *cim.Bin) []Section {
	// Byte position in the file
	var offset int
	var sections []Section

	// create a empty anonymous function
	var loop func(string, reflect.Value)

	// Define the function with references to itself via anonymous function
	loop = func(prefix string, valueOf reflect.Value) {
		for i := 0; i < valueOf.NumField(); i++ {
			tag := valueOf.Type().Field(i).Tag.Get("bin")
			if tag == "-" || tag == "" {
				continue
			}

			field := valueOf.Field(i)
			fieldName := valueOf.Type().Field(i).Name
			typeName := valueOf.Type().Field(i).Type.String()

			switch typeName {
			// time.Time is a struct and we don't want to itterate over any sub structures in it
			case "time.Time":
			default:
				// Recurse any structures
				if field.Kind() == reflect.Struct {
					loop(fieldName, field)
					continue
				}
			}
			// Keep track of previous bin len value for multi dimension arrays
			var previous int
			for _, p := range strings.Split(tag, ",") {
				var length, nlength int
				if _, err := fmt.Sscanf(p, "[len:%d]", &nlength); err == nil {
					structLen := (previous*nlength - previous)
					offset += structLen
					sections[len(sections)-1].Length = previous * nlength
					previous = 0
					continue
				}
				if _, err := fmt.Sscanf(p, "len:%d", &length); err != nil {
					continue
				}
				previous = length
				fname := genFieldName(prefix, fieldName)
				sections = append(sections, Section{
					ID:       fname,
					Start:    offset,
					Length:   length,
					Type:     typeName,
					Checksum: strings.Contains(fname, "CHECKSUM"),
				})
				offset += length
			}
		}
	}
	// start the recursing
	loop("", reflect.ValueOf(*fw))

	return sections
}

func genFieldName(prefix, name string) string {
	var fname string
	if prefix == "" {
		fname = name

	} else {
		fname = fmt.Sprintf("%s_%s", prefix, name)
		// Trim bank numbers at end so we get joined hilight sections
		fname = strings.TrimSuffix(fname, "1")
		fname = strings.TrimSuffix(fname, "2")
	}
	return strings.ToUpper(fname)
}
