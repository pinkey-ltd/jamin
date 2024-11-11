package main

import (
	"encoding/json"
	"github.com/pinkey-ltd/jamin/services/mqtt_client/migrations"
	"gofr.dev/pkg/gofr"
)

// ENV
//
// APP_NAME=mqtt_client
//
// DB_HOST=10.0.32.144
// DB_USER=<USER>
// DB_PASSWORD=<PASSWORD>
// DB_NAME=mqtt
// DB_PORT=5432
// DB_DIALECT=postgres
// DB_TABLE_NAME=<NAME>
//
// PUBSUB_BACKEND=MQTT
// MQTT_HOST=<HOST>
// MQTT_PORT=1883
// MQTT_PROTOCOL=tcp
// MQTT_USER=<USERNAME>
// MQTT_PASSWORD=<PASSWORD>
// MQTT_TOPIC=<TOPIC>
//
// METRICS_PORT=2122

var tableName string

func main() {
	app := gofr.New()

	// Custom ENV
	tableName = app.Config.GetOrDefault("DB_TABLE_NAME", "mqtt_telemetries")
	topic := app.Config.GetOrDefault("MQTT_TOPIC", "openapi-test/1")
	// Add migrations to run
	app.Migrate(migrations.All(tableName))

	// MQTT subscribe
	app.Subscribe(topic, subscriber)

	app.Run()
}

type RawJSON struct {
	json.RawMessage
}

func subscriber(ctx *gofr.Context) (err error) {
	telemetryRaw := &RawJSON{}
	if err = ctx.Bind(telemetryRaw); err != nil {
		ctx.Logger.Error(err, "failed to bind telemetry schema")
	}
	telemetry, err := json.Marshal(telemetryRaw)
	if err != nil {
		ctx.Logger.Error(err, "failed to marshal telemetry schema")
		return err
	}
	ctx.Logger.Infof("Received telemetry schema: %+v", string(telemetry))

	_, err = ctx.SQL.ExecContext(ctx, "INSERT INTO "+tableName+" (time,data) VALUES (NOW(),$1);", telemetry)
	if err != nil {
		ctx.Logger.Error(err, "failed to insert telemetry schema")
	}
	return nil
}
