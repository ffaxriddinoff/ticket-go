package models

import (
	"ticket-backend/app/dto"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type Watch struct {
	orm.Model

	UserID        uint
	From          string
	To            string
	TripDate      *carbon.Date
	Filters       dto.WatchFiltersDTO `gorm:"type:jsonb;serializer:json"`
	IsActive      bool
	LastCheckedAt *carbon.DateTime

	User *User `gorm:"foreignKey:UserID"`
}
