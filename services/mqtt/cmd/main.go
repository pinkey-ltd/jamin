package main

import (
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/adapter"
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/repo"
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/repo/migrations"
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/service"
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/usecase"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	// TODO: remove it
	app.Logger().Debug("SCHEMA is:", app.Config.Get("DB_SCHEMA_NAME"))
	// Add migrations to run
	app.Migrate(migrations.All())

	// Device REST
	err := app.AddRESTHandlers(&usecase.DeviceTemplateEntity{})
	if err != nil {
		app.Logger().Error(err, "failed to add Device REST handlers")
	}
	err = app.AddRESTHandlers(&usecase.ConstructionCURDEntity{})
	if err != nil {
		app.Logger().Error(err, "failed to add Construction REST handlers")
	}
	// Construction Actions
	podmanClient, err := adapter.NewPodmanAdapter()
	if err != nil {
		app.Logger().Error(err, "failed to initialize podman client")
		panic(err)
	}
	constructionService := service.NewConstruction(&repo.ConstructionPGRepository{}, nil, podmanClient)
	app.GET("/construction/{id}/test", constructionService.TestHandle)
	app.GET("/construction/{id}/start", constructionService.StartHandle)
	app.GET("/construction/{id}/stop", constructionService.StopHandle)
	app.GET("/construction/{id}/status", constructionService.StatusHandle)
	app.GET("/construction/{id}/remove", constructionService.RemoveHandle)
	app.GET("/construction/{id}/logs", constructionService.LogsHandle)

	err = app.AddRESTHandlers(&usecase.DeviceEntity{})
	if err != nil {
		app.Logger().Error(err, "failed to add Construction REST handlers")
	}

	// Docker driver handle
	app.GET("/cmd/docker/images", adapter.ImagesHandler)
	app.GET("/cmd/docker/containers", adapter.ContainersHandler)

	app.Run()
}
