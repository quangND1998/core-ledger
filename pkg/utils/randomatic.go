package utils

import (
	"crypto/rand"

	"fmt"
	"math/big"
	"strings"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Randomatic(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func generateRand(index int, prefix string) string {
	if index == 0 {
		index = 1
	}
	if prefix == "" {
		prefix = "W"
	}

	now := time.Now()
	time := now.Format("01022006") // MM/DD/YYYY format without slashes
	timestamp := now.UnixMilli()

	var indexString string
	if index < 10 {
		indexString = fmt.Sprintf("0%d", index)
	} else {
		indexString = fmt.Sprintf("%d", index)
	}

	timestampStr := fmt.Sprintf("%d", timestamp)
	lastTwo := timestampStr[len(timestampStr)-2:]

	result := fmt.Sprintf("%s%s%s%s", strings.ToUpper(prefix), time, indexString, lastTwo)
	return strings.ToUpper(result)
}

func GenerateDefaultPassword() string {
	return ""
}
func GenerateCustomerID() string {
	return generateRand(10, "C")
}
