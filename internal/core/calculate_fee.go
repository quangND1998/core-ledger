package core

import (
	"core-ledger/model/enum"
	model "core-ledger/model/wealify"
	"core-ledger/pkg/constants"
	"errors"
	"strings"
)

func CalculateFee(feeType string, feeValue float64, amount float64) float64 {
	switch feeType {
	case "PERCENT":
		return feeValue * amount
	case "FIXED":
		return feeValue
	default:
		return feeValue
	}
}

func CalculateTopUpFee(tier enum.Tier, platform *model.PlatformFee) (float64, error) {
	platformIncomplete := platform == nil ||
		platform.TopUpDiamondFee == nil ||
		platform.TopUpGoldFee == nil ||
		platform.TopUpSilverFee == nil ||
		platform.TopUpStandardFee == nil
	if platformIncomplete {
		return 0, ErrInvalidPlatform
	}
	switch tier {
	case enum.TierDiamond:
		return *platform.TopUpDiamondFee, nil
	case enum.TierGold:
		return *platform.TopUpGoldFee, nil
	case enum.TierSilver:
		return *platform.TopUpSilverFee, nil
	case enum.TierStandard:
		return *platform.TopUpStandardFee, nil
	default:
		return 0, errors.New("invalid input")
	}
}

func IsOutdatedBIDVBank(vaNumber string) bool {
	return strings.HasPrefix(vaNumber, constants.NotSupportBidvNumberHead)
}
