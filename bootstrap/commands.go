package bootstrap

import (
	"github.com/goravel/framework/contracts/console"

	"ticket-backend/app/console/commands"
)

func Commands() []console.Command {
	return []console.Command{
		&commands.CheckWatch{},
		&commands.CheckWatches{},
	}
}
