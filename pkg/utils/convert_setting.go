package utils

import (
	"core-ledger/model/jsonfield"
	"encoding/json"

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
