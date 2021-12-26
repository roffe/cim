package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roffe/cim/pkg/cim"
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
		//fw.Pretty()
		fmt.Println(filename)
		jesus(fw.UnknownData14)
		fmt.Println("----")
		//fw.Dump()
		return
	}

	fmt.Println("open http://localhost:8080")
	if err := serve(); err != nil {
		log.Fatal(err)
	}

}

func jesus(v interface{}) {
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

var templateHelpers = template.FuncMap{
	"printHex": func(v interface{}) template.HTML {
		return template.HTML(fmt.Sprintf("%X", v))
	},
	"print": func(v interface{}) template.HTML {
		return template.HTML(fmt.Sprintf("%s", v))
	},
	"isoDate": func(t time.Time) template.HTML {
		return template.HTML(t.Format(cim.IsoDate))
	},
	"boolChecked": func(b bool) template.HTML {
		if b {
			return template.HTML("checked")
		}
		return template.HTML("")
	},
	"keyOffset": func(factor int) template.HTML {
		return template.HTML(fmt.Sprintf("%d", 259+(4*factor)))
	},
}

func serve() error {
	r := gin.Default()
	// Load templates
	//r.LoadHTMLGlob("templates/*.tmpl")
	if tmpl, err := template.New("projectViews").Funcs(templateHelpers).ParseGlob("templates/*.tmpl"); err == nil {
		r.SetHTMLTemplate(tmpl)
	} else {
		return err
	}

	// Set upload limit for multipart form
	r.MaxMultipartMemory = 1 << 20

	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "upload.tmpl", nil) })
	r.POST("/save", saveHandler)
	r.POST("/", uploadHandler)

	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// Handle file uploads
func uploadHandler(c *gin.Context) {
	buf, filename, n, err := getFileFromCtx(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// We only accept exactly 512 bytes
	if n < 512 || n > 512 {
		c.String(http.StatusInternalServerError, "invalid bin size")
		return
	}

	fw, err := cim.LoadBytes(filename, buf)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := fw.Validate(); err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fwBytes, err := fw.Bytes()
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	hexRows := strings.Builder{}
	asciiColumns := strings.Builder{}

	pos := 0
	offset := 0
	width := 36

	for _, bb := range fwBytes {
		if pos == 0 {
			hexRows.WriteString(`<div class="hexRow">` + "\n" +
				"\t" + `<div class="addrColumn"><b>` + fmt.Sprintf("%03X", offset) + "</b></div>\n" +
				"\t" + `<div class="hexColumns">` + "\n")
		}

		hexRows.WriteString(fmt.Sprintf("\t\t"+`<div class="hexByte byte-%d" data-i="%d">%02X</div>`+"\n", offset, offset, bb))
		asciiColumns.WriteString(fmt.Sprintf(`<div class="asciiByte byte-%d" data-i="%d">%s</div>`+"\n", offset, offset, psafe(bb)))
		if pos == width {
			hexRows.WriteString("</div>\n")
			hexRows.WriteString(`<div class="asciiColumns">` + "\n")
			hexRows.WriteString(asciiColumns.String())
			hexRows.WriteString("</div>\n")
			hexRows.WriteString("</div>\n")
			asciiColumns.Reset()
			pos = 0
			offset++
			continue
		}
		pos++
		offset++
	}
	// Handle the tail that didn't fill a full width
	if pos <= width {
		for i := pos; i <= width; i++ {
			hexRows.WriteString(`<div class="fillByte">&nbsp;&nbsp;</div>` + "\n")
		}
		hexRows.WriteString("</div>")
		hexRows.WriteString(`<div class="asciiColumns">` + "\n")
		hexRows.WriteString(asciiColumns.String())
		hexRows.WriteString("</div>\n")
		hexRows.WriteString("</div>\n")
		asciiColumns.Reset()
	}

	hexRows.WriteString("</div>")

	b64 := base64.StdEncoding.EncodeToString(fwBytes)

	c.HTML(http.StatusOK, "hex.tmpl", gin.H{
		"filename": filepath.Base(filename),
		"fw":       fw,
		"B64":      b64,
		"Hexview":  template.HTML(hexRows.String()),
	})
}

func saveHandler(c *gin.Context) {
	file := c.PostForm("file")
	filename := c.PostForm("filename")

	if file == "" || filename == "" {
		c.String(http.StatusBadRequest, "missing file or filename")
		return
	}

	b, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	fw, err := cim.LoadBytes(filename, b)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	bs, err := fw.Bytes()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	xorBytes(bs)

	contentLength := int64(len(bs))
	contentType := "application/octet-stream"

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + filepath.Base(filename) + `"`,
	}
	c.DataFromReader(http.StatusOK, contentLength, contentType, bytes.NewReader(bs), extraHeaders)
}

func getFileFromCtx(c *gin.Context) ([]byte, string, int64, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return nil, "", 0, fmt.Errorf("getFileFromCtx err 1: %s", err.Error())
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	n, err := io.Copy(buf, file)
	if err != nil {
		return nil, "", 0, fmt.Errorf("getFileFromCtx err 2: %s", err.Error())
	}
	return buf.Bytes(), header.Filename, n, nil
}

func xorBytes(b []byte) []byte {
	for i, bb := range b {
		b[i] = bb ^ 0xFF
	}
	return b
}

// sanitize binary so we don't print controll characters
func psafe(b byte) string {
	a := uint8(b)
	if a == 0x00 {
		return "&centerdot;"
	}
	if a == 0x20 {
		return "&nbsp"
	}
	if a == 0xFF {
		return "&fflig;"
	}
	if a <= 0x20 {
		return "&#9618;"
	}

	if a >= 0x7F {
		return "&block;"
	}

	if a == 0x3c || a == 0x3e {
		return "˟"
	}

	return fmt.Sprintf("%s", []byte{b})
}
