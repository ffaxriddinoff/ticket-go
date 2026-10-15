package controllers

import (
	"context"
	"ticket-backend/app/facades"
	"time"

	"github.com/goravel/framework/contracts/http"
)

type HealthController struct {
	// Dependent services
}

func NewHealthController() *HealthController {
	return &HealthController{
		// Inject services
	}
}

func (r *HealthController) Health(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusOK, http.Json{
		"status": "ok",
	})
}

func (r *HealthController) Ready(ctx http.Context) http.Response {
	queryCtx, cancel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cancel()

	var ready bool

	err := facades.Orm().WithContext(queryCtx).Query().
		Raw(`
            SELECT to_regclass('users') IS NOT NULL
                AND to_regclass('plans') IS NOT NULL
        `).
		Scan(&ready)

	if err != nil || !ready {
		return ctx.Response().Json(http.StatusServiceUnavailable, http.Json{
			"status": "not_ready",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status": "ready",
	})
}
