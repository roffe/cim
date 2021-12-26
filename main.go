package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roffe/cim/pkg/cim"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	/*
		if len(os.Args) < 2 {
			log.Fatal("missing input filename")
		}
	*/
}

func main() {
	/*
		filename := os.Args[1]
		fw, err := cim.Load(filename)
		if err != nil {
			log.Fatal(err)
		}
		if err := fw.Validate(); err != nil {
			log.Fatal(err)
		}
	*/
	//fw.Pretty()
	//fw.Dump()
	//log.Println(fw.Vin.Value)
	fmt.Println("open http://localhost:8080")
	web()

}

func web() {
	r := gin.Default()
	r.MaxMultipartMemory = 1 << 20
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.tmpl", nil)
	})

	r.POST("/save", func(c *gin.Context) {
		file := c.PostForm("file")
		filename := c.PostForm("filename")

		b, err := base64.StdEncoding.DecodeString(file)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		for i, bb := range b {
			b[i] = bb ^ 0xff
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
		for i, b := range bs {
			bs[i] = b ^ 0xFF
		}

		contentLength := int64(len(bs))
		contentType := "application/octet-stream"

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="` + filepath.Base(filename) + `"`,
		}
		r := bytes.NewReader(bs)
		c.DataFromReader(http.StatusOK, contentLength, contentType, r, extraHeaders)
	})

	r.POST("/", func(c *gin.Context) {
		buf, filename, n, err := getFileFromCtx(c)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		//fmt.Printf("%X\n", buf)

		log.Println("bin uploaded, size:", n)
		if n < 512 || n > 512 {
			c.String(http.StatusInternalServerError, "invalid bin size")
			return
		}

		fw, err := cim.LoadBytes(filename, buf)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		bs, err := fw.Bytes()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		hexRow := strings.Builder{}
		asciiColumns := strings.Builder{}

		pos := 0
		offset := 0
		for _, bb := range bs {
			if pos == 0 {
				hexRow.WriteString(`
				<div class="hexRow">
                <div class="addrColumn">` + fmt.Sprintf("%03X", offset) + ` </div>
                <div class="hexColumns">
				`)
			}
			//hexRow.WriteString(fmt.Sprintf("%02X&nbsp;", bb))

			hexRow.WriteString(fmt.Sprintf(`<div class="hexByte byte-%d" data-i="%d">%02X</div>`+"\n", offset, offset, bb))
			asciiColumns.WriteString(fmt.Sprintf(`<div class="asciiByte byte-%d" data-i="%d">%s</div>`+"\n", offset, offset, ps(bb)))
			if pos == 26 {
				hexRow.WriteString("</div>")
				hexRow.WriteString(`<div class="asciiColumns">`)
				hexRow.WriteString(asciiColumns.String())
				hexRow.WriteString(`</div>`)
				hexRow.WriteString("</div>")
				asciiColumns.Reset()
				pos = 0
				offset++
				continue
			}
			pos++
			offset++
		}
		if pos <= 26 {
			hexRow.WriteString(`<div class="fillByte">&nbsp;&nbsp;</div>`)
			hexRow.WriteString("</div>")
			hexRow.WriteString(`<div class="asciiColumns">`)
			hexRow.WriteString(asciiColumns.String())
			hexRow.WriteString(`</div>`)
			hexRow.WriteString("</div>")
			asciiColumns.Reset()
		}

		hexRow.WriteString(`</div>`)

		b64 := base64.StdEncoding.EncodeToString(bs)

		c.HTML(http.StatusOK, "hex.tmpl", gin.H{
			"filename": filepath.Base(filename),
			"B64":      b64,
			"Bytes":    template.HTML(hexRow.String()),
		})
	})

	if err := r.Run(); err != nil {
		log.Fatal(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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

func ps(b byte) string {
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
	return fmt.Sprintf("%s", []byte{b})
}
