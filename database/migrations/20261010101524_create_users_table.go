package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ticket-backend/app/facades"
)

type M20261010101524CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20261010101524CreateUsersTable) Signature() string {
	return "20261010101524_create_users_table"
}

// Up Run the migrations.
func (r *M20261010101524CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			table.BigInteger("telegram_id")
			table.Unique("telegram_id")

			table.String("username", 100).Nullable()
			table.String("display_name", 200).Default("")

			table.BigInteger("plan_id").Nullable()
			table.Foreign("plan_id").References("id").On("plans")
			table.Index("plan_id").Name("users_plan_id_idx")

			table.TimestampTz("subscription_expires_at").Nullable()

			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261010101524CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
