package dto

import "time"

type ProfileResponse struct {
	ID                    uint       `json:"id"`
	TelegramID            int64      `json:"telegram_id"`
	Username              string     `json:"username"`
	DisplayName           string     `json:"display_name"`
	PlanID                *uint      `json:"plan_id"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
	HasActiveSubscription bool       `json:"has_active_subscription"`
	CreatedAt             time.Time  `json:"created_at"`
}
