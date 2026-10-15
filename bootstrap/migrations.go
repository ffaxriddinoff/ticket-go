package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ticket-backend/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20261010100619CreatePlansTable{},
		&migrations.M20261010101524CreateUsersTable{},
		&migrations.M20261010151813CreateWatchesTable{},
		&migrations.M20261015114449AddMonitoringToWatchesTable{},
	}
}
