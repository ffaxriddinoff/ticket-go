package main

import (
	"ticket-backend/bootstrap"
	_ "time/tzdata"
)

func main() {
	app := bootstrap.Boot()

	app.Start()
}
