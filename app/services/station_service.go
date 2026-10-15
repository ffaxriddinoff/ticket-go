package services

import (
	"slices"
	"ticket-backend/app/catalogs"
)

type StationService struct{}

func NewStationService() *StationService {
	return &StationService{}
}

func (s *StationService) List() []catalogs.Station {
	return slices.Clone(catalogs.Stations)
}
