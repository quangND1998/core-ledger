package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type RateField struct {
	Type   RateValueType `json:"type"`
	Amount float64       `json:"value"`
}

func (r *RateField) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

func (r *RateField) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("RateField: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, r)
}
