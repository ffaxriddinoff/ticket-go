package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ticket-backend/app/facades"
)

type M20261015114449AddMonitoringToWatchesTable struct{}

// Signature The unique signature for the migration.
func (r *M20261015114449AddMonitoringToWatchesTable) Signature() string {
	return "20261015114449_add_monitoring_to_watches_table"
}

// Up Run the migrations.
func (r *M20261015114449AddMonitoringToWatchesTable) Up() error {
	return facades.Schema().Table("watches", func(table schema.Blueprint) {
		table.Boolean("is_active").Default(true)
		table.TimestampTz("last_checked_at").Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20261015114449AddMonitoringToWatchesTable) Down() error {
	return facades.Schema().Table("watches", func(table schema.Blueprint) {
		table.DropColumn("is_active", "last_checked_at")
	})
}
