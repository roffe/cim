package server

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
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

type updateRequest struct {
	Vin             string   `json:"vin"`
	SpsCount        string   `json:"sps_count"`
	Pin             string   `json:"pin"`
	Sas             string   `json:"sas"`
	KeyCount        string   `json:"keycount"`
	Key             []string `json:"key"`
	IskHi           string   `json:"isk_hi"`
	IskLo           string   `json:"isk_lo"`
	Sync            []string `json:"sync"`
	ProgID          []string `json:"prog_id"`
	Snsticker       string   `json:"snsticker"`
	Partno1         string   `json:"partno1"`
	Pnbase1         string   `json:"pnbase1"`
	Pndelphi        string   `json:"pndelphi"`
	Partno          string   `json:"partno"`
	ConfVer         string   `json:"conf_ver"`
	FpDate          string   `json:"fp_date"`
	ProgrammingDate string   `json:"programming_date"`
	File            string   `json:"file_update"`
	Filename        string   `json:"filename"`
}

func updateHandler(c *gin.Context) {
	var u updateRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	b, err := base64.StdEncoding.DecodeString(u.File)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	fw, err := cim.MustLoadBytes(u.Filename, b)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := fw.Pin.Set(u.Pin); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := updateVin(fw, u); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	if u.Sas == "on" {
		fw.SetSasOpt(true)
	} else {
		fw.SetSasOpt(false)
	}

	if err := updateKeys(fw, u); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid key value: %v", err))
		return
	}

	if err := updateSync(fw, u); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid sync data: %v", err))
		return
	}

	for i, s := range u.ProgID {
		if err := fw.SetProgrammingID(i, s); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("invalid programming id %d value: %s: %v", i, s, err))
			return
		}
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
	//c.String(200, hexRows)
	/*
		sections := generateSections(fw)
		styles := generateStyles(sections)
		jsSections := jsSections(sections)

		c.HTML(http.StatusOK, "view.tmpl", gin.H{
			"filename": filepath.Base(u.Filename),
			"fw":       fw,
			"B64":      base64.StdEncoding.EncodeToString(fwBytes),
			"Hexview":  template.HTML(hexRows),
			"sections": template.JS(jsSections),
			"styles":   styles,
		})
	*/

	c.JSON(http.StatusOK, gin.H{
		"B64":     base64.StdEncoding.EncodeToString(fwBytes),
		"hexview": hexRows,
	})
}

func updateVin(fw *cim.Bin, u updateRequest) error {
	if err := fw.Vin.Set(u.Vin); err != nil {
		return fmt.Errorf("filed to set vin: %v", err)
	}
	if n, err := strconv.ParseUint(u.SpsCount, 0, 8); err == nil {
		fw.Vin.SetSpsCount(uint8(n))
	} else {
		return fmt.Errorf("failed to parse sps count: %q %s", u.SpsCount, err.Error())
	}
	return nil
}

func updateSync(fw *cim.Bin, u updateRequest) error {
	for i, opt := range u.Sync {
		syncData, err := hex.DecodeString(opt)
		if err != nil {
			return fmt.Errorf("failed to decode sync data %d: %v", i, err)
		}
		fw.Sync.SetData(uint8(i), syncData)

	}
	return nil
}

func updateKeys(fw *cim.Bin, u updateRequest) error {

	hi, err := hex.DecodeString(u.IskHi)
	if err != nil {
		return fmt.Errorf("failed to parse ISK High: %X: %s", u.IskHi, err.Error())
	}
	lo, err := hex.DecodeString(u.IskLo)
	if err != nil {
		return fmt.Errorf("failed to parse ISK Low: %X: %s", u.IskHi, err.Error())
	}

	if err := fw.Keys.SetIsk(hi, lo); err != nil {
		return fmt.Errorf("failed to set ISK: %v", err)
	}

	if n, err := strconv.ParseUint(u.KeyCount, 0, 8); err == nil {
		fw.Keys.SetKeyCount(uint8(n))
	} else {
		return fmt.Errorf("failed to parse key count: %q %s", u.KeyCount, err.Error())
	}

	for i, k := range u.Key {
		b, err := hex.DecodeString(k)
		if err != nil {
			return err
		}
		if err := fw.Keys.SetKey(uint8(i), b); err != nil {
			return err
		}
	}
	return nil
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
