package core

import (
	"core-ledger/model/enum"
	model "core-ledger/model/wealify"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type StatsTransactionInfo struct {
	CustomerID        int64   `json:"customer_id"`
	TransactionType   string  `json:"transaction_type"`
	CurrencySymbol    string  `json:"currency_symbol"`
	TransactionStatus string  `json:"transaction_status"`
	TotalAmount       float64 `json:"total_amount"`
	TotalReceived     float64 `json:"total_received"`
	TotalFee          float64 `json:"total_fee"`
}

const (
	calculateTotalAmountBeforeFeeQuery = `
    COALESCE(
        SUM(
            CASE 
                WHEN t.transaction_type = 'TOP_UP' 
                    THEN t.amount * CAST(JSON_UNQUOTE(JSON_EXTRACT(t.rate, '$.value')) AS DECIMAL(18,4))
                WHEN t.transaction_type = 'WITHDRAWAL' 
                    THEN t.amount
				WHEN t.transaction_type = 'ADJUSTMENT' 
                    THEN t.amount
                ELSE 0
            END
        ), 
        0
    ) AS total_amount`

	calculateTotalReceivedAfterFeeAndRateQuery = `
	COALESCE(
          SUM(
            CASE
              WHEN JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.type')) = 'FIXED' THEN 
                (t.amount - CAST(JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.value')) AS DECIMAL(18,4)))
                * CAST(JSON_UNQUOTE(JSON_EXTRACT(t.rate, '$.value')) AS DECIMAL(18,4))

              WHEN JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.type')) = 'PERCENT' AND t.transaction_type = 'TOP_UP' THEN 
                (t.amount - (t.amount * CAST(JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.value')) AS DECIMAL(18,4))))
                * CAST(JSON_UNQUOTE(JSON_EXTRACT(t.rate, '$.value')) AS DECIMAL(18,4))

              ELSE 
                t.amount
            END
          ), 
        0) AS total_received`
	calculateTotalFeeQuery = `
        COALESCE(
          SUM(
            CASE
              WHEN JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.type')) = 'FIXED' THEN 
                CAST(JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.value')) AS DECIMAL(18,4))
                * CAST(JSON_UNQUOTE(JSON_EXTRACT(t.rate, '$.value')) AS DECIMAL(18,4))

              WHEN JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.type')) = 'PERCENT' AND t.transaction_type = 'TOP_UP' THEN 
                (t.amount * CAST(JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.value')) AS DECIMAL(18,4)))
                * CAST(JSON_UNQUOTE(JSON_EXTRACT(t.rate, '$.value')) AS DECIMAL(18,4))

              WHEN JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.type')) = 'PERCENT' AND t.transaction_type = 'WITHDRAWAL' THEN 
                (t.amount * CAST(JSON_UNQUOTE(JSON_EXTRACT(t.fee, '$.value')) AS DECIMAL(18,4)))

              ELSE 
                0
            END
          ), 
        0) AS total_fee`
)

func CalculateStatsBalanceResult(infoRecords []StatsTransactionInfo) StatsBalanceResult {
	var result StatsBalanceResult
	for _, info := range infoRecords {
		switch info.TransactionType {
		case "TOP_UP":
			if info.TransactionStatus == "APPROVED" {
				result.TotalTopUpSuccessBeforeFee += info.TotalAmount
				result.TotalReceived += info.TotalReceived
				result.TotalFee += info.TotalFee
			}
			if info.TransactionStatus == "PENDING" {
				result.TotalTopUpPending += info.TotalAmount
			}
		case "WITHDRAWAL":
			if info.TransactionStatus == "PENDING" {
				result.TotalWithdrawalPending += info.TotalAmount
			}
			if info.TransactionStatus == "APPROVED" {
				result.TotalWithdrawalSuccess += info.TotalAmount
			}
			result.TotalWithdrawal += info.TotalAmount
		case enum.TransactionTypeAdjustment.String():
			result.TotalTopUpSuccessBeforeFee += info.TotalAmount
			result.TotalReceived += info.TotalReceived
		}
	}
	result.DisplayBalance = result.TotalReceived - result.TotalWithdrawal
	return result
}

func GetStatsTransactionInfo(db *gorm.DB,
	customerIDs []int64,
	walletType string,
	transactionStatuses []string,
	from, to string) (map[int64][]StatsTransactionInfo, error) {
	var rows []StatsTransactionInfo
	query := db.
		Table("transactions AS t").Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)}).
		Select("t.customer_id, t.transaction_type, t.currency_symbol, c.code AS currency_code, t.transaction_status", calculateTotalAmountBeforeFeeQuery, calculateTotalReceivedAfterFeeAndRateQuery, calculateTotalFeeQuery).
		Joins("INNER JOIN currencies c  ON c.symbol  = t.currency_symbol").
		Joins("LEFT JOIN wallets sent_wallet on t.sent_wallet_id = sent_wallet.id").
		Joins("LEFT JOIN wallets received_wallet on t.received_wallet_id = received_wallet.id").
		Where("sent_wallet.`type` = ? OR received_wallet.`type`  = ?", walletType, walletType).
		Where("t.customer_id IN (?)", customerIDs).
		Where("t.transaction_status IN (?)", transactionStatuses).
		Group("t.customer_id, t.transaction_type, c.symbol, c.code, t.transaction_status")
	if from != "" {
		query.Where("t.created_at >= ?", from)
	}
	if to != "" {
		query.Where("t.created_at < ?", to)
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[int64][]StatsTransactionInfo)

	for i, v := range rows {
		result[v.CustomerID] = append(result[v.CustomerID], rows[i])
	}
	return result, nil
}

