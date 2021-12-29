package server

import (
	"bytes"
	"embed"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roffe/cim/pkg/cim"
)

// embed templates into binary
//go:embed templates/*.tmpl
var tp embed.FS

func Run(enableShutdown bool, prefix string) error {
	r, err := setupRouter(enableShutdown, prefix)
	if err != nil {
		return err
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		openbrowser("http://localhost:8080")
	}()

	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Println("failed to open browser for you:", err)
	}

}

func setupRouter(enableShutdown bool, path string) (*gin.Engine, error) {
	r := gin.Default()
	// Load templates
	if err := loadTemplates(r); err != nil {
		return nil, err
	}

	// Set upload limit for multipart form
	r.MaxMultipartMemory = 1 << 20

	r.GET(p(path, "/"), func(c *gin.Context) { c.HTML(http.StatusOK, "upload.tmpl", nil) })
	r.POST(p(path, "/save"), saveHandler)
	r.POST(p(path, "/"), uploadHandler)
	r.POST(p(path, "/update"), updateHandler)
	r.GET(p(path, "/favicon.ico"), faviconHandler)

	if enableShutdown {
		r.GET("/shutdown", func(c *gin.Context) {
			go func() {
				time.Sleep(300 * time.Millisecond)
				os.Exit(0)
			}()
			c.String(200, "ok")
		})
	}
	return r, nil
}

func p(prefix, path string) string {
	var o strings.Builder
	if prefix != "" {
		o.WriteString(prefix)
	}
	o.WriteString(path)
	return o.String()
}

var bootOrder = []string{"83", "1B", "57", "AF", "C3", "C7", "F3", "FD", "147", "160", "176", "1A2", "1B0", "1B8", "1BF", "1E5", "83", "B9", "DD", "122", "18C", "1A8", "1C6"}

// Load templates from embed fs and add helper funcs to them
func loadTemplates(r *gin.Engine) error {
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
		"bootOrder": func() template.HTML {
			var out strings.Builder
			for _, b := range bootOrder {
				by, err := hex.DecodeString(fmt.Sprintf("%08s", b))
				if err != nil {
					log.Fatal(err)
				}
				u32 := binary.BigEndian.Uint32(by)
				_, err = out.WriteString(
					fmt.Sprintf(`<span data-i="%d" class="field byte-%d">0x%s</span> `, u32, u32, b),
				)
				if err != nil {
					log.Fatal(err)
				}
			}
			return template.HTML(out.String())
		},
	}

	if tmpl, err := template.New("views").Funcs(templateHelpers).ParseFS(tp, "templates/*.tmpl"); err == nil {
		r.SetHTMLTemplate(tmpl)
	} else {
		return err
	}
	return nil
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
	switch a {
	case 0x00:
		return "&centerdot;"
	case 0x20:
		return "&nbsp;"
	case 0xFF:
		return "&fflig;"
	case 0x3c, 0x3e:
		return "˟"
	}
	if a <= 0x20 {
		return "&#9618;"
	}
	if a >= 0x7F {
		return "&block;"
	}

	return fmt.Sprintf("%s", []byte{b})
}
