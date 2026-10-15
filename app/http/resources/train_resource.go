package resources

import (
	"ticket-backend/app/integrations/eticket"
)

type TrainResource struct {
	Number        string `json:"number"`
	Brand         string `json:"brand"`
	DepartureDate string `json:"departure_date"`
	ArrivalDate   string `json:"arrival_date"`
}

func NewTrainResource(train eticket.Train) *TrainResource {
	return &TrainResource{
		Number:        train.Number,
		Brand:         train.Brand,
		DepartureDate: train.DepartureDate,
		ArrivalDate:   train.ArrivalDate,
	}
}

func NewTrainCollection(trains []eticket.Train) []*TrainResource {
	response := make([]*TrainResource, len(trains))

	for i, train := range trains {
		response[i] = NewTrainResource(train)
	}

	return response
}
