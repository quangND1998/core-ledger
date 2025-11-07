package utils

import (
	"core-ledger/internal/module/constants"
	"core-ledger/model/enum"
	"core-ledger/pkg/repo"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GenerateWealifyID(index int64, prefix string) string {
	if prefix == "" {
		prefix = "W"
	}

	// Format ngày: MMDDYYYY (như toLocaleDateString('en-US').replace(/\//g, ''))
	now := time.Now()
	datePart := fmt.Sprintf("%02d%02d%04d", now.Month(), now.Day(), now.Year())

	// Lấy timestamp cuối 2 số
	timestamp := now.UnixMilli()
	timestampLast2 := fmt.Sprintf("%02d", timestamp%100)

	// Format index
	indexStr := fmt.Sprintf("%02d", index)

	// Ghép lại và viết hoa
	return strings.ToUpper(fmt.Sprintf("%s%s%s%s", prefix, datePart, indexStr, timestampLast2))
}

func GenerateWealifyIDVer2(prefix string) string {
	if prefix == "" {
		prefix = "W"
	}

	now := time.Now()

	// Get MMDDYY format
	month := fmt.Sprintf("%02d", now.Month())
	day := fmt.Sprintf("%02d", now.Day())
	year := fmt.Sprintf("%02d", now.Year()%100) // Last 2 digits

	// Get timestamp and last 6 digits
	timestamp := fmt.Sprintf("%d", now.UnixNano()/1e6) // Milliseconds
	lastSix := ""
	if len(timestamp) >= 6 {
		lastSix = timestamp[len(timestamp)-6:]
	} else {
		lastSix = timestamp
	}

	// Combine all parts
	return fmt.Sprintf("%s%s%s%s%s", prefix, month, day, year, lastSix)
}

func GenerateWealifyIDVer3(index int64, prefix string) string {
	if prefix == "" {
		prefix = "W"
	}
	now := time.Now()
	random := rand.Intn(1000) // 0–999
	return fmt.Sprintf("%s%s%d%d%03d", prefix, now.Format("20060102150405"), now.Nanosecond()/1e6, index, random)
}

func GenerateTransactionCode(tx *gorm.DB, transactionType enum.TransactionType) (string, error) {
	transactionRepo := repo.NewTransactionRepo(tx)
	var (
		index  int64
		prefix string
		err    error
	)
	if transactionType == enum.TransactionTypeTopUp {
		index, err = transactionRepo.Count(map[string]interface{}{
			"transaction_type": enum.TransactionTypeTopUp,
		})
		if err != nil {
			return "", errors.New("_generateId top up error")
		}
		prefix = "T"
	} else if transactionType == enum.TransactionTypeWithdrawal {
		index, err = transactionRepo.Count(map[string]interface{}{
			"transaction_type": enum.TransactionTypeWithdrawal,
		})
		if err != nil {
			return "", errors.New("_generateId withdraw error")
		}
		prefix = "W"
	} else {
		return "", errors.New("invalid transaction type")
	}
	return GenerateWealifyIDVer3(index, prefix), nil
}

func FormatAmount(amount float64, decimalPlaces int, groupDigits int) string {
	if decimalPlaces < 0 {
		decimalPlaces = 0
	}
	if groupDigits <= 0 {
		groupDigits = 3
	}

	isNegative := amount < 0
	if isNegative {
		amount = -amount
	}

	formatted := fmt.Sprintf("%.*f", decimalPlaces, amount)
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	length := len(integerPart)
	if length <= groupDigits {
		if decimalPart != "" {
			return integerPart + "." + decimalPart
		}
		return integerPart
	}

	var result []rune
	for i, char := range integerPart {
		if i > 0 && (length-i)%groupDigits == 0 {
			result = append(result, ',')
		}
		result = append(result, char)
	}
	finalResult := string(result)
	if decimalPart != "" {
		finalResult += "." + decimalPart
	}

	// for isNegative
	if isNegative {
		finalResult = "-" + finalResult
	}

	return finalResult
}

// AreApproximatelyEqual checks if two float64 numbers are approximately equal
func AreApproximatelyEqual(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}

// TruncateAmount truncates a number to the specified number of decimal places
func TruncateAmount(number float64, digits int) float64 {
	if digits < 0 {
		digits = 2 // default value
	}

	// Check if number is approximately equal to its rounded value
	if AreApproximatelyEqual(number, math.Round(number)) {
		return math.Round(number)
	}

	// Convert to string and split by decimal point
	numStr := fmt.Sprintf("%.15f", number) // Use high precision to avoid rounding
	parts := strings.Split(numStr, ".")

	if len(parts) == 1 {
		// No decimal part
		return number
	}

	integerPart := parts[0]
	decimalPart := parts[1]

	// Truncate decimal part to specified digits
	if len(decimalPart) > digits {
		decimalPart = decimalPart[:digits]
	}

	// Reconstruct the number
	if digits == 0 || decimalPart == "" {
		result, _ := strconv.ParseFloat(integerPart, 64)
		return result
	}

	truncatedStr := integerPart + "." + decimalPart
	result, _ := strconv.ParseFloat(truncatedStr, 64)
	return result
}

// TruncateAmountOnCurrency truncates amount based on currency type
func TruncateAmountOnCurrency(amount float64, currency string) float64 {
	var fixedPlaces int

	switch strings.ToUpper(currency) {
	case constants.VND_CODE, constants.SYSTEM_UNIT:
		fixedPlaces = 0
	default:
		fixedPlaces = 2
	}

	return TruncateAmount(amount, fixedPlaces)
}

func FormatTruncateAmountOnCurrency(amount float64, currency string, decimalPlaces int, groupDigits int) string {
	// Set default values
	if decimalPlaces < 0 {
		decimalPlaces = 2
	}
	if groupDigits <= 0 {
		groupDigits = 3
	}

	// Adjust decimal places based on currency
	switch strings.ToUpper(currency) {
	case constants.VND_CODE, constants.SYSTEM_UNIT:
		decimalPlaces = 0
	default:
		// Keep the provided decimalPlaces value
	}

	// First truncate the amount based on currency, then format it
	truncatedAmount := TruncateAmountOnCurrency(amount, currency)
	return FormatAmount(truncatedAmount, decimalPlaces, groupDigits)
}

func ChunkAmount(amount float64, minimumAmount float64, maximumAmount float64) []float64 {
	// Set default values if not provided (Go doesn't support default parameters)
	if minimumAmount <= 0 {
		minimumAmount = 50000
	}
	if maximumAmount <= 0 {
		maximumAmount = 290000000
	}

	if amount < maximumAmount {
		return []float64{amount}
	}

	var result []float64

	// Split amount into maximum chunks
	for amount > maximumAmount {
		amount = amount - maximumAmount
		result = append(result, maximumAmount)
	}

	// Handle remaining amount
	if amount >= minimumAmount {
		result = append(result, amount)
	} else {
		offset := minimumAmount - amount
		// Remove the last element
		if len(result) > 0 {
			result = result[:len(result)-1]
		}
		result = append(result, maximumAmount-offset)
		result = append(result, minimumAmount)
	}

	return result
}

// FormatAmountOnCurrency formats amount based on currency type
func FormatAmountOnCurrency(amount float64, currency string) string {
	var decimalPlaces int

	switch strings.ToUpper(currency) {
	case constants.VND_CODE, constants.SYSTEM_UNIT:
		decimalPlaces = 0
	default:
		decimalPlaces = 2
	}

	return FormatAmount(amount, decimalPlaces, 3)
}

// AmountStr returns formatted amount string based on transaction type
func GetAmountStr(amount float64, transactionType, currencySymbol string) string {
	switch enum.TransactionType(transactionType) {
	case enum.TransactionTypeWithdrawal, enum.TransactionTypeInternal:
		return FormatTruncateAmountOnCurrency(amount, constants.SYSTEM_UNIT, 2, 3)
	case enum.TransactionTypeTopUp:
		if currencySymbol == "" {
			currencySymbol = constants.SYSTEM_UNIT // default
		}
		return FormatTruncateAmountOnCurrency(amount, currencySymbol, 2, 3)
	default:
		return FormatTruncateAmountOnCurrency(amount, constants.SYSTEM_UNIT, 2, 3)
	}
}

// PaymentAmountStr returns formatted payment amount string based on transaction type
func GetPaymentAmountStr(amount float64, transactionType, currencySymbol string) string {
	switch enum.TransactionType(transactionType) {
	case enum.TransactionTypeWithdrawal, enum.TransactionTypeInternal:
		return FormatAmountOnCurrency(amount, constants.SYSTEM_UNIT)
	case enum.TransactionTypeTopUp:
		if currencySymbol == "" {
			currencySymbol = constants.SYSTEM_UNIT // default
		}
		return FormatAmountOnCurrency(amount, currencySymbol)
	default:
		return FormatAmountOnCurrency(amount, constants.SYSTEM_UNIT)
	}
}
