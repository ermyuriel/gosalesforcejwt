package gosalesforcejwt

import (
	"encoding/base64"
	"strings"
)

func EncodeBase64URL(data []byte) string {

	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")

}
