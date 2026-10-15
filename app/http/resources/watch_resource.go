package resources

import (
	"ticket-backend/app/dto"
	"ticket-backend/app/models"
	"time"
)

type WatchResource struct {
	ID        uint                `json:"id"`
	UserID    uint                `json:"user_id"`
	From      string              `json:"from"`
	To        string              `json:"to"`
	TripDate  string              `json:"trip_date"`
	Filters   dto.WatchFiltersDTO `json:"filters"`
	IsActive  bool                `json:"is_active"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

func NewWatchResource(watch models.Watch) *WatchResource {
	response := &WatchResource{
		ID:       watch.ID,
		UserID:   watch.UserID,
		From:     watch.From,
		To:       watch.To,
		Filters:  watch.Filters,
		IsActive: watch.IsActive,
	}

	if watch.TripDate != nil {
		response.TripDate = watch.TripDate.ToDateString()
	}

	if watch.CreatedAt != nil {
		response.CreatedAt = watch.CreatedAt.StdTime().UTC()
	}

	if watch.UpdatedAt != nil {
		response.UpdatedAt = watch.UpdatedAt.StdTime().UTC()
	}

	return response
}

func NewWatchCollection(watches []models.Watch) []*WatchResource {
	response := make([]*WatchResource, len(watches))

	for i := range watches {
		response[i] = NewWatchResource(watches[i])
	}

	return response
}
