package ruleValue

import (
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"

	"gorm.io/gorm"
)

type RuleCateogySerive struct {
	db         *gorm.DB
	
	logger     logger.CustomLogger
	dispatcher queue.Dispatcher
}

func NewRuleCateogySerive(dispatcher queue.Dispatcher, db *gorm.DB) *RuleCateogySerive {
	return &RuleCateogySerive{
		db:         db,
		logger:     logger.NewSystemLog("RuleCateogySerive"),
		dispatcher: dispatcher,
	}
}
