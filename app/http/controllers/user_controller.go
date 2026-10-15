package controllers

import (
	"errors"
	"ticket-backend/app/facades"
	"ticket-backend/app/http/middleware"
	"ticket-backend/app/http/resources"
	"ticket-backend/app/services"

	"github.com/goravel/framework/contracts/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

func (r *UserController) Index(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}

func (c *UserController) Me(ctx http.Context) http.Response {
	facades.Log().WithContext(ctx.Context()).
		With(map[string]any{
			"method": ctx.Request().Method(),
			"path":   ctx.Request().Path(),
			"input":  ctx.Request().All(),
		}).
		Info("Incoming request")

	identity, ok := middleware.TelegramIdentityFromContext(ctx)
	if !ok {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": http.Json{
				"code":    "unauthorized",
				"message": "Telegram authentication is required.",
			},
		})
	}

	user, err := c.userService.Profile(ctx.Context(), identity)
	if errors.Is(err, services.ErrInvalidIdentity) {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": http.Json{
				"code":    "unauthorized",
				"message": "Invalid Telegram identity.",
			},
		})
	}

	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Load user profile: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": http.Json{
				"code":    "internal_error",
				"message": "Could not load user profile.",
			},
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"data": resources.NewUserResource(*user),
	})
}
