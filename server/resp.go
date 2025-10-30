package server

import (
	"bytes"
	"io"

	"github.com/tidwall/resp"
)

func parseResp(raw string) ([]string, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	var result []string
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		if v.Type() == resp.Array {
			for _, v := range v.Array() {
				result = append(result, v.String())
			}
		} else {
			result = append(result, v.String())
		}
	}
	return result, nil
}
