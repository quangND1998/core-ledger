package wv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func StructToReader(data interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	fmt.Println("body:", string(b))
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
