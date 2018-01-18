package helpers

import (
	"encoding/base64"

	"github.com/kpango/glg"
)

func Encode64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		glg.Fatal(err)
	}
	return b
}
