package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strconv"
	"time"

	internal "go-web-scrapper/framework/error"
)

const (
	Limit   = "limit"
	Page    = "page"
	Keyword = "keyword"
)

func Encrypt(s string) string {
	passwordBytes := []byte(s)
	hash := sha256.Sum256(passwordBytes)
	res := hex.EncodeToString(hash[:])
	return res
}

func FormatTime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}

func GetLimit(v string) (int, *internal.Error) {
	if v == "" {
		return 10, nil
	} else {
		res, err := strconv.Atoi(v)
		if err != nil {
			return res, internal.NewError(500, err)
		}
		return res, nil
	}
}

func GetTargetPage(v string) (int, *internal.Error) {
	if v == "" {
		return 0, nil
	} else {
		res, err := strconv.Atoi(v)
		if err != nil {
			return res, internal.NewError(500, err)
		}
		return res, nil
	}
}

func GetPage(page int) int {
	if page > 0 {
		return page
	} else {
		return 1
	}
}

func GetTotalPage(limit int, total int) int {
	if total > 0 && limit > 0 {
		res := int(math.Ceil(float64(total) / float64(limit)))
		return res
	} else {
		return 1
	}
}
