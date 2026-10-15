package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type User struct {
	orm.Model

	TelegramID            int64
	Username              *string
	DisplayName           string
	PlanID                *uint
	SubscriptionExpiresAt *carbon.DateTime

	Plan *Plan `gorm:"foreignKey:PlanID"`
}
