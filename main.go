package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
		fw.Pretty()
		//fw.Dump()
		return
	}

	fmt.Println("open http://localhost:8080")
	if err := serve(); err != nil {
		log.Fatal(err)
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
	width := 40

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

	xorBytes(b)

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
		return "·"
	}
	if a == 0x20 {
		return "&nbsp"
	}
	if a == 0xFF {
		return "Ʃ"
	}
	if a <= 0x20 {
		return "˟"
	}

	if a >= 0x7F {
		return "˟"
	}

	if a == 0x3c || a == 0x3e {
		return "˟"
	}

	return fmt.Sprintf("%s", []byte{b})
}
