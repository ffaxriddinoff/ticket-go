package services

import (
	"context"
	"errors"
	"ticket-backend/app/catalogs"
	"ticket-backend/app/dto"
	"ticket-backend/app/integrations/eticket"

	"github.com/goravel/framework/support/carbon"
)

var ErrInvalidTrainRoute = errors.New("invalid train route")

type TrainService struct {
	eticket *eticket.Client
}

func NewTrainService() *TrainService {
	return &TrainService{
		eticket: eticket.NewClient(),
	}
}

func (s *TrainService) Search(
	ctx context.Context,
	dto dto.SearchTrainsDTO,
) ([]eticket.Train, error) {
	if dto.From == dto.To ||
		!stationCodeExists(dto.From) ||
		!stationCodeExists(dto.To) {
		return nil, ErrInvalidTrainRoute
	}

	tripDate := carbon.ParseByFormat(
		dto.TripDate,
		"Y-m-d",
	)

	if !tripDate.IsValid() || tripDate.ToDateString() != dto.TripDate {
		return nil, ErrInvalidTripDate
	}

	today := carbon.Now().StartOfDay()

	if tripDate.Lt(today) {
		return nil, ErrInvalidTripDate
	}

	switch dto.Lang {
	case "uz", "ru", "en":
	default:
		dto.Lang = "uz"
	}

	request := eticket.SearchDto{
		Directions: eticket.DirectionsDto{
			Forward: eticket.Direction{
				Date:           dto.TripDate,
				DepStationCode: dto.From,
				ArvStationCode: dto.To,
			},
		},
		RouteType: "INTERCITY",
	}

	return s.eticket.SearchTrains(
		ctx,
		request,
		dto.Lang,
	)
}

func stationCodeExists(code string) bool {
	for _, station := range catalogs.Stations {
		if string(station.Code) == code {
			return true
		}
	}

	return false
}
