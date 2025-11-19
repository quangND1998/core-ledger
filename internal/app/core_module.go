package app

import (
	"fmt"
	"core-ledger/internal/logging"
	"core-ledger/pkg/database"
	"core-ledger/pkg/queue"

	"go.uber.org/fx"
	"gorm.io/gorm"
	// "core-ledger/internal/app/core/mail"
	// "core-ledger/internal/app/core/sqs"
)

var CoreModule = fx.Module("core",
	fx.Provide(
		database.Instance,
	),
	// Register logging callbacks after database and dispatcher are available
	fx.Invoke(func(db *gorm.DB, dispatcher queue.Dispatcher) {
		if dispatcher == nil {
			fmt.Println("⚠️  [WARN] Dispatcher is nil, logging callbacks may not work properly")
		} else {
			fmt.Println("✅ [INFO] Registering logging callbacks with dispatcher")
		}
		logging.RegisterCallbacks(db, dispatcher)
	}),
)
