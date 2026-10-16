package routes

import (
	"ticket-backend/app/http/middleware"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/support"

	"ticket-backend/app/facades"
	"ticket-backend/app/http/controllers"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Static("public", "./public")

	userController := controllers.NewUserController()
	trainController := controllers.NewTrainController()
	stationController := controllers.NewStationController()
	healthController := controllers.NewHealthController()
	watchController := controllers.NewWatchController()

	facades.Route().Get("/api/health", healthController.Health)
	facades.Route().Get("/api/ready", healthController.Ready)

	facades.Route().Get("/users", userController.Index)

	facades.Route().Middleware(&middleware.TelegramAuth{}).Prefix("/api/v1").Group(func(router route.Router) {
		router.Get("/me", userController.Me)

		router.Get("/stations", stationController.Index)
		router.Get("/trains", trainController.Index)

		router.Get("/watches", watchController.Index)
		router.Post("/watches", watchController.Store)
		router.Patch("/watches/{id}/pause", watchController.Pause)
		router.Patch("/watches/{id}/resume", watchController.Resume)
		router.Delete("/watches/{id}", watchController.Destroy)
	})

	//facades.Route().
	//	Middleware(&middleware.TelegramAuth{}).
	//	Get("/api/v1/me", userController.Me)
	//
	//facades.Route().
	//	Middleware(&middleware.TelegramAuth{}).
	//	Get("/api/v1/watches", watchController.Index)
	//
	//facades.Route().
	//	Middleware(&middleware.TelegramAuth{}).
	//	Delete("/api/v1/watches/{id}", watchController.Destroy)
}
