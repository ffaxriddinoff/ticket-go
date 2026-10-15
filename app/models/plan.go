package models

import (
	"github.com/goravel/framework/database/orm"
)

type Plan struct {
	orm.Model

	Code         string
	Name         string
	DurationDays int
	MaxWatches   int
	AutoBooking  bool
	PriceUZS     *int64
	PriceStars   *int64
	IsActive     bool
}
