package main

import (
	"log"
	"os"

	"github.com/roffe/cim/pkg/cim"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatal("missing input filename")
	}
}

func main() {
	filename := os.Args[1]
	fw, err := cim.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	fw.Dump()
}
