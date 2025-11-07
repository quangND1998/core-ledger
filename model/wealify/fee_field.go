package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type FeeField struct {
	Type   FeeType `json:"type"`
	Amount float64 `json:"value"`
}

func (f *FeeField) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f)
}

func (f *FeeField) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("FeeField: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, f)
}
