package helper

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func GenerateSecureNumber() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	num := binary.BigEndian.Uint64(b)
	return fmt.Sprintf("%d", num)
}
