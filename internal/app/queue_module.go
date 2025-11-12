package app

import (
	"context"
	config "core-ledger/configs"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/handlers"
	"fmt"
	"reflect"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

// QueueModule: cung cấp QueueConfig, Worker và đăng ký handler + lifecycle start/stop
var QueueModule = fx.Module("queue",
	fx.Provide(
		// Cấp phát QueueConfig từ env (có validate)
		config.GetQueueConfigWithValidation,
		// Tạo worker theo config
		func(cfg *config.QueueConfig) *queue.Worker {
			return queue.NewWorkerWithRedis(asynq.RedisClientOpt{
				Addr:     cfg.RedisAddr,
				Password: cfg.RedisPassword,
				DB:       cfg.RedisDB,
			}, cfg.Concurrency, cfg.Queues)
		},
		queue.NewDispatcher,
		// Cấp phát handler có DI repo bên trong
		handlers.NewDataProcessHandler,
		handlers.NewMyJobHandler,
		handlers.NewImportCoaAccountHandler,

		fx.Annotate(handlers.NewDataProcessRegistration,
			fx.ResultTags(`group:"queue-registrations"`),
		),
		fx.Annotate(handlers.NewImportCoaAccountHandlerRegistration,
			fx.ResultTags(`group:"queue-registrations"`),
		),
		// Cấp phát registration theo group để dễ mở rộng nhiều job/handler
		fx.Annotate(handlers.NewMyJobHandlerRegistration,
			fx.ResultTags(`group:"queue-registrations"`),
		),
	),
	// Đăng ký routes của worker và khởi chạy theo lifecycle
	fx.Invoke(func(lc fx.Lifecycle, w *queue.Worker, in struct {
		fx.In
		Registrations []queue.Registration `group:"queue-registrations"`
	}) {
		// đăng ký toàn bộ job/handler đã provide vào group
		for _, r := range in.Registrations {
			fmt.Println("Type:", r.Type, "Template:", r.Template, "Handler:", reflect.TypeOf(r.Handler))
			w.RegisterJob(r.Type, r.Template, r.Handler)
		}

		// khởi chạy/dừng worker theo lifecycle
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					_ = w.Start()
				}()
				return nil
			},
			OnStop: func(_ context.Context) error {
				w.Stop()
				return nil
			},
		})
	}),
)
