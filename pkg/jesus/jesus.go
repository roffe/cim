package jesus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"reflect"
)

func Dump(v interface{}) {
	fmt.Println(v)
	fmt.Printf("%X\n", v)
	fmt.Printf("%q\n", v)
	reverseSlice(v)
	fmt.Println(v)
	fmt.Printf("%X\n", v)
	fmt.Printf("%q\n", v)
	reverseSlice(v)

	switch t := v.(type) {
	case []byte:
		var u16le uint16
		if err := binary.Read(bytes.NewReader(t), binary.LittleEndian, &u16le); err != nil {
			log.Fatal(err)
		}
		fmt.Println("u16le", u16le)
		var u16be uint16
		if err := binary.Read(bytes.NewReader(t), binary.BigEndian, &u16be); err != nil {
			log.Fatal(err)
		}
		fmt.Println("u16be", u16be)

		if len(t) == 3 {
			t = append(t, 0x00)
		}

		var u32le uint32
		if err := binary.Read(bytes.NewReader(t), binary.LittleEndian, &u32le); err != nil {
			log.Fatal(err)
		}
		fmt.Println("u32le", u32le)
		var u32be uint32
		if err := binary.Read(bytes.NewReader(t), binary.BigEndian, &u32be); err != nil {
			log.Fatal(err)
		}
		fmt.Println("u32be", u32be)

		if len(t) == 5 {
			b := append([]byte{0x00, 0x00, 0x00, 0x00}, t...)
			var u64le uint64
			if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &u64le); err != nil {
				log.Fatal(err)
			}
			fmt.Println("5u64le", u64le)
			var u64be uint64
			if err := binary.Read(bytes.NewReader(b), binary.BigEndian, &u64be); err != nil {
				log.Fatal(err)
			}
			fmt.Println("5u64be", u64be)

			b2 := append(t, []byte{0x00, 0x00, 0x00, 0x00}...)
			var u64le2 uint64
			if err := binary.Read(bytes.NewReader(b2), binary.LittleEndian, &u64le2); err != nil {
				log.Fatal(err)
			}
			fmt.Println("5u64leb", u64le)
			var u64be2 uint64
			if err := binary.Read(bytes.NewReader(b2), binary.BigEndian, &u64be2); err != nil {
				log.Fatal(err)
			}
			fmt.Println("5u64beb", u64be)

		} else {
			var u64le uint64
			if err := binary.Read(bytes.NewReader(t), binary.LittleEndian, &u64le); err != nil {
				log.Fatal(err)
			}
			fmt.Println("u64le", u64le)
			var u64be uint64
			if err := binary.Read(bytes.NewReader(t), binary.BigEndian, &u64be); err != nil {
				log.Fatal(err)
			}
			fmt.Println("u64be", u64be)
		}
		/*
			for i := byte(0x00); i < 0xff; i++ {
				var os []byte
				for _, bbb := range t {
					os = append(os, bbb^i)
				}
				fmt.Printf("%d %X\n", i, os)
				//fmt.Printf("%s\n", os)
				//fmt.Printf("%d %d\n", i, os)

			}
		*/
	}
}

func reverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}
