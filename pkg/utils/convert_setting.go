package utils

import (
	"core-ledger/model/jsonfield"
	"encoding/json"
	"strconv"

	"gorm.io/datatypes"
)

func GetEarlyApproveVATopUpCustomerIDs(data datatypes.JSON) ([]int64, error) {
	var parsed []jsonfield.VaEarlyApproveTopUpSettingDataItem
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	var res []int64
	for _, datum := range parsed {
		res = append(res, datum.ID)
	}
	return res, nil
}

func ParseIntIdParam(key string) (int64, error) {
	id, err := strconv.ParseInt(key, 10, 64)
	return id, err
}
