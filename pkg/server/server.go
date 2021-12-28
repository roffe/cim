package server

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roffe/cim/pkg/cim"
)

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

// embed templates into binary
//go:embed templates/*.tmpl
var tp embed.FS

// embed favicon.ico
//go:embed favicon.ico
var favicon []byte

func Run() error {
	r := gin.Default()

	// Load templates
	if tmpl, err := template.New("views").Funcs(templateHelpers).ParseFS(tp, "templates/*.tmpl"); err == nil {
		r.SetHTMLTemplate(tmpl)
	} else {
		return err
	}

	// Set upload limit for multipart form
	r.MaxMultipartMemory = 1 << 20

	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "upload.tmpl", nil) })
	r.POST("/save", saveHandler)
	r.POST("/", uploadHandler)
	r.POST("/update", updateHandler)
	r.GET("/favicon.ico", func(c *gin.Context) {
		if _, err := c.Writer.Write(favicon); err != nil {
			c.String(http.StatusInternalServerError, "failed to load favicon.ico")
			return
		}
		c.Status(200)
	})

	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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
