package controllers

import (
	"strings"
	"ticket-backend/app/http/resources"
	"ticket-backend/app/services"

	"github.com/goravel/framework/contracts/http"
)

type StationController struct {
	stationService *services.StationService
}

func NewStationController() *StationController {
	return &StationController{
		stationService: services.NewStationService(),
	}
}

func (c *StationController) Index(ctx http.Context) http.Response {
	lang := strings.ToLower(
		strings.TrimSpace(ctx.Request().Query("lang", "uz")),
	)

	stations := c.stationService.List()

	return ctx.Response().Json(http.StatusOK, http.Json{
		"data": resources.NewStationCollection(stations, lang),
	})
}
