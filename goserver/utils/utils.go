package utils

import (
	"encoding/base64"
)

func EncodeB64(frame []byte) string {
	return base64.StdEncoding.EncodeToString(frame)
}
