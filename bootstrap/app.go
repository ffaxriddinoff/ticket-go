package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"

	"ticket-backend/config"
	"ticket-backend/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithCommands(Commands).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Web()
			routes.Grpc()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
