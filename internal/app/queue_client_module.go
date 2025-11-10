package app

import (
	config "core-ledger/configs"
	"core-ledger/pkg/queue"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

// QueueClientModule: dành cho app producer (HTTP/API) để gửi job
var QueueClientModule = fx.Module("queue-client",
	fx.Provide(
		config.GetQueueConfigWithValidation,
		// asynq.Client
		func(cfg *config.QueueConfig) (*asynq.Client, error) {
			client := asynq.NewClient(asynq.RedisClientOpt{
				Addr:     cfg.RedisAddr,
				Password: cfg.RedisPassword,
				DB:       cfg.RedisDB,
			})
			return client, nil
		},
		// Dispatcher abstraction
		queue.NewDispatcher,
	),
)


