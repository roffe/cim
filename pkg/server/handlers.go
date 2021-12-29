package server

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roffe/cim/pkg/cim"
)

// embed favicon.ico
//go:embed favicon.ico
var favicon []byte

func faviconHandler(c *gin.Context) {
	if _, err := c.Writer.Write(favicon); err != nil {
		c.String(http.StatusInternalServerError, "failed to load favicon.ico")
		return
	}
	c.Status(200)
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
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := fw.Validate(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fwBytes, err := fw.Bytes()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	hexRows, err := buildHexview(fw)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	sections := generateSections(fw)
	styles := generateStyles(sections)
	jsSections := jsSections(sections)

	c.HTML(http.StatusOK, "view.tmpl", gin.H{
		"filename": filepath.Base(filename),
		"fw":       fw,
		"B64":      base64.StdEncoding.EncodeToString(fwBytes),
		"Hexview":  template.HTML(hexRows),
		"sections": template.JS(jsSections),
		"styles":   styles,
	})
}

func jsSections(sections []Section) string {
	js := strings.Builder{}
	js.WriteString(`var sections = [`)
	for i, s := range sections {
		js.WriteString(s.String())
		if i == len(sections)-1 {
			break
		}
		js.WriteString(",\n")
	}
	js.WriteString(`]`)
	return js.String()
}

func buildHexview(fw *cim.Bin) (string, error) {
	fwBytes, err := fw.Bytes()
	if err != nil {
		return "", err
	}
	hexRows := strings.Builder{}
	asciiColumns := strings.Builder{}

	pos := 0
	offset := 0
	width := 36

	for _, bb := range fwBytes {
		if pos == 0 {
			hexRows.WriteString(fmt.Sprintf(`<div class="hexRow"><div class="addrColumn"><b>%03X</b></div><div class="hexColumns">`, offset))
		}
		hexRows.WriteString(fmt.Sprintf(`<div class="hexByte byte-%d" data-i="%d">%02X</div>`+"\n", offset, offset, bb))
		asciiColumns.WriteString(fmt.Sprintf(`<div class="asciiByte byte-%d" data-i="%d">%s</div>`+"\n", offset, offset, psafe(bb)))
		if pos == width {
			hexRows.WriteString(`</div><div class="asciiColumns">` + "\n" + asciiColumns.String() + "</div></div>")
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
			hexRows.WriteString(`<div class="fillByte">&nbsp;&nbsp;</div>`)
		}
		hexRows.WriteString(`</div><div class="asciiColumns">` + "\n" + asciiColumns.String() + "</div></div>\n")
		asciiColumns.Reset()
	}

	hexRows.WriteString("</div>")
	return hexRows.String(), nil
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

	fw, err := cim.MustLoadBytes(filename, b)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	bs, err := fw.XORBytes()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	contentLength := int64(len(bs))
	contentType := "application/octet-stream"

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + filepath.Base(filename) + `"`,
	}
	c.DataFromReader(http.StatusOK, contentLength, contentType, bytes.NewReader(bs), extraHeaders)
}
