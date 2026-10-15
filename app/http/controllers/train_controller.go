package controllers

import (
	"context"
	"errors"
	"log"
	"strings"
	"ticket-backend/app/dto"
	"ticket-backend/app/facades"
	"ticket-backend/app/http/resources"
	"ticket-backend/app/services"

	"github.com/goravel/framework/contracts/http"
)

type TrainController struct {
	service *services.TrainService
}

func NewTrainController() *TrainController {
	return &TrainController{
		service: services.NewTrainService(),
	}
}

func (r *TrainController) Index(ctx http.Context) http.Response {
	query := dto.SearchTrainsDTO{
		From:     strings.TrimSpace(ctx.Request().Query("from")),
		To:       strings.TrimSpace(ctx.Request().Query("to")),
		TripDate: strings.TrimSpace(ctx.Request().Query("trip_date")),
		Lang:     strings.TrimSpace(ctx.Request().Query("lang", "uz")),
	}

	log.Printf("BEFORE Search: %+v", query)

	trains, err := r.service.Search(ctx, query)
	log.Printf(
		"AFTER Search: count=%d err=%v",
		len(trains),
		err,
	)
	if errors.Is(err, services.ErrInvalidTripDate) {
		return trainError(
			ctx,
			http.StatusUnprocessableEntity,
			"invalid_trip_date",
			"Select a valid travel date.",
		)
	}

	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Search trains: %v", err)

		if errors.Is(err, context.DeadlineExceeded) {
			return trainError(
				ctx,
				http.StatusGatewayTimeout,
				"train_search_timeout",
				"Train search timed out. Please try again.",
			)
		}

		return trainError(
			ctx,
			http.StatusBadGateway,
			"train_search_failed",
			"Could not load trains. Please try again.",
		)
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"data": resources.NewTrainCollection(trains),
	})
}

func trainError(
	ctx http.Context,
	status int,
	code string,
	message string,
) http.Response {
	return ctx.Response().Json(status, http.Json{
		"error": http.Json{
			"code":    code,
			"message": message,
		},
	})
}
