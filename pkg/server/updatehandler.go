package server

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/roffe/cim/pkg/cim"
)

type updateRequest struct {
	Vin             string   `json:"vin"`
	VinValue        string   `json:"vin_value"`
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
	Partno1Rev      string   `json:"partno1rev"`
	Pnbase1         string   `json:"pnbase1"`
	Pnbase1Rev      string   `json:"pnbase1rev"`
	Pndelphi        string   `json:"pndelphi"`
	Partno          string   `json:"partno"`
	ConfVer         string   `json:"conf_ver"`
	FpDate          string   `json:"fp_date"`
	ProgrammingDate string   `json:"programming_date"`
	PSKHi           string   `json:"psk_hi"`
	PSKLo           string   `json:"psk_lo"`
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
		return
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

	confVer, err := strconv.ParseUint(u.ConfVer, 0, 32)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	fw.SetConfVer(uint32(confVer))

	for i, s := range u.ProgID {
		if err := fw.SetProgrammingID(i, s); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("invalid programming id %d value: %s: %v", i, s, err))
			return
		}
	}
	if err := updatePSK(fw, u); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid PSK data: %v", err))
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

	c.JSON(http.StatusOK, gin.H{
		"md5":     fw.MD5(),
		"crc32":   fw.CRC32(),
		"B64":     base64.StdEncoding.EncodeToString(fwBytes),
		"hexview": hexRows,
	})
}

func updatePSK(fw *cim.Bin, u updateRequest) error {
	lo, err := hex.DecodeString(u.PSKLo)
	if err != nil {
		return fmt.Errorf("failed to decode PSK Low: %v", err)
	}
	hi, err := hex.DecodeString(u.PSKHi)
	if err != nil {
		return fmt.Errorf("failed to decode PSK High: %v", err)
	}
	if err := fw.PSK.SetHigh(hi); err != nil {
		return fmt.Errorf("invalid PSK high data: %v", err)
	}
	if err := fw.PSK.SetLow(lo); err != nil {
		return fmt.Errorf("invalid PSK low data: %v", err)
	}
	return nil
}

func updateVin(fw *cim.Bin, u updateRequest) error {
	if err := fw.Vin.Set(u.Vin); err != nil {
		return fmt.Errorf("failed to set vin: %v", err)
	}

	if n, err := strconv.ParseUint(u.VinValue, 0, 8); err == nil {
		fw.Vin.SetValue(uint8(n))
	} else {
		return fmt.Errorf("failed to parse vin value: %q %s", u.VinValue, err.Error())
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
