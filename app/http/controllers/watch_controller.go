package controllers

import (
	"errors"
	"strconv"
	"ticket-backend/app/dto"
	"ticket-backend/app/facades"
	"ticket-backend/app/http/middleware"
	"ticket-backend/app/http/requests"
	"ticket-backend/app/http/resources"
	"ticket-backend/app/services"

	"github.com/goravel/framework/contracts/http"
)

type WatchController struct {
	userService  *services.UserService
	watchService *services.WatchService
}

func NewWatchController() *WatchController {
	return &WatchController{
		userService:  services.NewUserService(),
		watchService: services.NewWatchService(),
	}
}

func (c *WatchController) Index(ctx http.Context) http.Response {
	identity, ok := middleware.TelegramIdentityFromContext(ctx)
	if !ok {
		return watchError(
			ctx,
			http.StatusUnauthorized,
			"unauthorized",
			"Telegram authentication is required.",
		)
	}
	//identity := dto.TelegramIdentity{
	//	TelegramID:  590924106,
	//	Username:    "qodirovfaxriddin",
	//	DisplayName: "Faxa",
	//}

	user, err := c.userService.Profile(ctx.Context(), identity)
	if errors.Is(err, services.ErrInvalidIdentity) {
		return watchError(
			ctx,
			http.StatusUnauthorized,
			"unauthorized",
			"Invalid Telegram identity.",
		)
	}

	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Resolve watch owner: %v", err)

		return watchError(
			ctx,
			http.StatusInternalServerError,
			"internal_error",
			"Could not load user profile.",
		)
	}

	watches, err := c.watchService.List(ctx.Context(), user.ID)
	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("List watches: %v", err)

		return watchError(
			ctx,
			http.StatusInternalServerError,
			"internal_error",
			"Could not load watches.",
		)
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"data": resources.NewWatchCollection(watches),
	})
}

func (c *WatchController) Store(ctx http.Context) http.Response {
	identity, ok := middleware.TelegramIdentityFromContext(ctx)
	if !ok {
		return watchError(
			ctx,
			http.StatusUnauthorized,
			"unauthorized",
			"Telegram authentication is required.",
		)
	}
	//identity := dto.TelegramIdentity{
	//	TelegramID:  590924106,
	//	Username:    "qodirovfaxriddin",
	//	DisplayName: "Faxa",
	//}

	var request requests.CreateWatchRequest
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return watchError(
			ctx,
			http.StatusBadRequest,
			"invalid_request",
			"Could not read request data.",
		)
	}

	if validationErrors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"error": http.Json{
				"code":    "validation_failed",
				"message": "Please check the submitted fields.",
				"fields":  validationErrors.All(),
			},
		})
	}

	user, err := c.userService.Profile(ctx.Context(), identity)
	if errors.Is(err, services.ErrInvalidIdentity) {
		return watchError(
			ctx,
			http.StatusUnauthorized,
			"unauthorized",
			"Invalid Telegram identity.",
		)
	}

	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Resolve watch owner: %v", err)

		return watchError(
			ctx,
			http.StatusInternalServerError,
			"internal_error",
			"Could not load user profile.",
		)
	}

	watch, err := c.watchService.Create(
		ctx.Context(),
		user.ID,
		dto.CreateWatchDTO{
			From:     request.From,
			To:       request.To,
			TripDate: request.TripDate,
			Filters:  request.Filters,
		},
	)

	if errors.Is(err, services.ErrInvalidTripDate) {
		return watchError(
			ctx,
			http.StatusUnprocessableEntity,
			"invalid_trip_date",
			"Invalid trip date.",
		)
	}

	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Create watch: %v", err)

		return watchError(
			ctx,
			http.StatusInternalServerError,
			"internal_error",
			"Could not create watch.",
		)
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"data": resources.NewWatchResource(*watch),
	})
}

func (c *WatchController) Pause(ctx http.Context) http.Response {
	identity, ok := middleware.TelegramIdentityFromContext(ctx)
	if !ok {
		return watchError(
			ctx,
			http.StatusUnauthorized,
			"unauthorized",
			"Telegram authentication is required.",
		)
	}
	//identity := dto.TelegramIdentity{
	//	TelegramID:  590924106,
	//	Username:    "qodirovfaxriddin",
	//	DisplayName: "Faxa",
	//}

	watchID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 63)
	if err != nil || watchID == 0 {
		return watchError(
			ctx,
			http.StatusBadRequest,
			"invalid_watch_id",
			"Invalid watch ID.",
		)
	}

	user, err := c.userService.Profile(ctx.Context(), identity)
	if errors.Is(err, services.ErrInvalidIdentity) {
		return watchError(
			ctx,
			http.StatusUnauthorized,
			"unauthorized",
			"Invalid Telegram identity.",
		)
	}

	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Resolve watch owner: %v", err)

		return watchError(
			ctx,
			http.StatusInternalServerError,
			"internal_error",
			"Could not load user profile.",
		)
	}

	err = c.watchService.Pause(
		ctx.Context(),
		user.ID,
		uint(watchID),
	)

	if errors.Is(err, services.ErrWatchNotFound) {
		return watchError(
			ctx,
			http.StatusNotFound,
			"watch_not_found",
			"Watch not found.",
		)
	}
	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Pause watch: %v", err)

		return watchError(
			ctx,
			http.StatusInternalServerError,
			"internal_error",
			"Could not pause watch.",
		)
	}

	return ctx.Response().NoContent(http.StatusNoContent)
}

func (c *WatchController) Destroy(ctx http.Context) http.Response {
	identity, ok := middleware.TelegramIdentityFromContext(ctx)
	if !ok {
		return watchError(
			ctx,
			http.StatusUnauthorized,
			"unauthorized",
			"Invalid Telegram authentication.",
		)
	}

	watchID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 63)
	if err != nil || watchID == 0 {
		return watchError(
			ctx,
			http.StatusBadRequest,
			"invalid_watch_id",
			"Invalid watch ID.",
		)
	}

	user, err := c.userService.Profile(ctx.Context(), identity)
	if errors.Is(err, services.ErrInvalidIdentity) {
		return watchError(
			ctx,
			http.StatusUnauthorized,
			"unauthorized",
			"Invalid Telegram authentication.",
		)
	}
	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Resolve watch owner: %v", err)

		return watchError(
			ctx,
			http.StatusInternalServerError,
			"internal_error",
			"Could not load user profile.",
		)
	}

	err = c.watchService.Delete(ctx.Context(), user.ID, uint(watchID))
	if errors.Is(err, services.ErrWatchNotFound) {
		return watchError(
			ctx,
			http.StatusNotFound,
			"watch_not_found",
			"Watch not found.",
		)
	}
	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Delete watch: %v", err)

		return watchError(
			ctx,
			http.StatusInternalServerError,
			"internal_error",
			"Could not delete watch.",
		)
	}

	return ctx.Response().NoContent(http.StatusNoContent)
}

func watchError(ctx http.Context, status int, code string, message string) http.Response {
	return ctx.Response().Json(status, http.Json{
		"error": http.Json{
			"code":    code,
			"message": message,
		},
	})
}
