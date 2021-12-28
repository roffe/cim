package main

import (
	"encoding/json"
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
	// if we pass a filename, print to the console instead of starting ui
	if len(os.Args) == 2 {
		filename := os.Args[1]
		fw, err := cim.MustLoad(filename)
		if err != nil {
			log.Fatal(err)
		}
		fw.Pretty()

		b, err := json.MarshalIndent(fw, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b[:]))
		return
	}

	// Run web ui
	fmt.Println("Server started @ http://localhost:8080")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
