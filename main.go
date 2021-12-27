package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/roffe/cim/pkg/cim"
	"github.com/roffe/cim/pkg/server"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// if we pass a filename, print to the console instead
	if len(os.Args) == 2 {
		filename := os.Args[1]
		fw, err := cim.Load(filename)
		if err != nil {
			log.Fatal(err)
		}

		if err := fw.Validate(); err != nil {
			log.Fatal(err)
		}
		fw.Pretty()

		return
	}

	// Run web ui
	fmt.Println("open http://localhost:8080")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
