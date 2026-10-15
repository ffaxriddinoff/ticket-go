package commands

import (
	"fmt"
	"strconv"
	"ticket-backend/app/services"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type CheckWatch struct {
}

// Signature The name and signature of the console command.
func (r *CheckWatch) Signature() string {
	return "watch:check"
}

// Description The console command description.
func (r *CheckWatch) Description() string {
	return "Check a single active watch"
}

// Extend The console command extend.
func (r *CheckWatch) Extend() command.Extend {
	return command.Extend{Category: "watch"}
}

// Handle Execute the console command.
func (r *CheckWatch) Handle(ctx console.Context) error {
	watchID, err := strconv.ParseUint(
		ctx.Argument(0),
		10,
		strconv.IntSize,
	)
	if err != nil || watchID == 0 {
		return fmt.Errorf(
			"musbat watch ID kiriting: ./artisan watch:check 1",
		)
	}

	ctx.Info(fmt.Sprintf(
		"Watch %d: tekshiruv boshlandi.",
		watchID,
	))

	monitor := services.NewWatchMonitorService()

	if err := monitor.CheckWatch(ctx, uint(watchID)); err != nil {
		return err
	}

	ctx.Success("Tekshiruv yakunlandi.")

	return nil
}
