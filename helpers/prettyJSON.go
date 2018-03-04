package helpers

import (
	"bytes"
	"encoding/json"
)

//PrettyJson indents json byte
func PrettyJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}
