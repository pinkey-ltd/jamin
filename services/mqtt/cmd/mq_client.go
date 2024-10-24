package main

import (
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/repo"
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/repo/migrations"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	// Add migrations to run
	app.Migrate(migrations.All())

	// MQTT subscribe
	app.Subscribe("openapi-test/299810258/tracker/telemetry", subscriber)

	app.Run()
}
func subscriber(ctx *gofr.Context) (err error) {
	telemetry := &repo.SchemaTest{}
	if err = ctx.Bind(telemetry); err != nil {
		ctx.Logger.Error(err, "failed to bind telemetry schema")
	}
	ctx.Logger.Infof("Received telemetry schema: %+v", telemetry)

	_, err = ctx.SQL.ExecContext(ctx, "INSERT INTO mqtt_telemetries (time,data) VALUES (NOW(),$1);", telemetry)
	if err != nil {
		ctx.Logger.Error(err, "failed to insert telemetry schema")
	}
	return nil
}
