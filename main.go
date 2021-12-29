package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roffe/cim/pkg/cim"
	"github.com/roffe/cim/pkg/server"
	flag "github.com/spf13/pflag"
)

var (
	outputMode     = "pretty"
	debugMode      = false
	enableShutdown = true
	httpPath       = ""
)

func init() {
	gin.SetMode(gin.ReleaseMode)

	flag.StringVarP(&outputMode, "output", "o", outputMode, "pretty|json|string")
	flag.BoolVarP(&debugMode, "debug", "d", debugMode, "true|false")
	flag.BoolVarP(&enableShutdown, "shutdown", "s", enableShutdown, "true|false enable shutdown api")
	flag.StringVar(&httpPath, "path", httpPath, "set http path")
	flag.Parse()

	if debugMode {
		cim.Debug = true
		gin.SetMode(gin.DebugMode)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	// if we pass a filename, print to the console instead of starting ui
	if len(flag.Args()) >= 1 {
		filename := flag.Args()[0]
		fw, err := cim.MustLoad(filename)
		if err != nil {
			log.Fatal(err)
		}
		switch strings.ToLower(outputMode) {
		case "string":
			fw.Dump()
		case "json":
			b, err := fw.Json()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(b[:]))
		case "pretty":
			fw.Pretty()
		default:
			fw.Pretty()
		}
		return
	}

	// Run web ui
	fmt.Println("Server started @ http://localhost:8080")
	if err := server.Run(enableShutdown, httpPath); err != nil {
		log.Fatal(err)
	}
}
