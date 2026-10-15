package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ticket-backend/app/facades"
)

type M20261010151813CreateWatchesTable struct{}

// Signature The unique signature for the migration.
func (r *M20261010151813CreateWatchesTable) Signature() string {
	return "20261010151813_create_watches_table"
}

// Up Run the migrations.
func (r *M20261010151813CreateWatchesTable) Up() error {
	if !facades.Schema().HasTable("watches") {
		return facades.Schema().Create("watches", func(table schema.Blueprint) {
			table.ID()

			table.BigInteger("user_id")
			table.Foreign("user_id").References("id").On("users")
			table.Index("user_id")

			table.String("from", 30)
			table.String("to", 30)
			table.Date("trip_date")
			table.Jsonb("filters")

			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261010151813CreateWatchesTable) Down() error {
	return facades.Schema().DropIfExists("watches")
}
