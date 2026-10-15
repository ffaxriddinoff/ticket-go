package commands

import (
	"ticket-backend/app/services"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type CheckWatches struct {
}

// Signature The name and signature of the console command.
func (r *CheckWatches) Signature() string {
	return "check:watches"
}

// Description The console command description.
func (r *CheckWatches) Description() string {
	return "Check all active watches"
}

// Extend The console command extend.
func (r *CheckWatches) Extend() command.Extend {
	return command.Extend{Category: "watch"}
}

// Handle Execute the console command.
func (r *CheckWatches) Handle(ctx console.Context) error {
	monitor := services.NewWatchMonitorService()

	if err := monitor.CheckActiveWatches(ctx); err != nil {
		return err
	}

	ctx.Success("Faol watchlar tekshiruvi yakunlandi.")

	return nil
}
