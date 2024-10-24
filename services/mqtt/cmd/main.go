package main

import "C"
import (
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/repo/migrations"
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/usecase"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()
	// Add migrations to run
	app.Migrate(migrations.All())

	// Device REST
	err := app.AddRESTHandlers(&usecase.DeviceEntity{})
	if err != nil {
		app.Logger().Error(err, "failed to add Device REST handlers")
	}

	err = app.AddRESTHandlers(&usecase.ConstructionEntity{})
	if err != nil {
		app.Logger().Error(err, "failed to add Construction REST handlers")
	}

	err = app.AddRESTHandlers(&usecase.DeviceEntity{})
	if err != nil {
		app.Logger().Error(err, "failed to add Construction REST handlers")
	}

	app.Run()
}
