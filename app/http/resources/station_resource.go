package resources

import (
	"ticket-backend/app/catalogs"
)

type StationResource struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func NewStationResource(
	station catalogs.Station,
	lang string,
) *StationResource {
	name := station.NameUZ

	switch lang {
	case "en":
		name = station.NameEN
	case "ru":
		name = station.NameRU
	}

	return &StationResource{
		Name: name,
		Code: string(station.Code),
	}
}

func NewStationCollection(
	stations []catalogs.Station,
	lang string,
) []*StationResource {
	response := make([]*StationResource, len(stations))

	for i, station := range stations {
		response[i] = NewStationResource(station, lang)
	}

	return response
}
