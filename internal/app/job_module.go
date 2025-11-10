package app

import (
	"core-ledger/pkg/queue/jobs"

	"go.uber.org/fx"
	// import thêm các repo khác
)

var JobModule = fx.Module("job",
	fx.Provide(
		jobs.NewDataProcessJob,
	),
)
