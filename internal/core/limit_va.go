package core

import "gorm.io/gorm"

func GetVALimitationByCustomer(db *gorm.DB, customerID int64) (effectiveLimit int, baseLimit int) {
	//get baseLimit

	return 0, 0
}

func GetVACapLimitationByCustomer(db *gorm.DB, customerID int64) (isPassCap bool, err error) {
	//getCapLimit

	return true, nil
}
