package util

import "encoding/base64"

func B64Encode(targetStr string) string {
	return base64.StdEncoding.EncodeToString([]byte(targetStr))
}
