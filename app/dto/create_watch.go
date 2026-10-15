package dto

type CreateWatchDTO struct {
	From     string
	To       string
	TripDate string
	Filters  WatchFiltersDTO
}