type StatsBalanceResult struct {
	TotalTopUpSuccessBeforeFee float64 //(only calculate top_up has been APPROVED)
	TotalReceived              float64 //(received = total success - total fee
	TotalWithdrawalSuccess     float64 //only calculate withdrawal has been APPROVED
	TotalWithdrawal            float64 //total withdrawal = SUM (PENDING withdrawal + APPROVED withdrawal)
	DisplayBalance             float64 //display balance = Total Received - Total Withdrawal
	TotalFee                   float64 //(total fee = TotalTopUpSuccessBeforeFee - TotalReceived
	TotalWithdrawalPending     float64 //only pending withdrawal tx
	TotalTopUpPending          float64 //only pending top up tx
}

type StatsVirtualAccountInfo struct {
	TotalCreated    int
	TotalActive     int
	TotalInactive   int
	TotalRestricted int
}

func GetStatsVirtualAccountInfo(db *gorm.DB, customerIDs []int64) (map[int64]StatsVirtualAccountInfo, error) {
	type statsVAInfo struct {
		CustomerID int64  `json:"customer_id"`
		Status     string `json:"status"`
		Total      int    `json:"total"`
	}
	var rows []statsVAInfo
	err := db.
		Table("`virtual-accounts` AS va").
		Select("va.customer_id, va.status, COUNT(va.status) AS total").
		Where("va.customer_id IN ?", customerIDs).
		Group("va.status, va.customer_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	dict := make(map[int64][]statsVAInfo)
	for i, row := range rows {
		dict[row.CustomerID] = append(dict[row.CustomerID], rows[i])
	}
	result := make(map[int64]StatsVirtualAccountInfo)

	for customerID, stats := range dict {
		var statsInfo StatsVirtualAccountInfo
		for _, v := range stats {
			switch v.Status {
			case string(model.VAStatusActive):
				statsInfo.TotalActive = v.Total
			case string(model.VAStatusInactive):
				statsInfo.TotalInactive = v.Total
			case string(model.VAStatusRestricted):
				statsInfo.TotalRestricted = v.Total
			}
			statsInfo.TotalCreated += v.Total
		}
		result[customerID] = statsInfo
	}
	return result, nil
}
