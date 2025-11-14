package helper

import (
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

func FormatJSONForExcel(value any) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case map[string]any:
		bytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("error: %v", err)
		}
		return string(bytes)
	case *datatypes.JSON:
		if v == nil {
			return ""
		}
		// Chuyển datatypes.JSON thành string
		return string(*v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
