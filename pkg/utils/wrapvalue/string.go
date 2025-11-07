package wv

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func StringToFloat64(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return 0
	}

	return f
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return 0
	}

	return i
}

func StringToInt64(str string) int64 {
	i64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int64:", err)
		return 0
	}
	return i64
}

func ToBytes(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
