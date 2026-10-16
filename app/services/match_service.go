package services

import (
	"slices"
	"ticket-backend/app/dto"
	"ticket-backend/app/integrations/eticket"
)

type MatchService struct{}

func NewMatchService() *MatchService {
	return &MatchService{}
}

func (s *MatchService) Filter(trains []eticket.Train, filters dto.WatchFiltersDTO) []eticket.Train {
	result := s.filterTrainsByNumber(
		trains,
		filters.Trains,
	)

	result = s.filterWagonsByType(
		result,
		filters.WagonTypes,
	)

	return result
}

func (s *MatchService) filterTrainsByNumber(trains []eticket.Train, numbers []string) []eticket.Train {
	result := make([]eticket.Train, 0, len(trains))
	for _, train := range trains {
		if len(numbers) == 0 || slices.Contains(numbers, train.Number) {
			result = append(result, train)
		}
	}

	return result
}

func (s *MatchService) filterWagonsByType(trains []eticket.Train, types []dto.WagonType) []eticket.Train {
	if len(types) == 0 {
		return trains
	}

	result := make([]eticket.Train, 0, len(trains))

	for _, train := range trains {
		cars := make([]eticket.Car, 0, len(train.Cars))

		for _, car := range train.Cars {
			wagonType, known := car.WagonType()

			if known && slices.Contains(types, wagonType) {
				cars = append(cars, car)
			}
		}

		if len(cars) > 0 {
			train.Cars = cars
			result = append(result, train)
		}

	}

	return result
}

func (s *MatchService) filterEmptySeats(trains []eticket.Train) []eticket.Train {
	result := make([]eticket.Train, 0, len(trains))

	for _, train := range trains {
		cars := make([]eticket.Car, 0, len(train.Cars))

		for _, car := range train.Cars {
			if car.FreeSeats > 0 {
				cars = append(cars, car)
			}
		}

		if len(cars) > 0 {
			train.Cars = cars
			result = append(result, train)
		}
	}

	return result
}

func (s *MatchService) filterPosition(trains []eticket.Train, position dto.Position) []eticket.Train {
	result := make([]eticket.Train, 0, len(trains))

	for _, train := range trains {
		cars := make([]eticket.Car, 0, len(train.Cars))

		for _, car := range train.Cars {
			if position == dto.PositionLower && car.SeatDetail.Down > 0 {
				cars = append(cars, car)
			}
			if position == dto.PositionUpper && car.SeatDetail.Up > 0 {
				cars = append(cars, car)
			}
			if car.SeatDetail.Up+car.SeatDetail.Down > 0 {
				cars = append(cars, car)
			}
		}

		if len(cars) > 0 {
			train.Cars = cars
			result = append(result, train)
		}
	}

	return result
}
