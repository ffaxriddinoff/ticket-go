package resources

import (
	"ticket-backend/app/dto"
	"ticket-backend/app/models"
	"time"
)

func NewUserResource(user models.User) *dto.ProfileResponse {
	response := &dto.ProfileResponse{
		ID:          user.ID,
		TelegramID:  user.TelegramID,
		DisplayName: user.DisplayName,
		PlanID:      user.PlanID,
	}

	if user.Username != nil {
		response.Username = *user.Username
	}

	if user.CreatedAt != nil {
		response.CreatedAt = user.CreatedAt.StdTime().UTC()
	}

	if user.SubscriptionExpiresAt != nil {
		expiresAt := user.SubscriptionExpiresAt.StdTime().UTC()
		response.SubscriptionExpiresAt = &expiresAt
		response.HasActiveSubscription =
			user.PlanID != nil && expiresAt.After(time.Now())
	}

	return response
}
