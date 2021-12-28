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
	VinOpt             string   `json:"vin_opt"`
	PinOpt             string   `json:"pin_opt"`
	SasOpt             string   `json:"sas_opt"`
	KeyCountOpt        string   `json:"keycount_opt"`
	KeyOpt             []string `json:"key_opt"`
	IskHiOpt           string   `json:"isk_hi_opt"`
	IskLoOpt           string   `json:"isk_lo_opt"`
	SyncOpt            []string `json:"sync_opt"`
	ProgIDOpt          string   `json:"prog_id_opt"`
	SnstickerOpt       string   `json:"snsticker_opt"`
	Partno1Opt         string   `json:"partno1_opt"`
	Pnbase1Opt         string   `json:"pnbase1_opt"`
	PndelphiOpt        string   `json:"pndelphi_opt"`
	PartnoOpt          string   `json:"partno_opt"`
	ConfVerOpt         string   `json:"conf_ver_opt"`
	FpDateOpt          string   `json:"fp_date_opt"`
	ProgrammingDateOpt string   `json:"programming_date_opt"`
	File               string   `json:"file"`
	Filename           string   `json:"filename"`
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

	fw.Vin.Set(u.VinOpt)
	fw.Pin.Set(u.PinOpt)

	if u.SasOpt == "on" {
		fw.SetSasOpt(true)
	} else {
		fw.SetSasOpt(false)
	}

	if err := updateKeys(fw, u); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid key value: %v", err))
		return
	}

	fw.Dump()
	c.String(200, "ok, not implemented yet")
}

func updateKeys(fw *cim.Bin, u updateRequest) error {
	for i, k := range u.KeyOpt {
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

	sections := generateSections(fw)
	styles := generateStyles(sections)

	out := strings.Builder{}
	out.WriteString(`var sections = [`)
	for i, s := range sections {
		out.WriteString(s.String())
		if i == len(sections)-1 {
			break
		}
		out.WriteString(",\n")
	}
	out.WriteString(`]`)

	c.HTML(http.StatusOK, "view.tmpl", gin.H{
		"filename": filepath.Base(filename),
		"fw":       fw,
		"B64":      base64.StdEncoding.EncodeToString(fwBytes),
		"Hexview":  template.HTML(hexRows.String()),
		"sections": template.JS(out.String()),
		"styles":   template.CSS(styles),
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
