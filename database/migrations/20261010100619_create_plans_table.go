package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ticket-backend/app/facades"
)

type M20261010100619CreatePlansTable struct{}

// Signature The unique signature for the migration.
func (r *M20261010100619CreatePlansTable) Signature() string {
	return "20261010100619_create_plans_table"
}

// Up Run the migrations.
func (r *M20261010100619CreatePlansTable) Up() error {
	if !facades.Schema().HasTable("plans") {
		return facades.Schema().Create("plans", func(table schema.Blueprint) {
			table.ID()

			table.String("code", 50)
			table.Unique("code")

			table.String("name", 100)
			table.Integer("duration_days")
			table.Integer("max_watches")
			table.Boolean("auto_booking").Default(false)
			table.BigInteger("price_uzs").Nullable()
			table.BigInteger("price_stars").Nullable()
			table.Boolean("is_active").Default(true)

			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261010100619CreatePlansTable) Down() error {
	return facades.Schema().DropIfExists("plans")
}
