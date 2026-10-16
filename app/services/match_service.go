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

	return result
}

func (s *MatchService) filterTrainsByNumber(trains []eticket.Train, numbers []string) []eticket.Train {
	result := make([]eticket.Train, 0, len(trains))
	for _, train := range trains {
		if len(numbers) > 0 && !slices.Contains(numbers, train.Number) {
			continue
		}

		result = append(result, train)
	}

	return result
}
